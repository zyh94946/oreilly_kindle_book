package lib

import (
    "fmt"
    "strings"
    "os"
    "os/exec"
    "path/filepath"
)

type book struct {
    Cover           string    `json:"cover"`
    ChapterList     string    `json:"chapter_list"`
    Toc             string    `json:"toc"`
    Identifier      string    `json:"identifier"`
    Title           string    `json:"title"`
    HasStylesheets  bool      `json:"has_stylesheets"`
    Isbn            string    `json:"isbn"`
    PageCount       int       `json:"pagecount"`
    OrderAbleAuthor string    `json:"orderable_author"`
    Language        string    `json:"language"`
    Description     string    `json:"description"`
    Issued          string    `json:"issued"`
    Authors         []authors `json:"authors"`
}

type authors struct {
    Name string `json:"name"`
}

var bookInfo = book{}

func (bk book) GenerateMobi() {
    tmpDir := GetTmpPath()
    mobiName := strings.Replace(bk.Title, " ", "_", -1) + ".mobi"
    cmd := exec.Command("kindlegen", tmpDir + "/build.opf", "-c1", "-o", mobiName, "-verbose")
    cmd.Stdout = os.Stdout
    fmt.Println(cmd.Args)

    if err := cmd.Start(); err != nil {
        fmt.Println("kindlegen err:", err)
        return
    }

    if err := cmd.Wait(); err != nil {
        //fmt.Println("generate mobi err:", err)
        //return
    }

    if isExist, _ := FileExists(tmpDir + "/" + mobiName); isExist == false {
        fmt.Println("generate mobi error!")
        return
    }

    moveDir, _ := filepath.Abs(".")
    if err := os.Rename(tmpDir + "/" + mobiName, moveDir + "/" + mobiName); err != nil {
        fmt.Println("move mobi err:", err)
        return
    }

    fmt.Println("successfully generated mobi to", moveDir + "/" + mobiName)
}

func (bk book) IsEmpty() bool {
    if "" == bk.Title {
        return true
    }
    return false
}

func GetBookInfo(bookId string) (book, error) {
    if !bookInfo.IsEmpty() {
        return bookInfo, nil
    }

    jc := jsonCus{url:"https://learning.oreilly.com/api/v1/book/" + bookId + "/", method:"GET"}
    if err := jc.getJson(&bookInfo); err != nil {
        return book{}, err
    }

    return bookInfo, nil
}
