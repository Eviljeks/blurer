package image

type BlurImage struct {
	UUID      string `json:"uuid"`
	ImageUUID string `json:"image_uuid"`
	X0        int    `json:"x_0"`
	Y0        int    `json:"y_0"`
	X1        int    `json:"x_1"`
	Y1        int    `json:"y_1"`
	TS        int64  `json:"ts"`
}
