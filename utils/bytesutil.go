// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

// 处理byte 数组或者slice的组件

package utils

import "strings"

var hexDict = []byte("0123456789ABCDEF")



// ToHexString byte数组转换成 16进制表示的string，upper 表示生成的十六进制是否是大写， true表示大写， false表示小写
// 与 hex.EncodeToString 相比它支持大小写输出
func ToHexString(data []byte, upper bool) string  {
	var builder strings.Builder
	builder.Grow(len(data)*2)
	for _, b := range  data {
		hc := hexDict[b>>4]
		lc := hexDict[b & 0x0F]

		if !upper {
			hc = toLow(hc)
			lc = toLow(lc)
		}
		builder.WriteByte(hc)
		builder.WriteByte(lc)
	}
	return builder.String()
}

// ToUpperHexString byte数组转换成大写的16进制表示的string
func ToUpperHexString(data []byte) string  {
	return ToHexString(data,true)
}

func toLow(c byte)  byte {
	if c >= 'A' && c <= 'F' {
		return c + 32
	}
	return c
}
