package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"bitbucket.org/deloz/zilang/sites"
	"bitbucket.org/deloz/zilang/utils"
)

func main() {
	rootDir, err := os.Getwd()
	utils.CheckError(err)
	downloadPath := path.Join(rootDir, "download")

	site := flag.String("site", "", "novel from which site")
	listURL := flag.String("url", "", "the list url")
	flag.Parse()
	if len(*listURL) == 0 {
		printHelp()
	}

	switch *site {
	case "zilang":
		sites.Zilang{}.ParseNovelList(*listURL, downloadPath)
		break
	case "lewen8":
		sites.Lewen8{}.ParseNovelList(*listURL, downloadPath)
		break
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println(`
	Usage examle => novel-downloader -site=zilang -url=http://www.zilang.net/170/170694/
	Now support sites :  zilang lewen8
	`)
	os.Exit(1)
}
