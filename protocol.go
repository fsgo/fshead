/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/19
 */

package fsprotocol

type Protocol interface {
	Bytes() []byte

	// 长度，定长,若不是定长则为-1
	Len() int

	// 识别该协议的最小长度 <= Len()
	DiscernLen() int

	// 解析
	Load(buf []byte) error

	// 使用DiscernLen()长度的[]byte 来判断是否是该协议
	Is(buf []byte) bool
}
