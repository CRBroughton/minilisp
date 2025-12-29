; Convert list to string (helper for HTML building)
(define list->string
  (lambda (lst)
    (string-join lst "")))

; Build attribute string from hash
(define attrs->string
  (lambda (attrs)
    (if (= attrs nil)
        ""
        (begin
          (define keys (hash-keys attrs))
          (define pairs
            (map (lambda (key)
                   (string-append key "=" (string-append (hash-get attrs key))))
                 keys))
          (string-join pairs " ")))))

; Core function: create HTML element with attributes
; Signature: (html-element tag attrs content)
(define html-element
  (lambda (tag attrs content)
    (define attr-str (attrs->string attrs))
    (define open-tag
      (if (= attrs nil)
          (string-append "<" tag ">")
          (string-append "<" tag " " attr-str ">")))

    (string-append open-tag content "</" tag ">")))


; Helper to detect if first arg is a hash (attrs)
(define hash?
  (lambda (x)
    (if (= x nil)
        nil
        (if (= (string? x) true) nil
            (if (= (number? x) true) nil
                (if (= (list? x) true) nil
                    true))))))

; Conditionally render a section of HTML (think v-if)                   
(define when-html
  (lambda (condition content)
          (if condition content "")))

(define <div>
  (lambda (first &rest rest)
    (cond
      ((= first nil) (html-element "div" nil ""))
      ((= rest nil) (html-element "div" nil first))
      ; If first is a hash, it's attrs, rest is children
      ((hash? first) (html-element "div" first (string-join rest "")))
      ; Otherwise first and rest are all children
      (true (html-element "div" nil (string-join (pair first rest) ""))))))

(define <h1>
  (lambda (first &rest rest)
    (cond
      ((= first nil) (html-element "h1" nil ""))
      ((= rest nil) (html-element "h1" nil first))
      ((hash? first) (html-element "h1" first (string-join rest "")))
      (true (html-element "h1" nil (string-join (pair first rest) ""))))))

(define <p>
  (lambda (first &rest rest)
    (cond
      ((= first nil) (html-element "p" nil ""))
      ((= rest nil) (html-element "p" nil first))
      ((hash? first) (html-element "p" first (string-join rest "")))
      (true (html-element "p" nil (string-join (pair first rest) ""))))))

(define <button>
  (lambda (first &rest rest)
    (cond
      ((= first nil) (html-element "button" nil ""))
      ((= rest nil) (html-element "button" nil first))
      ((hash? first) (html-element "button" first (string-join rest "")))
      (true (html-element "button" nil (string-join (pair first rest) ""))))))