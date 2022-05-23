package configstore

type Config struct {
	Id      string            `json:"id"`
	Entries map[string]string `json:"entries"`
	Version string            `json:"version"`
}

type Group struct {
	Id      string              `json:"id"`
	Configs []map[string]string `json:"configs"`
	Version string              `json:"version"`
}
