package main

type Issue struct {
	Number    int      `json:"number,omitempty"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Assignees []string `json:"assignees,omitempty"`
	Labels    []*Label `json:"labels"`
	State     string   `json:"state,omitempty"`
	User      *User    `json:",omitempty"`
}

type User struct {
	Login   string
	HTMLUrl string `json:"html_url"`
}

type Label struct {
	Name string `json:"name"`
}
