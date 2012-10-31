// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// Character Mappings (cmap).
package cmapt

// Fundamental type for character mappings.  This type is fairly central for
// mapping of characters to unicode and to character widths, so I assume an
// own package is justified.

type CMapRangeT struct {
  Prev           *CMapRangeT
  From, To, Dest int
}

type CMapT struct {
  Basic    [256]int    // basic characters
  Extended map[int]int // characters > 255
  Ranges   *CMapRangeT // ranges with growing content
  DRanges  *CMapRangeT // fixed ranges
}

//
// HINT: Better _NOT_ mix Ranges and DRanges above char(255)!
//

func New() *CMapT {
  r := new(CMapT)
  for k := range r.Basic {
    r.Basic[k] = -1
  }
  r.Extended = make(map[int]int)
  return r
}

func (m *CMapT) Code(s int) int {
  if s < 256 {
    return m.Basic[s]
  }
  if r, ok := m.Extended[s]; ok {
    return r
  }
  for t := m.Ranges; t != nil; t = t.Prev {
    if s >= t.From && s < t.To {
      return t.Dest + (s - t.From)
    }
  }
  for t := m.DRanges; t != nil; t = t.Prev {
    if s >= t.From && s < t.To {
      return t.Dest
    }
  }
  return -1
}

func (m *CMapT) AddRange(from, to, dest int) {
  for k := from; k < 256 && k < to; k++ {
    m.Basic[k] = dest + (k - from)
  }
  if from < 256 {
    if from = 256; from >= to {
      return
    }
  }
  if to-from < 32 {
    for k := from; k < to; k++ {
      m.Extended[k] = dest + (k - from)
    }
    return
  }
  for k := range m.Extended {
    if k >= from && k < to {
      delete(m.Extended, k)
    }
  }
  r := new(CMapRangeT)
  r.From, r.To, r.Dest, r.Prev = from, to, dest, m.Ranges
  m.Ranges = r
}

func (m *CMapT) AddDef(from, to, dest int) {
  for k := from; k < 256 && k < to; k++ {
    m.Basic[k] = dest
  }
  if from < 256 {
    if from = 256; from >= to {
      return
    }
  }
  if to-from < 32 {
    for k := from; k < to; k++ {
      m.Extended[k] = dest
    }
    return
  }
  for k := range m.Extended {
    if k >= from && k < to {
      delete(m.Extended, k)
    }
  }
  r := new(CMapRangeT)
  r.From, r.To, r.Dest = from, to, dest
  r.Prev, m.DRanges = m.DRanges, r
}

func (m *CMapT) Add(k, dest int) {
  if k < 256 {
    m.Basic[k] = dest
  } else {
    m.Extended[k] = dest
  }
}
