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

    chapterList, err := lib.GetAllChapter("https://learning.oreilly.com/api/v1/book/" + bookNum + "/chapter/")
    if err != nil {
        fmt.Println(err)
        return
    }

    err = lib.SaveImage("https://learning.oreilly.com/library/cover/" + bookNum + "/", "cover.jpg")
    if err != nil {
        fmt.Println(err)
        return
    }

    wg := sync.WaitGroup{}

    fmt.Println("start!")
    for _, val := range chapterList {
        wg.Add(1)
        go func(ci lib.ChapterItem){
            defer wg.Done()
            ci.Down()
        }(val)
    }

    wg.Wait()

    fmt.Println("success!")
}

