[![Go Report Card](https://goreportcard.com/badge/github.com/zyh94946/oreilly_kindle_book)](https://goreportcard.com/report/github.com/zyh94946/oreilly_kindle_book) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/zyh94946/oreilly_kindle_book/blob/master/LICENSE)

# Generate mobi file for o'reilly book

## Install
Only support mac and linux.

**Install [kindlegen](https://www.amazon.com/gp/feature.html?ie=UTF8&docId=1000765211.) first.**

Assuming you already have a recent version of Go installed, pull down the code with go get:

```
$ go get github.com/zyh94946/oreilly_kindle_book
```

```
$ go install github.com/zyh94946/oreilly_kindle_book
```

```
$ oreilly_kindle_book
Usage of oreilly_kindle_book:
  -email string
        you login email of https://www.oreilly.com/member/
  -n string
        the num of https://learning.oreilly.com/library/view/BOOK-NAME/***
  -p string
        you login password of https://www.oreilly.com/member/
```
