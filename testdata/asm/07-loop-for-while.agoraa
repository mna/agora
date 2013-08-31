//
// a := 5
// sum := 0
// for a > 0 { 
//   sum += a
//   a-- // implicit constant 1
// }
// return sum
//
[f]
07-loop-for-while.agora
2
0
2
0
6
[k]
sa
i5
ssum
i0
i1
[i]
PUSH K 1
POP V 0
PUSH K 3
POP V 2
PUSH V 0 // Loop start, compare
PUSH K 3
GT _ 0
TEST J 9 // Exit condition, if false
PUSH V 2 // Loop body
PUSH V 0
ADD _ 0
POP V 2
PUSH V 0 // Decrement a
PUSH K 4
SUB _ 0
POP V 0
JMPB J 12 // Jump back to loop start
PUSH V 2 // Outside loop, return
DUMP S 1
RET _ 0
