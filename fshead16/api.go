/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/21
 */

package fshead16

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// ErrMagicNumNotMatch 协议不匹配
var ErrMagicNumNotMatch = errors.New("magicNum not match")

// ErrHeaderLengthWrong 读取到的header长度不对，期望是32字节
var ErrHeaderLengthWrong = errors.New("header bytes length wrong, expect is 32")

// Load 解析header头
// 并使用 指定magicNum 来校验协议是否匹配(若该值为0，则使用默认值校验)
func Load(buf []byte, magicNumWant uint32) (*Head, error) {
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

	if got == want || (want == 0 && got == DefaultMagicNum) {
		return true
	}
	return false
}

// Is 是否该协议
func Is(buf []byte, magicNum uint32) bool {
	if len(buf) < DiscernLength {
		return false
	}
	magicNumGot := binary.LittleEndian.Uint32(buf[0:DiscernLength])
	return CheckMagicNum(magicNumGot, magicNum)
}
