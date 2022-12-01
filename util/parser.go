package util

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
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

func TrimBookParagraph(node *html.Node) string {
	output := htmlquery.OutputHTML(node, true)
	outputWithLine := strings.ReplaceAll(output, "</p>", "</p>\n")
	body, _ := htmlquery.Parse(strings.NewReader(outputWithLine))
	return TrimParagraph(htmlquery.InnerText(body))
}

func TrimParagraph(info string) string {
	var data strings.Builder

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
