package util

import (
	"mouban/consts"
	"reflect"
	"testing"
	"time"
)

func TestParseDoubanUid(t *testing.T) {
	type args struct {
		thumbnail string
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"1", args{thumbnail: "https://img2.doubanio.com/icon/up162448367-3.jpg"}, 162448367},
		{"2", args{thumbnail: "https://img2.doubanio.com/icon/up162448367.jpg"}, 162448367},
		{"3", args{thumbnail: "https://img2.doubanio.com/icon/u162448367-3.jpg"}, 162448367},
		{"4", args{thumbnail: "https://img2.doubanio.com/icon/u162448367.jpg"}, 162448367},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseDoubanUid(tt.args.thumbnail); got != tt.want {
				t.Errorf("ParseDoubanUid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNumber(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"1", args{number: "25 本"}, 25},
		{"2", args{number: "总共有 25 本"}, 25},
		{"3", args{number: "21474836380 本"}, 21474836380},
		{"4", args{number: "总共(123)"}, 123},
		{"5", args{number: "总12共(123)"}, 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseNumber(tt.args.number); got != tt.want {
				t.Errorf("ParseNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	type args struct {
		date string
	}

	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"1", args{date: "  2017-09-01 读过"}, time.Date(2017, 9, 1, 0, 0, 0, 0, time.Local)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseDate(tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trimLine(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{text: "   112 \n 21321    \n \n   a "}, "112 21321 a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimLine(tt.args.text); got != tt.want {
				t.Errorf("TrimLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimParagraph(t *testing.T) {
	type args struct {
		info string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{info: "\n 1 dal \n \n fsd  ds \n  "}, "1 dal\nfsd  ds"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimParagraph(tt.args.info); got != tt.want {
				t.Errorf("TrimParagraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseItem(t *testing.T) {
	type args struct {
		info string
	}
	type result struct {
		id uint64
		t  consts.Type
	}
	tests := []struct {
		name string
		args args
		want result
	}{
		{"1", args{info: "sitemap10.xml:    <loc>https://book.douban.com/subject/1503201/</loc>"}, result{id: 1503201, t: consts.TypeBook}},
		{"2", args{info: "sitemap10.xml:    <loc>https://book.douban.com/subject/1503201"}, result{id: 1503201, t: consts.TypeBook}},
		{"3", args{info: "sitemap10.xml:    <loc>https://music.douban.com/subject/1503937/</loc>"}, result{id: 1503937, t: consts.TypeSong}},
		{"4", args{info: "sitemap10.xml:    <loc>https://movie.douban.com/subject/1506676/</loc>"}, result{id: 1506676, t: consts.TypeMovie}},
		{"5", args{info: "sitemap10.xml:    <loc>https://m3ovie.douban.com/subject/1506676/</loc>"}, result{id: 0, t: consts.TypeUser}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if id, ty := ParseItem(tt.args.info); id != tt.want.id || ty != tt.want.t {
				t.Errorf("ParseItem() = %v, want %v", id, tt.want.id)
			}
		})
	}
}
