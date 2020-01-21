/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/19
 */

package fshead16

import (
	"encoding/binary"
	"io"
)

var magicNumDefaultBytes = make([]byte, 4)
var bytesPadding []byte
var zeroUint16LEBytes = make([]byte, 2)
var zeroUint32LEBytes = make([]byte, 4)

func init() {
	binary.LittleEndian.PutUint32(magicNumDefaultBytes, DefaultMagicNum)
	binary.LittleEndian.PutUint16(zeroUint16LEBytes, 0)
	binary.LittleEndian.PutUint32(zeroUint32LEBytes, 0)
	bytesPadding = make([]byte, _ClientNameLength)
	for i, x := range bytesPadding {
		bytesPadding[i] = byte(x)
	}
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
