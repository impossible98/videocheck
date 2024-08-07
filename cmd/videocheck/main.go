package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var flagvar2 string
var flagvar3 string

func init() {
	flag.StringVar(&flagvar2, "file", "", "help message for flagname")
	flag.StringVar(&flagvar3, "url", "", "help message for flagname")
}

func main() {
	flag.Parse()

	// 正则获取BV号
	reg := `BV[a-zA-Z0-9]{10}`
	bvid := regexp.MustCompile(reg).FindString(flagvar3)
	res, err := http.Get("https://api.bilibili.com/x/web-interface/view?bvid=" + bvid)
	if err != nil {
		// handle error
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		// handle error
	}
	// body is []byte 转化为json
	//

	type VideoData struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Bvid      string `json:"bvid"`
			Aid       int    `json:"aid"`
			Videos    int    `json:"videos"`
			Title     string `json:"title"`
			Pubdate   int    `json:"pubdate"`
			Ctime     int    `json:"ctime"`
			Dynamic   string `json:"dynamic"`
			Cid       int    `json:"cid"`
			Dimension struct {
				Width  int `json:"width"`
				Height int `json:"height"`
				Rotate int `json:"rotate"`
			} `json:"dimension"`
			Premiere interface{} `json:"premiere"`
			Pages    []struct {
				Cid       int    `json:"cid"`
				Page      int    `json:"page"`
				From      string `json:"from"`
				Part      string `json:"part"`
				Duration  int    `json:"duration"`
				Vid       string `json:"vid"`
				Weblink   string `json:"weblink"`
				Dimension struct {
					Width  int `json:"width"`
					Height int `json:"height"`
					Rotate int `json:"rotate"`
				} `json:"dimension"`
			} `json:"pages"`
		} `json:"data"`
	}

	var video VideoData
	err = json.Unmarshal(body, &video)
	if err != nil {
		fmt.Println("解析 JSON 数据出错:", err)
		return
	}

	// 这里可以根据您的需求获取具体的值
	fmt.Println("视频的 bvid:", video.Data.Bvid)
	fmt.Println("视频的标题:", video.Data.Title)
	fmt.Println("视频的发布时间:", video.Data.Pubdate)

	// 复制文件
	newFileName := video.Data.Title + filepath.Ext(flagvar2)
	err = copyFile(flagvar2, newFileName)

	// 修改文件的创建时间
	err = os.Chtimes(newFileName, time.Unix(int64(video.Data.Pubdate), 0), time.Unix(int64(video.Data.Pubdate), 0))

	// 删除文件
	os.Remove(flagvar2)
}

func copyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
