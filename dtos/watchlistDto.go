package dtos

type WatchlistAppendDto struct {
	Ticker string `form:"ticker" json:"ticker" xml:"ticker" binding:"required"`
}
