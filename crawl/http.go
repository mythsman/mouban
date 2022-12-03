package crawl

import (
	"crypto/tls"
	"fmt"
	cookiemonster "github.com/MercuryEngineering/CookieMonster"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

var userAgent = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/105.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:106.0) Gecko/20100101 Firefox/106.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:105.0) Gecko/20100101 Firefox/105.0",
	"Mozilla/5.0 (Windows NT 10.0; rv:105.0) Gecko/20100101 Firefox/105.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36 Edg/106.0.1370.42",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36 Edg/106.0.1370.47",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36 Edg/106.0.1370.52",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0",
}[rand.Intn(13)]

var client http.Client
var DefaultLimiter *rate.Limiter
var BackLimiter *rate.Limiter

func init() {
	jar, _ := cookiejar.New(nil)
	client = http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 0 && req.Header.Get("cookie") == "" {
				req.Header.Set("Cookie", via[len(via)-1].Header.Get("Cookie"))
			}
			return nil
		},
		Timeout: time.Duration(viper.GetInt("http.timeout")) * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 10 * time.Second,
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
	DefaultLimiter = rate.NewLimiter(rate.Every(time.Duration(viper.GetInt("http.default_interval"))*time.Second), 1)
	BackLimiter = rate.NewLimiter(rate.Every(time.Duration(viper.GetInt("http.back_interval"))*time.Second), 1)
}

func Get(url string, limiter *rate.Limiter) (*string, int, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*120)
	err := limiter.Wait(ctx)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("GET", url, nil)
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

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	fmt.Printf(" code is %d for %s\n", resp.StatusCode, url)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	bodyStr := string(body)
	return &bodyStr, resp.StatusCode, err
}
