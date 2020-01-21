/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/14
 */

package fshead16

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFsHead(t *testing.T) {
	h := &Head{
		ClientName: "demodemodemo",
		MagicNum:   1234,
		MetaLen:    3238,
		BodyLen:    198765,
	}

	got := h.Bytes()
	if gotLen := len(got); gotLen != 16 {
		t.Errorf("h.Bytes() want length=16, but gotLen=%d", gotLen)
	}

	if !h.Is(got) {
		t.Fatalf("h.Is(got) failed")
	}

	if h.Is(got[:3]) {
		t.Fatalf("h.Is(got) want failed")
	}

	want := []byte("\xd2\x04\x00\x00demode\xa6\fm\b\x03\x00")

	if !bytes.Equal(got, want) {
		t.Errorf("not eq, got=%q\n, want=%q", got, want)
	}

	h2, err := Load(got, 1234)
	if err != nil {
		t.Fatalf("parser Load failed, err=%s", err)
	}

	// clientName已被截取
	if h2.ClientName == "demode" {
		h.ClientName = "demode"
	}

	if !reflect.DeepEqual(h, h2) {
		t.Fatalf("Load not eq,got=%v,want=%v", h2, h)
	}

}

func TestFsHead2(t *testing.T) {
	h := &Head{
		ClientName: "demo",
		MetaLen:    3238,
		BodyLen:    198765,
		MagicNum:   DefaultMagicNum,
	}

	got := h.Bytes()
	if gotLen := len(got); gotLen != 16 {
		t.Errorf("h.Bytes() want length=16, but gotLen=%d", gotLen)
	}

	if !h.Is(got) {
		t.Fatalf("h.Is(got) failed")
	}
	{
		h3 := &Head{}
		if !h3.Is(got) {
			t.Fatalf("h2.Is(got) failed")
		}

		if err := h3.Load(got); err != nil {
			t.Fatalf("h3.Load(got) with error: %s", err)
		} else {
			if !reflect.DeepEqual(h, h3) {
				t.Fatalf("h!=h3, h=%v, h3=%v", h, h3)
			}
		}
	}

	want := []byte("x\r\x1d\xcedemo\x00\x00\xa6\fm\b\x03\x00")

	if !bytes.Equal(got, want) {
		t.Errorf("not eq, got=%q\n, want=%q", got, want)
	}

	h2, err := Load(got, 0)
	if err != nil {
		t.Fatalf("parser Load failed, err=%s", err)
	}

	h.MagicNum = DefaultMagicNum

	if !reflect.DeepEqual(h, h2) {
		t.Fatalf("Load not eq,got=%v,want=%v", h2, h)
	}

}

var buf []byte

func BenchmarkFsHead_Bytes(b *testing.B) {
	h := &Head{
		ClientName: "demo",
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
