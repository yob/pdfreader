// Stacks of different types.
package stacks

type Stack interface {
  Push([]byte)
  Pop() []byte
  Drop(int) (st [][]byte)
  Dump() [][]byte
  Depth() int
  Index(p int) []byte
  Clear()
}

type StackT struct {
  st    [][]byte
  sp    int
  limit bool
}

func (st *StackT) Push(s []byte) {
  if st.sp >= len(st.st) && !st.limit {
    t := make([][]byte, len(st.st)+1024)
    st.sp = 0
    for k := range st.st {
      t[k] = st.st[k]
      st.sp++
    }
    st.st = t
  }
  st.st[st.sp] = s
  st.sp++
}
func (st *StackT) Drop(n int) [][]byte {
  st.sp -= n
  return st.st[st.sp : st.sp+n]
}
func (st *StackT) Pop() []byte {
  st.sp--
  return st.st[st.sp]
}
func (st *StackT) Dump() [][]byte     { return st.st[0:st.sp] }
func (st *StackT) Depth() int         { return st.sp }
func (st *StackT) Index(p int) []byte { return st.st[st.sp-p] }
func (st *StackT) Clear()             { st.sp = 0 }
func NewStack(n int) *StackT {
  r := new(StackT)
  if r.limit = n >= 0; !r.limit {
    n = 1024
  }
  r.st = make([][]byte, n)
  return r
}

type StrStack interface {
  Push(string)
  Pop() string
  Drop(int) (st []string)
  Dump() []string
  Depth() int
  Index(p int) string
  Clear()
}

type StrStackT struct {
  st    []string
  sp    int
  limit bool
}

func (st *StrStackT) Push(s string) {
  if st.sp >= len(st.st) && !st.limit {
    t := make([]string, len(st.st)+1024)
    st.sp = 0
    for k := range st.st {
      t[k] = st.st[k]
      st.sp++
    }
    st.st = t
  }
  st.st[st.sp] = s
  st.sp++
}
func (st *StrStackT) Drop(n int) []string {
  st.sp -= n
  return st.st[st.sp : st.sp+n]
}
func (st *StrStackT) Pop() string {
  st.sp--
  return st.st[st.sp]
}
func (st *StrStackT) Dump() []string     { return st.st[0:st.sp] }
func (st *StrStackT) Depth() int         { return st.sp }
func (st *StrStackT) Index(p int) string { return st.st[st.sp-p] }
func (st *StrStackT) Clear()             { st.sp = 0 }
func NewStrStack(n int) *StrStackT {
  r := new(StrStackT)
  if r.limit = n >= 0; !r.limit {
    n = 1024
  }
  r.st = make([]string, n)
  return r
}
