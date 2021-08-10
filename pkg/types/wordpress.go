package types

type Plugin struct {
	Name    string `json:"name" yaml:"name"`
	Status  string `json:"status" yaml:"status"`
	Update  string `json:"update" yaml:"update"`
	Version string `json:"version" yaml:"version"`
	Url     string `json:"url" yaml:"url"`
}

type Theme struct {
	Name    string `json:"name" yaml:"name"`
	Status  string `json:"status" yaml:"status"`
	Update  string `json:"update" yaml:"update"`
	Version string `json:"version" yaml:"version"`
	Url     string `json:"url" yaml:"url"`
}

type Post struct {
	Id     int    `json:"ID"`
	Title  string `json:"post_title"`
	Name   string `json:"post_name"`
	Date   string `json:"post_date"`
	Status string `json:"post_status"`
}
