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

### Using the REPL

```bash
go build
./minilisp
```

The REPL starts with the standard library loaded (thread macros, when, cond, factorial, sum):

```
> (-> 5 (* 2) (+ 3))
13
> (factorial 5)
120
> (load "fetch.lisp")
```

### Running a file

```bash
./minilisp < file.lisp
```

## Examples

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
Here is an example of fetching data:

```lisp
(define get-github-user
  (lambda (username)
    (->> username
         (string-append "https://api.github.com/users/")
         (fetch)
         (json-parse))))

(define user (get-github-user "crbroughton"))
(-> user (hash-get "login") (print))
(-> user (hash-get "public_repos") (print))
```

Here is an example combining fetch with Result type and cond for error handling:

```lisp
(load "std/macro.lisp")
(load "std/result.lisp")

(define get-github-user
  (lambda (username)
    (cond
      ((= username "") (err "Username cannot be empty"))
      (true (ok (->> username
                     (string-append "https://api.github.com/users/")
                     (fetch) ; TODO - make the fetch return a Result type
                     (json-parse)))))))

(define print-user-info
  (lambda (result)
    (cond
      ((ok? result)
        (begin
          (define user (unwrap result))
          (print (string-append "User: " (hash-get user "login")))
          (print (string-append "Repos: " (->string (hash-get user "public_repos"))))))
      ((err? result)
        (print (string-append "Error: " (unwrap-err result)))))))

(->> "crbroughton" (get-github-user) (print-user-info))
```

Here is an example of conditionals:
```lisp
(load "std/macro.lisp")

(define x 5)
(print (cond
  ((< x 0) "negative")
  ((= x 0) "zero")
  ((< 0 x) "positive"))) ; need to add a > operator
```
