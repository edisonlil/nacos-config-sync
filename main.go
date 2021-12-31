package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type NaocsConfigInfo struct {

	FileName string

	Content string

}

func main()  {

	fmt.Println("开始同步配置文件...")

	InitConfig()

	startUploadConfig(config.GetString("sync-config.config-addr"))
}

func startUploadConfig(configDir string)  {


	loginResult,_ := naocsLogin()

	files, err := GetAllFiles(configDir)

	if err != nil {
		fmt.Errorf("error : %s", err)
		return
	}

	for i := range files {

		var fileName = GetFileName(files[i])

		if fileName == ".gitlab-ci.yml" {
			continue
		}
		bytes, err := ioutil.ReadFile(files[i])

		if err != nil {
			fmt.Errorf("error : %s", err)
			return
		}

		//上传配置文件到naocsConfig
		uploadConfig(loginResult["accessToken"].(string), NaocsConfigInfo{
			FileName: fileName,
			Content: string(bytes),
		})

		fmt.Println(files[i])
	}
}

func uploadConfig(accessToken string,configInfo NaocsConfigInfo)  {
	naocsSaveConfigUrl := config.GetString("sync-config.naocs.configUrl")+"?accessToken="+ accessToken
	naocsUrl := config.GetString("sync-config.naocs.url")
	naocsAddr := config.GetString("sync-config.naocs.addr")
	client := &http.Client{}

	values := url.Values{}
	values.Set("dataId", configInfo.FileName)
	values.Set("group", config.GetString("sync-config.naocs.group"))
	values.Set("content", configInfo.Content)
	values.Set("tenant", config.GetString("sync-config.naocs.namespace"))
	values.Set("type",  config.GetString("sync-config.naocs.file-extension"))

	reqBody := ioutil.NopCloser(strings.NewReader(values.Encode()))

	req, err := http.NewRequest("POST", naocsSaveConfigUrl, reqBody)

	if err != nil {
		fmt.Println("error：" + err.Error())
	}

	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", naocsUrl)
	req.Header.Add("Referer", naocsAddr)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36")
	req.Header.Add("accessToken", accessToken)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error：" + err.Error())
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(string(respBody))
}

func naocsLogin () (map[string]interface{},error){

	values := url.Values{}
	values.Set("username",config.GetString("sync-config.naocs.username"))
	values.Set("password",config.GetString("sync-config.naocs.password"))

	reqBody := ioutil.NopCloser(strings.NewReader(values.Encode()))

	client := &http.Client{}


	loginUrl := config.GetString("sync-config.naocs.loginUrl")
	naocsUrl := config.GetString("sync-config.naocs.url")
	naocsAddr := config.GetString("sync-config.naocs.addr")
	req,err := http.NewRequest("POST",loginUrl,reqBody)

	req.Header.Add("Accept","application/json,text/plain,*/*")
	req.Header.Add("Accept-Encoding","gzip,deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", naocsUrl)
	req.Header.Add("Referer", naocsAddr)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36")

	resp,err := client.Do(req)


	defer resp.Body.Close()

	respBody,err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error",err)
	}

	var result map[string]interface{}

	err = json.Unmarshal([]byte(respBody),&result)

	fmt.Println(string(respBody))

	return result,nil
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

func GetFileName(file string) string  {

	var separator string

	if strings.Contains(file,"\\") {
		separator = "\\"
	}else {
		separator = "/"
	}

	splits := strings.Split(file,separator)

	return splits[len(splits)-1]

}