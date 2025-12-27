# MiniLisp

An attempt to create a very small Lisp dialect with an
expandable 'metaprogramming' macro core. Also included a std file
for some typical functions you might want.

What is currently included is:

Numbers, symbols, pairs, lambdas, macros, `&rest` parameters for variadic functions, file loading so you can pull in code from other files (see below example). I've also started to
include some basic helpers for web development, such as a fetch function, basic JSON support etc.

Will also create a standard lib file using Minilisp, to help prove the language can
load in external function definitions (that's a TODO).

## Getting started

```bash
go run .
./minilisp < file.lisp
```

## Example

```lisp
; Loads the included thread macros
(load "std/macro.lisp")

(print (-> (-> 5 (* 2)) (+ 3)))  ; 13
(print (sum 1 2 3 4 5))          ; 15

(define x 5)
(print (cond
  ((< x 0) "negative")
  ((= x 0) "zero")
  ((< 0 x) "positive")))         ; "positive"
```