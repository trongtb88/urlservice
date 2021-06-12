package entity

import "time"

type HttpRequestUrl struct {
	Url string 			`json:"url"`
	Shortcode string 	`json:"shortcode"`
}

type HttpResponseUrl struct {
	Shortcode string 	`json:"shortcode"`
}

type HttpResponseStatsUrl struct {
	StartDate string 		`json:"startDate"`
	RedirectCount int64 	`json:"redirectCount"`
	LastSeenDate string 	`json:"lastSeenDate,omitempty"`
}

type Url struct {
	OriginUrl  string    `gorm:"type:varchar(200);" json:"url"`
	ShortCode  string    `gorm:"primaryKey;type:varchar(6);json:"shortcode"`
	RedirectCount int64  `gorm:""DEFAULT:0;type:int(20)"" json:"redirect_count"`
	CreatedAt time.Time  `gorm:""DEFAULT:current_timestamp; type:timestamp"" json:"created_at"`
	UpdatedAt time.Time  `gorm:""DEFAULT:current_timestamp;type:timestamp"" json:"updated_at"`
	LastSeenAt time.Time `gorm:""type:timestamp"" json:"last_seen_at"`
}

