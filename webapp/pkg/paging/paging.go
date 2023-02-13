package paging

type PageResponse struct {
	MetaInfo MetaInfo    `json:"metaInfo"`
	Items    interface{} `json:"items"`
}

func New(items interface{}, pageMetaInfo MetaInfo) *PageResponse {
	return &PageResponse{pageMetaInfo, items}
}
