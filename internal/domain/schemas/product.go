package schemas

type Product struct {
	UUID     string `json:"-"`
	Name     string `json:"name"`
	Size     string `json:"size"`
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
}
