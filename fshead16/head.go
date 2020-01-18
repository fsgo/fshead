/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/14
 */

package fshead16

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// MagicNumDefault 默认的magicNum，用于校验协议是否匹配
// uint32 = 3458010488
const MagicNumDefault uint32 = 0xce1d0d78

// Length 协议头长度，固定长16字节
const Length = 16

const _ClientNameLength = 6

// Head 协议头
type Head struct {
	// 魔法变量 用于校验协议是否匹配
	// 若为0，则使用默认值
	MagicNum uint32 // 0-4

	// 调用方名称,前6个字节
	ClientName string // 4-10

	// 后面的元数据长度
	// 消息完整格式为：{Head:固定长度}{Meta}{Body}
	MetaLen uint16 // 10-12

	// 消息体的长度
	BodyLen uint32 // 12-16

}

var magicNumDefaultBytes = make([]byte, 4)
var bytesPadding []byte

func init() {
	binary.LittleEndian.PutUint32(magicNumDefaultBytes, MagicNumDefault)

	bytesPadding = make([]byte, _ClientNameLength)
	for i, x := range bytesPadding {
		bytesPadding[i] = byte(x)
	}
}

// Bytes 将协议头对象装换为可以传输的bytes
func (h *Head) Bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, Length))

	if h.MagicNum == 0 {
		buf.Write(magicNumDefaultBytes) // 0-4
	} else {
		binary.Write(buf, binary.LittleEndian, h.MagicNum)
	}

	//  4-10
	if len(h.ClientName) >= _ClientNameLength {
		binary.Write(buf, binary.LittleEndian, []byte(h.ClientName[0:_ClientNameLength]))
	} else {
		binary.Write(buf, binary.LittleEndian, []byte(h.ClientName))
		buf.Write(bytesPadding[:_ClientNameLength-len(h.ClientName)])
	}

	writeBinaryUint16(buf, h.MetaLen) // 10-12
	writeBinaryUint32(buf, h.BodyLen) // 12-16
	return buf.Bytes()
}

var zeroUint16LEBytes = make([]byte, 2)
var zeroUint32LEBytes = make([]byte, 4)

func init() {
	binary.LittleEndian.PutUint16(zeroUint16LEBytes, 0)
	binary.LittleEndian.PutUint32(zeroUint32LEBytes, 0)
}

func writeBinaryUint16(w io.Writer, num uint16) {
	if num == 0 {
		w.Write(zeroUint16LEBytes)
		return
	}
	binary.Write(w, binary.LittleEndian, num)
}

func writeBinaryUint32(w io.Writer, num uint32) {
	if num == 0 {
		w.Write(zeroUint32LEBytes)
		return
	}
	binary.Write(w, binary.LittleEndian, num)
}

// ErrMagicNumNotMatch 协议不匹配
var ErrMagicNumNotMatch = errors.New("magicNum not match")

// ErrHeaderLengthWrong 读取到的header长度不对，期望是32字节
var ErrHeaderLengthWrong = errors.New("header bytes length wrong, expect is 32")

// ParserBytes 解析header头
// 并使用 指定magicNum 来校验协议是否匹配(若该值为0，则使用默认值校验)
func ParserBytes(buf []byte, magicNumWant uint32) (*Head, error) {
	if len(buf) != Length {
		return nil, ErrHeaderLengthWrong
	}

	magicNumGot := binary.LittleEndian.Uint32(buf[0:4])
	if !CheckMagicNum(magicNumGot, magicNumWant) {
		return nil, ErrMagicNumNotMatch
	}
	header := &Head{
		MagicNum:   magicNumGot,
		ClientName: string(bytes.TrimRight(buf[4:10], "\x00")),
		MetaLen:    binary.LittleEndian.Uint16(buf[10:12]),
		BodyLen:    binary.LittleEndian.Uint32(buf[12:16]),
	}
	return header, nil
}

// CheckMagicNum 检查magicNum 是否匹配
// 若 want为0，则使用默认值进行校验
func CheckMagicNum(got uint32, want uint32) bool {
	if got == 0 {
		return false
	}

	if got == want || (want == 0 && got == MagicNumDefault) {
		return true
	}
	return false
}
