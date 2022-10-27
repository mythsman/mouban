package do

type Game struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Platform    string `json:"platform"`
	Alias       string `json:"alias"`
	Developer   string `json:"developer"`
	Publisher   string `json:"publisher"`
	PublishDate string `json:"publish_date"`
	Intro       string `json:"intro"`
	Thumbnail   string `json:"thumbnail"`
}
