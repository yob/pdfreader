// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// PS top-down parser.
package ps

import (
  "github.com/yob/pdfreader/fancy"
  "github.com/yob/pdfreader/hex"
)

func SkipLE(f fancy.Reader) {
  for {
    c, err := f.ReadByte()
    if err != nil {
      return
    }
    if c > 32 {
      f.UnreadByte()
      return
    }
    if c == 13 {
      c, err = f.ReadByte()
      if err == nil && c != 10 {
        f.UnreadByte()
      }
      return
    }
    if c == 10 {
      return
    }
  }
}

func skipSpaces(f fancy.Reader) byte {
  for {
    c, err := f.ReadByte()
    if err != nil {
      break
    }
    if c > 32 {
      return c
    }
  }
  return 0
}

func skipToDelim(f fancy.Reader) byte {
  for {
    c, err := f.ReadByte()
    if err != nil {
      break
    }
    if c < 33 {
      return c
    }
    switch c {
    case '<', '>', '(', ')', '[', ']', '/', '%', '{', '}':
      return c
    }
  }
  return 255
}

func skipString(f fancy.Reader) {
  for depth := 1; depth > 0; {
    c, err := f.ReadByte()
    if err != nil {
      break
    }
    switch c {
    case '(':
      depth++
    case ')':
      depth--
    case '\\':
      f.ReadByte()
    }
  }
}

func skipComment(f fancy.Reader) {
  for {
    c, err := f.ReadByte()
    if err != nil || c == 13 || c == 10 {
      break
    }
  }
}

func skipComposite(f fancy.Reader) {
  for depth := 1; depth > 0; {
    switch skipToDelim(f) {
    case '<', '[', '{':
      depth++
    case '>', ']', '}':
      depth--
    case '(':
      skipString(f)
    case '%':
      skipComment(f)
    }
  }
}

func fpos(f fancy.Reader) int64 {
  r, _ := f.Seek(0, 1)
  return r
}

func Token(f fancy.Reader) ([]byte, int64) {
again:
  c := skipSpaces(f)
  if c == 0 {
    return []byte{}, -1
  }
  p := fpos(f) - 1
  switch c {
  case '%':
    skipComment(f)
    goto again
  case '<', '[', '{':
    skipComposite(f)
  case '(':
    skipString(f)
  default:
    if skipToDelim(f) != 255 {
      f.UnreadByte()
    }
  }
  n := int(fpos(f) - p)
  f.Seek(p, 0)
  return f.Slice(n), p
}

func String(s []byte) []byte {
  if s[0] == '<' {
    r := hex.Decode(string(s[1 : len(s)-1]))
    return r
  }
  if s[0] != '(' {
    return s
  }
  r := make([]byte, len(s))
  q := 0
  for p := 1; p < len(s)-1; p++ {
    if s[p] == '\\' {
      p++
      switch s[p] {
      case 13:
        if s[p+1] == 10 {
          p++
        }
        q--
      case 10:
        q--
      case 'n':
        r[q] = 10
      case 'r':
        r[q] = 13
      case 't':
        r[q] = 9
      case 'b':
        r[q] = 8
      case 'f':
        r[q] = 12
      case '0', '1', '2', '3', '4', '5', '6', '7':
        a := s[p] - '0'
        if s[p+1] >= '0' && s[p+1] <= '7' {
          p++
          a = (a << 3) + (s[p] - '0')
          if s[p+1] >= '0' && s[p+1] <= '7' {
            p++
            a = (a << 3) + (s[p] - '0')
          }
        }
        r[q] = a
      default:
        r[q] = s[p]
      }
    } else {
      r[q] = s[p]
    }
    q++
  }
  return r[0:q]
}

func StrIntL(s []byte) (r, l int) {
  for k := range s {
    r <<= 8
    r += int(s[k])
  }
  return r, len(s)
}

func StrInt(s []byte) int {
  r, _ := StrIntL(s)
  return r
}
