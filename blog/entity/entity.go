package entity

import (
	"database/sql"
	"time"
)

type Article struct {
	Id            int
	Title         string
	Content       string
	Subtract      string
	ClickNumber   int
	PraisedNumber int
	CreateTime    int
	UpdateTime    int
	IsShow        bool
	AuthorId      int
}

type User struct {
	Id         int
	Name       string
	Salt       string
	Password   string
	Email      sql.NullString
	DelFlag    bool
	CreateTime int
	Avatar     sql.NullString
}

type History struct {
	Id         int
	CreateTime time.Time
	UserAgent  string
	FullIp     string
	TrimIp     string
	VisitPath  string
	UserHash   int
}
