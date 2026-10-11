package domain

// Status of an item in the inventory
type Status int

// Status constants
const (
	Unknown Status = iota
	Available
	NeedsRepair
)

// String returns the string representation of the Status
func (status Status) NameStatus() string {
	switch status {
	case Unknown:
		return "Unknown"
	case Available:
		return "Available"
	case NeedsRepair:
		return "Needs Repair"
	default:
		return "Invalid Status"
	}
}

// Return status
func (status Status) String() string {
	return status.NameStatus()
}
