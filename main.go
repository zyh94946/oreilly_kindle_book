package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/zyh94946/oreilly_kindle_book/lib"
)

func main() {
	flag.Parse()

	if lib.Config.Version {
		fmt.Printf("oreilly_kindle_book %s (built: %s, Git SHA: %s, Go Version: %s)\n", Version, Built, GitSHA, runtime.Version())
		return
	}

	if lib.Config.ArgsIsEmpty() || lib.Config.Help {
		flag.Usage()
		return
	}

	if err := lib.InitCheck(); err != nil {
		log.Println(err)
		return
	}
	defer lib.TmpClear()

	if err := lib.Config.Login(); err != nil {
		log.Println(err)
		return
	}
	log.Println("login success!")

	bookInfo, err := lib.Config.GetBookInfo()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("book name:", bookInfo.Title)

	// Build toc.ncx.
	if err := lib.BuildToc(bookInfo.Toc); err != nil {
		log.Println(err)
		return
	}
	log.Println("build toc success!")

	log.Println("get chapter:")
	chapterList, err := lib.GetAllChapter(bookInfo.ChapterList)
	if err != nil {
		log.Println(err)
		return
	}

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
	if err := lib.BuildOpenPackagingFormat(chapterList); err != nil {
		log.Println(err)
		return
	}
	log.Println("build opf file success!")

	log.Println("generate mobi:")
	bookInfo.GenerateMobi()

	return
}
