package reqs

import (
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
)

const KALSHI_API_URL = "trading-api.kalshi.com"

type Token struct {
	Token        string `json:"token"`
	CreationTime int64  `json:"creation_time"`
}

type HttpMethod int

const (
	GET HttpMethod = iota
	POST
	PUT
	DELETE
)

type HttpRequestTemplate struct {
	Path   string
	Method HttpMethod
}

var httpClient *http.Client

func getHttpClient() *http.Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return httpClient
}

func KalshiRequest(reqTemplate HttpRequestTemplate, token Token, body string) (string, error) {
	var req *http.Request
	var err error
	fullPath := "https://" + path.Join(KALSHI_API_URL, reqTemplate.Path)
	switch reqTemplate.Method {
	case GET:
		req, err = http.NewRequest("GET", fullPath, nil)
		if err != nil {
			return "", err
		}
	case POST:
		req, err = http.NewRequest("POST", fullPath, nil)
		if err != nil {
			return "", err
		}
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	if token.Token != "" {
		req.Header.Add("Authorization", "Bearer "+token.Token)
	}
	if body != "" {
		closer := io.NopCloser(strings.NewReader(body))
		defer closer.Close()
		req.Body = closer
	}
	log.Debug().
		Str("method", req.Method).
		Str("path", req.URL.Path).
		Msgf("Sending request to %s", fullPath)

	client := getHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}
