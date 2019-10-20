package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Tags map[string]interface{}
type Body map[string]interface{}

func main() {

	file, err := os.Open("playerids.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		playerID := strings.TrimSpace(scanner.Text())
		if playerID != "" {
			err := SetAsApproved(playerID)
			if err != nil {
				fmt.Printf("playerID: %s failed with error: %s\n", playerID, err.Error())
			}
			fmt.Println("-------------------------------------")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func SetAsApproved(playerID string) error {
	resp, err := SendTag(playerID, Tags{
		"isApproved": true,
	})
	if err != nil {
		return err
	}
	fmt.Printf("playerID: %s success with response: %s\n", playerID, resp)
	return nil
}

func SendTag(playerID string, tags Tags) (string, error) {
	body := Body{}
	body["tags"] = tags
	return NewRequest("PUT", "/players/"+playerID, body)
}

func NewRequest(method, subUrl string, body Body) (string, error) {

	body["app_id"] = os.Getenv("APP_ID")
	jsonObj, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	fmt.Println("request json", string(jsonObj))
	req, err := http.NewRequest(method, "https://onesignal.com/api/v1"+subUrl, bytes.NewBuffer(jsonObj))
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

	b, err := ioutil.ReadAll(resp.Body.(io.Reader))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(
			fmt.Sprintf("status code: %d, response: %s", resp.StatusCode, string(b)))
	}

	return string(b), nil
}
