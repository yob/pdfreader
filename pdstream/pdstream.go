// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// Example program for pdfread.go
package main

import (
  "os"
  "fmt"
  "github.com/yob/pdfreader/pdfread"
  "github.com/yob/pdfreader/util"
)

// The program takes a PDF file and an object reference of a stream.
// The output are the decoded stream contents.
//
// Example:
//  ./pdstream.go foo.pdf "9 0 R"

func main() {
  pd := pdfread.Load(os.Args[1])
  _, d := pd.DecodedStream(util.Bytes(os.Args[2]))
  fmt.Printf("%s", d)

  /*
     a := cmapi.Read(fancy.SliceReader(d));
     fmt.Printf("\n%v\n%v\n%v\n", a, a.Ranges, a.Uni);
  */
}
