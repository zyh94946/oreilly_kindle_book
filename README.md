[![Go Report Card](https://goreportcard.com/badge/github.com/zyh94946/oreilly_kindle_book)](https://goreportcard.com/report/github.com/zyh94946/oreilly_kindle_book) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/zyh94946/oreilly_kindle_book/blob/master/LICENSE)

# Generate mobi file for o'reilly book

**Technical learning only.**

## Install
Only support mac and linux.

**Install [kindlegen](https://www.amazon.com/gp/feature.html?ie=UTF8&docId=1000765211.) first.**

Assuming you already have a recent version of Go installed.

Use make & make install:

```
$ git clone https://github.com/zyh94946/oreilly_kindle_book.git
$ cd oreilly_kindle_book
$ make && make install
Building oreilly_kindle_book...
Building success...
Installing oreilly_kindle_book...
Install success to /usr/local/bin/oreilly_kindle_book.
Usage of oreilly_kindle_book:
  -email string
    	you login email of https://www.oreilly.com/member/
  -help
    	help
  -n string
    	the num of https://learning.oreilly.com/library/view/BOOK-NAME/***
  -p string
    	you login password of https://www.oreilly.com/member/
  -version
    	print version and exit
```

Or pull down the code with go get:

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
  -help
    	help
  -n string
    	the num of https://learning.oreilly.com/library/view/BOOK-NAME/***
  -p string
    	you login password of https://www.oreilly.com/member/
  -version
    	print version and exit
```

## Usage

To generate mobi use:

```
$ oreilly_kindle_book -n BOOK_NUM -email YOU_EMAIL -p YOU_PASSWORD
``` 
