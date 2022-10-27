package do

type User struct {
	DoubanUid    uint64 `json:"douban_uid"`
	Domain       string `json:"domain"`
	Name         string `json:"name"`
	Thumbnail    string `json:"thumbnail"`
	BookWish     uint32 `json:"book_wish"`
	BookDo       uint32 `json:"book_do"`
	BookCollect  uint32 `json:"book_collect"`
	GameWish     uint32 `json:"game_wish"`
	GameDo       uint32 `json:"game_do"`
	GameCollect  uint32 `json:"game_collect"`
	MusicWish    uint32 `json:"music_wish"`
	MusicDo      uint32 `json:"music_do"`
	MusicCollect uint32 `json:"music_collect"`
	MovieWish    uint32 `json:"movie_wish"`
	MovieDo      uint32 `json:"movie_do"`
	MovieCollect uint32 `json:"movie_collect"`
}

func (User) TableName() string {
	return "user"
}
