package download

type DowloadFileParams struct {
	Filename string `params:"filename" validate:"required,gte=1,lte=255"`
}
