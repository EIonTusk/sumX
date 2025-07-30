package db

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
    bun.BaseModel `bun:"table:users,alias:u"`

    ID        int64     `bun:",pk,autoincrement"`
    UserName  string    `bun:",notnull"`
    Email     string    `bun:",unique,notnull"`
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`

	Summaries []*Summary `bun:"rel:has-many,join:id=user_id"`
}

type XUser struct {
	bun.BaseModel `bun:"table:xusers,alias:xu"`

    ID        int64     `bun:",pk"`
    UserName  string    `bun:",notnull"`
}

type SummaryData struct {
	Heading string `json:"heading"`
	Text    string `json:"text"`
}

type Summary struct {
    bun.BaseModel `bun:"table:summaries,alias:s"`

    ID        int64     `bun:",pk,autoincrement"`
    UserID   int64     `bun:",notnull"` // FK to users.id
    XUserID   int64     `bun:",notnull"` // FK to users.id
    From      time.Time `bun:",nullzero"`
    To        time.Time `bun:",nullzero"`
    Limit     int16     `bun:",nullzero"`
	Summary 	[]SummaryData `bun:",nullzero"`
	Tweets []string	`bun:",nullzero"`	
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`

    XUser     *XUser    `bun:"rel:belongs-to,join:x_user_id=id"`
    User      *User    `bun:"rel:belongs-to,join:user_id=id"`
}
