package do

type Movie struct {
	DoubanId  uint64 `json:"douban_id"`
	Title     string `json:"title"`
	Director  string `json:"director"`
	Writer    string `json:"writer"`
	Actor     string `json:"actor"`
	Style     string `json:"style"`
	Site      string `json:"site"`
	Country   string `json:"country"`
	Language  string `json:"language"`
	PublishAt string `json:"publish_at"`
	Season    uint32 `json:"season"`
	Episode   uint32 `json:"episode"`
	Duration  uint32 `json:"duration"`
	Alias     string `json:"alias"`
	IMDb      string `json:"imdb"`
	Intro     string `json:"intro"`
	Thumbnail string `json:"thumbnail"`
}

func (Movie) TableName() string {
	return "movie"
}
