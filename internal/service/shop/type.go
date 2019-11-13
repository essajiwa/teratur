package shop

// Shop is shop data
type Shop struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// DataResp hold shop response from API
type DataResp struct {
	Data Shop `json:"data"`
}
