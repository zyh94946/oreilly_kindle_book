package lib

import (
    "bytes"
    "fmt"
    "io/ioutil"
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
    //Url         string `json:"url"`
}

var chapterList = make([]ChapterItem, 0)
var cssList = make([]string, 0)

func GetAllChapter(url string) ([]ChapterItem, error) {
    chapterRes := chapter{}

    err := GetJson(url, &chapterRes)
    if err != nil {
        return []ChapterItem{}, err
    }

    fmt.Println(chapterRes.Next)

    for _, val := range chapterRes.Results {
        chapterList = append(chapterList, val)
    }

    if chapterRes.Next != "" {
        _, err = GetAllChapter(chapterRes.Next)
        if err != nil {
            return []ChapterItem{}, err
        }
    }
    return chapterList, nil

}

func (ci ChapterItem) Down() {
    err := ci.saveHtml(ci.Content, ci.FullPath)
    if err != nil {
        fmt.Println(err)
        return
    }
    if len(ci.Images) > 0 {
        for _, imgUrl := range ci.Images {
            err := SaveHttpFile(ci.AssetBaseUrl+imgUrl, imgUrl)
            if err != nil {
                fmt.Println(err)
            }
        }
    }
}

func (ci ChapterItem) saveHtml(useUrl string, fullPath string) error {
    fileBody, err := HttpGet(useUrl)
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

    err = SaveFile(fullPath, html.Bytes())
    if err != nil {
        return err
    }

    return nil
}

func (ci ChapterItem) getCssHtml() (string, error) {
    styleHtml := ""
    if len(ci.Stylesheets) == 0 {
        return styleHtml, nil
    }

    for _, styleVal := range ci.Stylesheets {
        err := styleVal.saveCss()
        if err != nil {
            return "", err
        }
        styleHtml += fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\" type=\"text/css\" />", styleVal.FullPath)
    }

    return styleHtml, nil
}

func (ss styleSheet) saveCss() error {
    for _, val := range cssList {
        if val == ss.FullPath {
            return nil
        }
    }

    err := SaveHttpFile(ss.OriginalUrl, ss.FullPath)
    if err != nil {
        return err
    }

    cssList = append(cssList, ss.FullPath)

    return nil
}
