package models

type (
	WarehouseProduct struct {
		WarehouseUUID    string
		ProductUUID      string
		Quantity         int
		ReservedQuantity int
	}
)
