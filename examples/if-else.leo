if a > b then a else b

PAREN IF IDENT OP IDENT THEN IDENT ELSE IDENT PAREN


if a > b then
   print a
   print b
else
   print c

(if a > b then
   (print a)
   (print b)
)(else
   (print c))


(if a > b then
   (print a)
   (print b)
else
   (print c))

if a > b then a
else 1 + 3


(if a > b then a)
(else 1 + 3)

if a then b else c

(if a then b else c)

Two forms:

(if cond st1 st2 ELSE st1 st2)

OR

(if cond st1 st2)(ELSE st1 st2)
