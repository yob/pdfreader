// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// hex encoder/decoder for PDF.
package hex

var deco [256]byte

func init() {
  for i := 0; i <= 255; i++ {
    deco[i] = 255
  }
  for i := '0'; i <= '9'; i++ {
    deco[i] = byte(i) - '0'
  }
  for i := 'A'; i <= 'F'; i++ {
    deco[i] = byte(i) - 'A' + 10
  }
  for i := 'a'; i <= 'f'; i++ {
    deco[i] = byte(i) - 'a' + 10
  }
}

func Decode(s string) []byte {
  r := make([]byte, (len(s)+1)/2)
  q := 0
  for p := 0; p < len(s); p++ {
    if c := deco[s[p]]; c != 255 {
      if q%2 == 0 {
        c <<= 4
      }
      r[q/2] += c
      q++
    } else if s[p] > 32 {
      if s[p] == '>' {
        break
      }
      return []byte{}
    }
  }
  return r[0 : (q+1)/2]
}

func IsHex(c byte) bool { return deco[c] != 255 }

const _hex = "0123456789ABCDEF"

func Encode(i []byte) []byte {
  r := make([]byte, len(i)*2)
  for k := range i {
    r[k*2] = _hex[i[k]>>4]
    r[k*2+1] = _hex[i[k]&15]
  }
  return r
}

func EncodeLen(i []byte) int { return len(i) * 2 }
