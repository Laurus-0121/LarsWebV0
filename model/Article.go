package model

import "time"

type Article struct {
	ID            string    `json:"id"`
	UserId        int64     `column:"user_id"json:"user"`
	Title         string    `json:"title"`
	Content       string    `column:"content" json:"body"`
	CreateTime    time.Time `column:"create_time"json:"create_time"`
	Zan_size      int64     `column:"zan_size"json:"zan_size"`
	Bookmark_size int64     `column:"bookmark_size"json:"collect"`
	Forward_size  int64     `column:"forward_size"json:"transmit"`
	Comment_size  int64     `column:"comment_size" json:"comment"`
	Image         string    `column:"image" json:"image"`
	View          string    `column:"view"json:"view"`
}
