package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

type NaocsContentInfo struct {
	FileName string

	Content string
}

func main() {

	log.Infoln("开始同步配置文件...")
	log.Infoln("配置文件地址：", env.GetString("project_addr"))
	err := startUploadConfig(env.GetString("project_addr"))

	if err != nil {
		log.Errorf("同步配置文件失败...")
	} else {
		log.Infoln("同步配置文件完成...")
	}

}

func startUploadConfig(configDir string) error {

	loginResult, err := naocsLogin()

	if err != nil {
		return err
	}

	files, err := GetAllFiles(configDir)

	if err != nil {
		log.Errorf("error : %s", err)
		return err
	}

	var wgr sync.WaitGroup

	for i := range files {

		var fileName = GetFileName(files[i])

		if fileName == ".gitlab-ci.yml" || fileName == "sync-config.yml" {
			continue
		}

		bytes, err := ioutil.ReadFile(files[i])

		if err != nil {
			log.Errorf("error : %s", err)
			return err
		}

		wgr.Add(1)

		//上传配置文件到naocsConfig
		go uploadConfig(loginResult["accessToken"].(string), NaocsContentInfo{
			FileName: fileName,
			Content:  string(bytes),
		})

		wgr.Done()

		log.Infoln(files[i])
	}

	wgr.Wait()

	return nil
}

func uploadConfig(accessToken string, info NaocsContentInfo) {

	nacosConfig := env.GetNacosConfig()

	naocsSaveConfigUrl := nacosConfig.ConfigUrl + "?accessToken=" + accessToken
	client := &http.Client{}

	values := url.Values{}
	values.Set("dataId", info.FileName)
	values.Set("group", nacosConfig.Group)
	values.Set("content", info.Content)
	values.Set("tenant", nacosConfig.Namespace)
	values.Set("type", nacosConfig.FileExtension)

	reqBody := ioutil.NopCloser(strings.NewReader(values.Encode()))

	req, err := http.NewRequest("POST", naocsSaveConfigUrl, reqBody)

	if err != nil {
		log.Errorf("error：%s", err.Error())
	}

	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", nacosConfig.Url)
	req.Header.Add("Referer", nacosConfig.Addr)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36")
	req.Header.Add("accessToken", accessToken)

	resp, err := client.Do(req)

	if err != nil {
		log.Infoln("error：" + err.Error())
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorf("error: %s", err)
	}

	log.Infoln(string(respBody))
}

func naocsLogin() (map[string]interface{}, error) {

	nacosConfig := env.GetNacosConfig()

	values := url.Values{}
	values.Set("username", nacosConfig.Username)
	values.Set("password", nacosConfig.Password)

	reqBody := ioutil.NopCloser(strings.NewReader(values.Encode()))

	client := &http.Client{}

	req, err := http.NewRequest("POST", nacosConfig.LoginUrl, reqBody)

	req.Header.Add("Accept", "application/json,text/plain,*/*")
	req.Header.Add("Accept-Encoding", "gzip,deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", nacosConfig.Url)
	req.Header.Add("Referer", nacosConfig.Addr)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36")

	resp, err := client.Do(req)

	if err != nil {
		log.Errorf("登录nacos失败,%s", err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorf("%s", err.Error())
	}

	var result map[string]interface{}

	err = json.Unmarshal([]byte(respBody), &result)

	log.Infoln(string(respBody))

	return result, nil
}

//GetAllFiles 获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".yaml") || strings.HasSuffix(fi.Name(), ".yml")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}

func GetFileName(file string) string {

	var separator string

	if strings.Contains(file, "\\") {
		separator = "\\"
	} else {
		separator = "/"
	}

	splits := strings.Split(file, separator)

	return splits[len(splits)-1]

}
