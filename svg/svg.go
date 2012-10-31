// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// Library to convert PDF pages to SVG.
package svg

import (
  "fmt"
  "os"
  "github.com/yob/pdfreader/fancy"
  "github.com/yob/pdfreader/pdfread"
  "github.com/yob/pdfreader/strm"
  "github.com/yob/pdfreader/svgdraw"
  "github.com/yob/pdfreader/svgtext"
  "github.com/yob/pdfreader/util"
)

func complain(err string) {
  fmt.Printf("%s", err)
  os.Exit(1)
}

func Page(pd *pdfread.PdfReaderT, page int) []byte {
  pg := pd.Pages()
  if page >= len(pg) {
    complain("Page does not exist!\n")
  }
  mbox := util.StringArray(pd.Arr(pd.Att("/MediaBox", pg[page])))
  drw := svgdraw.NewTestSvg()
  svgtext.New(pd, drw).Page = page
  w := strm.Mul(strm.Sub(mbox[2], mbox[0]), "1.25")
  h := strm.Mul(strm.Sub(mbox[3], mbox[1]), "1.25")
  drw.Write.Out(
    "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"no\"?>\n"+
      "<svg\n"+
      "   xmlns:svg=\"http://www.w3.org/2000/svg\"\n"+
      "   xmlns=\"http://www.w3.org/2000/svg\"\n"+
      "   version=\"1.0\"\n"+
      "   width=\"%s\"\n"+
      "   height=\"%s\">\n"+
      "<g transform=\"matrix(1.25,0,0,-1.25,%s,%s)\">\n",
    w, h,
    strm.Mul(mbox[0], "-1.25"),
    strm.Mul(mbox[3], "1.25"))
  cont := pd.ForcedArray(pd.Dic(pg[page])["/Contents"])
  _, ps := pd.DecodedStream(cont[0])
  drw.Interpret(fancy.SliceReader(ps))
  drw.Draw.CloseDrawing()
  drw.Write.Out("</g>\n</svg>\n")
  return drw.Write.Content
}
