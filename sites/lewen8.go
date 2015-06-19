package sites

import (
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"bitbucket.org/deloz/zilang/utils"
)

type Lewen8 struct{}

func (l Lewen8) ParseNovelList(listURL, downloadPath string) {
	if !strings.HasPrefix(listURL, "/") {
		listURL += "/"
	}
	doc, err := utils.FetchPage("utf-8", listURL)
	utils.CheckError(err)

	bookName := doc.Find(".kfml .infot h1").Text()
	author := strings.Replace(doc.Find(".kfml .infot span").Text(), "/著", "", -1)

	filename := bookName + "--" + author + ".txt"

	saveFilePath := path.Join(downloadPath, filename)
	f, err := os.Create(saveFilePath)
	utils.CheckError(err)
	_, err = f.WriteString(bookName + "\n" + author + "\n\n\n\n")
	utils.CheckError(err)

	doc.Find("#defaulthtml dd").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		url := s.Find("a").AttrOr("href", "")
		if len(url) == 0 {
			return
		}

		url = utils.FixURL(listURL, url)

		log.Println("title: ", title)
		log.Println("url: ", url)

		content := l.downloadArticle(title, url)
		_, err := f.WriteString(strings.TrimSpace(title) + "\n\n" + strings.TrimSpace(content) + "\n\n\n\n")
		utils.CheckError(err)
	})
}

func (l Lewen8) downloadArticle(title, articleURL string) string {
	defer utils.Un(utils.Trace("download lewen8 article: " + title + " url=>> " + articleURL))
	doc, err := utils.FetchPage("utf-8", articleURL)
	utils.CheckError(err)

	html, err := doc.Find("#content").Html()
	utils.CheckError(err)
	html = strings.TrimSpace(html)
	re := regexp.MustCompile(`</p>\s*<p[^>]*>|<br\s*/?>\s*<br\s*/?>|<br\s*/?>`)
	html = re.ReplaceAllString(html, "\n")
	re = regexp.MustCompile(`<!--.*?-->|<script[^>]*>.*?</script>|<a[^>]*>.*?</a>|<div[^>]*>|</div>|<p[^>]*>|</p>|&nbsp;|@@`)
	html = re.ReplaceAllString(html, "")
	html = strings.Replace(html, "\xEF\xBB\xBF", "", -1)

	return html
}
