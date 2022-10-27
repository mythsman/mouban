package do

type Rating struct {
	Total  uint32  `json:"total"`
	Rating float32 `json:"rating"`
	Star5  float32 `json:"star5"`
	Star4  float32 `json:"star4"`
	Star3  float32 `json:"star3"`
	Star2  float32 `json:"star2"`
	Star1  float32 `json:"star1"`
	Status uint8   `json:"status"` //comment:0-normal,1-not enough,2-can not rate
}

func (Rating) TableName() string {
	return "rating"
}
