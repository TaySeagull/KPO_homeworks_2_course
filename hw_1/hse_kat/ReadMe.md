# Homework 1
Консольное приложение для учёта парка сервиса кикшеринга «ВышКат»

Console application for fleet management of the "HSE_kat" kick-sharing service

---

# Launch

Для запуска необходим Go версии 1.25 или совместимой версии

Запустить приложение:

```bash
go run ./cmd
```

---

Go version 1.25 or a compatible version is required to run the application

Run the application:

```bash
go run ./cmd
```

# Tests

Запуск всех тестов:

```bash
go test ./...
```

---

Run the tests:

```bash
go test ./...
```

## Structure

```
hse_kat/
├── cmd/
│   ├── input.go
│   ├── main.go
│   ├── menu.go
│   ├── print.go
│   └── registry.go
│
├── container/
│   └── container.go
│
├── internal/
│   ├── domain/
│   │   ├── energy.go
│   │   ├── inventory.go
│   │   ├── items.go
│   │   └── status.go
│   │
│   └── entities/
│       ├── bike.go
│       ├── helmet.go
│       └── scooter.go
│
├── service/
│   ├── inspector.go
│   ├── park.go
│   └── service_center.go
│
├── tests/
│   └── park_test.go
│
├── go.mod
└── ReadMe.md
```

## Details

`internal/domain` содержит / contains:

- `InventoryItem` - инвентарный интерфейс / inventory interface
- `EnergyConsumer` - иинтерфейс потребления энергии / energy consumption interface
- `Vehicle` - транспорт / transport
- `Equipment` - оборудование / equipment
- `Status` - состояние объекта / object status

`internal/entities` содержит конкретные сущности / contains user's entities:

- `Scooter` - самокат / scooter
- `Bike` - велосипед / bike
- `Helmet` - шлем / helmet

`service` содержит бизнес-логику / contains business logic:

- `Inspector` и `BasicInspector` - технический осмотр / technical inspection
- `ServiceCenter` - использование инспектора / uses the inspector (for DIP)
- `Park` - сам парк-станция / station for transpost

`container` - собирает зависимости / Composition Root

`cmd` - консольный интерфейс пользователя / console module

`tests` - тесты парка / test for the whole park


---

#  Architecture

В качестве основной модели использованы интерфейсы, описывающие поведение объектов

Interfaces describing object behavior are used as the primary model

## Inventory items

`internal/domain/inventory.go`

```go
type InventoryItem interface {
	Name() string
	InventoryNumber() string
}
```

Он позволяет работать с разными объектами через единую абстракцию инвентарного объекта

Например, и `Scooter`, и `Bike`, и `Helmet` являются `InventoryItem`

Благодаря этому парк может сформировать единый список инвентаря:

```go
[]domain.InventoryItem
```

без необходимости знать конкретный тип каждого объекта

---

It allows working with various objects through a unified inventory object abstraction

For example, `Scooter`, `Bike`, `Helmet` are all `InventoryItem`

Because of this, the park can compile a unified inventory list:

```go
[]domain.InventoryItem
```

without needing to know the specific type of each object

## Transport

`internal/domain/items.go`

```go
type Vehicle interface {
	InventoryItem
	Simplicity() int
	Status() Status
}
```

Транспорт должен иметь инвентарные данные, характеристику простоты для новичка и состояние

Сейчас реализованы:

- `Scooter`
- `Bike`

Добавление нового вида транспорта не требует изменения `Park`, если новый тип реализует `domain.Vehicle`

---

The vehicle must have inventory data, a rating of how easy it is for a beginner to use, and a status

The following have been implemented:

- `Scooter`
- `Bike`

Adding a new type of transport does not require modifying `Park` if the new type implements `domain.Vehicle`

## Equipment

`internal/domain/items.go`

```go
type Equipment interface {
	InventoryItem
	Status() Status
}
```

Сейчас реализованы:

- `Helmet`

---

The following have been implemented:

- `Helmet`

## Energy consumers

`internal/domain/energy.go`

```go
type EnergyConsumer interface {
	EnergyConsumption() float64
}
```

Это позволяет не заставлять все сущности хранить энергопотребление

Например, `Scooter` реализует `EnergyConsumer`, а `Bike` в текущей модели не реализует его

При расчёте энергии парк проверяет, поддерживает ли конкретный объект данный интерфейс
Таким образом, добавление нового энергопотребляющего устройства не требует изменения алгоритма расчёта энергии

---

This avoids the need for all entities to store energy consumption data

For example, `Scooter` implements `EnergyConsumer`, whereas `Bike` does not implement it in the current model

When calculating energy, the fleet checks whether a specific object supports this interface
Thus, adding a new energy-consuming device does not require modifying the energy calculation alorithm


## Techical inspector

Перед добавлением транспорта в парк вызывается технический осмотр

Цепочка выглядит следующим образом:

```
Park
  ↓
ServiceCenter
  ↓
Inspector
  ↓
BasicInspector
```

`Park` не выполняет проверку самостоятельно

Вместо этого он обращается к:

```go
park.serviceCenter.Inspect(vehicle)
```

`ServiceCenter` использует абстракцию:

```go
type Inspector interface {
	Inspect(vehicle domain.Vehicle) bool
}
```

Текущая реализация `BasicInspector` принимает транспорт, если он имеет статус:

```go
domain.Available
```

Транспорт со статусом:

```go
domain.NeedsRepair
```

отклоняется

---

A technical inspection is triggered before a vehicle is added to the fleet

The chain looks like this:

```
Park
  ↓
ServiceCenter
  ↓
Inspector
  ↓
BasicInspector
```

`Park` does not perform the inspection itself

Instead, it calls:

```go
park.serviceCenter.Inspect(vehicle)
```

`ServiceCenter` uses an abstraction:

```go
type Inspector interface {
	Inspect(vehicle domain.Vehicle) bool
}
```

The current implementation of `BasicInspector` accepts a vehicle if it has the status:

```go
domain.Available
```

A vehicle with the status:

```go
domain.NeedsRepair
```

is rejected

# Dependency Injection

В проекте используется Dependency Injection

`ServiceCenter` получает `Inspector` через конструктор:

```go
func NewServiceCenter(inspector Inspector) *ServiceCenter
```

При этом `ServiceCenter` зависит не от конкретного `BasicInspector`, а от интерфейса `Inspector`.

Парк получает готовый `ServiceCenter`:

```go
func NewPark(serviceCenter *ServiceCenter) *Park
```

Сборка зависимостей происходит в `container/container.go`:

```go
inspector := &service.BasicInspector{}
serviceCenter := service.NewServiceCenter(inspector)
park := service.NewPark(serviceCenter)
```

Таким образом, классы и сервисы не создают свои зависимости самостоятельно

`container` является Composition Root

---

The project uses Dependency Injection

`ServiceCenter` receives an `Inspector` via its constructor:

```go
func NewServiceCenter(inspector Inspector) *ServiceCenter
```

Notably, `ServiceCenter` depends on the `Inspector` interface rather than a concrete `BasicInspector`

`Park` receives a fully initialized `ServiceCenter`:

```go
func NewPark(serviceCenter *ServiceCenter) *Park
```

Dependency assembly takes place in `container/container.go`:

```go
inspector := &service.BasicInspector{}
serviceCenter := service.NewServiceCenter(inspector)
park := service.NewPark(serviceCenter)
```

Thus, classes and services do not instantiate their own dependencies

The `container` acts as the Composition Root

# SOLID

## Single Responsibility Principle

Основные обязанности разделены:

- `Park` управляет парком и выполняет операции над его содержимым
- `ServiceCenter` отвечает за взаимодействие с техническим осмотром
- `BasicInspector` содержит конкретное правило технического осмотра
- сущности `Scooter`, `Bike` и `Helmet` хранят состояние конкретных объектов
- `cmd` отвечает за консольный интерфейс
- `container` отвечает за сборку зависимостей

Таким образом, изменение правила технического осмотра не требует изменения логики `Park`

---

Key responsibilities are distributed as follows:

- `Park` manages the fleet and performs operations on its contents
- `ServiceCenter` handles interactions related to technical inspections
- `BasicInspector` encapsulates a specific technical inspection rule
- The `Scooter`, `Bike`, and `Helmet` entities store the state of specific objects
- `cmd` handles the console interface
- `container` manages dependency assembly

Thus, modifying a technical inspection rule does not require changing the `Park` logic

## Open/Closed Principle

Парк открыт для добавления новых типов транспорта. `Park` при этом изменять не потребуется.

То же самое относится к `EnergyConsumer`:
новый объект может реализовать этот интерфейс и автоматически участвовать в расчёте энергопотребления

---

The park is open to the addition of new types of transport. `Park` class itself does not need to be modified.

The same applies to `EnergyConsumer`:
a new object can implement this interface and automatically participate in energy consumption calculations.

## Liskov Substitution Principle

`Park` работает с `domain.Vehicle`, а не с конкретными `Scooter` или `Bike`

Поэтому любой корректный `Vehicle` может быть передан в:

```go
AddVehicle(vehicle domain.Vehicle)
```

без изменения поведения парка.

Аналогично различные сущности могут использоваться как `InventoryItem`

---

`Park` works with `domain.Vehicle` rather than specific `Scooter` or `Bike` types

Therefore, any valid `Vehicle` can be passed to:

```go
AddVehicle(vehicle domain.Vehicle)
```

without altering the park's behavior

Similarly, various entities can be used as `InventoryItem`

## Interface Segregation Principle

Интерфейсы сделаны небольшими и разделены по назначению

Объект не обязан реализовывать методы, которые ему не нужны

---

Interfaces are designed to be small and are separated by purpose

An object is not required to implement methods it does not need

## Dependency Inversion Principle

Любимый принцип

`ServiceCenter` зависит от:

```go
Inspector
```

а не от:

```go
BasicInspector
```

Конкретная реализация передаётся извне через конструктор

---

`ServiceCenter` depends on:

```go
Inspector
```

rather than on:

```go
BasicInspector
```

The concrete implementation is passed in via the constructor

# 6. Which principle is applied the most?

Наиболее явно в проекте применяется **Dependency Inversion Principle**.

Это видно в связке:

```text
ServiceCenter → Inspector
```

Основная идея DIP заключается в том, что высокоуровневые компоненты не должны напрямую зависеть от конкретных реализаций низкоуровневых компонентов. Вместо этого они должны зависеть от абстракций. Конкретные реализации связываются с системой отдельно, в точке сборки приложения

Это прослеживается в архитектуре всего приложения

--- 

The **Dependency Inversion Principle** is most clearly applied in the project

This is evident in the following relationship:

```text
ServiceCenter → Inspector
```

The core idea of ​​DIP is that high-level components should not depend directly on concrete implementations of low-level components; instead, they should depend on abstractions. Concrete implementations are integrated into the system separately, at the application assembly point

This pattern is reflected throughout the application's architecturе

---

# Which principle is deliberately limited?

В проекте Dependency Inversion Principle реализован не полностью

`Park` хранит конкретный `ServiceCenter`:

```go
type Park struct {
    vehicles      []domain.Vehicle
    equipment     []domain.Equipment
    serviceCenter *ServiceCenter
}
```

Теоретически можно было бы создать ещё один интерфейс для сервисного центра. Но в программе используется один сервисный центр с одной реализацией, и требований к нескольким разным сервисным центрам нет. Поэтому было принято решение ограничиться таким решением 

---

The Dependency Inversion Principle is not fully implemented in the project

`Park` holds a concrete `ServiceCenter`:

```go
type Park struct {
    vehicles      []domain.Vehicle
    equipment     []domain.Equipment
    serviceCenter *ServiceCenter
}
```

Theoretically, one could create an additional interface for the service center. However, the program uses a single service center with a single implementation, and there are no requirements for multiple different service centers. Therefore, the decision was made to stick with this approach

# Expansion of the app

## New vehicle

Достаточно создать новую сущность, реализующую:

```go
domain.Vehicle
```

Если транспорт потребляет энергию, он дополнительно может реализовать:

```go
domain.EnergyConsumer
```

`Park` изменять не потребуется

--- 

Сreate a new entity that implements:

```go
domain.Vehicle
```

If the vehicle consumes energy, it can additionally implement:

```go
domain.EnergyConsumer
```

There is no need to modify `Park`

## New employee

Сотрудник не является инвентарным объектом и в текущей модели не относится к `Vehicle` или `Equipment`

Поэтому потребуется создать отдельную сущность, написать логику, добавить соответствующий сервис, расширить Composition Root

--- 

An employee is not an inventory item and, in the current model, does not fall under `Vehicle` or `Equipment`

Therefore, it will be necessary to create a separate entity, implement the logic, add the corresponding service, and extend the Composition Root

## New spare parts warehouse

Склад также является отдельной предметной областью

Можно создать отдельный сервис:

```text
Warehouse
```

и соответствующие сущности для запчастей

Текущая логика `Park` при этом не должна измениться, если склад не требует изменения существующего поведения парка

Новая зависимость будет добавлена в Composition root

---

The warehouse also represents a distinct domain area

A separate service can be created:

```text
Warehouse
```

along with the corresponding entities for spare parts

The existing `Park` logic should remain unchanged, provided the warehouse does not require modifications to the park's current behavior

The new dependency will be added to the Composition root

# Console app

Приложение предоставляет следующие возможности:

1. Просмотр всего инвентаря
2. Добавление транспортного средства
3. Добавление оборудования
4. Расчёт суммарного суточного энергопотребления
5. Просмотр транспорта, подходящего для новичков
6. Завершение работы программы

При добавлении транспорта выполняется технический осмотр

Если транспорт доступен, он добавляется в парк

Если транспорт требует ремонта, он отклоняется

Для новичков используется правило:

```text
Simplicity >= 6
```

---
The app is in russian

The application provides the following features:

1. View entire inventory
2. Add a vehicle
3. Add equipment
4. Calculate total daily energy consumption
5. View vehicles suitable for beginners
6. Exit the program

A technical inspection is performed when adding a vehicle:

If the vehicle is in good condition, it is added to the fleet

If the vehicle requires repairs, it is rejected

The following rule applies to beginners:

```text
Simplicity >= 6
```

# Test details

tests/park_test.go

- TestAddVehicleAvailable — проверка добавления транспорта, успешно прошедшего техосмотр

- TestAddVehicleNeedsRepair — проверка отказа в добавлении транспорта, требующего ремонта

- TestAddEquipment — проверка добавления оборудования в парк

- TestCalculateDailyEnergy — проверка расчёта суммарного суточного энергопотребления транспорта

- TestGetBeginnerVehicles — проверка получения транспорта, подходящего для новичков

- TestGetInventory — проверка получения полного списка инвентаря

- TestCalculateDailyEnergyNilPark — проверка поведения при расчёте энергопотребления для nil-парка

- TestGetBeginnerVehiclesBoundary — проверка граничного значения простоты 6

- TestGetBeginnerVehiclesBelowBoundary — проверка транспорта с простотой ниже порога 6

- TestCalculateDailyEnergyEmptyPark — проверка нулевого энергопотребления пустого парка

- TestGetBeginnerVehiclesEmptyPark — проверка отсутствия транспорта для новичков в пустом парке

- TestGetInventoryEmptyPark — проверка пустого инвентаря нового парка

- TestRejectedVehicleDoesNotAffectEnergy — проверка того, что отклонённый транспорт не учитывается при расчёте энергопотребления

- TestBikeDoesNotAffectEnergy — проверка того, что транспорт, не реализующий EnergyConsumer, не учитывается при расчёте энергопотребления

- TestGetSeveralBeginnerVehicles — проверка получения нескольких транспортных средств, подходящих для новичков

- TestGetInventoryNilPark — проверка возврата nil для инвентаря nil-парка

- TestGetBeginnerVehiclesNilPark — проверка возврата nil для списка транспорта nil-парка

- TestAddVehicleWithFakeInspector — проверка подмены инспектора через Dependency Injection с использованием fakeInspector

---

tests/park_test.go

- TestAddVehicleAvailable — verifies the addition of a vehicle that has successfully passed inspection

- TestAddVehicleNeedsRepair — verifies the rejection of a vehicle requiring repair

- TestAddEquipment — verifies the addition of equipment to the fleet

- TestCalculateDailyEnergy — verifies the calculation of total daily energy consumption for vehicles

- TestGetBeginnerVehicles — verifies the retrieval of vehicles suitable for beginners

- TestGetInventory — verifies the retrieval of the complete inventory list

- TestCalculateDailyEnergyNilPark — verifies behavior when calculating energy consumption for a nil fleet

- TestGetBeginnerVehiclesBoundary — verifies the boundary value of 6 for the simplicity rating

- TestCalculateDailyEnergyEmptyPark — verify zero energy consumption for an empty park

- TestGetBeginnerVehiclesEmptyPark — verify the absence of beginner vehicles in an empty park

- TestGetInventoryEmptyPark — verify that a new park has an empty inventory

- TestRejectedVehicleDoesNotAffectEnergy — verify that rejected vehicles are not included in energy consumption calculations

- TestBikeDoesNotAffectEnergy — verify that vehicles not implementing EnergyConsumer are not included in energy consumption calculations

- TestGetBeginnerVehicles — verifies the retrieval of multiple vehicles suitable for beginners

- TestGetInventoryNilPark — verifies that `nil` is returned for the inventory of a `nil` park

- TestGetBeginnerVehiclesNilPark — verifies that `nil` is returned for the vehicle list of a `nil` park

- TestAddVehicleWithFakeInspector — verifies inspector substitution via Dependency Injection using a `fakeInspector`


# Tesring Dependency Injection

Зависимость технического осмотра представлена интерфейсом:

```go
type Inspector interface {
	Inspect(vehicle domain.Vehicle) bool
}
```

Благодаря этому в тесте можно передать вместо `BasicInspector` собственную тестовую реализацию

Например, fake-инспектор может всегда возвращать `true` или `false`

Это позволяет отдельно проверить поведение `Park` при разных результатах технического осмотра, не привязывая тест к реализации `BasicInspector`

---

The vehicle inspection dependency is defined by the following interface:

```go
type Inspector interface {
	Inspect(vehicle domain.Vehicle) bool
}
```

This allows you to pass a custom test implementation instead of `BasicInspector` during testing

For instance, a fake inspector could simply return `true` or `false`

This makes it possible to test the behavior of `Park` under different inspection outcomes without coupling the test to the `BasicInspector` implementation

# The use of AI

При разработке проекта использовался ИИ-помощник

ИИ использовался для:

- обсуждения структуры проекта
- поиска возможных вариантов организации кода
- помощи при написании и расширении юнит-тестов
- проверки отдельных решений и объяснения ошибок
- объяснения нюансов языка Golang
- оформление ReadMe

---

An AI assistant was used during the development of this project.

AI was used for:

- discussing the project structure
- exploring code organization options
- assisting with writing and expanding unit tests
- verifying specific solutions and explaining errors
- explaining nuances of the Golang
- drafting the README
