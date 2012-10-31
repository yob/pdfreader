// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// Decoder for pfb fonts.
package pfb

import (
  "github.com/yob/pdfreader/hex"
)

func Decode(b []byte) []byte {
  r := make([]byte, len(b)*2)[0:0]
  for {
    if b[0] != 128 {
      break
    }
    if b[1] == 3 {
      break
    }
    l := int(b[2]) + (int(b[3]) << 8) + (int(b[4]) << 16) + 6
    if b[1] == 1 {
      r = append(r, b[6:l]...)
    } else {
      r = append(r, hex.Encode(b[6:l])...)
    }
    b = b[l:]
  }
  return r
}
