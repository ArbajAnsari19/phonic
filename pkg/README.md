# Shared Packages

This directory contains reusable Go packages shared across services.

## Packages
- `utils/` - Common utility functions
- `middleware/` - gRPC and HTTP middleware
- `models/` - Shared data models and structs

## Usage
Import packages using the full module path:
```go
import "github.com/ArbajAnsari19/phonic/pkg/utils"
import "github.com/ArbajAnsari19/phonic/pkg/middleware"
```
