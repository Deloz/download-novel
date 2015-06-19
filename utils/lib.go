package utils

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/qiniu/iconv"
)

const LANG = "utf-8"

func FetchPage(sourceLang, url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")
	req.Header.Add("Referer", url)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("fetch fail")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var input io.Reader
	if strings.ToLower(sourceLang) != LANG {
		cd, err := iconv.Open(LANG, sourceLang)
		CheckError(err)
		defer cd.Close()
		input = bytes.NewReader([]byte(cd.ConvString(string(body))))
	} else {
		input = bytes.NewReader(body)
	}

	return goquery.NewDocumentFromReader(input)
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func FixURL(baseURL string, path string) string {
	b, err := url.Parse(baseURL)
	CheckError(err)

	p, err := b.Parse(path)
	CheckError(err)

	return p.String()
}

func Trace(s string) (string, time.Time) {
	log.Println("START: ", s)
	return s, time.Now()
}

func Un(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println(" END: ", s, "ElapsedTime in seconds:", endTime.Sub(startTime), "\n")
}
