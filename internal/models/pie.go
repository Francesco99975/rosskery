package models

type PieItem struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type Pie struct {
	Title string    `json:"title"`
	Items []PieItem `json:"items"`
}
