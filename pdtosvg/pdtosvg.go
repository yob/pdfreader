// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// Convert PDF-pages to SVG.
package main

import (
  "fmt"
  "os"
  "github.com/yob/pdfreader/pdfread"
  "github.com/yob/pdfreader/strm"
  "github.com/yob/pdfreader/svg"
)

// The program takes a PDF file and converts a page to SVG.

func complain(err string) {
  fmt.Printf("%susage: pdtosvg foo.pdf [page] >foo.svg\n", err)
  os.Exit(1)
}

func main() {
  if len(os.Args) == 1 || len(os.Args) > 3 {
    complain("")
  }
  page := 0
  if len(os.Args) > 2 {
    page = strm.Int(os.Args[2], 1) - 1
    if page < 0 {
      complain("Bad page!\n\n")
    }
  }
  pd := pdfread.Load(os.Args[1])
  if pd == nil {
    complain("Could not load pdf file!\n\n")
  }
  fmt.Printf("%s", svg.Page(pd, page))
}
