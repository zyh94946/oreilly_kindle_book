package main

import (
	"fmt"
	"oreilly_kindle_book/lib"
	"sync"
)

func main() {
	bookNum := "9781491926291"

	err := lib.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	bookInfo, err := lib.GetBookInfo(bookNum)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Build toc.html, toc.ncx.
	err = lib.BuildToc(bookInfo.Toc)
	if err != nil {
		fmt.Println(err)
		return
	}

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

	// Save chapter html, images, css files.
	wg := sync.WaitGroup{}
	for _, val := range chapterList {
		wg.Add(1)
		go func(ci lib.ChapterItem) {
			defer wg.Done()
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

	fmt.Println("success!")
}
