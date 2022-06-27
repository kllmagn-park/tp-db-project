package models

type Forum struct {
	Title   string `json:"title" db:"title"`
	User    string `json:"user" db:"user"` 
	Slug    string `json:"slug" db:"slug"`    
	Posts   int64  `json:"posts" db:"posts"`     
	Threads int32  `json:"threads" db:"threads"` 
}

type CreateForumForm struct {
	Title string `json:"title" db:"title"`
	User  string `json:"user" db:"user"`
	Slug  string `json:"slug" db:"slug"`
}

type GetUsersForm struct {
	Slug  string `json:"slug" db:"slug"`
	Limit int32  `json:"limit" db:"limit"`
	Since string `json:"since" db:"since"` 
	Desc  bool   `json:"desc" db:"desc"`
}

type GetThreadsForm struct {
	Limit int32  `json:"limit,omitempty"`
	Since string `json:"since,omitempty"`
	Desc  bool   `json:"desc,omitempty"`
}
