<h1 align="center">Nali Go</h1>

<p align="center">给IP加上地理信息的命令行小工具</p>

<p align="center">
<a title="Release" target="_blank" href="https://github.com/Mikubill/nali-go/releases"><img src="https://img.shields.io/github/release/Mikubill/nali-go.svg?style=flat-square&hash=c7"></a>
<a title="Go Report Card" target="_blank" href="https://goreportcard.com/report/github.com/Mikubill/nali-go"><img src="https://goreportcard.com/badge/github.com/Mikubill/nali-go?style=flat-square"></a>
</p>

本项目支持IPv4（纯真IP数据库）和IPv6（ZX公网IPv6库）。

## 编译 说明

安装 Go，然后设置一下GOPROXY：

```bash
➜ export GOPROXY=https://goproxy.cn
```

clone本项目，然后编译即可。

```bash
➜ go build -ldflags "-w -s -X main.githash=`git rev-parse HEAD`" .
```

## 下载/运行 说明

Go语言程序, 可直接在[发布页](https://github.com/Mikubill/nali-go/releases)下载使用，也可以使用下面的命令安装:

```bash
go get -u https://github.com/Mikubill/nali-go
```

Query simple IP address:

```bash
➜ nali 2.3.6.7 1.1.2.5

2.3.6.7 [法国 Orange]
1.1.2.5 [福建省 电信]
```

Query IP addresses from `stdin`:

```bash
➜  dig github.io +short | nali 

185.199.110.153 [美国 GitHub+Fastly节点]
185.199.111.153 [美国 GitHub+Fastly节点]
185.199.109.153 [美国 GitHub+Fastly节点]
185.199.108.153 [美国 GitHub+Fastly节点]
```

Use Nali CLI built-in tools shortcut:

```bash
nali nslookup ip.sb

Server:         1.0.0.1 [美国 APNIC&CloudFlare 公共 DNS 服务器]
Address:        1.0.0.1 [美国 APNIC&CloudFlare 公共 DNS 服务器]#53

Non-authoritative answer:
Name:   ip.sb
Address: 119.9.95.61 [香港 Rackspace Hosting公司]
```

Update IP Database:

```bash
nali update
```

Delete IP Database:

```bash
nali delete
```

Check Version:

```bash
nali version

NaLi-Go 
Source: https://github.com/Mikubill/nali-go
Git Commit Hash: 61e7869a02dc88c28093fac5a5aa35c06ef18333
```

## Related

- [Sukka - NaLi-CLI](https://github.com/sukkaw/nali-cli)
- [CZ88 QQIP 数据库](http://www.cz88.net/fox/ipdat.shtml) 纯真网络提供的免费离线 IP 数据库
- [QQWry Mirror](https://qqwry.mirror.noc.one) Just a mirror of qqwry ipdb

