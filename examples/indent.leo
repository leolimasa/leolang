import
html
    ident-with-no-open-paren
    head
        script (src test)
	    do-something
	    do-something else
    body
    	p (text "hello world!")

import LE 
html INDENT
indent-with-no-open-paren LE
head INDENT
script (src test) INDENT
do-something LE
do-something else DEDENT 
body INDENT 
p (text "hello world!")

WITH PARENS

(html
	indent-with-no-open-paren
	(head
		(script (src test)
			do-something
			(do-something else))
		(body
			(p (text "hello world")))

	
# New rules

* Lines end when there is an LE or an INDENT or a DEDENT

* Auto-grouping
    * Lines that start with open paren will never be auto-grouped
    * Lines that have a single token will never be auto-grouped
    * Lines that end in INDENT will never be auto-grouped

* Lines that end with an indent will add an open parens at the beginning
* Lines that end in dedent add closing parens equivalent to the dedent levels (indepdendently of auto grouping)


Auto group if a line DOESN'T start with an open paren or line has more than one token



