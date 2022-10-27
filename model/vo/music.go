package do

type Music struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Actor       string `json:"actor"`
	Style       string `json:"style"`
	Media       string `json:"media"`
	Genre       string `json:"genre"`
	PublishDate string `json:"publish_date"`
	Publisher   string `json:"publisher"`
	Barcode     string `json:"barcode"`
	ISRC        string `json:"isrc"`
	Alias       string `json:"alias"`
	Thumbnail   string `json:"thumbnail"`
	Intro       string `json:"intro"`
}

func (Music) TableName() string {
	return "music"
}
