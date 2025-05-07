# Trading System

A simple trading system with API server and worker process.

## Features

- POST /trades - Submit new trades
- GET /stats/{acc} - Get account statistics
- GET /healthz - Health check endpoint

## How to run

### Local development

1. Build and run:
```bash
make run
```

2. Test with curl:
```bash
# Submit trade
curl -X POST http://localhost:8080/trades \
    -H 'Content-Type: application/json' \
    -d '{"account":"123","symbol":"EURUSD","volume":1.0,"open":1.1000,"close":1.1050,"side":"buy"}'

# Get stats
curl http://localhost:8080/stats/123
```

### Docker

1. Start services:
```bash
make docker-up
```

2. Stop services:
```bash
make docker-down
```

## Testing

Run tests:
```bash
make test
```

## CI/CD

GitLab CI pipeline will:
- Run go vet and tests
- Check race conditions
- Verify API endpoints