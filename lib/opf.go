package lib

import (
	"fmt"
	"strings"
)

var opfTemplate = `<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="2.0" unique-identifier="BookId">
    <metadata xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:opf="http://www.idpf.org/2007/opf">
        <!-- Title [mandatory]: The title of the publication. This is the title that will appear on the "Home" screen. -->
        <dc:title>%s</dc:title>

        <!-- Language [mandatory]: the language of the publication. The language codes used are the same as in XML
        and HTML. The full list can be found here: http://www.w3.org/International/articles/language-tags/
        Some common language strings are:
        "en-us" English - USA
        "fr"    French
        "de"    German
        "es"    Spanish
        -->
        <dc:language>%s</dc:language>

        <!-- Cover [mandatory]. The cover image must be specified in <manifest> and referenced from
        this <meta> element with a name="cover" attribute.
        -->
        <meta name="cover" content="My_Cover" />

        <!-- The ISBN of your book goes here -->
        <dc:identifier id="BookId" opf:scheme="ISBN">%s</dc:identifier>

        <!-- The author of the book. For multiple authors, use multiple <dc:Creator> tags.
        Additional contributors whose contributions are secondary to those listed in
        creator  elements should be named in contributor elements.
        -->
        %s

        <!-- Publisher: An entity responsible for making the resource available -->
        <dc:publisher>O'Reilly Media, Inc.</dc:publisher>

        <!-- Subject: A topic of the content of the resource. Typically, Subject will be
        expressed as keywords, key phrases or classification codes that describe a topic
        of the resource. The BASICCode attribute should contain the subject code
        according to the BISG specification:
        http://www.bisg.org/what-we-do-20-73-bisac-subject-headings-2008-edition.php
        -->
        <dc:subject></dc:subject>

        <!-- Date: Date of publication in YYYY-MM-DD format. (Days and month can be omitted).
        Standard to follow: http://www.w3.org/TR/NOTE-datetime
        -->
        <dc:date>%s</dc:date>

        <!-- Description: A short description of the publication's content. -->
        <dc:description>%s</dc:description>

    </metadata>

    <manifest>
        <!-- HTML content files [mandatory] -->
        %s

        <!-- table of contents [mandatory] -->
        <item id="My_Table_of_Contents" media-type="application/x-dtbncx+xml" href="toc.ncx"/>

        <!-- cover image [mandatory] -->
        <item id="My_Cover" media-type="image/jpeg" href="cover.jpg"/>
    </manifest>


    <spine toc="My_Table_of_Contents">
        <!-- the spine defines the linear reading order of the book -->
        %s
    </spine>

    <guide>
        %s
    </guide>

</package>
`

func BuildOpenPackagingFormat(cl []ChapterItem) error {
	bookInfo, _ := Config.GetBookInfo()
	authors := strings.Builder{}
	for _, val := range bookInfo.Authors {
		authors.WriteString(`<dc:creator>` + val.Name + `</dc:creator>`)
	}

	manifestItem := strings.Builder{}
	spineItem := strings.Builder{}
	for _, val := range cl {
		manifestItem.WriteString(`<item id="` + val.FullPath + `" media-type="text/x-oeb1-document" href="` + val.FullPath + `"></item>`)
		spineItem.WriteString(`<itemref idref="` + val.FullPath + `"/>`)
	}

	rangeChapterImage(func(k interface{}, v interface{}) bool {
		manifestItem.WriteString(`<item id="` + k.(string) + `" media-type="" href="` + k.(string) + `"></item>`)
		return true
	})

	firstPage := fmt.Sprintf("<reference type=\"text\" title=\"%s\" href=\"%s\"></reference>", chapterList[0].Title, chapterList[0].FullPath)

	opf := fmt.Sprintf(opfTemplate, bookInfo.Title, bookInfo.Language, bookInfo.Isbn, authors.String(), bookInfo.Issued, bookInfo.Description, manifestItem.String(), spineItem.String(), firstPage)
	if err := saveFile("build.opf", []byte(opf)); err != nil {
		return err
	}

	return nil

}
