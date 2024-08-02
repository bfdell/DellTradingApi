package dtos

type PortfolioUpdateRequestDto struct {
	TickerRequestDto
	Shares uint `form:"shares" json:"shares" xml:"shares" binding:"required"`
}
