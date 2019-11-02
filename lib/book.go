package lib

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func (bk *book) GenerateMobi() {
	tmpDir := GetTmpPath()
	mobiName := strings.Replace(bk.Title, " ", "_", -1) + ".mobi"
	cmd := exec.Command("kindlegen", tmpDir+"/build.opf", "-c1", "-o", mobiName, "-verbose")
	cmd.Stdout = os.Stdout
	log.Println(cmd.Args)

	if err := cmd.Start(); err != nil {
		log.Println("kindlegen err:", err)
		return
	}

	if err := cmd.Wait(); err != nil {
		//log.Println("generate mobi err:", err)
		//return
	}

	if isExist, _ := fileExists(tmpDir + "/" + mobiName); isExist == false {
		log.Println("generate mobi error!")
		return
	}

	moveDir, _ := filepath.Abs(".")
	if err := os.Rename(tmpDir+"/"+mobiName, moveDir+"/"+mobiName); err != nil {
		log.Println("move mobi err:", err)
		return
	}

	log.Println("successfully generated mobi to", moveDir+"/"+mobiName)
}

func (bk *book) IsEmpty() bool {
	if "" == bk.Title {
		return true
	}
	return false
}
