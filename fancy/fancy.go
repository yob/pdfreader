// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// Enhanced input.
package fancy

import "os"
import "io"
import "io/ioutil"
import "bufio"

type Reader interface {
  ReadAt(buf []byte, pos int64) (n int, err error)
  Read(b []byte) (n int, err error)
  Slice(n int) []byte
  Seek(off int64, whence int) (ret int64, err error)
  ReadByte() (c byte, err error)
  UnreadByte() error
  Size() int64
}

// ------------------------------------------------------------------

const (
  _SECTOR_SIZE  = 512
  _SECTOR_COUNT = 32
)

type SecReaderT struct {
  cache     map[int64][]byte
  age       map[int64]int
  ticker    int
  pos, size int64
  f         io.ReaderAt
}

func min(a, b int64) int64 {
  if a < b {
    return a
  }
  return b
}

func (sr *SecReaderT) access(pos int64) (sl []byte, p int) {
  p = int(pos % _SECTOR_SIZE)
  pos /= _SECTOR_SIZE
  if s, ok := sr.cache[pos]; ok {
    if sr.age[pos] != sr.ticker {
      sr.ticker++
      sr.age[pos] = sr.ticker
    }
    return s, p
  }
  if len(sr.cache) >= _SECTOR_COUNT {
    a := sr.ticker
    old := int64(0)
    for k := range sr.age {
      if sr.age[k] < a {
        old = k
        a = sr.age[k]
      }
    }
    delete(sr.cache, old)
    delete(sr.age, 0)
  }
  sr.ticker++
  sl = make([]byte, min(sr.size-pos*_SECTOR_SIZE, _SECTOR_SIZE))
  sr.f.ReadAt(sl, pos*_SECTOR_SIZE)
  sr.cache[pos] = sl
  sr.age[pos] = sr.ticker
  return
}

func (sr *SecReaderT) ReadAt(buf []byte, pos int64) (n int, err error) {
  if pos >= sr.size {
    return 0, io.EOF
  }
  b, p := sr.access(pos)
  for ; p < _SECTOR_SIZE && n < len(buf); p++ {
    buf[n] = b[p]
    n++
  }
  if secs := (len(buf) - n) / _SECTOR_SIZE; secs > 0 {
    sr.f.ReadAt(buf[n:n+secs*_SECTOR_SIZE], pos+int64(n))
    n += secs * _SECTOR_SIZE
  }
  if len(buf)-n > 0 {
    b, p = sr.access(pos + int64(n))
    for ; n < len(buf); p++ {
      buf[n] = b[p]
      n++
    }
  }
  if pos+int64(n) >= sr.size {
    n -= int(pos + int64(n) - sr.size)
  }
  return
}

func (sr *SecReaderT) Read(b []byte) (n int, err error) {
  n, err = sr.ReadAt(b, sr.pos)
  sr.pos += int64(n)
  return
}

func (sr *SecReaderT) Seek(off int64, whence int) (ret int64, err error) {
  ret = sr.pos
  switch whence {
  case 0:
    sr.pos = 0
  case 2:
    sr.pos = sr.size
  }
  sr.pos += off
  return
}

func (sr *SecReaderT) ReadByte() (c byte, err error) {
  if sr.pos < sr.size {
    b, p := sr.access(sr.pos)
    c = b[p]
    sr.pos++
  } else {
    err = io.EOF
  }
  return
}

func (sr *SecReaderT) UnreadByte() error {
  sr.pos--
  return nil
}

func (sr *SecReaderT) Size() int64 { return sr.size }

func (sr *SecReaderT) Slice(n int) []byte {
  r := make([]byte, n)
  sr.Read(r)
  return r
}

// grmpf: Next is for AUTOGENERATE!
// The thing here is only (!) for convenience.
func (sr *SecReaderT) ReadBytes(delim byte) ([]byte, error) {
  return bufio.NewReader(sr).ReadBytes(delim)
}

func SecReader(f io.ReaderAt, size int64) Reader {
  sr := new(SecReaderT)
  sr.f = f
  sr.size = size
  sr.cache = make(map[int64][]byte)
  sr.age = make(map[int64]int)
  return sr
}

// ------------------------------------------------------------------

type SliceReaderT struct {
  bin []byte
  pos int64
}

func (sl *SliceReaderT) ReadAt(b []byte, off int64) (n int, err error) {
  for n := 0; n < len(b); n++ {
    if off >= int64(len(sl.bin)) {
      if n > 0 {
        break
      }
      return n, io.EOF
    }
    b[n] = sl.bin[off]
    off++
  }
  return len(b), nil
}

func (sl *SliceReaderT) Read(b []byte) (n int, err error) {
  n, err = sl.ReadAt(b, sl.pos)
  sl.pos += int64(n)
  return
}

func (sl *SliceReaderT) Seek(off int64, whence int) (ret int64, err error) {
  ret = sl.pos
  switch whence {
  case 0:
    sl.pos = 0
  case 2:
    sl.pos = int64(len(sl.bin))
  }
  sl.pos += off
  return
}

func (sl *SliceReaderT) Size() int64 { return int64(len(sl.bin)) }

func (sl *SliceReaderT) ReadByte() (c byte, err error) {
  if sl.pos < int64(len(sl.bin)) {
    c = sl.bin[sl.pos]
    sl.pos++
  } else {
    err = io.EOF
  }
  return
}

func (sl *SliceReaderT) UnreadByte() error {
  sl.pos--
  return nil
}

func (sl *SliceReaderT) Slice(n int) []byte {
  sl.pos += int64(n)
  return sl.bin[sl.pos-int64(n) : sl.pos]
}

// grmpf: Next is for AUTOGENERATE!
// The thing here is only (!) for convenience.
func (sl *SliceReaderT) ReadBytes(delim byte) ([]byte, error) {
  return bufio.NewReader(sl).ReadBytes(delim)
}

func SliceReader(bin []byte) Reader {
  r := new(SliceReaderT)
  r.bin = bin
  return r
}

// ------------------------------------------------------------------

func ReadAndClose(f io.ReadCloser, err error) []byte {
  if err != nil {
    return []byte{}
  }
  r, _ := ioutil.ReadAll(f)
  f.Close()
  return r
}

func FileReader(fn string) Reader {
  dir, err := os.Stat(fn)
  if err != nil {
    return nil
  }
  fil, err := os.Open(fn)
  if err != nil {
    return nil
  }
  return SecReader(fil, int64(dir.Size()))
}

func ReadAll(f io.Reader) []byte {
  r, _ := ioutil.ReadAll(f)
  return r
}

// ------------------------------------------------------------------
