# echallan-services

Go implementation of the DIGIT echallan-services microservice. Structured according to the Go Gold Standard Clean Architecture Guidelines.

## 🏗️ Repository Layout

- `cmd/`: Entrypoints.
- `configs/`: App configurations.
- `internal/`: Encapsulated codebase.
  - `domain/`: Business entities and API request/response contracts.
  - `repository/`: Database interactions.
  - `service/`: Core application use-cases.
  - `transport/`: Network layers (HTTP/Kafka).
  - `util/`: Helper utilities.
  - `validator/`: Payload validation.
- `deployments/`: Deploy configurations (Docker, Kubernetes, Helm).
- `docs/`: Architecture design guides.
- `scripts/`: Dev & automation scripts.
- `test/`: Integration and unit test fixtures.

## 🛠️ Usage

### Build
```bash
make build
```

### Run
```bash
make run
```

### Test
```bash
make test
```
