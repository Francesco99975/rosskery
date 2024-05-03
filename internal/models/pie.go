package models

type PieItem struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Color int     `json:"color"`
}

type Pie struct {
	Title string    `json:"title"`
	Items []PieItem `json:"items"`
}
