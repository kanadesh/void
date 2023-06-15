package types

type Audio struct {
	File  string
	Start float64
}

type Option struct {
	Fps        int    `json:"fps"`
	Frames     int    `json:"frames"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	CacheDir   string `json:"cacheDir"`
	ResultFile string `json:"resultFile"`
	Number     int    `json:"number"`
	Audios     []struct {
		Link  string  `json:"link"`
		Start float64 `json:"start"`
	} `json:"audios"`
}
