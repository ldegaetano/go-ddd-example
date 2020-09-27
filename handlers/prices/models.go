package prices

type (
	priceCreate struct {
		ItemCode  string `json:"item_code" binding:"required,max=5"`
		ItemPrice float64 `json:"item_price" binding:"required"`
	}

	pricesResponse struct {
		Items []item `json:"items"`
	}

	item struct {
		ItemCode  string  `json:"item_code"`
		ItemPrice float64 `json:"item_price"`
	}
)
