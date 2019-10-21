package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Tags map[string]interface{}
type Body map[string]interface{}

func (tags *Tags) Parse(jsonStr string) error {
	err := json.Unmarshal([]byte(jsonStr), tags)
	if err != nil {
		return err
	}
	return nil
}

const (
	ColorRed   = "\033[0;31m"
	ColorGreen = "\033[0;32m"
	ColorNone  = "\033[0m"
)

func main() {

	checkArgs()

	tags := Tags{}
	err := tags.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("playerids.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if playerID := strings.TrimSpace(scanner.Text()); playerID != "" {
			resp, err := SendTags(playerID, tags)
			if err != nil {
				printError(playerID, err)
			} else {
				printSuccess(playerID, resp)
			}
			fmt.Println("------------------------------------")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func SendTags(playerID string, tags Tags) (string, error) {
	body := Body{}
	body["tags"] = tags
	return NewRequest("PUT", "/players/"+playerID, body)
}

func NewRequest(method, subUrl string, body Body) (string, error) {

	reqUrl := "https://onesignal.com/api/v1" + subUrl
	fmt.Println("url: ", reqUrl)

	body["app_id"] = os.Getenv("APP_ID")
	jsonObj, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	fmt.Println("request json: ", string(jsonObj))

	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(jsonObj))
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(
			fmt.Sprintf("status code: %d, response: %s", resp.StatusCode, string(b)))
	}

	return string(b), nil
}

func printSuccess(playerID, resp string) {
	fmt.Print(ColorGreen)
	fmt.Printf("PlayerID: %s success with response: %s\n", playerID, resp)
	fmt.Print(ColorNone)
}

func printError(playerID string, err error) {
	fmt.Print(ColorRed)
	fmt.Printf("PlayerID: %s failed with error: %s\n", playerID, err.Error())
	fmt.Print(ColorNone)
}

func checkArgs() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ", os.Args[0], "<json tag>")
		os.Exit(-1)
	}
}
