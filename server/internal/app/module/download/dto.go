package download

type DowloadFileParams struct {
	Filename string `params:"filename" validate:"required,gte=1,lte=255"`
}

type DowloadFileQuery struct {
	Query string `query:"query" validate:"omitempty,oneof=info download"`
}
