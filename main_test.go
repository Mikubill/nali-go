package main

import (
	"encoding/binary"
	"math/rand"
	"net"
	"testing"
	"time"
)

func BenchmarkAnalyseV4(b *testing.B) {
	Analyse("255.255.255.255")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Analyse(genIPv4())
	}
}

func BenchmarkAnalyseV6(b *testing.B) {
	Analyse("FFFF:FFFF:FFFF:FFFF::")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Analyse(genIPv6())
	}
}

func genIPv4() string {
	ip := net.IPv4zero
	binary.BigEndian.PutUint32(ip.To4(), rand.New(rand.NewSource(time.Now().UnixNano())).Uint32())
	return ip.String()
}

func genIPv6() string {
	ip := net.IPv6zero
	binary.BigEndian.PutUint64(ip.To16(), rand.New(rand.NewSource(time.Now().UnixNano())).Uint64())
	binary.BigEndian.PutUint64(ip.To16(), rand.New(rand.NewSource(time.Now().UnixNano())).Uint64())
	return ip.String()
}
