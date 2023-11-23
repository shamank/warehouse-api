package schemas

type (
	WarehouseCounter struct {
		WarehouseUUID string
		Count         int
	}

	ProductWarehouseSplitted struct {
		ProductArticle string
		WarehouseData  []WarehouseCounter
	}
)
