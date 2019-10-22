package lib

import "fmt"

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

        <item id="toc" media-type="application/xhtml+xml" href="toc.html"></item>

        <!-- table of contents [mandatory] -->
        <item id="My_Table_of_Contents" media-type="application/x-dtbncx+xml" href="toc.ncx"/>

        <!-- cover image [mandatory] -->
        <item id="My_Cover" media-type="image/jpeg" href="cover.jpg"/>
    </manifest>


    <spine toc="My_Table_of_Contents">
        <!-- the spine defines the linear reading order of the book -->
        <itemref idref="toc"/>
        %s
    </spine>

    <guide>
        <reference type="toc" title="Table of Contents" href="toc.html"></reference>
        %s
    </guide>

</package>
`

func BuildOpenPackagingFormat(cl []ChapterItem) error {
	bookInfo, _ := GetBookInfo("")
	authors := ""
	for _, val := range bookInfo.Authors {
		authors += fmt.Sprintf("<dc:creator>%s</dc:creator>", val.Name)
	}

	manifestItem := ""
	spineItem := ""
	for _, val := range cl {
		manifestItem += fmt.Sprintf("<item id=\"%s\" media-type=\"%s\" href=\"%[1]s\"></item>\r\n", val.FullPath, "application/xhtml+xml")
		spineItem += fmt.Sprintf("<itemref idref=\"%s\"/>\r\n", val.FullPath)
		if len(val.Images) > 0 {
			for _, imgVal := range val.Images {
				manifestItem += fmt.Sprintf("<item id=\"%s\" media-type=\"%s\" href=\"%[1]s\"></item>\r\n", imgVal, "")
			}
		}
	}

	firstItem := GetFirstItem()

	next := fmt.Sprintf("<reference type=\"text\" title=\"%s\" href=\"%s\"></reference>", firstItem.Label, firstItem.Href)

	opf := fmt.Sprintf(opfTemplate, bookInfo.Title, bookInfo.Language, bookInfo.Isbn, authors, bookInfo.Issued, bookInfo.Description, manifestItem, spineItem, next)
	err := SaveFile("build.opf", []byte(opf))
	if err != nil {
		return err
	}

	return nil

}
