package main

import (
	"flag"
	"os"
	"log"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
}

/*
运行程序的时候加一个--cpuprofile参数，比如fabonacci --cpuprofile=fabonacci.prof
这样程序运行的时候的cpu信息就会记录到XXX.prof中了。
下一步就可以使用这个prof信息做出性能分析图了（需要安装graphviz）。
使用go tool pprof (应用程序) （应用程序的prof文件）
进入到pprof，使用web命令就会在/tmp下生成svg文件，svg文件是可以在浏览器下看的
*/
