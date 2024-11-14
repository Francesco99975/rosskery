package storage

type Setting string

const (
	Online    Setting = "online"
	Operative Setting = "operative"
	Message   Setting = "message"
)

type Settings struct {
	Online    bool   `json:"online"`
	Operative bool   `json:"operative"`
	Message   string `json:"message"`
}
