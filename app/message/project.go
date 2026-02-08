package message

import "time"

type Project struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	Repo      string    `json:"repo"`
	Media     []Media   `json:"media"`
	Links     []Link    `json:"links"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Media struct {
	Id        int    `json:"id"`
	ProjectId int    `json:"projectId"`
	Type      string `json:"type"`
	URL       string `json:"url"`
}

type Link struct {
	Id        int    `json:"id"`
	ProjectId int    `json:"projectId"`
	Name      string `json:"name"`
	URL       string `json:"url"`
}

type CreateProjectRequest struct {
	Name   string   `json:"name"`
	Desc   string   `json:"desc"`
	Repo   string   `json:"repo"`
	Photos []string `json:"photos"`
	Videos []string `json:"videos"`
	Links  []Link   `json:"links"`
}

type UpdateProjectRequest struct {
	Name   string   `json:"name"`
	Desc   string   `json:"desc"`
	Repo   string   `json:"repo"`
	Photos []string `json:"photos"`
	Videos []string `json:"videos"`
	Links  []Link   `json:"links"`
}
