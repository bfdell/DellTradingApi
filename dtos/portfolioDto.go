package dtos

type PortfolioEntryDto struct {
	TickerRequestDto
	Shares uint `form:"shares" json:"shares" xml:"shares" binding:"required"`
}
