// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// "crush" bytes into bits - variable length.
package crush

type BitT struct {
  s    []byte
  p, b int
}

var mask = [9]int{0, 1, 3, 7, 15, 31, 63, 127, 255}

func (x *BitT) Get(n int) (r int) {
  if x.b == 0 {
    x.b = 8
    x.p++
  }
  if x.b >= n {
    x.b -= n
    r = int(x.s[x.p]>>uint8(x.b)) & mask[n]
  } else {
    n -= x.b
    r = x.Get(x.b) << uint8(n)
    r += x.Get(n)
  }
  return
}

func NewBits(s []byte) *BitT {
  r := new(BitT)
  r.s = s
  r.b = 8
  return r
}
