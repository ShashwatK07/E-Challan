# Architecture Overview

Flow of data inside the echallan-calculator microservice:

```
Request Info (Challan JSON)
        ↓
Controller (Gin handler binding)
        ↓
CalculationService (Generates estimate amount per tax head)
        ↓
DemandService (Prepares billing-service Demand Payload)
        ↓
Repository (HTTP call to /billing-service/demands/_create)
```
