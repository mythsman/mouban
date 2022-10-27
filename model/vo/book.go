package do

type Book struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Author      string `json:"author"`
	Translator  string `json:"translator"`
	Publisher   string `json:"publisher"`
	PublishDate string `json:"publish_date"`
	ISBN        string `json:"isbn"`
	Page        string `json:"page"`
	Price       uint32 `json:"price"`
	Intro       string `json:"intro"`
	Thumbnail   string `json:"thumbnail"`
}
