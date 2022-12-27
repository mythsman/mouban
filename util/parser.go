package util

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"mouban/consts"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var doubanUidParser = regexp.MustCompile(`.*/icon/[a-z]+(\d+)(-.*)?\..*`)
var dateParser = regexp.MustCompile(`.*(\d{4}-\d{2}(-\d{2})?).*`)
var numberParser = regexp.MustCompile(`\D*(\d+).*`)
var floatParser = regexp.MustCompile(`([0-9]*\.?[0-9]+)`)
var domainParser = regexp.MustCompile(`.*people/(.*)/`)
var bookItemParser = regexp.MustCompile("https://book.douban.com/subject/[0-9]*")
var movieItemParser = regexp.MustCompile("https://movie.douban.com/subject/[0-9]*")
var songItemParser = regexp.MustCompile("https://music.douban.com/subject/[0-9]*")
var gameItemParser = regexp.MustCompile("https://www.douban.com/game/[0-9]*")

func ParseDoubanUid(thumbnail string) uint64 {
	result := doubanUidParser.FindStringSubmatch(thumbnail)
	if len(result) == 0 {
		return 0
	}
	doubanUid, err := strconv.ParseUint(result[1], 10, 64)
	if err != nil {
		return 0
	}
	return doubanUid
}

func ParseDate(date string) time.Time {
	result := dateParser.FindStringSubmatch(date)
	if len(result) == 0 {
		return time.Time{}
	}
	dateTime, err := time.ParseInLocation("2006-01-02", result[1], time.Local)
	if err != nil {
		return time.Time{}
	}
	return dateTime
}

func ParseNumber(number string) uint64 {
	result := numberParser.FindStringSubmatch(number)
	if result == nil {
		return 0
	}
	num, err := strconv.ParseUint(result[1], 10, 64)
	if err != nil {
		return 0
	}
	return num
}

func ParseFloat(float string) float32 {
	result := floatParser.FindStringSubmatch(float)
	if result == nil {
		return 0
	}
	f, err := strconv.ParseFloat(result[1], 32)
	if err != nil {
		return 0
	}
	return float32(f)
}

func ParseUidOrDomain(link string) string {
	result := domainParser.FindStringSubmatch(link)
	if len(result) == 0 {
		return ""
	}
	return result[1]
}

func ParseDomain(doubanUid uint64, link string) string {
	result := domainParser.FindStringSubmatch(link)
	if len(result) == 0 || result[1] == strconv.FormatUint(doubanUid, 10) {
		return ""
	}
	return result[1]
}

func ParseNewUsers(doc *html.Node) *[]string {
	var newUsers []string
	newUserSet := make(map[string]bool)

	newUserNodes := htmlquery.Find(doc, "//a[contains(@href,'www.douban.com/people/')]")
	for _, node := range newUserNodes {
		userLink := htmlquery.SelectAttr(node, "href")
		uidOrDomain := ParseUidOrDomain(userLink)
		if uidOrDomain != "" {
			if !newUserSet[uidOrDomain] {
				newUsers = append(newUsers, uidOrDomain)
				newUserSet[uidOrDomain] = true
			}
		}
	}
	return &newUsers
}

func ParseNewItems(doc *html.Node, t consts.Type) *[]uint64 {
	var newItems []uint64
	newItemSet := make(map[uint64]bool)

	var newItemsNodes []*html.Node

	switch t {
	case consts.TypeBook:
		newItemsNodes = htmlquery.Find(doc, "//a[contains(@href,'book.douban.com/subject/')]")
		break
	case consts.TypeMovie:
		newItemsNodes = htmlquery.Find(doc, "//a[contains(@href,'movie.douban.com/subject/')]")
		break
	case consts.TypeGame:
		newItemsNodes = htmlquery.Find(doc, "//a[contains(@href,'www.douban.com/game/')]")
		break
	case consts.TypeSong:
		newItemsNodes = htmlquery.Find(doc, "//a[contains(@href,'music.douban.com/subject/')]")
		break
	}

	for _, node := range newItemsNodes {
		itemLink := htmlquery.SelectAttr(node, "href")
		doubanId := ParseNumber(itemLink)
		if doubanId > 0 {
			if !newItemSet[doubanId] {
				newItems = append(newItems, doubanId)
				newItemSet[doubanId] = true
			}
		}
	}
	return &newItems
}

func TrimBookParagraph(node *html.Node) string {
	output := htmlquery.OutputHTML(node, true)
	outputWithLine := strings.ReplaceAll(output, "</p>", "</p>\n")
	body, _ := htmlquery.Parse(strings.NewReader(outputWithLine))
	return TrimParagraph(htmlquery.InnerText(body))
}

func TrimParagraph(info string) string {
	var data strings.Builder

	info = strings.ReplaceAll(info, "<br>", "\n")
	info = strings.ReplaceAll(info, "<br/>", "\n")

	for _, p := range strings.Split(info, "\n") {
		t := strings.TrimSpace(p)
		if t != "" {
			data.WriteString(t)
			data.WriteString("\n")
		}
	}
	return strings.Trim(data.String(), "\n")
}

func TrimInfo(info string) map[string]string {
	list := strings.Split(info, "<br/>")
	result := make(map[string]string)

	for _, s := range list {
		parse, err := htmlquery.Parse(strings.NewReader(s))
		if err != nil {
			return nil
		}
		text := htmlquery.InnerText(parse)
		line := TrimLine(text)
		colonIndex := strings.Index(line, ":")
		if colonIndex != -1 {
			result[line[0:colonIndex]] = line[colonIndex+1:]
		}
	}

	return result
}

func TrimLine(text string) string {
	var data strings.Builder

	mark := true
	for _, ss := range strings.ReplaceAll(text, "\n", " ") {
		if ss == ' ' {
			if mark {
				continue
			} else {
				data.WriteString(" ")
				mark = true
			}
		} else {
			if mark {
				data.WriteString(string(ss))
				mark = false
			} else {
				data.WriteString(string(ss))
			}
		}
	}

	return strings.TrimSpace(data.String())
}

func ParseItem(line string) (uint64, consts.Type) {
	matches := bookItemParser.FindStringSubmatch(line)
	if matches != nil {
		return ParseNumber(matches[0]), consts.TypeBook
	}

	matches = movieItemParser.FindStringSubmatch(line)
	if matches != nil {
		return ParseNumber(matches[0]), consts.TypeMovie
	}

	matches = songItemParser.FindStringSubmatch(line)
	if matches != nil {
		return ParseNumber(matches[0]), consts.TypeSong
	}

	matches = gameItemParser.FindStringSubmatch(line)
	if matches != nil {
		return ParseNumber(matches[0]), consts.TypeGame
	}
	
	return 0, consts.TypeUser
}
