package main

// Status of an item in the inventory
type Status int

// Status constants
const (
	Unknown Status = iota
	Available
	NeedsRepair
)

// String returns the string representation of the Status
func (status Status) String() string {
	switch status {
	case Available:
		return "Available"

	case NeedsRepair:
		return "Needs Repair"

	default:
		return "Invalid Status"
	}
}
