package types

type WpUser struct {
	ID             int    `json:"ID"`
	Roles          string `json:"roles"`
	DisplayName    string `json:"display_name"`
	UserLogin      string `json:"user_login"`
	UserEmail      string `json:"user_email"`
	UserStatus     string `json:"user_status"`
	UserPassword   string `json:"user_pass"`
	UserRegistered string `json:"user_registered"`
}

type WpPlugin struct {
	Name    string `json:"name" yaml:"name"`
	Status  string `json:"status" yaml:"status"`
	Update  string `json:"update" yaml:"update"`
	Version string `json:"version" yaml:"version"`
	Url     string `json:"url" yaml:"url"`
}

type WpTheme struct {
	Name    string `json:"name" yaml:"name"`
	Status  string `json:"status" yaml:"status"`
	Update  string `json:"update" yaml:"update"`
	Version string `json:"version" yaml:"version"`
	Url     string `json:"url" yaml:"url"`
}

type WpPost struct {
	Id     int    `json:"ID"`
	Title  string `json:"post_title"`
	Name   string `json:"post_name"`
	Date   string `json:"post_date"`
	Status string `json:"post_status"`
}
