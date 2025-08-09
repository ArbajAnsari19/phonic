# Services

This directory contains all microservices for the Phonic AI Calling Agent.

## Services
- `gateway/` - WebRTC gateway handling browser/telephony connections
- `stt-client/` - Client for connecting to Moshi STT server
- `tts-client/` - Client for connecting to Moshi TTS server
- `orchestrator/` - Main coordination service (STT → LLM → TTS)
- `session/` - Session management and state persistence

Each service contains:
- `main.go` - Service entry point
- `server.go` - gRPC server implementation
- `config.go` - Service configuration
- `handlers/` - Request handlers
- `client/` - gRPC client code (if needed)
