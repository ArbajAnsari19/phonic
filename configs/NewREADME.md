# Roadmap for building a production-grade AI Calling Agent backend:

## Phase 1: Foundation & Environment Setup (Steps 1-10)
1. **Initialize Go workspace** - Set up Go module and basic project structure
2. **Install development dependencies** - Go, protoc, Docker, Make tools
3. **Create monorepo structure** - Set up `/proto`, `/services`, `/infra`, `/configs`, `/docs` directories
4. **Initialize Git repository** - Version control with proper .gitignore
5. **Create Makefile** - Build automation and common commands
6. **Set up Docker development environment** - Base Dockerfile and docker-compose structure
7. **Configure logging infrastructure** - Structured logging with zap
8. **Set up configuration management** - Environment-based config with viper
9. **Create basic health check endpoints** - HTTP health checks for all services
10. **Test foundation setup** - Verify all tools and basic structure work

## Phase 2: Protocol Definitions & Code Generation (Steps 11-20)
11. **Define STT service protobuf** - Streaming speech-to-text gRPC definitions
12. **Define TTS service protobuf** - Streaming text-to-speech gRPC definitions
13. **Define Session service protobuf** - Session management and state protobuf
14. **Define Orchestrator protobuf** - Main orchestration service definitions
15. **Install protoc-gen-go plugins** - gRPC code generation tools
16. **Generate Go code from protos** - Create gRPC client/server stubs
17. **Create shared gRPC utilities** - Common connection, retry, and error handling
18. **Set up proto validation** - Input validation for all gRPC messages
19. **Create proto documentation** - Auto-generated API documentation
20. **Test proto compilation** - Verify all generated code compiles

## Phase 3: Database & Storage Setup (Steps 21-28)
21. **Set up PostgreSQL schema** - Database for call metadata and user data
22. **Create database migrations** - Version-controlled schema management
23. **Set up Redis configuration** - Session state and caching layer
24. **Configure MinIO/S3 setup** - Audio file storage and logging
25. **Create database models** - Go structs for data persistence
26. **Implement database connections** - Connection pooling and health checks
27. **Create data access layer** - Repository pattern for database operations
28. **Test database connectivity** - Verify all storage systems work

## Phase 4: STT Service Implementation (Steps 29-36)
29. **Create STT service skeleton** - Basic gRPC server structure
30. **Implement Kyutai Unmute integration** - Connect to Rust STT container
31. **Set up audio streaming pipeline** - WebRTC audio → STT processing
32. **Implement partial transcription streaming** - Real-time transcription results
33. **Add audio format conversion** - Handle different audio codecs
34. **Implement barge-in detection** - Voice activity detection for interruptions
35. **Add STT error handling** - Graceful degradation and retry logic
36. **Test STT service** - End-to-end audio transcription testing

## Phase 5: TTS Service Implementation (Steps 37-44)
37. **Create TTS service skeleton** - Basic gRPC server structure
38. **Implement Kyutai Unmute TTS integration** - Connect to Rust TTS container
39. **Set up text streaming pipeline** - Text → audio generation
40. **Implement audio streaming output** - Real-time audio generation
41. **Add voice configuration** - Voice selection and parameters
42. **Implement audio quality optimization** - Compression and format optimization
43. **Add TTS error handling** - Fallback and retry mechanisms
44. **Test TTS service** - End-to-end text-to-speech testing

## Phase 6: WebRTC Gateway Service (Steps 45-52)
45. **Create Gateway service skeleton** - HTTP server with WebRTC support
46. **Implement WebRTC peer connection** - Browser audio input handling
47. **Set up SDP negotiation** - WebRTC signaling protocol
48. **Implement ICE candidate handling** - NAT traversal and connectivity
49. **Create audio track processing** - Extract audio from WebRTC streams
50. **Implement audio output streaming** - Send processed audio back to client
51. **Add connection state management** - Handle WebRTC lifecycle
52. **Test WebRTC connectivity** - Browser to Gateway audio flow

## Phase 7: Session Management Service (Steps 53-58)
53. **Create Session service skeleton** - Session lifecycle management
54. **Implement session creation** - New call session initialization
55. **Add session state persistence** - Redis-backed session storage
56. **Implement session middleware** - Request routing and state management
57. **Add session cleanup** - Automatic session expiration and cleanup
58. **Test session management** - Session CRUD operations

## Phase 8: Orchestrator Service (Steps 59-66)
59. **Create Orchestrator service skeleton** - Main coordination service
60. **Implement STT client integration** - Connect to STT service
61. **Implement TTS client integration** - Connect to TTS service
62. **Create mock LLM service** - Echo bot for initial testing
63. **Implement conversation flow** - STT → LLM → TTS pipeline
64. **Add barge-in handling** - Interrupt and resume logic
65. **Implement streaming coordination** - Real-time audio processing
66. **Test orchestrator pipeline** - End-to-end conversation flow

## Phase 9: Observability & Monitoring (Steps 67-70)
67. **Set up Prometheus metrics** - Service metrics and monitoring
68. **Implement OpenTelemetry tracing** - Distributed request tracing
69. **Create monitoring dashboards** - Grafana dashboard configuration
70. **Add performance monitoring** - Latency and throughput metrics

## Phase 10: Final Integration & Testing (Steps 71-73)
71. **Complete Docker composition** - All services running together
72. **End-to-end system testing** - Full pipeline validation
73. **Production readiness checklist** - Security, performance, deployment prep

---


## Key Architecture Overview

Before we start, here's what we're building:

```
Browser/Phone → WebRTC Gateway → STT Service (Kyutai) → Orchestrator → TTS Service (Kyutai) → Audio Output
                       ↓                                        ↓
                 Session Manager ←→ Redis                Mock LLM (Echo Bot)
                       ↓
                 PostgreSQL + MinIO
```

**Why this architecture?**
- **Microservices**: Each component can scale independently
- **Streaming**: Sub-300ms latency through real-time audio streaming
- **Barge-in**: STT can interrupt TTS when user speaks
- **Stateful**: Sessions persist across network issues
- **Observable**: Full metrics, logging, and tracing

