package crawl

import (
	"crypto/tls"
	"github.com/MercuryEngineering/CookieMonster"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"
)

var userAgent = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/106.0.1370.52",
}[rand.Intn(3)]

var retryClient *retryablehttp.Client
var UserLimiter *rate.Limiter
var BookLimiter *rate.Limiter
var MovieLimiter *rate.Limiter
var GameLimiter *rate.Limiter
var DiscoverLimiter *rate.Limiter

func init() {
	jar, _ := cookiejar.New(nil)
	retryClient = retryablehttp.NewClient()
	retryClient.RetryMax = viper.GetInt("http.retry_max")
	retryClient.Logger = nil
	retryClient.RetryWaitMin = time.Duration(1) * time.Second
	retryClient.RetryWaitMax = time.Duration(60) * time.Second
	retryClient.HTTPClient = &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 0 && req.Header.Get("cookie") == "" {
				req.Header.Set("Cookie", via[len(via)-1].Header.Get("Cookie"))
			}
			return nil
		},
		Timeout: time.Duration(viper.GetInt("http.timeout")) * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: time.Duration(viper.GetInt("http.timeout")) * time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				},
			},
		}}
	UserLimiter = rate.NewLimiter(rate.Every(time.Duration(viper.GetInt("http.interval.user"))*time.Second), 1)
	BookLimiter = rate.NewLimiter(rate.Every(time.Duration(viper.GetInt("http.interval.book"))*time.Second), 1)
	MovieLimiter = rate.NewLimiter(rate.Every(time.Duration(viper.GetInt("http.interval.movie"))*time.Second), 1)
	GameLimiter = rate.NewLimiter(rate.Every(time.Duration(viper.GetInt("http.interval.game"))*time.Second), 1)
	DiscoverLimiter = rate.NewLimiter(rate.Every(time.Duration(viper.GetInt("http.interval.discover"))*time.Second), 1)
}

func Get(url string, limiter *rate.Limiter) (*string, int, error) {
	err := limiter.Wait(context.Background())
	if err != nil {
		return nil, 0, err
	}

	req, err := retryablehttp.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.douban.com/")

	cookies, err := cookiemonster.ParseFile("./cookie.txt")
	if err != nil {
		cookies, err = cookiemonster.ParseFile("../cookie.txt")
		if err != nil {
			panic(err)
		}
	}
	for _, c := range cookies {
		c.Value = strings.Trim(c.Value, "\"")
		req.AddCookie(c)
	}

	if err != nil {
		return nil, 0, err
	}

	resp, err := retryClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	log.Println("code is " + strconv.Itoa(resp.StatusCode) + " for " + url)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	bodyStr := string(body)
	return &bodyStr, resp.StatusCode, err
}
