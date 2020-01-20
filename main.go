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
	version = "NaLi-Go 1.5.0\n" +
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
	refresh = "Usage: %s update <command> [options] \n" +
		"\nOptions:" +
		"\n  --force        强制执行，不询问\n" +
		"\nCommands:" +
		"\n  ipv4           更新IPv4数据" +
		"\n  ipv6           更新IPv6数据\n"
	remove = "Usage: %s delete <command> [options] \n" +
		"\nOptions:" +
		"\n  --force        强制执行，不询问\n" +
		"\nCommands:" +
		"\n  ipv4           删除IPv4数据" +
		"\n  ipv6           删除IPv6数据\n"
)

var (
	commands = []string{"dig", "ping", "traceroute", "tracepath", "nslookup"}
	help     = []string{"help", "--help", "-h", "h"}
	ver      = []string{"version", "--version", "-v", "v"}
	githash  = ""
	v4Data   = fileData{}
	v6Data   = fileData{}
	v4db     = pointer{Data: &v4Data}
	v6db     = pointer{Data: &v6Data}
)

func main() {

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if args := os.Args; len(args) > 1 {

		cmd()

		for i := range args {
			item := args[i]
			if strings.Contains(item, os.Args[0]) == false {
				fmt.Println(analyse(item))
			}
		}

		//fmt.Printf("\n")
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
	//fmt.Printf("\n")
}

func analyse(item string) string {
	// ipv4, https://stackoverflow.com/questions/53497/regular-expression-that-matches-valid-ipv6-addresses/17871737#17871737
	re4 := regexp.MustCompile(`((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9])`)
	if ip := re4.FindStringSubmatch(item); len(ip) != 0 {
		if v4db.Data.Data == nil {
			res := v4Data.InitIPData("https://qqwry.mirror.noc.one/qqwry.rar", "ipv4.dat", 5252)
			if v, ok := res.(error); ok {
				panic(v)
			}
		}
		res := v4db.findv4(ip[0])
		result := ip[0] + " " + "\x1b[0;0;36m[" + res.Country + res.Area + "]\x1b[0m"
		return strings.ReplaceAll(item, ip[0], result)
	}

	// ipv6, https://github.com/lilydjwg/winterpy/blob/master/pyexe/ipmarkup
	re6 := regexp.MustCompile(`fe80:(:[0-9a-fA-F]{1,4}){0,4}(%\w+)?|([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?`)
	if ip := re6.FindStringSubmatch(item); len(ip) != 0 {
		if v6db.Data.Data == nil {
			res := v6db.Data.InitIPData("https://cdn.jsdelivr.net/gh/Mikubill/nali-go@1.3.0/ipv6wry.db", "ipv6.dat", 1951)
			if v, ok := res.(error); ok {
				panic(v)
			}
		}
		res := v6db.findv6(ip[0])
		result := res.IP + " " + "\x1b[0;0;36m[" + res.Country + res.Area + "]\x1b[0m"
		return strings.ReplaceAll(item, ip[0], result)
	}
	return item
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
		_version()
		os.Exit(0)
	case os.Args[1] == "update":
		_update(false, contains(os.Args, "--force"))
		os.Exit(0)
	case os.Args[1] == "delete":
		_update(true, contains(os.Args, "--force"))
		os.Exit(0)
	}
}

func _version() {
	fmt.Printf(version, githash)
	if _, err := os.Stat("ipv4.dat"); err == nil || os.IsExist(err) {
		ver := strings.ReplaceAll(analyse("255.255.255.255"), "255.255.255.255", "")
		fmt.Printf("IPv4 Version： %s\n", ver)
	} else {
		fmt.Printf("IPv4 Version： Database Not Found.\n")
	}

	if _, err := os.Stat("ipv6.dat"); err == nil || os.IsExist(err) {
		ver := strings.ReplaceAll(analyse("FFFF:FFFF:FFFF:FFFF::"), "FFFF:FFFF:FFFF:FFFF::", "")
		fmt.Printf("IPv6 Version： %s\n", ver)
	} else {
		fmt.Printf("IPv6 Version： Database Not Found.\n")
	}
}

func ipv4Update(force bool, del bool) {
	filename := "ipv4.dat"
	if _, err := os.Stat(filename); err == nil || os.IsExist(err) {
		if force != true {
			if del != true {
				question("删除现有 IPv4 数据库", "正在删除 IPv4 数据库...")
			} else {
				question("更新现有 IPv4 数据库", "")
			}
		}
		if err := os.Remove(filename); err == nil {
			log.Printf("数据文件 %s 已清理。", filename)
		}
	} else {
		log.Printf("没有待清理的数据文件。")
	}
	if del != true {
		res := v4Data.InitIPData("https://qqwry.mirror.noc.one/qqwry.rar", filename, 5252)
		if v, ok := res.(error); ok {
			panic(v)
		}
	}
}

func ipv6Update(force bool, del bool) {
	filename := "ipv6.dat"
	if _, err := os.Stat(filename); err == nil || os.IsExist(err) {
		if force != true {
			if del == true {
				question("删除现有 IPv6 数据库", "正在删除 IPv6 数据库...")
			} else {
				question("更新现有 IPv6 数据库", "")
			}
		}
		if err := os.Remove(filename); err == nil {
			log.Printf("数据文件 %s 已清理。", filename)
		}
	} else {
		log.Printf("没有待清理的数据文件。")
	}
	if del != true {
		res := v6Data.InitIPData("https://cdn.jsdelivr.net/gh/Mikubill/nali-go@1.3.0/ipv6wry.db", "ipv6.dat", 1951)
		if v, ok := res.(error); ok {
			panic(v)
		}
	}
}

func allUpdate(force bool, del bool) {
	_, er1 := os.Stat("ipv4.dat")
	_, er2 := os.Stat("ipv6.dat")
	if er1 == nil || er2 == nil {
		if force != true {
			if del == true {
				question("删除现有所有 IP 数据库", "正在删除数据库...")
			} else {
				question("更新现有所有 IP 数据库", "")
			}
		}
		if err := os.Remove("ipv4.dat"); err == nil {
			log.Printf("数据文件 %s 已清理。", "ipv4.dat")
		}
		if err := os.Remove("ipv6.dat"); err == nil {
			log.Printf("数据文件 %s 已清理。", "ipv6.dat")
		}
	} else {
		log.Printf("没有待清理的数据文件。")
	}
	if del != true {
		res := v4Data.InitIPData("https://qqwry.mirror.noc.one/qqwry.rar", "ipv4.dat", 5252)
		if v, ok := res.(error); ok {
			panic(v)
		}

		res = v6Data.InitIPData("https://cdn.jsdelivr.net/gh/Mikubill/nali-go@1.3.0/ipv6wry.db", "ipv6.dat", 1951)
		if v, ok := res.(error); ok {
			panic(v)
		}
	}
}

func updateTip(del bool) {
	if del != true {
		fmt.Printf(refresh, os.Args[0])
	} else {
		fmt.Printf(remove, os.Args[0])
	}
}

func _update(del bool, force bool) {
	switch {
	case len(os.Args) < 3:
		updateTip(del)
	case os.Args[2] == "ipv4":
		ipv4Update(force, del)
	case os.Args[2] == "ipv6":
		ipv6Update(force, del)
	case os.Args[2] == "all":
		allUpdate(force, del)
	default:
		updateTip(del)
	}
	os.Exit(0)

}

func question(action string, result string) {
	var str string
	var err error
	fmt.Printf("确定要 %s 嘛（此操作会影响原有数据）? [Y/n]", action)
	_, err = fmt.Scanln(&str)
	if err != nil || (str != "Y" && str != "y") {
		fmt.Println("Cancelled.")
		os.Exit(0)
	}
	if err == nil {
		log.Printf("%s", result)
	}
}
