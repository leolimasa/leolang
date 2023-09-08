html
    head
        script (src test)
	    do-something
	    do-something else
    body
    	p (text "hello world!")

html LE
INDENT head
INDENT script (src test) LE
INDENT do-something LE do-something else LE
DEDENT body LE
INDENT p (text "hello world!") LE

(html
	(head
		(script (src test)
			(do-something)
			(do-something else))
		(body
			(p (text "hello world")))

if just LE -> add closing paren and open paren
if LE IDENT -> do NOT add closing paren, add opening paren
if DEDENT -> add extra closing paren
		

