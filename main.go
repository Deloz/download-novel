package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/deloz/zilang/sites"
)

func main() {
	site := flag.String("site", "", "novel from which site")
	listURL := flag.String("url", "", "the list url")
	if len(*listURL) == 0 {
		printHelp()
	}
	flag.Parse()

	switch *site {
	case "zilang":
		sites.Zilang{}.ParseNovelList(*listURL)
		break
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println(`
	Usage examle => downloader -site=zilang -url=http://www.zilang.net/170/170694/
	Now support sites :  zilang
	`)
	os.Exit(1)
}
