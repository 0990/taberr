package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/davyxu/golog"
	"github.com/davyxu/tabtoy/data"
	"github.com/davyxu/tabtoy/exportor"
	"github.com/davyxu/tabtoy/exportorv2"
)

var log *golog.Logger = golog.New("main")

// 开启调试信息
var paramDebugLevel = flag.Int("debug", 0, "show debug info")

// 并发导出,提高导出速度, 输出日志会混乱
var paramPara = flag.Bool("para", false, "parallel export by your cpu count")

// 显示版本号
var paramVersion = flag.Bool("version", false, "Show version")

// 工作模式
var paramMode = flag.String("mode", "", "mode: xls2pbt")

// 出现错误时暂停
var paramHaltOnError = flag.Bool("haltonerr", false, "halt on error")

// 输入协议二进制描述文件,通过protoc配合github.com/davyxu/pbmeta/protoc-gen-meta插件导出
var paramPbFile = flag.String("pb", "PB", "input protobuf binary descript file, export by protoc-gen-meta plugins")

// 输入电子表格文件
var paramXlsFile = flag.String("xls", "XLS", "input excel file, use ',' splited file list by multipy files")

// 输出文件夹
var paramOutDir = flag.String("outdir", "OUT_DIR", "output directory")

// 补丁文件
var paramPatch = flag.String("patch", "", "patch input files then output")

// 输出文件格式
var paramFormat = flag.String("fmt", "pbt", "output file format, support 'pbt', 'lua' ")

func main() {

	flag.Parse()

	// 版本
	if *paramVersion {
		fmt.Println("tabtoy 1.2.3")
		return
	}

	// 调试信息挂接命令行
	data.DebuggingLevel = *paramDebugLevel

	switch *paramMode {
	case "xls2pbt", "exportorv1":
		if !exportor.Run(exportor.Parameter{
			InputFileList: flag.Args(),
			PBFile:        *paramPbFile,
			PatchFile:     *paramPatch,
			Format:        *paramFormat,
			ParaMode:      *paramPara,
			OutDir:        *paramOutDir,
		}) {
			goto Err
		}
	case "exporortv2":
		if !exportorv2.Run(exportorv2.Parameter{
			InputFileList: flag.Args(),
			Format:        *paramFormat,
			ParaMode:      *paramPara,
			OutDir:        *paramOutDir,
		}) {
			goto Err
		}
	default:
		fmt.Println("--mode not specify")
		goto Err
	}

	return

Err:

	if *paramHaltOnError {
		halt()
	}

	os.Exit(1)
	return

}

func halt() {
	reader := bufio.NewReader(os.Stdin)

	reader.ReadLine()
}
