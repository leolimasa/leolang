import
html
    head
        script (src test)
	    do-something
	    do-something else
    body
    	p (text "hello world!")

import LE
html
INDENT head
INDENT script (src test)
INDENT do-something LE do-something else
DEDENT body
INDENT p (text "hello world!")

(html
	(head
		(script (src test)
			(do-something)
			(do-something else))
		(body
			(p (text "hello world")))

if just LE -> add closing paren and open paren if another token comes
if INDENT -> do NOT add closing paren, add opening paren if another token comes
if DEDENT -> add as many closing parens as the dedent level. Remove one dedent if next token is an else.
		

