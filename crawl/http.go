package crawl

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	_ "mouban/common"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
)

var clients []*retryablehttp.Client
var UserLimiter *rate.Limiter
var ItemLimiter *rate.Limiter
var DiscoverLimiter *rate.Limiter
var clientIndex = uint32(0)

func init() {
	UserLimiter = rate.NewLimiter(rate.Every(time.Duration(viper.GetInt("http.interval.user"))*time.Millisecond), 1)
	ItemLimiter = rate.NewLimiter(rate.Every(time.Duration(viper.GetInt("http.interval.item"))*time.Millisecond), 1)
	DiscoverLimiter = rate.NewLimiter(rate.Every(time.Duration(viper.GetInt("http.interval.discover"))*time.Millisecond), 1)

	authString := viper.GetString("http.auth")
	authList := strings.Split(authString, ";")
	for _, authKey := range authList {
		if authKey == "" {
			continue
		}
		entry := strings.Split(authKey, ",")
		if len(entry) == 1 {
			client := initClient(entry[0], nil)
			clients = append(clients, client)
		} else if len(entry) == 2 {
			proxy, _ := url.ParseRequestURI(entry[1])
			client := initClient(entry[0], proxy)
			clients = append(clients, client)
		}
	}
	logrus.Infoln(len(clients), "user auth initialized")
}

func initClient(dbcl2 string, proxy *url.URL) *retryablehttp.Client {
	var retryClient = retryablehttp.NewClient()
	retryClient.RetryMax = viper.GetInt("http.retry_max")
	retryClient.Logger = nil
	retryClient.RetryWaitMin = time.Duration(1) * time.Second
	retryClient.RetryWaitMax = time.Duration(60) * time.Second
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		shouldRetry, e := retryablehttp.DefaultRetryPolicy(ctx, resp, err)
		if shouldRetry && strings.Contains(err.Error(), "too many redirects") {
			return false, e
		}
		return shouldRetry, e
	}

	retryClient.HTTPClient = &http.Client{
		Jar: initCookieJar(dbcl2),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				logrus.Infoln("too many redirects found for", req.URL.String())
				return errors.New("too many redirects for " + req.URL.String())
			}
			if len(via) > 0 && req.Header.Get("cookie") == "" {
				req.Header.Set("Cookie", via[len(via)-1].Header.Get("Cookie"))
			}
			return nil
		},
		Timeout: time.Duration(viper.GetInt("http.timeout")) * time.Millisecond,
		Transport: &http.Transport{
			TLSHandshakeTimeout: time.Duration(viper.GetInt("http.timeout")) * time.Millisecond,
			Proxy:               http.ProxyURL(proxy),
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

	go func() {
		for range time.NewTicker(time.Hour * 1).C {
			retryClient.HTTPClient.Jar = initCookieJar(dbcl2)
			logrus.Infoln("client cookie jar refreshed")
		}
	}()
	return retryClient
}

func initCookieJar(dbcl2 string) http.CookieJar {
	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   "dbcl2",
		Value:  dbcl2,
		Path:   "/",
		Domain: ".douban.com",
	}
	cookies = append(cookies, cookie)
	u, _ := url.Parse("https://douban.com")
	jar.SetCookies(u, cookies)
	return jar
}

func Get(url string, limiter *rate.Limiter) (*string, int, error) {
	err := limiter.Wait(context.Background())
	if err != nil {
		return nil, 0, err
	}

	startTime := time.Now()
	req, err := retryablehttp.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.douban.com/")

	clientIdx := int(atomic.AddUint32(&clientIndex, 1)-1) % len(clients)
	retryClient := clients[clientIdx]

	resp, err := retryClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
		duration := time.Now().Sub(startTime).Milliseconds()
		logrus.Infoln("code is", strconv.Itoa(resp.StatusCode), "in", duration, "ms", "at user", 1+clientIdx, "for", url)
		resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	bodyStr := string(body)

	if resp.StatusCode == 403 && strings.Contains(bodyStr, "error code: 004") {
		return nil, 0, errors.New("IP is probably banned (error code: 004)")
	}

	return &bodyStr, resp.StatusCode, err
}
