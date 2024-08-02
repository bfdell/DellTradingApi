package dtos

import (
	"strconv"
	"time"
)

type StockResponseDto struct {
	Ticker        string      `form:"ticker" json:"symbol" xml:"symbol" binding:"required"`
	Price         CustomFloat `form:"price" json:"close" xml:"close" binding:"required"`
	Name          string      `form:"name" json:"name" xml:"name" binding:"required"`
	DateTime      CustomDate  `form:"datetime" json:"timestamp" xml:"timestamp" binding:"required"`
	PercentChange CustomFloat `form:"percent_change" json:"percent_change" xml:"percent_change" binding:"required"`
}

// These two custom types are made to convert the data recieved from the api into the proper type and format
type CustomFloat float64

func (cf *CustomFloat) UnmarshalJSON(b []byte) error {
	// Convert JSON bytes to string
	str := string(b)

	// Remove quotes around the string (if present)
	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// Convert string to float64
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	// Set the float value
	*cf = CustomFloat(value)
	return nil
}

type CustomDate time.Time

func (cd CustomDate) String() string {
	return time.Time(cd).Format(time.RFC3339)
}

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	//convert unix timestamp from byte array into int64
	timeStr := string(b)
	timestamp, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return err
	}
	//convert timestamp into date
	value := time.Unix(timestamp, 0)

	*cd = CustomDate(value)
	return nil
}
