(define html-page
  (lambda (body-content)
          (string-append
           "<!DOCTYPE html>
                <html lang=en>
                    <head>
                        <meta charset=UTF-8>
                        <meta name=viewport content=width=device-width,initial-scale=1.0>
                        <link rel=stylesheet href=https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css>
                        <script src=https://unpkg.com/htmx.org@2.0.4></script>
                        <title>MiniLisp App</title>
                    </head>
                    <body>
                        <main class=container>"
                            body-content
           "</main></body></html>")))

(define card
  (lambda (props)
    ; slots
    (define title (hash-get props "title"))
    (define content (hash-get props "content"))

    ; template
    (<div>
      (hash "class" "card")
      (when-html title
                (<h1> (hash "class" "card-title") title))
      (<p> (hash "class" "card-content") content))))

(define router
  (lambda (request)
          (define path (hash-get request "path"))
          (cond
            ((= path "/") (home-handler request))
            ((= path "/api") (api-handler request))
            ((= path "/about") (about-handler request))
            ((= path "/counter") (counter-handler request))
            ((= path "/get-latest-count") (get-latest-count-handler request))
            (true (not-found-handler request)))))

(define home-handler
  (lambda (request)
          (hash "status" 200
                "headers" (hash "Content-Type" "text/html")
                "body" (html-page
                        (<div>
                         (<h1> "HTMX Counter Demo")
                         (<div> (hash
                         "id" "counter-display"
                         "hx-get" "/counter"
                         "hx-trigger" "load") "Loading..."))))))

(define app-state (hash "counter" 0))
(define counter-handler
  (lambda (request)
          (define current (hash-get app-state "counter"))
          (hash "status" 200
                "headers" (hash "Content-Type" "text/html")
                "body" (string-append
                        (<p> (hash "id" "count") (@string current))
                        (<button> (hash
                                   "hx-post" "/get-latest-count"
                                   "hx-target" "#count"
                                   "hx-swap" "outerHTML")
                                   "Increment Counter")))))

(define get-latest-count-handler
  (lambda (request)
          (define current (hash-get app-state "counter"))
          (define new-count (+ current 1))
          (hash-set app-state "counter" new-count)
          (hash "status" 200
                "headers" (hash "Content-Type" "text/html")
                "body" (<p> (hash "id" "count") (@string new-count)))))

(define not-found-handler
  (lambda (request)
          (hash "status" 404
                "headers" (hash "Content-Type" "text/plain")
                "body" "Not Found")))

(http-server 3000 router)