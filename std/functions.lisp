; ============================================
; Common Functions for MiniLisp
; ============================================

; Factorial function
(define factorial
  (lambda (n)
    (if (= n 0)
        1
        (* n (factorial (- n 1))))))

; Sum function - variadic addition with tail-recursive helper
(define sum-helper
  (lambda (nums acc)
    (if (= nums nil)
        acc
        (sum-helper (tail nums) (+ acc (head nums))))))

(define sum
  (lambda (&rest nums)
    (sum-helper nums 0)))