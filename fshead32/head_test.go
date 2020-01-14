/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/14
 */

package fshead32

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFsHead(t *testing.T) {
	h := &Head{
		Version:    1,
		ClientName: "demodemodemo",
		MagicNum:   1234,
		UserID:     6789,
		LogID:      34556,
		MetaLen:    3238,
		BodyLen:    198765,
	}

	got := h.Bytes()
	if gotLen := len(got); gotLen != 32 {
		t.Errorf("h.Bytes() want length=32, but gotLen=%d", gotLen)
	}

	want := []byte("\x01\x00demodemo\x85\x1a\x00\x00\xfc\x86\x00\x00\x00\x00\x00\x00\xa6\fm\b\x03\x00\xd2\x04\x00\x00")

	if !bytes.Equal(got, want) {
		t.Errorf("not eq, got=%q\n, want=%q", got, want)
	}

	h2, err := ParserBytes(got, 1234)
	if err != nil {
		t.Fatalf("parser ParserBytes failed, err=%s", err)
	}

	// clientName已被截取
	if h2.ClientName == "demodemo" {
		h.ClientName = "demodemo"
	}

	if !reflect.DeepEqual(h, h2) {
		t.Fatalf("ParserBytes not eq,got=%v,want=%v", h2, h)
	}

}

func TestFsHead2(t *testing.T) {
	h := &Head{
		Version:    1,
		ClientName: "demo",
		UserID:     6789,
		LogID:      34556,
		MetaLen:    3238,
		BodyLen:    198765,
	}

	got := h.Bytes()
	if gotLen := len(got); gotLen != 32 {
		t.Errorf("h.Bytes() want length=32, but gotLen=%d", gotLen)
	}

	want := []byte("\x01\x00demo\x00\x00\x00\x00\x85\x1a\x00\x00\xfc\x86\x00\x00\x00\x00\x00\x00\xa6\fm\b\x03\x00x\r\x1d\xce")

	if !bytes.Equal(got, want) {
		t.Errorf("not eq, got=%q\n, want=%q", got, want)
	}

	h2, err := ParserBytes(got, 0)
	if err != nil {
		t.Fatalf("parser ParserBytes failed, err=%s", err)
	}

	h.MagicNum = MagicNumDefault

	if !reflect.DeepEqual(h, h2) {
		t.Fatalf("ParserBytes not eq,got=%v,want=%v", h2, h)
	}

}

var buf []byte

func BenchmarkFsHead_Bytes(b *testing.B) {
	h := &Head{
		Version:    1,
		ClientName: "demo",
		UserID:     6789,
		LogID:      34556,
		MetaLen:    3238,
		BodyLen:    198765,
	}
	b.ResetTimer()
	b.ReportAllocs()

	var bf []byte
	for i := 0; i < b.N; i++ {
		bf = h.Bytes()
	}
	buf = bf
}
func BenchmarkFsHead_BytesManyZero(b *testing.B) {
	h := &Head{
		ClientName: "demo",
	}
	b.ResetTimer()
	b.ReportAllocs()

	var bf []byte
	for i := 0; i < b.N; i++ {
		bf = h.Bytes()
	}
	buf = bf
}
