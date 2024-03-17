package model

import "time"

type Article struct {
	ID            string    `json:"id"`
	User          User      `json:"user"`
	CreateTime    time.Time `json:"create_time"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Image         string    `json:"image"`
	View          string    `json:"view"`
	Zan_size      int64     `json:"like"`
	Bookmark_size int64     `json:"collect"`
	Forward_size  int64     `json:"transmit"`
	Comment_size  int64     `json:"comment"`
}
