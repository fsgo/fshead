/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/10
 */

package fshead

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// MagicNumDefault 默认的magicNum，用于校验协议是否匹配
// uint32 = 3458010488
const MagicNumDefault uint32 = 0xce1d0d78

// Length 协议头长度，固定长32字节
const Length = 32

// FsHead 协议头
type FsHead struct {

	// 协议版本
	Version uint16

	// 调用方名称,前8个字节
	ClientName string

	// 调用方ID，若不需要，可以传0
	// server端也可以依次做身份校验
	UserID uint32

	// 日志ID
	LogID uint32

	// 预留字段，业务可以扩展使用
	Reserve uint32

	// 后面的元数据长度
	// 消息完整格式为：{FsHead:固定长度}{Meta}{Body}
	MetaLen uint16

	// 消息体的长度
	BodyLen uint32

	// 魔法变量 用于校验协议是否匹配
	// 若为0，则使用默认值
	MagicNum uint32
}

// Bytes 将协议头对象装换为可以传输的bytes
func (h *FsHead) Bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, Length))
	binary.Write(buf, binary.LittleEndian, h.Version) // 0-2

	if len(h.ClientName) >= 8 {
		binary.Write(buf, binary.LittleEndian, []byte(h.ClientName[0:8])) // 2-10
	} else {
		binary.Write(buf, binary.LittleEndian, []byte(h.ClientName))
		cnamePad := make([]byte, 8-len(h.ClientName))
		binary.Write(buf, binary.LittleEndian, cnamePad)
	}
	binary.Write(buf, binary.LittleEndian, h.UserID)  // 10-14
	binary.Write(buf, binary.LittleEndian, h.LogID)   // 14-18
	binary.Write(buf, binary.LittleEndian, h.Reserve) // 18-22
	binary.Write(buf, binary.LittleEndian, h.MetaLen) // 22-24
	binary.Write(buf, binary.LittleEndian, h.BodyLen) // 24-28
	if h.MagicNum == 0 {
		binary.Write(buf, binary.LittleEndian, MagicNumDefault) // 28-32
	} else {
		binary.Write(buf, binary.LittleEndian, h.MagicNum)
	}
	return buf.Bytes()
}

// ErrMagicNumNotMatch 协议不匹配
var ErrMagicNumNotMatch = errors.New("magicNum not match")

// ErrHeaderLengthWrong 读取到的header长度不对，期望是32字节
var ErrHeaderLengthWrong = errors.New("header bytes length wrong, expect is 32")

// ParserBytes 解析header头
// 并使用 指定magicNum 来校验协议是否匹配(若该值为0，则使用默认值校验)
func ParserBytes(buf []byte, magicNumWant uint32) (*FsHead, error) {
	if len(buf) != Length {
		return nil, ErrHeaderLengthWrong
	}

	magicNumGot := binary.LittleEndian.Uint32(buf[28:32])
	if !CheckMagicNum(magicNumGot, magicNumWant) {
		return nil, ErrMagicNumNotMatch
	}
	header := &FsHead{
		Version:    binary.LittleEndian.Uint16(buf[0:2]),
		ClientName: string(bytes.TrimRight(buf[2:10], "\x00")),
		UserID:     binary.LittleEndian.Uint32(buf[10:14]),
		LogID:      binary.LittleEndian.Uint32(buf[14:18]),
		Reserve:    binary.LittleEndian.Uint32(buf[18:22]),
		MetaLen:    binary.LittleEndian.Uint16(buf[22:24]),
		BodyLen:    binary.LittleEndian.Uint32(buf[24:28]),
		MagicNum:   magicNumGot,
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
