package dto

type PaginationQuery struct {
	Page       *int32  `query:"page" validate:"omitempty,gte=0,lte=100"`
	Limit      *int32  `query:"limit" validate:"omitempty,gte=1,lte=100"`
	NextCursor *string `query:"cursor"`
}

func (p *PaginationQuery) UseCursor() bool {
	return p.NextCursor != nil && *p.NextCursor != ""
}

func (p *PaginationQuery) SetDefaults() {
	defaultLimit := int32(10)
	if p.Limit == nil {
		p.Limit = &defaultLimit
	}
	if p.Page == nil {
		defaultPage := int32(1)
		p.Page = &defaultPage
	}
}
