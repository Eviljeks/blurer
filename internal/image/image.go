package image

type Image struct {
	UUID string `json:"uuid"`
	Path string `json:"path"`
	TS   int64  `json:"ts"`
}
