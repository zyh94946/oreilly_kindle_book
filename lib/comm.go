package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var tmpDir = ""

// Init tmp dir.
func InitCheck() error {

	err := error(nil)

	switch runtime.GOOS {
	case "linux":
	case "darwin":
	default:
		return errors.New("emmm.. only supports linux and macos.")
	}

	if _, err := exec.LookPath("kindlegen"); err != nil {
		return errors.New("please install kindlegen first. \r\nfrom https://www.amazon.com/gp/feature.html?ie=UTF8&docId=1000765211. \r\nerr :" + err.Error())
	}

	if tmpDir, err = ioutil.TempDir("", "oreilly_kindle_book."); err != nil {
		return errors.New("create temp dir error! err: " + err.Error())
	}

	if err := os.Mkdir(tmpDir+"/assets", 0777); err != nil {
		return errors.New("create assets dir error! err: " + err.Error())
	}

	return nil
}

func GetTmpPath() string {
	return tmpDir
}

func TmpClear() {
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Println("tmp dir remove error! err:", err)
	}
}

type jsonCus struct {
	url string
}

func (jc *jsonCus) getJson(structVal interface{}) error {
	body, err := httpGet(jc.url)

	if err != nil {
		return errors.New("getJson error: " + err.Error())
	}

	if err := json.NewDecoder(body).Decode(&structVal); err != nil {
		return errors.New("json decode " + jc.url + " error: " + err.Error())
	}

	return nil
}

func saveHttpFile(baseUrl string, saveName string) error {
	body, err := httpGet(baseUrl)
	if err != nil {
		return err
	}

	fileData, _ := ioutil.ReadAll(body)
	if err := saveFile(saveName, fileData); err != nil {
		return err
	}

	return nil

}

func saveFile(fullPath string, fileData []byte) error {
	return ioutil.WriteFile(tmpDir+"/"+fullPath, fileData, 0644)
}

func httpGet(useUrl string) (io.ReadCloser, error) {
	httpCl := http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequest("GET", useUrl, nil)
	if err != nil {
		return nil, errors.New("newRequest err:" + err.Error())
	}

	// set header
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Add("accept-encoding", "deflate, br")
	req.Header.Add("accept-language", "en-US;q=0.8,en;q=0.7,ja;q=0.6,zh-TW;q=0.5,st;q=0.4,sk;q=0.3,ko;q=0.2")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh)")

	// set cookie
	for _, cookieItem := range Config.userCookie {
		req.AddCookie(cookieItem)
	}

	res, err := httpCl.Do(req)
	if err != nil {
		return nil, errors.New("get " + useUrl + " err: " + err.Error())
	}
	defer res.Body.Close()

	cpByte, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return nil, errors.New("get " + useUrl + " status code: " + res.Status + " err: " + string(cpByte))
	}

	body := ioutil.NopCloser(bytes.NewBuffer(cpByte))

	return body, nil

}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
