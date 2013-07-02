package main

import "os"
import ("bufio";"io";"strings")
type dfa struct {
  acc []bool
  f []func(rune) int
  id int
}
type family struct {
  a []dfa
  endcase int
}
var a0 [13]dfa
var a []family
func init() {
a = make([]family, 1)
{
var acc [4]bool
var fun [4]func(rune) int
fun[0] = func(r rune) int {
  switch(r) {
  case 47: return 1
  case 10: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[1] = func(r rune) int {
  switch(r) {
  case 47: return 2
  case 10: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[2] = func(r rune) int {
  switch(r) {
  case 47: return 3
  case 10: return -1
  default:
    switch {
    default: return 3
    }
  }
  panic("unreachable")
}
acc[3] = true
fun[3] = func(r rune) int {
  switch(r) {
  case 47: return 3
  case 10: return -1
  default:
    switch {
    default: return 3
    }
  }
  panic("unreachable")
}
a0[0].acc = acc[:]
a0[0].f = fun[:]
a0[0].id = 0
}
{
var acc [2]bool
var fun [2]func(rune) int
fun[0] = func(r rune) int {
  switch(r) {
  default:
    switch {
    case 48 <= r && r <= 57: return 1
    default: return -1
    }
  }
  panic("unreachable")
}
acc[1] = true
fun[1] = func(r rune) int {
  switch(r) {
  default:
    switch {
    case 48 <= r && r <= 57: return 1
    default: return -1
    }
  }
  panic("unreachable")
}
a0[1].acc = acc[:]
a0[1].f = fun[:]
a0[1].id = 1
}
{
var acc [28]bool
var fun [28]func(rune) int
fun[2] = func(r rune) int {
  switch(r) {
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return 16
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[6] = true
fun[6] = func(r rune) int {
  switch(r) {
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[23] = true
fun[23] = func(r rune) int {
  switch(r) {
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[13] = true
fun[13] = func(r rune) int {
  switch(r) {
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[17] = func(r rune) int {
  switch(r) {
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return 18
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[21] = func(r rune) int {
  switch(r) {
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return 22
  case 115: return -1
  case 108: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[22] = func(r rune) int {
  switch(r) {
  case 101: return 23
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[0] = func(r rune) int {
  switch(r) {
  case 114: return 1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return 2
  case 116: return -1
  case 110: return -1
  case 102: return 3
  case 105: return 4
  case 112: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[5] = func(r rune) int {
  switch(r) {
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return 7
  case 114: return -1
  case 117: return -1
  case 99: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[7] = func(r rune) int {
  switch(r) {
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return 8
  case 97: return -1
  case 103: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[19] = func(r rune) int {
  switch(r) {
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return -1
  case 116: return 24
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[26] = func(r rune) int {
  switch(r) {
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return 27
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[3] = func(r rune) int {
  switch(r) {
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return 11
  case 99: return -1
  case 109: return -1
  case 111: return 12
  case 97: return -1
  case 103: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[9] = func(r rune) int {
  switch(r) {
  case 101: return -1
  case 116: return 10
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[8] = func(r rune) int {
  switch(r) {
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return 9
  case 117: return -1
  case 99: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[12] = func(r rune) int {
  switch(r) {
  case 114: return 13
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[11] = func(r rune) int {
  switch(r) {
  case 101: return -1
  case 116: return -1
  case 110: return 14
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[16] = func(r rune) int {
  switch(r) {
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return 17
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[24] = func(r rune) int {
  switch(r) {
  case 114: return -1
  case 117: return 25
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[25] = func(r rune) int {
  switch(r) {
  case 114: return 26
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[27] = true
fun[27] = func(r rune) int {
  switch(r) {
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[4] = func(r rune) int {
  switch(r) {
  case 109: return 5
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return 6
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[10] = true
fun[10] = func(r rune) int {
  switch(r) {
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[15] = true
fun[15] = func(r rune) int {
  switch(r) {
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[18] = true
fun[18] = func(r rune) int {
  switch(r) {
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[20] = func(r rune) int {
  switch(r) {
  case 101: return -1
  case 116: return -1
  case 110: return 21
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[1] = func(r rune) int {
  switch(r) {
  case 101: return 19
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return -1
  case 109: return -1
  case 111: return -1
  case 97: return 20
  case 103: return -1
  case 115: return -1
  case 108: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[14] = func(r rune) int {
  switch(r) {
  case 101: return -1
  case 116: return -1
  case 110: return -1
  case 102: return -1
  case 105: return -1
  case 112: return -1
  case 114: return -1
  case 117: return -1
  case 99: return 15
  case 109: return -1
  case 111: return -1
  case 97: return -1
  case 103: return -1
  case 115: return -1
  case 108: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
a0[2].acc = acc[:]
a0[2].f = fun[:]
a0[2].id = 2
}
{
var acc [10]bool
var fun [10]func(rune) int
fun[0] = func(r rune) int {
  switch(r) {
  case 116: return 1
  case 114: return -1
  case 117: return -1
  case 101: return -1
  case 102: return 2
  case 97: return -1
  case 108: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[1] = func(r rune) int {
  switch(r) {
  case 116: return -1
  case 114: return 7
  case 117: return -1
  case 101: return -1
  case 102: return -1
  case 97: return -1
  case 108: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[2] = func(r rune) int {
  switch(r) {
  case 116: return -1
  case 114: return -1
  case 117: return -1
  case 101: return -1
  case 102: return -1
  case 97: return 3
  case 108: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[3] = func(r rune) int {
  switch(r) {
  case 116: return -1
  case 114: return -1
  case 117: return -1
  case 101: return -1
  case 102: return -1
  case 97: return -1
  case 108: return 4
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[4] = func(r rune) int {
  switch(r) {
  case 116: return -1
  case 114: return -1
  case 117: return -1
  case 101: return -1
  case 102: return -1
  case 97: return -1
  case 108: return -1
  case 115: return 5
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[5] = func(r rune) int {
  switch(r) {
  case 116: return -1
  case 114: return -1
  case 117: return -1
  case 101: return 6
  case 102: return -1
  case 97: return -1
  case 108: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[6] = true
fun[6] = func(r rune) int {
  switch(r) {
  case 116: return -1
  case 114: return -1
  case 117: return -1
  case 101: return -1
  case 102: return -1
  case 97: return -1
  case 108: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[7] = func(r rune) int {
  switch(r) {
  case 116: return -1
  case 114: return -1
  case 117: return 8
  case 101: return -1
  case 102: return -1
  case 97: return -1
  case 108: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[8] = func(r rune) int {
  switch(r) {
  case 116: return -1
  case 114: return -1
  case 117: return -1
  case 101: return 9
  case 102: return -1
  case 97: return -1
  case 108: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[9] = true
fun[9] = func(r rune) int {
  switch(r) {
  case 116: return -1
  case 114: return -1
  case 117: return -1
  case 101: return -1
  case 102: return -1
  case 97: return -1
  case 108: return -1
  case 115: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
a0[3].acc = acc[:]
a0[3].f = fun[:]
a0[3].id = 3
}
{
var acc [3]bool
var fun [3]func(rune) int
fun[0] = func(r rune) int {
  switch(r) {
  case 95: return 1
  default:
    switch {
    case 48 <= r && r <= 57: return -1
    case 65 <= r && r <= 90: return 1
    case 97 <= r && r <= 122: return 1
    default: return -1
    }
  }
  panic("unreachable")
}
acc[1] = true
fun[1] = func(r rune) int {
  switch(r) {
  case 95: return 2
  default:
    switch {
    case 48 <= r && r <= 57: return 2
    case 65 <= r && r <= 90: return 2
    case 97 <= r && r <= 122: return 2
    default: return -1
    }
  }
  panic("unreachable")
}
acc[2] = true
fun[2] = func(r rune) int {
  switch(r) {
  case 95: return 2
  default:
    switch {
    case 48 <= r && r <= 57: return 2
    case 65 <= r && r <= 90: return 2
    case 97 <= r && r <= 122: return 2
    default: return -1
    }
  }
  panic("unreachable")
}
a0[4].acc = acc[:]
a0[4].f = fun[:]
a0[4].id = 4
}
{
var acc [29]bool
var fun [29]func(rune) int
fun[0] = func(r rune) int {
  switch(r) {
  case 58: return 1
  case 61: return 2
  case 33: return 3
  case 43: return 4
  case 37: return 5
  case 124: return 6
  case 45: return 7
  case 42: return 8
  case 47: return 9
  case 94: return 10
  case 62: return 11
  case 60: return 12
  case 38: return 13
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[11] = true
fun[11] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return 16
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[15] = true
fun[15] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[25] = true
fun[25] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[4] = true
fun[4] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return 24
  case 33: return -1
  case 43: return 25
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[21] = true
fun[21] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[1] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return 28
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[8] = true
fun[8] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return 19
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[13] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return 14
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[14] = true
fun[14] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[26] = true
fun[26] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[27] = true
fun[27] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[28] = true
fun[28] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[3] = true
fun[3] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return 26
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[5] = true
fun[5] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return 23
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[12] = true
fun[12] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return 15
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[18] = true
fun[18] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[19] = true
fun[19] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[20] = true
fun[20] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[24] = true
fun[24] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
fun[6] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return 22
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[16] = true
fun[16] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[22] = true
fun[22] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[23] = true
fun[23] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[2] = true
fun[2] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return 27
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[7] = true
fun[7] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return 20
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return 21
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[9] = true
fun[9] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return 18
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[10] = true
fun[10] = func(r rune) int {
  switch(r) {
  case 58: return -1
  case 61: return 17
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[17] = true
fun[17] = func(r rune) int {
  switch(r) {
  case 45: return -1
  case 42: return -1
  case 47: return -1
  case 94: return -1
  case 62: return -1
  case 60: return -1
  case 38: return -1
  case 58: return -1
  case 61: return -1
  case 33: return -1
  case 43: return -1
  case 37: return -1
  case 124: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
a0[5].acc = acc[:]
a0[5].f = fun[:]
a0[5].id = 5
}
{
var acc [3]bool
var fun [3]func(rune) int
fun[0] = func(r rune) int {
  switch(r) {
  case 40: return 1
  case 41: return 2
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[1] = true
fun[1] = func(r rune) int {
  switch(r) {
  case 40: return -1
  case 41: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[2] = true
fun[2] = func(r rune) int {
  switch(r) {
  case 40: return -1
  case 41: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
a0[6].acc = acc[:]
a0[6].f = fun[:]
a0[6].id = 6
}
{
var acc [3]bool
var fun [3]func(rune) int
fun[0] = func(r rune) int {
  switch(r) {
  case 123: return 1
  case 125: return 2
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[1] = true
fun[1] = func(r rune) int {
  switch(r) {
  case 123: return -1
  case 125: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[2] = true
fun[2] = func(r rune) int {
  switch(r) {
  case 123: return -1
  case 125: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
a0[7].acc = acc[:]
a0[7].f = fun[:]
a0[7].id = 7
}
{
var acc [3]bool
var fun [3]func(rune) int
fun[0] = func(r rune) int {
  switch(r) {
  case 91: return 1
  case 93: return 2
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[1] = true
fun[1] = func(r rune) int {
  switch(r) {
  case 91: return -1
  case 93: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[2] = true
fun[2] = func(r rune) int {
  switch(r) {
  case 91: return -1
  case 93: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
a0[8].acc = acc[:]
a0[8].f = fun[:]
a0[8].id = 8
}
{
var acc [5]bool
var fun [5]func(rune) int
fun[0] = func(r rune) int {
  switch(r) {
  case 44: return 1
  case 46: return 2
  case 59: return 3
  case 58: return 4
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[1] = true
fun[1] = func(r rune) int {
  switch(r) {
  case 44: return -1
  case 46: return -1
  case 59: return -1
  case 58: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[2] = true
fun[2] = func(r rune) int {
  switch(r) {
  case 44: return -1
  case 46: return -1
  case 59: return -1
  case 58: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[3] = true
fun[3] = func(r rune) int {
  switch(r) {
  case 44: return -1
  case 46: return -1
  case 59: return -1
  case 58: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[4] = true
fun[4] = func(r rune) int {
  switch(r) {
  case 44: return -1
  case 46: return -1
  case 59: return -1
  case 58: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
a0[9].acc = acc[:]
a0[9].f = fun[:]
a0[9].id = 9
}
{
var acc [8]bool
var fun [8]func(rune) int
fun[1] = func(r rune) int {
  switch(r) {
  case 34: return 2
  case 92: return 3
  default:
    switch {
    default: return 4
    }
  }
  panic("unreachable")
}
fun[3] = func(r rune) int {
  switch(r) {
  case 34: return 5
  case 92: return 6
  default:
    switch {
    default: return 7
    }
  }
  panic("unreachable")
}
fun[4] = func(r rune) int {
  switch(r) {
  case 34: return 2
  case 92: return 3
  default:
    switch {
    default: return 4
    }
  }
  panic("unreachable")
}
fun[0] = func(r rune) int {
  switch(r) {
  case 34: return 1
  case 92: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[2] = true
fun[2] = func(r rune) int {
  switch(r) {
  case 34: return -1
  case 92: return -1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[5] = true
fun[5] = func(r rune) int {
  switch(r) {
  case 34: return 2
  case 92: return 3
  default:
    switch {
    default: return 4
    }
  }
  panic("unreachable")
}
fun[6] = func(r rune) int {
  switch(r) {
  case 34: return 5
  case 92: return 6
  default:
    switch {
    default: return 7
    }
  }
  panic("unreachable")
}
fun[7] = func(r rune) int {
  switch(r) {
  case 34: return 2
  case 92: return 3
  default:
    switch {
    default: return 4
    }
  }
  panic("unreachable")
}
a0[10].acc = acc[:]
a0[10].f = fun[:]
a0[10].id = 10
}
{
var acc [2]bool
var fun [2]func(rune) int
fun[0] = func(r rune) int {
  switch(r) {
  case 32: return 1
  case 9: return 1
  case 10: return 1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
acc[1] = true
fun[1] = func(r rune) int {
  switch(r) {
  case 32: return 1
  case 9: return 1
  case 10: return 1
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
a0[11].acc = acc[:]
a0[11].f = fun[:]
a0[11].id = 11
}
{
var acc [2]bool
var fun [2]func(rune) int
fun[0] = func(r rune) int {
  switch(r) {
  default:
    switch {
    default: return 1
    }
  }
  panic("unreachable")
}
acc[1] = true
fun[1] = func(r rune) int {
  switch(r) {
  default:
    switch {
    default: return -1
    }
  }
  panic("unreachable")
}
a0[12].acc = acc[:]
a0[12].f = fun[:]
a0[12].id = 12
}
a[0].endcase = 13
a[0].a = a0[:]
}
func getAction(c *frame) int {
  if -1 == c.match { return -1 }
  c.action = c.fam.a[c.match].id
  c.match = -1
  return c.action
}
type frame struct {
  atEOF bool
  action, match, matchn, n int
  buf []rune
  text string
  in *bufio.Reader
  state []int
  fam family
}
func newFrame(in *bufio.Reader, index int) *frame {
  f := new(frame)
  f.buf = make([]rune, 0, 128)
  f.in = in
  f.match = -1
  f.fam = a[index]
  f.state = make([]int, len(f.fam.a))
  return f
}
type Lexer []*frame
func NewLexer(in io.Reader) Lexer {
  stack := make([]*frame, 0, 4)
  stack = append(stack, newFrame(bufio.NewReader(in), 0))
  return stack
}
func (stack Lexer) isDone() bool {
  return 1 == len(stack) && stack[0].atEOF
}
func (stack Lexer) nextAction() int {
  c := stack[len(stack) - 1]
  for {
    if c.atEOF { return c.fam.endcase }
    if c.n == len(c.buf) {
      r,_,er := c.in.ReadRune()
      switch er {
      case nil: c.buf = append(c.buf, r)
      case io.EOF:
	c.atEOF = true
	if c.n > 0 {
	  c.text = string(c.buf)
	  return getAction(c)
	}
	return c.fam.endcase
      default: panic(er.Error())
      }
    }
    jammed := true
    r := c.buf[c.n]
    for i, x := range c.fam.a {
      if -1 == c.state[i] { continue }
      c.state[i] = x.f[c.state[i]](r)
      if -1 == c.state[i] { continue }
      jammed = false
      if x.acc[c.state[i]] {
	if -1 == c.match || c.matchn < c.n+1 || c.match > i {
	  c.match = i
	  c.matchn = c.n+1
	}
      }
    }
    if jammed {
      a := getAction(c)
      if -1 == a { c.matchn = c.n + 1 }
      c.n = 0
      for i, _ := range c.state { c.state[i] = 0 }
      c.text = string(c.buf[:c.matchn])
      copy(c.buf, c.buf[c.matchn:])
      c.buf = c.buf[:len(c.buf) - c.matchn]
      return a
    }
    c.n++
  }
  panic("unreachable")
}
func (stack Lexer) push(index int) Lexer {
  c := stack[len(stack) - 1]
  return append(stack,
      newFrame(bufio.NewReader(strings.NewReader(c.text)), index))
}
func (stack Lexer) pop() Lexer {
  return stack[:len(stack) - 1]
}
func (stack Lexer) Text() string {
  c := stack[len(stack) - 1]
  return c.text
}
func main() { 
  lex := NewLexer(os.Stdin) 
  txt := func() string { 
		return lex.Text() 
	} 
  func(yylex Lexer) {
  for !yylex.isDone() {
    switch yylex.nextAction() {
    case -1:
    case 0:  //\/\/[^\n]+/
{ println("TOK_COMMENT:", txt())}
    case 1:  //[0-9]+/
{ println("TOK_INT_LIT:", txt()) }
    case 2:  //return|func|import|for|range|if|else/
{ println("TOK_KEYWORD:", txt()) }
    case 3:  //true|false/
{ println("TOK_BOOL_LIT:", txt()) }
    case 4:  //[a-zA-Z_][a-zA-Z_0-9]*/
{ println("TOK_IDENT:", txt()) }
    case 5:  //:=|==|=|!=|\+\+|\-\-|\+=|\-=|\*=|\/=|%=|\^=|\+|\-|\*|\/|%|\^|!|>=|<=|<|>|&&|\|\|/
{ println("TOK_OP:", txt()) }
    case 6:  //\(|\)/
{ println("TOK_PAREN:", txt()) }
    case 7:  //\{|\}/
{ println("TOK_BRACE:", txt()) }
    case 8:  //\[|\]/
{ println("TOK_INDEX:", txt()) }
    case 9:  //,|\.|;|:/
{ println("TOK_PUNC:", txt()) }
    case 10:  //"(\\.|[^"])*"/
{ println("TOK_STR_LIT:", txt()) }
    case 11:  //[ \t\n]+/
{ /* eat up whitespace */ }
    case 12:  //./
{ println("Unrecognized character:", txt()) }
    case 13:  ///
// [END]
    }
  }
  }(lex) 
}
