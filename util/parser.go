package util

import (
	"regexp"
	"strconv"
	"time"
)

var doubanUidParser = regexp.MustCompile(`.*/icon/[a-z]+(\d+)(-.*)?\..*`)
var dateParser = regexp.MustCompile(`.*(\d{4}-\d{2}(-\d{2})?).*`)
var numberParser = regexp.MustCompile(`\D*(\d+).*`)
var domainParser = regexp.MustCompile(`.*people/(.*)/`)

func ParseDoubanUid(thumbnail string) uint64 {
	result := doubanUidParser.FindStringSubmatch(thumbnail)
	doubanUid, err := strconv.ParseUint(result[1], 10, 64)
	if err != nil {
		return 0
	}
	return doubanUid
}

func ParseDate(date string) time.Time {
	result := dateParser.FindStringSubmatch(date)
	dateTime, err := time.Parse("2006-01-02", result[1])
	if err != nil {
		return time.Time{}
	}
	return dateTime
}

func ParseNumber(number string) uint64 {
	result := numberParser.FindStringSubmatch(number)
	num, err := strconv.ParseUint(result[1], 10, 64)
	if err != nil {
		return 0
	}
	return num
}

func ParseDomain(link string) string {
	result := domainParser.FindStringSubmatch(link)
	return result[1]
}
