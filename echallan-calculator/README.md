# echallan-calculator

Go implementation of the DIGIT echallan-calculator microservice. Structured according to the Go Gold Standard Clean Architecture Guidelines.

## 🏗️ Repository Layout

- `cmd/`: Application entry point.
- `configs/`: App configurations.
- `internal/`: Encapsulated codebase.
  - `domain/`: Business entities and API request/response contracts.
  - `repository/`: Outbound HTTP repository interfaces to other DIGIT services.
  - `service/`: Tax estimation calculations.
  - `transport/`: Gin routers & controllers mapping requests.
  - `util/`: Stateless helper utilities.

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
