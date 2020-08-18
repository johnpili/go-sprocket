package api

// ResponsePayload ...
type ResponsePayload struct {
	Reference  string      `json:"reference"`
	Success    bool        `json:"success"`
	Result     interface{} `json:"result"`
	Errors     []string    `json:"errors"`
	Messages   []string    `json:"messages"`
	ResultInfo *ResultInfo `json:",omitempty"`
}

// ResultInfo ...
type ResultInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Count      int `json:"count"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}
