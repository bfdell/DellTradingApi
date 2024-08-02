package dtos

// todo: update title of dto make it more generic?
type TickerRequestDto struct {
	Ticker string `form:"ticker" json:"ticker" xml:"ticker" binding:"required"`
}
