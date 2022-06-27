package models

import "time"

type Post struct {
	Id       int64     `json:"id"`
	Parent   int64     `json:"parent"`
	Author   string    `json:"author"` 
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"` 
	Forum    string    `json:"forum"`  
	Thread   int32     `json:"thread"`
	Created  time.Time `json:"created"`
}

type CreatePostForm struct {
	Parent  int64  `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type GetPostForm struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}

type GetPostReqForm struct {
	Related []string `json:"related,omitempty"`
}

type PostUpdate struct {
	Message string `json:"message,omitempty"`
}

const (
	RelatedUser   = "user"
	RelatedThread = "thread"
	RelatedForum  = "forum"
)
