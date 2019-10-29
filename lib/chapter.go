package lib

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

type chapter struct {
	Next    string        `json:"next"`
	Results []ChapterItem `json:"results"`
}

type ChapterItem struct {
	FullPath     string       `json:"full_path"`
	Content      string       `json:"content"`
	Title        string       `json:"title"`
	AssetBaseUrl string       `json:"asset_base_url"`
	Images       []string     `json:"images"`
	Stylesheets  []styleSheet `json:"stylesheets"`
}

type styleSheet struct {
	FullPath    string `json:"full_path"`
	OriginalUrl string `json:"original_url"`
}

var chapterList = make([]ChapterItem, 0)
var chapterCssList = sync.Map{}
var chapterImageList = sync.Map{}

func GetAllChapter(url string) ([]ChapterItem, error) {
	chapterRes := chapter{}

	jc := jsonCus{url: url}
	if err := jc.getJson(&chapterRes); err != nil {
		return []ChapterItem{}, err
	}

	log.Println(chapterRes.Next)

	for _, val := range chapterRes.Results {
		chapterList = append(chapterList, val)
	}

	if chapterRes.Next != "" {
		if _, err := GetAllChapter(chapterRes.Next); err != nil {
			return []ChapterItem{}, err
		}
	}
	return chapterList, nil

}

func (ci *ChapterItem) Down() {
	log.Println("download", ci.FullPath)
	if err := ci.saveHtml(ci.Content, ci.FullPath); err != nil {
		log.Fatalln(err)
	}
	if len(ci.Images) > 0 {
		for _, imgUrl := range ci.Images {
			if _, isExist := chapterImageList.Load(imgUrl); isExist {
				continue
			}

			log.Println("download", imgUrl)
			chapterImageList.Store(imgUrl, true)
			if err := saveHttpFile(ci.AssetBaseUrl+imgUrl, imgUrl); err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func (ci *ChapterItem) saveHtml(useUrl string, fullPath string) error {
	fileBody, err := httpGet(useUrl)
	if err != nil {
		return err
	}

	styleHtml, err := ci.getCssHtml()
	if err != nil {
		return err
	}

	head := "<!doctype html><html lang=\"en\"><head><meta charset=\"utf-8\" /><title></title>" + styleHtml + "</head><body>"
	foot := "</body></html>"

	body, err := ioutil.ReadAll(fileBody)
	if err != nil {
		return err
	}
	html := bytes.Buffer{}
	html.WriteString(head)
	html.Write(body)
	html.WriteString(foot)

	if err = saveFile(fullPath, html.Bytes()); err != nil {
		return err
	}

	return nil
}

func (ci *ChapterItem) getCssHtml() (string, error) {
	styleHtml := strings.Builder{}
	if len(ci.Stylesheets) == 0 {
		return "", nil
	}

	for _, styleVal := range ci.Stylesheets {
		if err := styleVal.saveCss(); err != nil {
			return "", err
		}
		styleHtml.WriteString(`<link rel="stylesheet" href="`)
		styleHtml.WriteString(styleVal.FullPath)
		styleHtml.WriteString(`" type="text/css" />`)
	}

	return styleHtml.String(), nil
}

func (ss *styleSheet) saveCss() error {
	if _, isExist := chapterCssList.Load(ss.FullPath); isExist {
		return nil
	}

	if err := saveHttpFile(ss.OriginalUrl, ss.FullPath); err != nil {
		return err
	}

	chapterCssList.Store(ss.FullPath, true)
	return nil
}

func rangeChapterImage(f func(k, v interface{}) bool) {
	chapterImageList.Range(f)
}
