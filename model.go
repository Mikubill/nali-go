package main

import (
	"encoding/binary"
	"golang.org/x/text/encoding/simplifiedchinese"
	"math/big"
	"net"
	"os"
)

var (
	header  []byte
	country []byte
	area    []byte
	v4ip    uint32
	v6ip    uint64
	offset  uint32
	start   uint32
	end     uint32
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

func (q *pointer) findv4(ip string) (res result) {

	res = result{}
	res.IP = ip
	q.Offset = 0

	v4ip = binary.BigEndian.Uint32(net.ParseIP(ip).To4())
	offset = q.searchIndexV4(v4ip)
	q.Offset = offset + q.ItemLen

	enc := simplifiedchinese.GBK.NewDecoder()
	country, area = q.getAddr()
	res.Country, _ = enc.String(string(country))
	res.Area, _ = enc.String(string(area))

	// Delete CZ88.NET (防止不相关的信息产生干扰）
	if res.Area == " CZ88.NET" || res.Area == "" {
		res.Area = ""
	} else {
		res.Area = " " + res.Area
	}

	return
}

func (q *pointer) findv6(ip string) (res result) {

	res = result{}
	res.IP = ip
	q.Offset = 0

	tp := big.NewInt(0)
	op := big.NewInt(0)
	tp.SetBytes(net.ParseIP(ip).To16())
	op.SetString("18446744073709551616", 10)
	op.Div(tp, op)
	tp.SetString("FFFFFFFFFFFFFFFF", 16)
	op.And(op, tp)

	v6ip = op.Uint64()
	offset = q.searchIndexV6(v6ip)
	q.Offset = offset

	country, area = q.getAddr()
	res.Country = string(country)
	res.Area = string(area)

	// Delete ZX (防止不相关的信息产生干扰）
	if res.Area == "ZX" || res.Area == "" {
		res.Area = ""
	} else {
		res.Area = " " + res.Area
	}

	return
}

func (q *pointer) getAddr() ([]byte, []byte) {
	mode := q.readData(1)[0]
	if mode == 0x01 {
		// [IP][0x01][国家和地区信息的绝对偏移地址]
		q.Offset = byteToUInt32(q.readData(3))
		return q.getAddr()
	}
	// [IP][0x02][信息的绝对偏移][...] or [IP][国家][...]
	_offset := q.Offset - 1
	c1 := q.readArea(_offset)
	if mode == 0x02 {
		q.Offset = 4 + _offset
	} else {
		q.Offset = _offset + uint32(1+len(c1))
	}
	c2 := q.readArea(q.Offset)
	return c1, c2
}

func (q *pointer) readArea(offset uint32) []byte {
	q.Offset = offset
	mode := q.readData(1)[0]
	if mode == 0x01 || mode == 0x02 {
		return q.readArea(byteToUInt32(q.readData(3)))
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
	header = q.Data.Data[0:8]
	start = binary.LittleEndian.Uint32(header[:4])
	end = binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, q.IndexLen)

	for {
		mid := start + q.IndexLen*(((end-start)/q.IndexLen)>>1)
		buf = q.Data.Data[mid : mid+q.IndexLen]
		_ip := binary.LittleEndian.Uint32(buf[:q.ItemLen])

		if end-start == q.IndexLen {
			if ip >= binary.LittleEndian.Uint32(q.Data.Data[end:end+q.ItemLen]) {
				buf = q.Data.Data[end : end+q.IndexLen]
			}
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

func (q *pointer) searchIndexV6(ip uint64) uint32 {

	q.ItemLen = 8
	q.IndexLen = 11

	header = q.Data.Data[8:24]
	start = binary.LittleEndian.Uint32(header[8:])
	counts := binary.LittleEndian.Uint32(header[:8])
	end = start + counts*q.IndexLen

	buf := make([]byte, q.IndexLen)

	for {
		mid := start + q.IndexLen*(((end-start)/q.IndexLen)>>1)
		buf = q.Data.Data[mid : mid+q.IndexLen]
		_ip := binary.LittleEndian.Uint64(buf[:q.ItemLen])

		if end-start == q.IndexLen {
			if ip >= binary.LittleEndian.Uint64(q.Data.Data[end:end+q.ItemLen]) {
				buf = q.Data.Data[end : end+q.IndexLen]
			}
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
func byteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
