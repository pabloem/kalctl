package auth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/pabloem/kalctl/reqs"
)

var DefaultCredsFile = filepath.Join(os.Getenv("HOME"), ".kalctl/auth.json")
var DefaultTokenFile = filepath.Join(os.Getenv("HOME"), ".kalctl/token.json")

const KALSHI_AUTH_PATH = "trade-api/v2/login"

type Creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func WriteToken(token reqs.Token) error {
	file, err := os.OpenFile(DefaultTokenFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	return enc.Encode(token)
}

func GetToken() (reqs.Token, error) {
	file, err := os.Open(DefaultTokenFile)
	if err != nil {
		return reqs.Token{}, err
	}
	defer file.Close()

	var token reqs.Token
	dec := json.NewDecoder(file)
	err = dec.Decode(&token)
	if err != nil {
		return reqs.Token{}, err
	}

	return token, nil
}

func WriteUserCredentials(creds Creds) error {
	file, err := os.OpenFile(DefaultCredsFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	return enc.Encode(creds)
}

func ImportUserCredentials() (Creds, error) {
	file, err := os.Open(DefaultCredsFile)
	if err != nil {
		return Creds{}, err
	}
	defer file.Close()

	var creds Creds
	dec := json.NewDecoder(file)
	err = dec.Decode(&creds)
	if err != nil {
		return Creds{}, err
	}

	return creds, nil
}

func ReadUserCredentials() (Creds, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return Creds{}, err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return Creds{}, err
	}

	password := string(bytePassword)
	return Creds{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}, nil
}

type expectedAuthResponse struct {
	MemberId string `json:"member_id"`
	Token    string `json:"token"`
}

func RunKalshiAuth(permAuth bool) error {
	// Check that $HOME/.kalctl exists. If it doesn't, create it
	if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".kalctl")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(os.Getenv("HOME"), ".kalctl"), 0700)
	}

	creds, err := ImportUserCredentials()

	// If the file doesn't exist, prompt the user for credentials
	if err != nil {
		creds, err := ReadUserCredentials()
		if err != nil {
			return err
		}
		if permAuth {
			// If the user wants to save the credentials, write them to the file
			return WriteUserCredentials(creds)
		}
	}
	res, err := reqs.KalshiRequest(
		reqs.HttpRequestTemplate{
			Path:   KALSHI_AUTH_PATH,
			Method: reqs.POST,
		},
		reqs.Token{}, // Empty token
		fmt.Sprintf(`{"email": "%s", "password": "%s"}`, creds.Username, creds.Password))
	if err != nil {
		return err
	}
	dec := json.NewDecoder(strings.NewReader(res))
	var authResponse expectedAuthResponse
	err = dec.Decode(&authResponse)
	if err != nil {
		return err
	}

	token := reqs.Token{
		Token:        authResponse.Token,
		CreationTime: time.Now().Unix(),
	}
	return WriteToken(token)
}
