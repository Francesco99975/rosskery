package models

type Dataset struct {
	Topic      string   `json:"topic"`
	Horizontal []string `json:"horizontal"`
	Vertical   []int    `json:"vertical"`
	Color      int      `json:"color"`
}
