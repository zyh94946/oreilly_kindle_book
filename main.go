package main

import (
	"flag"
	"fmt"
	"oreilly_kindle_book/lib"
	"sync"
)

var bookNum *string
var userEmail *string
var userPasswd *string

func init() {
	bookNum = flag.String("n", "", "the num of https://learning.oreilly.com/library/view/BOOK-NAME/***")
	userEmail = flag.String("email", "", "you login email of https://www.oreilly.com/member/")
	userPasswd = flag.String("p", "", "you login password of https://www.oreilly.com/member/")
	flag.Parse()
}

func main() {

	if *bookNum == "" || *userEmail == "" || *userPasswd == "" {
		flag.Usage()
		return
	}

	err := lib.InitCheck()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lib.TmpClear()

	if err := lib.Login(*userEmail, *userPasswd); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("login success!")

	bookInfo, err := lib.GetBookInfo(*bookNum)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("book name:", bookInfo.Title)

	// Build toc.html, toc.ncx.
	err = lib.BuildToc(bookInfo.Toc)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("build toc success!")

	fmt.Println("get chapter:")
	chapterList, err := lib.GetAllChapter(bookInfo.ChapterList)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = lib.SaveHttpFile(bookInfo.Cover, "cover.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("get cover success!")

	// Save chapter html, images, css files.
	maxPro := make(chan bool, 4)
	wg := sync.WaitGroup{}
	for _, val := range chapterList {
		maxPro <- true
		wg.Add(1)
		go func(ci lib.ChapterItem) {
			defer func() {
				<-maxPro
				wg.Done()
			}()
			ci.Down()
		}(val)
	}
	wg.Wait()

	// Build opf file.
	err = lib.BuildOpenPackagingFormat(chapterList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("build opf file success!")

	fmt.Println("generate mobi:")
	bookInfo.GenerateMobi()

}
