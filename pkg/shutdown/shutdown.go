// Package shutdown provides graceful shutdown functionality for Phonic AI Calling Agent
package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ArbajAnsari19/phonic/pkg/logger"
)

// Hook represents a cleanup function to run during shutdown
type Hook func(ctx context.Context) error

// Manager manages graceful shutdown of the service
type Manager struct {
	hooks   []Hook
	timeout time.Duration
	logger  *logger.Logger
	mu      sync.Mutex
}

// NewManager creates a new shutdown manager
func NewManager(timeout time.Duration, log *logger.Logger) *Manager {
	return &Manager{
		hooks:   make([]Hook, 0),
		timeout: timeout,
		logger:  log,
	}
}

// AddHook adds a shutdown hook
func (m *Manager) AddHook(hook Hook) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hooks = append(m.hooks, hook)
}

// Listen starts listening for shutdown signals
func (m *Manager) Listen() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigChan
	m.logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	m.Shutdown()
}

// Shutdown executes all shutdown hooks with timeout
func (m *Manager) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	m.logger.Info("Starting graceful shutdown", zap.Duration("timeout", m.timeout))

	m.mu.Lock()
	hooks := make([]Hook, len(m.hooks))
	copy(hooks, m.hooks)
	m.mu.Unlock()

	// Execute hooks in reverse order (LIFO)
	for i := len(hooks) - 1; i >= 0; i-- {
		hook := hooks[i]
		if err := hook(ctx); err != nil {
			m.logger.Error("Shutdown hook failed", zap.Error(err), zap.Int("hook_index", i))
		}
	}

	m.logger.Info("Graceful shutdown completed")
}

// WaitForShutdown waits for shutdown signal and executes hooks
func (m *Manager) WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigChan
	m.logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	m.Shutdown()
	os.Exit(0)
}
