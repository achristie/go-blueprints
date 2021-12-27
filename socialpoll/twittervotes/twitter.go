package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joeshaw/envdecode"
)

type ts struct {
	AccessToken string `env:"TWTR_TOKEN,required"`
}

var (
	authSetupOnce sync.Once
	httpClient    *http.Client
	twtr          ts
)

func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
	authSetupOnce.Do(func() {
		setupTwitterAuth()
		httpClient = &http.Client{Transport: &http.Transport{
			Dial: dial,
		},
		}
	})
	formEnc := params.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	req.Header.Set("Authorization", "Bearer "+twtr.AccessToken)

	return httpClient.Do(req)
}

func setupTwitterAuth() {
	if err := envdecode.Decode(&twtr); err != nil {
		log.Fatalln(err)
	}
}

var conn net.Conn

func dial(netw, addr string) (net.Conn, error) {
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

type tweet struct {
	Data struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
}

func readFromTwitter(votes chan<- string) {
	options, err := loadOptions()
	if err != nil {
		log.Println("failed to load options:", err)
		return
	}

	u, err := url.Parse("https://api.twitter.com/2/tweets/search/stream")
	if err != nil {
		log.Println("creating filter request failed:", err)
		return
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Println("creating filter request failed:", err)
		return
	}
	resp, err := makeRequest(req, nil)
	if err != nil {
		log.Println("making request failed:", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalln("Failed to connect to Streaming API", string(body), resp.StatusCode)
	}
	reader := resp.Body
	decoder := json.NewDecoder(reader)
	for {
		var t tweet
		if err := decoder.Decode(&t); err != nil {
			break
		}
		log.Println(t)
		for _, option := range options {
			if strings.Contains(
				strings.ToLower(t.Data.Text),
				strings.ToLower(option),
			) {
				log.Println("vote:", option)
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
			case <-stopchan:
				log.Println("stopping twitter")
				return
			default:
				log.Println("Querying twitter")
				readFromTwitter(votes)
				log.Println("waiting...")
				time.Sleep(10 * time.Second)

			}
		}
	}()
	return stoppedchan
}
