package main

import (
	"encoding/binary"
	"golang.org/x/text/encoding/simplifiedchinese"
	"math/big"
	"net"
	"os"
	"strings"
)

var (
	country []byte
	area    []byte
	off     uint32
	numip   uint32
	offset  uint32
	start   uint32
	end     uint32
	_cmp    int
)

type result struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Area    string `json:"area"`
}

type fileData struct {
	Data []byte
	Path *os.File
}

type pointer struct {
	Data     *fileData
	Offset   uint32
	ItemLen  uint32
	IndexLen uint32
}

func (q *pointer) readData(length uint32) (rs []byte) {
	end := q.Offset + length
	dataNum := uint32(len(q.Data.Data))
	if q.Offset > dataNum {
		return nil
	}

	if end > dataNum {
		end = dataNum
	}
	rs = q.Data.Data[q.Offset:end]
	q.Offset = end
	return rs
}

func (q *pointer) find(ip string) (res result) {

	res = result{}
	res.IP = ip
	q.Offset = 0

	if strings.Contains(ip, ":") {
		tp := big.NewInt(0)
		op := big.NewInt(0)
		tp.SetBytes(net.ParseIP(ip).To16())
		op.SetString("18446744073709551616", 10)
		op.Div(tp, op)
		tp.SetString("FFFFFFFFFFFFFFFF", 16)
		op.And(op, tp)

		offset = q.searchIndexV6(op) - q.ItemLen
	} else {
		numip = binary.BigEndian.Uint32(net.ParseIP(ip).To4())
		offset = q.searchIndexV4(numip)
	}

	if offset <= 0 {
		return
	}

	q.Offset = offset + q.ItemLen
	mode := q.readData(1)[0]
	if mode == 0x01 {
		q.Offset = byteToUInt32(q.readData(3))
		mode = q.readData(1)[0]
		if mode == 0x02 {
			off = q.Offset + 3
			q.Offset = byteToUInt32(q.readData(3))
			country = q.readString()
		} else {
			q.Offset--
			off = q.Offset
			country = q.readString()
			off += uint32(len(country) + 1)
		}
		area = q.readArea(off)
	} else if mode == 0x02 {
		q.Offset = byteToUInt32(q.readData(3))
		country = q.readString()
		area = q.readArea(offset + 8)
	} else {
		q.Offset = offset + 4
		country = q.readString()
		area = q.readArea(offset + uint32(5+len(country)))
	}

	if strings.Contains(ip, ":") {
		res.Country = string(country)
		res.Area = string(area)
	} else {
		enc := simplifiedchinese.GBK.NewDecoder()
		res.Country, _ = enc.String(string(country))
		res.Area, _ = enc.String(string(area))
	}

	// Delete CZ88.NET (防止不相关的信息产生干扰）
	if res.Area == " CZ88.NET" || res.Area == "" {
		res.Area = ""
	} else {
		res.Area = " " + res.Area
	}

	return
}

func (q *pointer) readArea(offset uint32) []byte {
	q.Offset = offset
	mode := q.readData(1)[0]
	if mode == 0x01 || mode == 0x02 {
		areaOffset := byteToUInt32(q.readData(3))
		if areaOffset == 0 {
			return []byte("")
		}
		q.Offset = areaOffset
		return q.readString()
	}
	q.Offset = offset
	return q.readString()
}

func (q *pointer) readString() []byte {
	data := make([]byte, 0)
	for {
		buf := q.readData(1)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

func (q *pointer) searchIndexV4(ip uint32) uint32 {

	q.ItemLen = 4
	q.IndexLen = 7
	start = binary.LittleEndian.Uint32(q.Data.Data[0:4])
	end = binary.LittleEndian.Uint32(q.Data.Data[4:8])

	buf := make([]byte, q.IndexLen)

	for {
		mid := start + q.IndexLen*(((end-start)/q.IndexLen)>>1)
		buf = q.Data.Data[mid : mid+q.IndexLen]
		_ip := binary.LittleEndian.Uint32(buf[:q.ItemLen])

		if end-start == q.IndexLen {
			return byteToUInt32(buf[q.ItemLen:])
		}

		if _ip > ip {
			end = mid
		} else if _ip < ip {
			start = mid
		} else if _ip == ip {
			return byteToUInt32(buf[q.ItemLen:])
		}
	}
}

func (q *pointer) searchIndexV6(ip *big.Int) uint32 {

	q.ItemLen = 8
	q.IndexLen = 11
	start = binary.LittleEndian.Uint32(q.Data.Data[16:23])
	end = start + binary.LittleEndian.Uint32(q.Data.Data[8:15])

	buf := make([]byte, q.IndexLen)
	_ip := big.NewInt(0)

	for {
		mid := start + q.IndexLen*(((end-start)/q.IndexLen)>>1)
		buf = q.Data.Data[mid : mid+q.IndexLen]
		_ip.SetBytes(buf[:q.ItemLen])
		_cmp = _ip.Cmp(ip)

		if _cmp == 1 {
			end = mid
		} else if _cmp == -1 {
			start = mid
		} else if _cmp == 0 {
			return byteToUInt32(buf[q.ItemLen:])
		}
	}
}
func byteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
