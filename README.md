# 🎵 Phonic AI Calling Agent

A production-grade AI calling agent backend built with Go and Rust, designed to compete with Vapi AI and Retell AI.

## ✨ Features

- **Sub-300ms Latency**: Real-time audio streaming with WebRTC
- **Barge-in Support**: Voice activity detection for natural conversations  
- **Enterprise Scale**: Microservices architecture with Kubernetes support
- **High-Performance STT/TTS**: Powered by Kyutai Moshi (Rust servers)
- **Production Ready**: Full observability, monitoring, and deployment configs

## 🏗️ Architecture

```
Browser/Phone → WebRTC Gateway → STT Client → Moshi STT Server
                      ↓                            ↑
                 Orchestrator ←→ Session Manager ←→ Redis
                      ↓                            ↓
              TTS Client → Moshi TTS Server → Audio Output
                      ↓
                 Mock LLM ←→ PostgreSQL + MinIO
```

## 📁 Project Structure

```
phonic/
├── proto/              # gRPC service definitions
├── services/           # Microservices (Go)
│   ├── gateway/        # WebRTC gateway
│   ├── stt-client/     # Moshi STT client
│   ├── tts-client/     # Moshi TTS client
│   ├── orchestrator/   # Main coordination service
│   └── session/        # Session management
├── pkg/                # Shared Go packages
├── cmd/                # Application binaries
├── infra/              # Docker & Kubernetes configs
├── configs/            # Environment configurations
├── docs/               # Documentation
├── scripts/            # Build and deployment scripts
└── tools/              # Development tools
```

## 🚀 Quick Start

```bash
# Start development environment
make up

# Build all services
make build

# Run tests
make test

# View logs
make logs
```

## 🛠️ Development

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- protoc with Go plugins
- Make

### Setup
```bash
git clone https://github.com/ArbajAnsari19/phonic.git
cd phonic
make setup
make up
```

## 📊 Monitoring

- **Metrics**: Prometheus + Grafana
- **Tracing**: OpenTelemetry  
- **Logging**: Structured JSON with Zap
- **Health Checks**: Built-in health endpoints

## 🔧 Configuration

Environment-specific configs in `configs/`:
- `dev/` - Development settings
- `staging/` - Staging environment  
- `prod/` - Production settings

## 📖 Documentation

- [API Documentation](docs/api/)
- [Architecture Overview](docs/architecture/)
- [Deployment Guide](docs/deployment.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
