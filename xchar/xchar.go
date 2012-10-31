// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// Encode UTF-8.
package xchar

var utconv = []int{
  0x7F, 0x00,
  0x7FF, 0xC0,
  0xFFFF, 0xE0,
  0x1FFFFF, 0xF0,
}

func Utf8(rune int) []byte {
  up := 3
  out := make([]byte, up+1)
  if rune < 0 || rune > 0x10FFFF {
    rune = 0xFFFD
  }
  uc := 0
  r := rune
  for ; utconv[uc] < rune; uc += 2 {
    out[up] = byte((r & 0x3F) | 0x80)
    up--
    r >>= 6
  }
  out[up] = byte(r | utconv[uc+1])
  return out[up:]
}

func EncodeRune(rune int, out []byte) int {
  r := Utf8(rune)
  copy(out, r)
  return len(r)
}
