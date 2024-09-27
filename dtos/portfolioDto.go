package dtos

type PortfolioEntryDto struct {
	TickerRequestDto
	Shares uint `form:"shares" json:"shares" xml:"shares" binding:"required"`
}

type PortfolioGraphDto struct {
	Date        string  `form:"date" json:"date" xml:"date" binding:"required"`
	StockAssets float64 `form:"stock_assets" json:"stock_assets" xml:"stock_assets" binding:"required"`
	Cash        float64 `form:"cash" json:"cash" xml:"cash" binding:"required"`
}
