package lib

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/pkg/errors"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
    "runtime"
    "time"
)

var tmpDir = ""
var userCookie []*http.Cookie

// Init tmp dir.
func InitCheck() error {

    err := error(nil)

    switch runtime.GOOS {
    case "linux":
    case "darwin":
    default:
        return errors.New(fmt.Sprintln("emmm.. only supports linux and macos."))
    }

    if _, err := exec.LookPath("kindlegen"); err != nil {
        return errors.New(fmt.Sprintln("please install kindlegen first. \r\nfrom https://www.amazon.com/gp/feature.html?ie=UTF8&docId=1000765211. \r\nerr:", err))
    }

    if tmpDir, err = ioutil.TempDir("", "go.test."); err != nil {
        return errors.New(fmt.Sprintln("create temp dir error! err:", err))
    }

    if err := os.Mkdir(tmpDir + "/assets", 0777); err != nil {
        fmt.Println("create assets dir error!")
        return errors.New(fmt.Sprintln("create assets dir error! err:", err))
    }

    return nil
}

func GetTmpPath() string {
    return tmpDir
}

func TmpClear() {
    if err := os.RemoveAll(tmpDir); err != nil {
        fmt.Sprintln("tmp dir remove error! err:", err)
    }
}

type jsonCus struct {
    url string
    method string
    body []byte
}

func (jc jsonCus) getJson(str interface{}) error {
    var body io.ReadCloser
    var err error

    if jc.method == "GET" {
        body, err = HttpGet(jc.url)
    } else if jc.method == "POST" {

    }

    if err != nil {
        return errors.New(fmt.Sprintln("getJson error:", err))
    }

    if err := json.NewDecoder(body).Decode(&str); err != nil {
        return errors.New(fmt.Sprintln("json decode", jc.url, "error:", err))
    }

    return nil
}

func Login(email string, passwd string) error {
    reqBody := map[string]string{"email":email, "password":passwd}
    reqBodyJson, _ := json.Marshal(reqBody)
    _, err := httpPost("https://www.oreilly.com/member/auth/login/", reqBodyJson)

    if err != nil {
        return err
    }

    return nil
}

func SaveHttpFile(baseUrl string, saveName string) error {
    body, err := HttpGet(baseUrl)
    if err != nil {
        return err
    }

    fileData, _ := ioutil.ReadAll(body)
    if err := SaveFile(saveName, fileData); err != nil {
        return err
    }

    return nil

}

func SaveFile(fullPath string, fileData []byte) error {
    return ioutil.WriteFile(tmpDir+"/"+fullPath, fileData, 0644)
}

func HttpGet(useUrl string) (io.ReadCloser, error) {
    httpCl := http.Client{
        Timeout: 60 * time.Second,
    }

    req, err := http.NewRequest("GET", useUrl, nil)
    if err != nil {
        return nil, errors.New(fmt.Sprintln("newRequest err:", err))
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
    for _, cookieItem := range userCookie {
        req.AddCookie(cookieItem)
    }

    res, err := httpCl.Do(req)
    if err != nil {
        return nil, errors.New(fmt.Sprintln("get", useUrl, "err:", err))
    }
    defer res.Body.Close()

    cpByte, _ := ioutil.ReadAll(res.Body)
    if res.StatusCode != 200 {
        return nil, errors.New(fmt.Sprintln("get", useUrl, "status code:", res.Status, "err:", string(cpByte)))
    }

    body := ioutil.NopCloser(bytes.NewBuffer(cpByte))

    return body, nil

}

func httpPost(useUrl string, reqBody []byte) (io.ReadCloser, error) {

    httpCl := http.Client{
        Timeout: 60 * time.Second,
    }

    req, err := http.NewRequest("POST", useUrl, bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, errors.New(fmt.Sprintln("newRequest err:", err))
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
        return nil, errors.New(fmt.Sprintln("get", useUrl, "err:", err))
    }
    defer res.Body.Close()

    userCookie = res.Cookies()
    cpByte, _ := ioutil.ReadAll(res.Body)
    if res.StatusCode != 200 {
        return nil, errors.New(fmt.Sprintln("The email address or password you've entered is incorrect.", "status code:", res.Status, "err:", string(cpByte)))
    }

    body := ioutil.NopCloser(bytes.NewBuffer(cpByte))

    return body, nil
}

func FileExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}