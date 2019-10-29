package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/zyh94946/oreilly_kindle_book/lib"
)

func main() {
	flag.Parse()

	if lib.Config.Version {
		fmt.Printf("oreilly_kindle_book %s (built: %s, Git SHA: %s, Go Version: %s)\n", Version, Built, GitSHA, runtime.Version())
		os.Exit(0)
	}

	if lib.Config.ArgsIsEmpty() || lib.Config.Help {
		flag.Usage()
		if lib.Config.Help {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	defer lib.TmpClear()
	if err := lib.InitCheck(); err != nil {
		log.Fatalln(err)
	}

	if err := lib.Config.Login(); err != nil {
		log.Fatalln(err)
	}
	log.Println("login success!")

	bookInfo, err := lib.Config.GetBookInfo()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("book name:", bookInfo.Title)

	// Build toc.html, toc.ncx.
	if err := lib.BuildToc(bookInfo.Toc); err != nil {
		log.Fatalln(err)
	}
	log.Println("build toc success!")

	log.Println("get chapter:")
	chapterList, err := lib.GetAllChapter(bookInfo.ChapterList)
	if err != nil {
		log.Fatalln(err)
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
		log.Fatalln(err)
	}
	log.Println("build opf file success!")

	log.Println("generate mobi:")
	bookInfo.GenerateMobi()

	os.Exit(0)
}
