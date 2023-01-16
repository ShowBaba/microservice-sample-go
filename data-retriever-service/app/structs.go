package app

type ListResult struct {
	Nodes      []interface{} `json:"nodes"`
	TotalCount int           `json:"totalCount"`
}

type GraphQLPayload struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}
