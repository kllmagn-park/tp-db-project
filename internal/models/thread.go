package models

import "time"

type Thread struct {
	Id      int32     `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"` 
	Forum   string    `json:"forum"` 
	Message string    `json:"message"`
	Votes   int32     `json:"votes"`   
	Slug    string    `json:"slug,omitempty"`
	Created time.Time `json:"created"`   
}

type CreateThreadForm struct {
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum,omitempty"`
	Message string    `json:"message"`
	Slug    string    `json:"slug,omitempty"`
	Created time.Time `json:"created"`
}

type ThreadUpdate struct {
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

type GetPostsForm struct {
	Limit int32  `json:"limit,omitempty"`
	Since int64  `json:"since"`
	Sort  string `json:"sort"`
	Desc  bool   `json:"desc"`
}

const (
	Flat       = "flat"
	Tree       = "tree"
	ParentTree = "parent_tree"
)
