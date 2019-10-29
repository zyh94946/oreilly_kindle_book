package lib

import (
	"fmt"
	"strings"
)

type tocItem struct {
	Href     string    `json:"href"`
	Filename string    `json:"filename"`
	Depth    int       `json:"depth"`
	Children []tocItem `json:"children"`
	Label    string    `json:"label"`
}

var tocList = make([]tocItem, 0)

var tocNum = 1
var tocDepth = 1
var firstTocItem = tocItem{}

var tocHtmlVar = strings.Builder{}
var tocHtmlTemplate = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
    <head><title>Table of Contents</title></head>
    <body>
        <div>
            <h1><b>TABLE OF CONTENTS</b></h1>
            <br />
            <div>%s</div>
        </div>
    </body>
</html>
`

var tocNcxVar = strings.Builder{}
var tocNcxTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE ncx PUBLIC "-//NISO//DTD ncx 2005-1//EN" "http://www.daisy.org/z3986/2005/ncx-2005-1.dtd">
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1" xml:lang="en-US">
    <head>
        <meta name="dtb:uid" content="%s"/>
        <meta name="dtb:depth" content="%d"/>
        <meta name="dtb:totalPageCount" content="%d"/>
        <meta name="dtb:maxPageNumber" content="%[3]d"/>
    </head>
    <docTitle><text>%s</text></docTitle>
    <docAuthor><text>%s</text></docAuthor>
    <navMap>
        <navPoint class="toc" id="toc" playOrder="1">
            <navLabel>
                <text>Table of Contents</text>
            </navLabel>
            <content src="toc.html"/>
        </navPoint>

        %s

    </navMap>
</ncx>
`

func (ti tocItem) IsEmpty() bool {
	if "" == ti.Label {
		return true
	}
	return false
}

func BuildToc(tocUrl string) error {

	bookInfo, _ := Config.GetBookInfo()

	if err := getToc(tocUrl); err != nil {
		return err
	}
	getTocVal(tocList)

	tocHtml := fmt.Sprintf(tocHtmlTemplate, tocHtmlVar)
	if err := saveFile("toc.html", []byte(tocHtml)); err != nil {
		return err
	}

	tocNcx := fmt.Sprintf(tocNcxTemplate, bookInfo.Isbn, tocDepth, bookInfo.PageCount, bookInfo.Title, bookInfo.OrderAbleAuthor, tocNcxVar)
	if err := saveFile("toc.ncx", []byte(tocNcx)); err != nil {
		return err
	}

	return nil
}

func getToc(url string) error {
	jc := jsonCus{url: url}
	if err := jc.getJson(&tocList); err != nil {
		return err
	}

	return nil
}

func getTocVal(tl []tocItem) {
	tocHtmlVar.WriteString("<ul>")
	for _, val := range tl {
		if tocDepth < val.Depth {
			tocDepth = val.Depth
		}

		if firstTocItem.IsEmpty() {
			firstTocItem = val
		}

		tocNum++
		tocHtmlVar.WriteString(`<li><a href="` + val.Href + `"><b>` + val.Label + `</b></a>`)
		tocNcxVar.WriteString(fmt.Sprintf("<navPoint class=\"chapter\" id=\"chapter_%d\" playOrder=\"%[1]d\"><navLabel><text>%s</text></navLabel><content src=\"%s\"/>", tocNum, val.Label, val.Href))
		if len(val.Children) > 0 {
			getTocVal(val.Children)
		}
		tocHtmlVar.WriteString("</li>")
		tocNcxVar.WriteString("</navPoint>")

	}
	tocHtmlVar.WriteString("</ul>")
}
