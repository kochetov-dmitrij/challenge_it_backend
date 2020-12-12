package database

import (
	"time"
)

type User struct {
	Id          int32  `pg:",pk"`
	Name        string `pg:"type:varchar(50),notnull"`
	Email       string `pg:"type:varchar(50),notnull"`
	EncPassword string `pg:"type:varchar(50),notnull"`
}

type Challenge struct {
	Id           int32  `pg:",pk"`
	Title        string `pg:"type:varchar(50),notnull"`
	AuthorId     int32  `pg:",notnull"`
	Author       *User  `pg:"rel:has-one"`
	Description  string `pg:"type:text"`
	Requirements string `pg:"type:text"`
	Rating       int32  `pg:",notnull"`
	Taken        int32  `pg:",notnull"`
}

const (
	Assigned int32 = iota + 1
	InProgress
	Completed
	Rejected
)

type UserChallenge struct {
	Id                int32     `pg:",pk"`
	UserId            int32     `pg:",notnull"`
	User              *User     `pg:"rel:has-one"`
	StartDate         time.Time `pg:",notnull"`
	Comment           string    `pg:"type:text"`
	Status            int32     `pg:",notnull"`
	Photo             []byte    `pg:"type:bytea"`
	ConclusionComment string    `pg:"type:text"`
}

type Tag struct {
	Id   int32  `pg:",pk"`
	Name string `pg:"type:varchar(50),notnull"`
}

type ChallengeTag struct {
	Id          int32      `pg:",pk"`
	ChallengeId int32      `pg:",notnull"`
	Challenge   *Challenge `pg:"rel:has-one"`
	TagId       int32      `pg:",notnull"`
	Tag         *Tag       `pg:"rel:has-one"`
}
