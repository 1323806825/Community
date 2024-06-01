package dto

type Paginate struct {
	PageNum  int `json:"page_num,omitempty" form:"page_num"`
	PageSize int `json:"page_size,omitempty" form:"page_size"`
}
