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
    "time"
)

var tmpDir = ""

func Init() error {

    err := error(nil)
    tmpDir, err = ioutil.TempDir("", "go.test.")
    if err != nil {
        return errors.New(fmt.Sprintln("create temp dir error! err:", err))
    }

    err = os.Mkdir(tmpDir + "/assets", 0777)
    if err != nil {
        fmt.Println("create assets dir error!")
        return errors.New(fmt.Sprintln("create assets dir error! err:", err))
    }

    return nil
}

func GetJson(url string, str interface{}) error {
    body, err := HttpGet(url)
    if err != nil {
        return errors.New(fmt.Sprintln("getJson error:", err))
    }

    err = json.NewDecoder(body).Decode(&str)
    if err != nil {
        return errors.New(fmt.Sprintln("json decode error:", err))
    }

    return nil
}

func SaveImage(baseUrl string, saveName string) error {
    body, err := HttpGet(baseUrl)
    if err != nil {
        return err
    }

    fileData, _ := ioutil.ReadAll(body)
    err = SaveFile(saveName, fileData)
    if err != nil {
        return err
    }

    return nil

}

func SaveFile(fullPath string, fileData []byte) error {
    return ioutil.WriteFile(tmpDir + "/" + fullPath, fileData, 0644)
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
    //req.Header.Add("", "")

    res, err :=httpCl.Do(req)
    if err != nil {
        return nil, errors.New(fmt.Sprintln("get", useUrl, "err:", err))
    }
    if res.StatusCode != 200 {
        return nil, errors.New(fmt.Sprintln("get", useUrl, "err:", res.Status))
    }
    defer res.Body.Close()

    cpByte, _ := ioutil.ReadAll(res.Body)
    body := ioutil.NopCloser(bytes.NewBuffer(cpByte))

    return body, nil

}