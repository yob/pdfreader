// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// Some utilities.
package util

import (
  "fmt"
  "github.com/yob/pdfreader/xchar"
)

var wrongUniCode = xchar.Utf8(-1)

// util.Bytes() is a dup of string.Bytes()
func Bytes(a string) []byte {
  r := make([]byte, len(a))
  for k := 0; k < len(a); k++ {
    r[k] = byte(a[k])
  }
  return r
}

func JoinStrings(a []string, c byte) []byte {
  if a == nil || len(a) == 0 {
    return []byte{}
  }
  l := 0
  for k := range a {
    l += len(a[k]) + 1
  }
  r := make([]byte, l)
  q := 0
  for k := range a {
    for i := 0; i < len(a[k]); i++ {
      r[q] = a[k][i]
      q++
    }
    r[q] = c
    q++
  }
  return r[0 : l-1]
}

func StringArray(i [][]byte) []string {
  r := make([]string, len(i))
  for k := range i {
    r[k] = string(i[k])
  }
  return r
}

func set(o []byte, q string) int {
  for k := 0; k < len(q); k++ {
    o[k] = q[k]
  }
  return len(q)
}

func ToXML(s []byte) []byte {
  l := len(s)
  for k := range s {
    switch s[k] {
    case '<', '>':
      l += 3
    case '&':
      l += 4
    case 0, 1, 2, 3, 4, 5, 6, 7, 8,
      11, 12, 14, 15, 16, 17, 18, 19, 20,
      21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
      31:
      l += len(wrongUniCode) - 1
    }
  }
  r := make([]byte, l)
  p := 0
  for k := range s {
    switch s[k] {
    case '<':
      p += set(r[p:p+4], "&lt;")
    case '>':
      p += set(r[p:p+4], "&gt;")
    case '&':
      p += set(r[p:p+5], "&amp;")
    case 10, 9, 13:
      r[p] = s[k]
      p++
    default:
      if s[k] < 32 {
        p += copy(r[p:], wrongUniCode)
      } else {
        r[p] = s[k]
        p++
      }
    }
  }
  return r
}

type OutT struct {
  Content []byte
}

func (t *OutT) Out(f string, args ... string) {
  p := fmt.Sprintf(f, args)
  q := len(t.Content)
  if cap(t.Content)-q < len(p) {
    n := make([]byte, cap(t.Content)+(len(p)/512+2)*512)
    copy(n, t.Content)
    t.Content = n[0:q]
  }
  t.Content = t.Content[0 : q+len(p)]
  for k := 0; k < len(p); k++ {
    t.Content[q+k] = p[k]
  }
}
