package storage

type Setting string

const (
	Online    Setting = "online"
	Operative Setting = "operative"
	Message   Setting = "message"
)

type Settings struct {
	Online    bool
	Operative bool
	Message   string
}
