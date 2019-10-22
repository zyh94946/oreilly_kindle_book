package lib

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

    err := GetJson("https://learning.oreilly.com/api/v1/book/"+bookId+"/", &bookInfo)
    if err != nil {
        return book{}, err
    }

    return bookInfo, nil
}
