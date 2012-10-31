// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// SVG driver for graf.go.
package svgdraw

import (
  "fmt"
  "github.com/yob/pdfreader/graf"
  "github.com/yob/pdfreader/stacks"
  "github.com/yob/pdfreader/strm"
  "github.com/yob/pdfreader/util"
)

type SvgT struct {
  Drw     *graf.PdfDrawerT
  drwpath stacks.StrStack
  p       int
  groups  int
}

func (s *SvgT) SvgPath() string {
  return fmt.Sprintf("path d=\"%s\"",
    util.JoinStrings(s.drwpath.Dump(), ' '))
}

func (s *SvgT) DropPath() { s.drwpath.Clear() }
func (s *SvgT) MoveTo(coord [][]byte) {
  s.drwpath.Push(fmt.Sprintf("M%s %s", coord[0], coord[1]))
}
func (s *SvgT) LineTo(coord [][]byte) {
  s.drwpath.Push(fmt.Sprintf("L%s %s", coord[0], coord[1]))
}

func (s *SvgT) CurveTo(coords [][]byte) {
  s.drwpath.Push(fmt.Sprintf("C%s %s %s %s %s %s",
    coords[0], coords[1],
    coords[2], coords[3],
    coords[4], coords[5]))
}

func (s *SvgT) Rectangle(coords [][]byte) {
  s.drwpath.Push(fmt.Sprintf("M%s %s V%s H%s V%s H%s Z",
    coords[0], coords[1],
    strm.Add(string(coords[1]), string(coords[3])),
    strm.Add(string(coords[0]), string(coords[2])),
    coords[1], coords[0]))
}

func (s *SvgT) ClosePath() { s.drwpath.Push("Z") }

func (s *SvgT) Stroke() {
  s.Drw.Write.Out("<%s fill=\"none\" stroke-width=\"%s\" stroke=\"%s\" />\n",
    s.SvgPath(), s.Drw.ConfigD.LineWidth, s.Drw.ConfigD.StrokeColor)
}

func (s *SvgT) Fill() {
  fill := s.Drw.ConfigD.FillColor
  if fill == "" {
    fill = "none"
  }
  s.Drw.Write.Out("<%s fill=\"%s\" stroke=\"none\" />\n",
    s.SvgPath(), fill)
}

func (s *SvgT) EOFill() { s.Fill() }

func (s *SvgT) FillAndStroke() {
  fill := s.Drw.ConfigD.FillColor
  if fill == "" {
    fill = "none"
  }
  s.Drw.Write.Out("<%s fill=\"%s\" stroke-width=\"%s\" stroke=\"%s\" />\n",
    s.SvgPath(), fill, s.Drw.ConfigD.LineWidth, s.Drw.ConfigD.StrokeColor)
}

func (s *SvgT) EOFillAndStroke() { s.FillAndStroke() }
func (s *SvgT) Clip()            {}
func (s *SvgT) EOClip()          {}

func (s *SvgT) Concat(m [][]byte) {
  s.Drw.Write.Out("<g transform=\"matrix(%s,%s,%s,%s,%s,%s)\">\n",
    string(m[0]), string(m[1]), string(m[2]), string(m[3]), string(m[4]), string(m[5]))
  s.groups++
}

func (s *SvgT) SetIdentity() {
  for s.groups > 0 {
    s.Drw.Write.Out("</g>\n")
    s.groups--
  }
}

func (s *SvgT) CloseDrawing() { s.SetIdentity() }

func (s *SvgT) Gray(a []byte) string {
  c := strm.Percent(a)
  return fmt.Sprintf("rgb(%s%%,%s%%,%s%%)", c, c, c)
}
func (s *SvgT) CMYK(cmyk [][]byte) string {
  return fmt.Sprintf("cmyk(%s%%,%s%%,%s%%,%s%%)",
    strm.Percent(cmyk[0]),
    strm.Percent(cmyk[1]),
    strm.Percent(cmyk[2]),
    strm.Percent(cmyk[3]))
}
func (s *SvgT) RGB(rgb [][]byte) string {
  return fmt.Sprintf("rgb(%s%%,%s%%,%s%%)",
    strm.Percent(rgb[0]),
    strm.Percent(rgb[1]),
    strm.Percent(rgb[2]))
}

func NewTestSvg() *graf.PdfDrawerT {
  t := new(SvgT)
  t.Drw = graf.NewPdfDrawer()
  t.Drw.ConfigD.SetColors(t)
  t.Drw.Draw = t
  t.drwpath = stacks.NewStrStack(-1)
  return t.Drw
}
