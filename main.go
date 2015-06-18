package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/qiniu/iconv"
)

const LANG = "utf-8"

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`
			Usage:
				zilang url   e.g.  zilang http://www.zilang.net/170/170694/
		`)
		os.Exit(1)
	}

	parseNovelList(os.Args[1])
}

func parseNovelList(targetURL string) {
	doc, err := fetchPage("gbk", targetURL)
	checkError(err)

	re := regexp.MustCompile(`《|》`)
	bookName := doc.Find(".book h1").Text()
	bookName = re.ReplaceAllString(bookName, "")
	author := doc.Find(".book .small span").First().Text()

	filename := bookName + "--" + author + ".txt"
	dir, err := os.Getwd()
	checkError(err)

	saveFilePath := path.Join(dir, "download", filename)
	f, err := os.Create(saveFilePath)
	checkError(err)
	_, err = f.WriteString(bookName + "\n" + author + "\n\n\n\n")
	checkError(err)

	doc.Find(".book .list ul li").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		url := s.Find("a").AttrOr("href", "")
		if len(url) == 0 {
			return
		}

		url = fixURL(targetURL, url)

		log.Println("title: ", title)
		log.Println("url: ", url)
		content := downloadArticle(title, url)
		_, err := f.WriteString(strings.TrimSpace(title) + "\n\n" + strings.TrimSpace(content) + "\n\n\n\n")
		checkError(err)
	})
}

func downloadArticle(title, articleURL string) string {
	defer un(trace("download article: " + title + " url=>> " + articleURL))
	doc, err := fetchPage("gbk", articleURL)
	checkError(err)

	html, err := doc.Find("#chapter_content").Html()
	checkError(err)
	html = strings.TrimSpace(html)
	re := regexp.MustCompile(`</p>\s*<p[^>]*>|<br\s*/?>`)
	html = re.ReplaceAllString(html, "\n")
	re = regexp.MustCompile(`<!--.*?-->|<script[^>]*>.*?</script>|<a[^>]*>.*?</a>|<div[^>]*>|</div>|<p[^>]*>|</p>|\(紫琅文学http://www\.zilang\.net\)`)
	html = re.ReplaceAllString(html, "")

	return html
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func fixURL(baseURL string, path string) string {
	b, err := url.Parse(baseURL)
	checkError(err)

	p, err := b.Parse(path)
	checkError(err)

	return p.String()
}

func gbk2UTF8(str string) string {
	cd, err := iconv.Open("utf-8", "gbk")
	checkError(err)
	defer cd.Close()

	utf8 := cd.ConvString(str)

	return utf8
}

func trace(s string) (string, time.Time) {
	log.Println("START: ", s)
	return s, time.Now()
}

func un(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println(" END: ", s, "ElapsedTime in seconds:", endTime.Sub(startTime))
}

func fetchPage(sourceLang, url string) (*goquery.Document, error) {
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
		checkError(err)
		defer cd.Close()
		input = bytes.NewReader([]byte(cd.ConvString(string(body))))
	} else {
		input = bytes.NewReader(body)
	}

	return goquery.NewDocumentFromReader(input)
}
