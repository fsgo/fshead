/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/14
 */

package fshead16

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/fsgo/fsprotocol"
)

// DefaultMagicNum 默认的magicNum，用于校验协议是否匹配
// uint32 = 3458010488
const DefaultMagicNum uint32 = 0xce1d0d78

// Length 协议头长度，固定长16字节
const Length = 16

// DiscernLength 识别该协议的最小长度
const DiscernLength = 4

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

// Is 判断是否当前协议
func (h *Head) Is(buf []byte) bool {
	if h == nil {
		return Is(buf, 0)
	}
	return Is(buf, h.MagicNum)
}

// Len 协议长度
func (h *Head) Len() int {
	return Length
}

// DiscernLen 确认协议的长度
func (h *Head) DiscernLen() int {
	return DiscernLength
}

// Load 加载解析协议
func (h *Head) Load(buf []byte) error {
	magicNum := DefaultMagicNum
	if h != nil && h.MagicNum != 0 {
		magicNum = h.MagicNum
	}
	head, err := Load(buf, magicNum)
	if err != nil {
		return err
	}
	// todo
	rv := reflect.ValueOf(h).Elem()
	rv.Set(reflect.ValueOf(head).Elem())
	return nil
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

var _ fsprotocol.Protocol = &Head{}
