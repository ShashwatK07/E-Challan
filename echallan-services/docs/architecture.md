# Architecture Overview

This microservice uses Go Gold Standard Layered Architecture:

1. **Transport Layer** (`internal/transport/`): Maps HTTP paths to logic.
2. **Validator Layer** (`internal/validator/`): Checks payload sanitization rules.
3. **Service Layer** (`internal/service/`): Orchestrates entities and use cases.
4. **Domain Layer** (`internal/domain/`): Standard business objects and model schemas.
5. **Repository Layer** (`internal/repository/`): Interface with DB or other API endpoints.
