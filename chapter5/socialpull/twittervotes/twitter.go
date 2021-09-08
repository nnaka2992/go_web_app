package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/garyburd/go-oauth/oauth"
)

var (
	conn net.Conn
)

func dial(netw, addr string) (net.Conn, error){
	if conn != nil {
		conn.Close()
		conn = nil
	}
	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	conn = netc
	return netc, nil
}

var reader io.ReadCloser
func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}

}

var (
	authClient *oauth.Client
	creds *oauth.Credentials
)

func setupTwitterAuth() {
	// read env file
	_ = godotenv.Load(".env")
	var ts struct {
		ConsumerKey string `env:"SP_TWITTER_KEY,required"`
		ConsumerSecretKey string `env:"SP_TWITTER_SECRET,required"`
		AccessToken string `env:"SP_TWITTER_ACCESSTOKEN,required"`
		AccessSecret string `env:"SP_TWITTER_ACCESSSECRET,required"`
	}
	if err := envdecode.Decode(&ts); err != nil {
		log.Fatal(err)
	}
	creds = &oauth.Credentials{
		Token: ts.AccessToken,
		Secret: ts.AccessSecret,
	}
	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token: ts.ConsumerKey,
			Secret: ts.ConsumerSecretKey,
		},
	}
}

var (
	authSetupOnce sync.Once
	httpClient *http.Client
)

func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
	authSetupOnce.Do(func() {
		setupTwitterAuth()
		httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					if conn != nil {
						conn.Close()
						conn = nil
					}
					netc, err := net.DialTimeout(netw, addr, 5*time.Second)
					if err != nil {
						return nil, err
					}
					conn = netc
					return netc, nil
				},
			},
		}
	})
	forEnc := params.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(forEnc)))
	req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", req.URL, params))
	return httpClient.Do(req)
}

type tweet struct {
	Text string
}

func readFromTwitter(votes chan<- string) {
	options, err := loadOptions()
	if err != nil {
		log.Println("Failed to load options: ", err)
		return
	}
	u, err := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
	if err != nil {
		log.Println("Failed to parse URL")
		return
	}
	query := make(url.Values)
	query.Set("track", strings.Join(options, ", "))
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(query.Encode()))
	if err != nil {
		log.Println("Failed to query request: ", err)
		return
	}
	resp, err := makeRequest(req, query)
	reader := resp.Body
	decoder := json.NewDecoder(reader)
	for {
		var tweet tweet
		if err := decoder.Decode(&tweet); err != nil {
			break
		}
		for _, option := range options {
			if strings.Contains(strings.ToLower(tweet.Text), strings.ToLower(option)) {
				log.Println("Poll: ", option)
				votes <- option
			}
		}
	}
}

func startTwitterStream(stopchan <-chan struct{}, votes chan<- string) <-chan struct{} {
	stoppedchan := make(chan struct{}, 1)
	go func() {
		defer func() {
			stoppedchan <- struct{}{}
		}()
		for {
			select {
			case <- stopchan:
				log.Println("Finish quering to Twitter")
				return
			default:
				log.Println("Quering to Twitter")
				readFromTwitter(votes)
				log.Println("Waiting for result...")
				time.Sleep(10 * time.Second)
			}
		}
	}()
	return stoppedchan
}
