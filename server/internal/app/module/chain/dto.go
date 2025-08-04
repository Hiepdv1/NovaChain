package chain

type GetBlocksDto struct {
	Page  *int32 `query:"Page" validate:"gte=0,lte=100"`
	Limit *int32 `query:"limit" validate:"gte=1,lte=100"`
}

func (d *GetBlocksDto) SetDefaults() {
	if d.Limit == nil {
		defaultLimit := int32(10)
		d.Limit = &defaultLimit
	}

	if d.Page == nil {
		defaultPage := int32(0)
		d.Page = &defaultPage
	}
}
