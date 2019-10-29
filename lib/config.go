package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"time"
)

type configType struct {
	bookNum    string
	userEmail  string
	userPasswd string
	Version    bool
	Help       bool
	userCookie []*http.Cookie
	bookInfo   book
}

var Config configType

func init() {
	flag.BoolVar(&Config.Help, "help", false, "help")
	flag.BoolVar(&Config.Version, "version", false, "print version and exit")
	flag.StringVar(&Config.bookNum, "n", "", "the num of https://learning.oreilly.com/library/view/BOOK-NAME/***")
	flag.StringVar(&Config.userEmail, "email", "", "you login email of https://www.oreilly.com/member/")
	flag.StringVar(&Config.userPasswd, "p", "", "you login password of https://www.oreilly.com/member/")
}

func (cf *configType) ArgsIsEmpty() bool {
	if cf.bookNum == "" || cf.userEmail == "" || cf.userPasswd == "" {
		return true
	}

	return false
}

func (cf *configType) Login() error {
	reqBody := map[string]string{"email": cf.userEmail, "password": cf.userPasswd}
	reqBodyJson, _ := json.Marshal(reqBody)
	useUrl := "https://www.oreilly.com/member/auth/login/"

	httpCl := http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequest("POST", useUrl, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return errors.New("newRequest err:" + err.Error())
	}

	// set header
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("origin", "https://www.oreilly.com")
	req.Header.Add("accept-encoding", "deflate, br")
	req.Header.Add("accept-language", "en-US;q=0.8,en;q=0.7,ja;q=0.6,zh-TW;q=0.5,st;q=0.4,sk;q=0.3,ko;q=0.2")
	req.Header.Add("cookie", "")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh)")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "*/*")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("authority", "www.oreilly.com")

	res, err := httpCl.Do(req)
	if err != nil {
		return errors.New("post " + useUrl + " err: " + err.Error())
	}
	defer res.Body.Close()

	cpByte, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return errors.New("The email address or password you've entered is incorrect." + "status code: " + res.Status + " err: " + string(cpByte))
	}

	cf.userCookie = res.Cookies()

	return nil
}

func (cf *configType) GetBookInfo() (book, error) {
	if !cf.bookInfo.IsEmpty() {
		return cf.bookInfo, nil
	}

	jc := jsonCus{url: "https://learning.oreilly.com/api/v1/book/" + cf.bookNum + "/"}
	if err := jc.getJson(&cf.bookInfo); err != nil {
		return book{}, err
	}

	if err := saveHttpFile(cf.bookInfo.Cover, "cover.jpg"); err != nil {
		return book{}, err
	}

	return cf.bookInfo, nil
}
