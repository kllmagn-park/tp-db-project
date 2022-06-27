package models

type Vote struct {
	Nickname string `json:"nickname" db:"nickname"`
	Thread   int64  `json:"thread" db:"thread"`
	Voice    int32  `json:"voice" db:"voice"`
}

type CreateVoteForm struct {
	Nickname string `json:"nickname" db:"nickname"`
	Voice    int32  `json:"voice" db:"voice"`    
}
