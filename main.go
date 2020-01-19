package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	version = "NaLi-Go 1.3.0\n" +
		"Source: https://github.com/Mikubill/nali-go\n" +
		"Git Commit Hash: %s\n"
	helper = "Usage: %s <command> [options] \n" +
		"\nOptions:" +
		"\n  -v, --version  版本信息" +
		"\n  -h, --help     帮助信息\n" +
		"\nCommands:" +
		"\n  IP Address     解析 stdin 或参数中的 IP 信息 (默认)" +
		"\n  update         更新 IP 库" +
		"\n  delete         删除 IP 库数据\n"
)

var (
	commands = []string{"dig", "ping", "traceroute", "tracepath", "nslookup", "mtr"}
	help     = []string{"help", "--help", "-h", "h"}
	ver      = []string{"version", "--version", "-v", "v"}
	githash  = ""
	v4Data   = fileData{}
	//v6Data = fileData{}
	v4db = pointer{Data: &v4Data}
	//v6db = pointer{Data: &IPv6Data}
)

func main() {

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	//IPv6 （未启用）
	//res = IPv6Data.InitIPData("https://github.com/Mikubill/nali-go/raw/master/ipv6wry.db", "ipv6.dat")
	//if v, ok := res.(error); ok { panic(v) }

	if args := os.Args; len(args) > 1 {

		cmd()

		for i := range args {
			item := args[i]
			if strings.Contains(item, os.Args[0]) == false {
				fmt.Println(analyse(item))
			}
		}
		os.Exit(0)
	}

	if (info.Mode() & os.ModeCharDevice) != 0 {
		self := os.Args[0]
		fmt.Printf(helper, self)
		os.Exit(0)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		fmt.Printf("%s", analyse(line))
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	//scanner := bufio.NewScanner(os.Stdin)
	//for scanner.Scan() {
	//	t := scanner.Text()
	//	_, _ = fmt.Fprint(os.Stdout, analyse(t, v4db, v6db))
	//}
	//
	//if err := scanner.Err(); err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
}

func analyse(item string) string {
	// ipv4
	if v4db.Data.Data == nil {
		res := v4db.Data.InitIPData("https://qqwry.mirror.noc.one/QQWry.Dat", "ipv4.dat")
		if v, ok := res.(error); ok {
			panic(v)
		}
	}
	re4 := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	if ip := re4.FindStringSubmatch(item); len(ip) != 0 {
		res := v4db.find(ip[0])
		result := ip[0] + " " + "\x1b[0;0;36m[" + res.Country + res.Area + "]\x1b[0m"
		return strings.ReplaceAll(item, ip[0], result)
	}
	return item

	//ipv6
	//re6 := regexp.MustCompile(`[a-fA-F0-9:]+`)
	//if ip := re6.FindStringSubmatch(item); len(ip) != 0 {
	//	res := v6db.find(ip[0])
	//	result := res.IP + " " + "[" +  res.Country + " " + res.Area + "]"
	//	fmt.Println(strings.ReplaceAll(item, ip[0], result))
	//} else {
	//	fmt.Println(item)
	//}
}

func contains(array []string, flag string) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == flag {
			return true
		}
	}
	return false
}

func execute(cmd string) {
	runner := exec.Command("sh", "-c", cmd)
	fmt.Println(runner.Args)
	stdout, err := runner.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = runner.Start()
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		fmt.Printf("%s", line)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}

func cmd() {
	switch {
	case contains(commands, os.Args[1]):
		execute(strings.Join(os.Args[1:], " ") + " | " + os.Args[0])
		os.Exit(0)
	case contains(help, os.Args[1]):
		fmt.Printf(helper, os.Args[0])
		os.Exit(0)
	case contains(ver, os.Args[1]):
		fmt.Printf(version, githash)
		os.Exit(0)
	case os.Args[1] == "update":
		update()
		os.Exit(0)
	case os.Args[1] == "delete":
		del()
		os.Exit(0)
	}
}

func update() {
	_, err := os.Stat("ipv4.dat")
	if err == nil || os.IsExist(err) {
		var str string
		fmt.Print("确定要更新数据库嘛（此操作会删除原有数据）? [Y/n]")
		_, err = fmt.Scanln(&str)
		if err != nil || (str != "Y" && str != "y") {
			fmt.Println("Cancelled.")
			return
		}
		err := os.Remove("ipv4.dat")
		if err == nil {
			log.Println("数据文件已清理。正在重新下载...")
		}
	}
	res := v4Data.InitIPData("https://qqwry.mirror.noc.one/QQWry.Dat", "ipv4.dat")
	if v, ok := res.(error); ok {
		panic(v)
	}
}

func del() {
	_, err := os.Stat("ipv4.dat")
	if err == nil || os.IsExist(err) {
		var str string
		fmt.Print("确定要删除数据库嘛? [Y/n]")
		_, err = fmt.Scanln(&str)
		if err != nil || (str != "Y" && str != "y") {
			fmt.Println("Cancelled.")
			return
		}
		err := os.Remove("ipv4.dat")
		if err == nil {
			log.Println("数据文件已清理。")
		}
	}
}
