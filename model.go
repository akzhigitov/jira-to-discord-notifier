package main

type Message struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Color       int    `json:"color"`
	Author      Author `json:"author"`
}

type Author struct {
	Name string `json:"name"`
	Icon string `json:"icon_url"`
}
