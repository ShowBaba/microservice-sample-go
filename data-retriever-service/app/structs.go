package app

type ListResult struct {
	Nodes      []interface{} `json:"nodes"`
	TotalCount int           `json:"totalCount"`
}
