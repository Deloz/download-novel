package sites

import (
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"bitbucket.org/deloz/zilang/utils"

	"github.com/PuerkitoBio/goquery"
)

type Zilang struct{}

func (z Zilang) ParseNovelList(targetURL string) {
	doc, err := utils.FetchPage("gbk", targetURL)
	utils.CheckError(err)

	re := regexp.MustCompile(`《|》`)
	bookName := doc.Find(".book h1").Text()
	bookName = re.ReplaceAllString(bookName, "")
	author := doc.Find(".book .small span").First().Text()

	filename := bookName + "--" + author + ".txt"
	dir, err := os.Getwd()
	utils.CheckError(err)

	saveFilePath := path.Join(dir, "download", filename)
	f, err := os.Create(saveFilePath)
	utils.CheckError(err)
	_, err = f.WriteString(bookName + "\n" + author + "\n\n\n\n")
	utils.CheckError(err)

	doc.Find(".book .list ul li").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		url := s.Find("a").AttrOr("href", "")
		if len(url) == 0 {
			return
		}

		url = utils.FixURL(targetURL, url)

		log.Println("title: ", title)
		log.Println("url: ", url)
		content := downloadArticle(title, url)
		_, err := f.WriteString(strings.TrimSpace(title) + "\n\n" + strings.TrimSpace(content) + "\n\n\n\n")
		utils.CheckError(err)
	})
}

func downloadArticle(title, articleURL string) string {
	defer utils.Un(utils.Trace("download article: " + title + " url=>> " + articleURL))
	doc, err := utils.FetchPage("gbk", articleURL)
	utils.CheckError(err)

	html, err := doc.Find("#chapter_content").Html()
	utils.CheckError(err)
	html = strings.TrimSpace(html)
	re := regexp.MustCompile(`</p>\s*<p[^>]*>|<br\s*/?>`)
	html = re.ReplaceAllString(html, "\n")
	re = regexp.MustCompile(`<!--.*?-->|<script[^>]*>.*?</script>|<a[^>]*>.*?</a>|<div[^>]*>|</div>|<p[^>]*>|</p>|\(紫琅文学http://www\.zilang\.net\)`)
	html = re.ReplaceAllString(html, "")

	return html
}
