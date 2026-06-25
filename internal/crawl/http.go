package crawl

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

var clientInitOnce sync.Once
var clientInitErr error

func initClientsFromConfig() error {
	userInterval := viper.GetInt("http.interval.user")
	if userInterval <= 0 {
		userInterval = 4000
	}
	itemInterval := viper.GetInt("http.interval.item")
	if itemInterval <= 0 {
		itemInterval = 4000
	}
	discoverInterval := viper.GetInt("http.interval.discover")
	if discoverInterval <= 0 {
		discoverInterval = 4000
	}

	UserLimiter = rate.NewLimiter(rate.Every(time.Duration(userInterval)*time.Millisecond), 1)
	ItemLimiter = rate.NewLimiter(rate.Every(time.Duration(itemInterval)*time.Millisecond), 1)
	DiscoverLimiter = rate.NewLimiter(rate.Every(time.Duration(discoverInterval)*time.Millisecond), 1)

	authString := strings.TrimSpace(viper.GetString("http.auth"))
	authList := strings.Split(authString, ";")
	for _, authKey := range authList {
		authKey = strings.TrimSpace(authKey)
		if authKey == "" {
			continue
		}
		entry := strings.Split(authKey, ",")
		if len(entry) == 1 {
			client := initClient(strings.TrimSpace(entry[0]), nil)
			clients = append(clients, client)
		} else if len(entry) == 2 {
			proxy, _ := url.ParseRequestURI(strings.TrimSpace(entry[1]))
			client := initClient(strings.TrimSpace(entry[0]), proxy)
			clients = append(clients, client)
		}
	}

	if len(clients) == 0 {
		return fmt.Errorf("http.auth is empty or invalid, no crawl client initialized")
	}

	logrus.Infoln(len(clients), "user auth initialized")
	return nil
}

func ensureClientsInitialized() error {
	clientInitOnce.Do(func() {
		clientInitErr = initClientsFromConfig()
	})
	return clientInitErr
}

// Bootstrap initializes crawler clients and dependencies.
func Bootstrap() error {
	if isCrawlEnabled() {
		if err := ensureClientsInitialized(); err != nil {
			return err
		}
		logrus.Infoln("crawl client enabled")
	} else {
		logrus.Infoln("crawl client disabled")
	}

	if viper.GetBool("storage.enable") {
		if err := ensureStorageInitialized(); err != nil {
			return err
		}
		logrus.Infoln("storage enabled")
	} else {
		logrus.Infoln("storage disabled")
	}

	return nil
}

func initClient(dbcl2 string, proxy *url.URL) *retryablehttp.Client {
	var retryClient = retryablehttp.NewClient()
	retryClient.RetryMax = viper.GetInt("http.retry_max")
	retryClient.Logger = nil
	retryClient.RetryWaitMin = time.Duration(1) * time.Second
	retryClient.RetryWaitMax = time.Duration(60) * time.Second
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		shouldRetry, e := retryablehttp.DefaultRetryPolicy(ctx, resp, err)
		if shouldRetry && err != nil && strings.Contains(err.Error(), "too many redirects") {
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

var (
	callOpsHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "mouban_call_ops_duration",
		Help:    "Histogram of the duration of HTTP call requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"code", "index"})
)

func Get(url string, limiter *rate.Limiter) (*string, int, error) {
	if !isCrawlEnabled() {
		return nil, 0, fmt.Errorf("crawl is disabled")
	}

	if err := ensureClientsInitialized(); err != nil {
		return nil, 0, err
	}

	err := limiter.Wait(context.Background())
	if err != nil {
		return nil, 0, err
	}

	startTime := time.Now()
	req, _ := retryablehttp.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.douban.com/")

	clientIdx := int(atomic.AddUint32(&clientIndex, 1)-1) % len(clients)
	retryClient := clients[clientIdx]

	resp, err := retryClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
		duration := time.Since(startTime).Milliseconds()
		logrus.WithField("code", strconv.Itoa(resp.StatusCode)).Infoln("code is", strconv.Itoa(resp.StatusCode), "in", duration, "ms", "at user", 1+clientIdx, "for", url)
		callOpsHistogram.WithLabelValues(strconv.Itoa(resp.StatusCode), strconv.Itoa(clientIdx)).Observe(float64(duration))
		resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	bodyStr := string(body)

	if resp.StatusCode == 403 && strings.Contains(bodyStr, "error code: 004") {
		return nil, 0, errors.New("IP is probably banned (error code: 004)")
	}

	return &bodyStr, resp.StatusCode, err
}
