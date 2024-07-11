package dtos

type WatchlistRequestDto struct {
	Ticker string `form:"ticker" json:"ticker" xml:"ticker" binding:"required"`
}
