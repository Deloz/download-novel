## Usage
  1. compile for *Linux amd64*

    `GOOS=linux GOARCH=amd64 go build -x -o novel-downloader main.go`

  2. select site, and download novel

    `./novel-downloader -site=zilang -url=http://zilang.com/xx/xx`

  3. while done the noval would save at direcotry `./download`
