package main

// An interface that defines the methods for an item
type InventoryItem interface {
	Name() string
	InventoryNumber() string
}
