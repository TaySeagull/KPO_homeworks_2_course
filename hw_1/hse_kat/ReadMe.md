# Homework 1

### Entities

#### Interface

**Vehicle**
```
├── Name
├── InventoryNumber
├── Simplicity
└── Status
```

**Equipment**
```
├── Name
├── InventoryNumber
└── Status
```

**Dock Station**
```
├── name
├── inventory number
├── energy consumption
├── vehicles
└── equipment
```

**Inventory Item**
```
├── name
└── inventory_number
```

**Energy Consumer**
```
├── name
└── energy_consumption
```

**Inspector**
```
└── Inspect
```

```
hse_kat/
├── go.mod
├── cmd/
│   ├── main.go
│   ├── menu.go
│   ├── input.go
│   ├── print.go
│   └── registry.go
├── container/
│   └── container.go
├── internal/
│   ├── domain/
│   │   ├── status.go
│   │   ├── inventory.go
│   │   ├── energy.go
│   │   └── items.go
│   └── entities/
│       ├── scooter.go
│       ├── bike.go
│       └── helmet.go
└── service/
    ├── park.go
    ├── inspector.go
    └── service_center.go
```