; ============================================
; Macro Library for MiniLisp
; ============================================

; ============================================
; Thread Macros
; ============================================

; Thread-first macro (->)
; Threads value through multiple forms, inserting as FIRST argument
; Example: (-> 5 (* 2) (+ 3)) expands to (+ (* 5 2) 3)
(defmacro -> (value &rest forms)
  (if (= forms nil)
      value
      (if (= (tail forms) nil)
          ; Single form - thread value as first arg
          (pair (head (head forms))
                (pair value (tail (head forms))))
          ; Multiple forms - recurse
          (pair '->
                (pair (pair (head (head forms))
                           (pair value (tail (head forms))))
                      (tail forms))))))

; Helper function for ->> macro
(define append-at-end
  (lambda (form val)
    (if (= form nil)
        (pair val nil)
        (pair (head form) (append-at-end (tail form) val)))))

; Thread-last macro (->>)
; Threads value through multiple forms, inserting as LAST argument
; Example: (->> 5 (* 2) (+ 3)) expands to (+ 3 (* 2 5))
(defmacro ->> (value &rest forms)
  (if (= forms nil)
      value
      (if (= (tail forms) nil)
          ; Single form - append value at end
          (append-at-end (head forms) value)
          ; Multiple forms - recurse
          (pair '->>
                (pair (append-at-end (head forms) value)
                      (tail forms))))))

; ============================================
; Control Flow Macros
; ============================================

; When macro - conditional execution without else branch
; Example: (when true (print 42))
(defmacro when (test body)
  (pair 'if (pair test (pair body (pair 'nil nil)))))

; Cond macro - multi-way conditional
; Takes multiple (test expr...) pairs and returns the result of the first true test
; Supports multiple expressions per clause
; Example: (cond ((= x 0) 100) ((< x 0) 200) (true 300))
; Example: (cond ((ok? r) (print "ok") (unwrap r)) ((err? r) (print "err")))
(defmacro cond (&rest clauses)
  (if (= clauses nil)
      nil
      (pair 'if
            (pair (head (head clauses))
                  (pair (if (= (tail (tail (head clauses))) nil)
                            ; Single expression
                            (head (tail (head clauses)))
                            ; Multiple expressions - wrap in begin
                            (pair 'begin (tail (head clauses))))
                        (pair (pair 'cond (tail clauses))
                              nil))))))

; Map function - applies a function to each element of a list
(define map
  (lambda (f lst)
    (if (null? lst)
        nil
        (pair (f (head lst))
              (map f (tail lst))))))