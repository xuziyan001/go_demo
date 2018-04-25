package main

import (
	"fmt"

	"github.com/go-pg/pg"
)

type User struct {
	Id        int64    `sql:"id,pk"`
	Name      string   `sql:"name"`
	Type      string   `sql:"type"`
	tableName struct{} `sql:"public.user"`
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.Id, u.Name, u.Type)
}

type Relationship struct {
	Id           int64    `sql:"id,pk"`
	UserID       int64    `sql:"userid"`
	RelateUserID int64    `sql:"relateuserid"`
	State        string   `sql:"state"`
	tableName    struct{} `sql:"public.relationship"`
}

func (r Relationship) String() string {
	return fmt.Sprintf("Relationship<%d %d %s>", r.UserID, r.RelateUserID, r.State)
}

var db = pg.Connect(&pg.Options{
	User:     "xuziyan",
	Database: "postgres",
})
