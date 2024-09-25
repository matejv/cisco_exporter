package inventory

type InventoryItem struct {
	Name         string
	Description  string
	PartNumber   string
	SerialNumber string
}

type TransceiverItem struct {
	Name         string
	Description  string
	Vendor       string
	PartNumber   string
	SerialNumber string
}
