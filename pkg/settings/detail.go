package settings

import (
	"time"
)

type Detail struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ReleaseNames []string `json:"releases"`
	Issues       []Issue  `json:"issues"`
}

type Issue struct {
	Id      int64      `json:"id"`
	Number  int        `json:"number"`
	Title   string     `json:"title"`
	Url     string     `json:"url"`
	Created *time.Time `json:"created"`
	Closed  *time.Time `json:"closed"`
}
