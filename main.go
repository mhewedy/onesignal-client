package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

const playerIDCharLen = 36

type Tags map[string]interface{}
type Body map[string]interface{}

type ErrorInfo struct {
	PlayerID string
	Error    string
}

func (tags *Tags) Parse(jsonStr string) error {
	err := json.Unmarshal([]byte(jsonStr), tags)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	checkArgs()

	errorInfo := make([]ErrorInfo, 0)
	var sCount int

	registerExitHook(&errorInfo, &sCount)

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

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	bar := progressbar.New(int(fileInfo.Size() / playerIDCharLen))
	_ = bar.RenderBlank()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if playerID := strings.TrimSpace(scanner.Text()); playerID != "" {
			_, err := SendTags(playerID, tags)
			if err != nil {
				errorInfo = append(errorInfo, ErrorInfo{PlayerID: playerID, Error: err.Error()})
			} else {
				sCount++
			}
			_ = bar.Add(1)
		}
	}
	_ = bar.Finish()

	err = printSummary(errorInfo, sCount)
	if err != nil {
		log.Fatal(err)
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

	body["app_id"] = os.Getenv("APP_ID")
	jsonObj, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

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

func printSummary(errorInfo []ErrorInfo, sCount int) error {

	fmt.Println("\n\nSummary report has been written to \"summary.log\"")

	file, err := os.OpenFile("summary.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_ = file.Truncate(0)
	_, _ = file.Seek(0, 0)
	_, _ = fmt.Fprintln(file, "************ Summary Report ************")

	if len(errorInfo) == 0 {
		_, _ = fmt.Fprintln(file, "All (", sCount, ") PlayerID succeed with no errors.")
	} else {
		_, _ = fmt.Fprintln(file, "We have (", sCount, ") PlayerIDs succeed and",
			"(", len(errorInfo), ") PlayerIDs Failed as follows:")
		for key := range errorInfo {
			_, _ = fmt.Fprintf(file, "PlayerID: %s failed with error: %s\n",
				errorInfo[key].PlayerID, errorInfo[key].Error)
		}
	}
	return nil
}

func registerExitHook(errorInfo *[]ErrorInfo, sCount *int) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = printSummary(*errorInfo, *sCount)
		os.Exit(-1)
	}()
}

func checkArgs() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<json tag>")
		fmt.Println("example:\nAPP_ID=myOneSignalAppId", os.Args[0], `'{"approved": true}'`)
		os.Exit(-1)
	}
}
