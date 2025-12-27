; Create an Ok result containing a value
(define ok
  (lambda (value)
    (hash "type" "ok" "value" value)))

; Create an Error result containing an error message
(define err
  (lambda (error)
    (hash "type" "err" "error" error)))

; Check if a result is Ok
(define ok?
  (lambda (result)
    (= (hash-get result "type") "ok")))

; Check if a result is an Error
(define err?
  (lambda (result)
    (= (hash-get result "type") "err")))

; Unwrap a result value (panics on error)
(define unwrap
  (lambda (result)
    (if (ok? result)
        (hash-get result "value")
        (hash-get result "error"))))

; Get error message from an Err result
(define unwrap-err
  (lambda (result)
    (hash-get result "error")))

; Unwrap a result, returning the value or a default
; Example: (unwrap-or (ok 42) 0) returns 42
; Example: (unwrap-or (err "failed") 0) returns 0
(define unwrap-or
  (lambda (result default)
    (if (ok? result)
        (hash-get result "value")
        default)))
