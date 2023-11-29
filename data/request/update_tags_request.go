package request

type UpdateTagsRequest struct {
	ID   int    `validate:"required"`
	Name string `validate:"required,max=200,min=1" json:"name"`
}
