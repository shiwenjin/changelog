package main

import (
	"log"
	"os/exec"
	"strings"

	//"bytes"
	"fmt"
	"os"
)

/*
feat： 新增 feature
fix: 修复 bug
docs: 仅仅修改了文档，比如 README, CHANGELOG, CONTRIBUTE等等
style: 仅仅修改了空格、格式缩进、逗号等等，不改变代码逻辑
refactor: 代码重构，没有加新功能或者修复 bug
perf: 优化相关，比如提升性能、体验
test: 测试用例，包括单元测试、集成测试等
chore: 改变构建流程、或者增加依赖库、工具等
revert: 回滚到上一个版本
*/

var keyworlds = []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "chore", "revert", "others"}
var logs = make(map[string]map[string][]string)

func main() {

	//var m map[string][]string
	//date:keywrld:values[]

	//sort.Strings(keyworlds)
	cmd := exec.Command("git", "log", "--date=format:'%Y-%m-%d'", "--pretty=format:'@%h@%cd@%B")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	strs := strings.Split(string(out), "\n")
	for i := range strs {
		if len(strs[i]) < 23 {
			continue
		}
		if strs[i][1] != '@' {
			continue
		}
		if strs[i][9] != '@' {
			continue
		}
		if strs[i][22] != '@' {
			continue
		}
		buildLog(strs[i][11:21], strs[i][23:])
	}
	writeFile()
}
func buildLog(date, title string) {
	index, logKey := getLogKey(title)
	dateLogs := logs[date]
	if dateLogs == nil {
		dateLogs = make(map[string][]string)
		logs[date] = dateLogs
	}
	logValues := dateLogs[logKey]
	if logValues == nil {
		logValues = []string{}
	}
	logValues = append(logValues, getLog(index, title))
	dateLogs[logKey] = logValues
}
func getLog(sindex int, title string) string {
	if sindex == -1 {
		return title
	}
	return title[sindex+1:]
}
func getLogKey(title string) (int, string) {
	index := strings.Index(title, ":")
	if index == -1 {
		return index, "others"
		//logs
	}
	title = title[:index]
	for _, v := range keyworlds {
		if v == title {
			return index, v
		}
	}
	index = -1
	return index, "others"
}
func writeFile() {
	filePath := "CHANGELOG.MD"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		fmt.Printf("%s", err)
	}
	for key, value := range logs {
		file.WriteString("## " + key)
		file.WriteString("\n")
		writeKeyWorldsLogDetails(value, file)
	}

}
func writeKeyWorldsLogDetails(logdetails map[string][]string, file *os.File) {
	for _, keyworld := range keyworlds {
		if len(logdetails[keyworld]) == 0 {
			continue
		}
		file.WriteString("### " + keyworld)
		file.WriteString("\n")
		writeLogDetails(logdetails[keyworld], file)

	}
	//for _, value := range logdetails {
	//	fmt.Println(value)
	//}
	//for key, value := range logdetails {
	//	file.WriteString("### " + key)
	//	file.WriteString("\n")
	//	writeLogDetails(value, file)
	//}
}
func writeLogDetails(logdetails []string, file *os.File) {
	for _, value := range logdetails {
		file.WriteString("- " + value)
		file.WriteString("\n")

	}
}
