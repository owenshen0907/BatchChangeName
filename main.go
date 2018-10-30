// BatchChangeName project main.go
package main

import (
	//	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/larspensjo/config"
)

var (
	configFile = flag.String("configfile", "Sender_config.ini", "General configuration file")
	Version    = "BatchChangeName V1.0.20180814 "
	Auther     = "Owen Shen"
)

func getArgs() {
	version := flag.Bool("v", false, "version")
	flag.Parse()
	if *version {
		fmt.Println("Version：", Version)
		fmt.Println("Auther:", Auther)
		return
	}
}
func main() {

	if len(os.Args) > 1 {
		getArgs()
	} else {
		body()
	}

}

func body() {
	var TOPIC = readconfigfile()

	//	logFile, _ := os.OpenFile(TOPIC["fulename"], os.O_RDWR|os.O_CREATE, 0666)
	fmt.Println("Hello World!")
	var pwd, filename string

	filename = TOPIC["filename"]
	pwd = getCurrentDirectory()
	fmt.Println(TOPIC["filename"])
	//	pwd = pwd + "\\" + TOPIC["source"]

	files, _ := ListDir(pwd+"\\"+TOPIC["source"], TOPIC["stuffix"])
	fmt.Println(len(files))
	fmt.Println(files)
	//	var tmp, tmpv []string
	if checkFileIsExist(pwd + "\\" + filename) {
		fmt.Println("文件已经存在，我就不删除了，你自己把已经有的内容清掉吧。哈哈！")
		//		del := os.Remove(filename)
		//		if del != nil {
		//			fmt.Println(del)
		//		}
	}

	//新增内容
	//f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0777)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	erro(err)
	SEEK_END, _ := f.Seek(0, os.SEEK_END)
	_, _ = f.WriteAt([]byte("\r\n"), SEEK_END)
	//ExportFileName := generate(logFile)
	defer f.Close()
	explain := "#####请删除掉此段文本以上的内容哦！##"
	explain1 := "  Please delete the content above, ok?"
	f.WriteString(explain + "\r\n")
	f.WriteString(explain1 + "\r\n")
	for _, v := range files {
		fmt.Println(v)
		tmp := strings.Split(v, "\\")

		fmt.Println(tmp[len(tmp)-1])
		f.WriteString(tmp[len(tmp)-1] + ",\r\n")
		//i := len(tmp) - 1
		//tmpv = append(tmpv, tmp[i])
		//tmp = append(tmp, v)
		//	w.Write(tmpv)
		//	tmpv = nil
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "\\", -1)
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10) //初始化file切片，预留十个位置
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)                                                     //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		//if err != nil { //忽略错误
		// return err
		//}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func readconfigfile() (TOPIC map[string]string) {
	TOPIC = make(map[string]string)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//set config file std
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		log.Fatalf("Fail to find", *configFile, err)
	}
	//set config file std End

	//Initialized topic from the configuration
	if cfg.HasSection("topicArr") {
		section, err := cfg.SectionOptions("topicArr")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("topicArr", v)
				if err == nil {
					TOPIC[v] = options
				}
			}
		}
	}
	//Initialized topic from the configuration END
	return TOPIC
}

func erro(err error) {
	if err != nil {
		panic(err)
	}
}
