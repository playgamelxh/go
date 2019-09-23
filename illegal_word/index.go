package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

//定义全局字符串类型切片
var words []string

func main() {
	fmt.Println("星矿违禁词过滤")

	fmt.Println("加载违禁词")
	//违禁词文件
	var path string = "违禁词.csv"
	//定义字符串类型切片
	//	var words []string
	//将文件读取到切片
	words = getFile(path)
	//格式化输出切片大小
	//	fmt.Printf("len = %d\n", len(words))

	//监听8092的post请求
	http.HandleFunc("/", parseHttp)
	err := http.ListenAndServe(":8092", nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}

//读取文件内容
func getFile(path string) []string {

	var words []string

	file, err := os.Open(path)
	if err != nil {
		return words
	}
	defer file.Close()
	fd := bufio.NewReader(file)
	for {
		str, err := fd.ReadString('\n')
		if err != nil {
			break
		}
		words = append(words, str)
	}
	return words
}

//处理http请求
func parseHttp(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	str := r.Form.Get("str")

	log.Println("Post form-urlencoded: str = %s", str)

	//全局违禁词切片
	fmt.Printf("len = %d\n", len(words))
	//遍历切片
	for i := 0; i < len(words); i++ {
		var list []string = strings.Split(words[i], "\t")
		//违禁词处理
		str = str_replace(str, list)
	}
	fmt.Println(str)

	fmt.Fprintf(w, `{"code":0,"data":{"str":"`+str+`"}}`)
}

//违禁词处理
func str_replace(str string, list []string) string {
	//违禁词
	var illegalWord string = strings.Replace(list[1], " ", "", -1)
	illegalWord = strings.Replace(illegalWord, "\t", "", -1)
	//相近词
	var variant string = list[2]
	//正选词
	var positiveSelectionWord = list[7]
	//反选词
	var antiSelectionWord = list[8]

	//相近词替换
	if variant != "" {
		var variantList []string = strings.Split(variant, ",")
		for i := 0; i < len(variantList); i++ {
			str = strings.Replace(str, variantList[i], str_replace_start(variantList[i]), -1)
		}
	}
	//违禁词+正选词替换
	if positiveSelectionWord != "" {
		var positiveList []string = strings.Split(positiveSelectionWord, ",")
		for i := 0; i < len(positiveList); i++ {
			str = strings.Replace(str, illegalWord+positiveList[i], str_replace_start(illegalWord+positiveList[i]), -1)
		}
	}
	//违禁词+反选词
	if antiSelectionWord != "" {
		var antiList []string = strings.Split(antiSelectionWord, ",")
		for i := 0; i < len(antiList); i++ {
			var index int = strings.Index(antiList[i], illegalWord)
			if index == -1 {
				str = strings.Replace(str, illegalWord+antiList[i], "/*"+string(i)+"*/", -1)
			} else {
				str = strings.Replace(str, antiList[i], "/*"+string(i)+"*/", -1)
			}
		}
	}

	//没有正选的违禁词替换
	if positiveSelectionWord == "" {
		str = strings.Replace(str, illegalWord, str_replace_start(illegalWord), -1)
	}
	//违禁词+反选词还原
	if antiSelectionWord != "" {
		var antiList []string = strings.Split(antiSelectionWord, ",")
		for i := 0; i < len(antiList); i++ {
			var index int = strings.Index(antiList[i], illegalWord)
			if index == -1 {
				str = strings.Replace(str, "/*"+string(i)+"*/", illegalWord+antiList[i], -1)
			} else {
				str = strings.Replace(str, "/*"+string(i)+"*/", antiList[i], -1)
			}
		}
	}

	return str
}

//获取字符串长度并返回对应长度的*号
func str_replace_start(str string) string {
	var strLen int = strings.Count(str, "")
	var returnStr string = ""
	for i := 1; i < strLen; i++ {
		returnStr += "*"
	}
	return returnStr
}
