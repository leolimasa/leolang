# Leo Lang LISP version

```
#idea #language
```

A TYPED LISP without the parenthesis. Focus on concurrency.

## Automatic parenthesis removal by indentation

* Each line is automatically enclosed in parenthesis.
* The last parenthesis is **not** added if the **next** line indentation level increases
* A closing parenthesis is added for each level of **dedent** of the next line
* Parenthesis are **not** added if the line has a single word **and** the next line **does not** have an ident level increase

```
(define (enumerate-list lst)
  (map 
    (lambda (num item) (string-append (number->string num) ") " item)) 
    (range 0 (length lst))
    lst))
```

Would become:

```
define (enumerate-list lst)
  map 
    lambda (num item) (string-append (number->string num) ") " item) 
    range 0 (length lst)
    lst
```

## Operators

The language has operators:

Assignment: =
Comparison: <,>,<=,>=,!=,==
Logical: and, not, or
Arithmetic: `+,-, * ,/, **, %`
Callback: v=

They are implemented as traits.

## Assignments

```
myfun = fn (a b)
  a + b
```

## IF statements

```
if (a > b)
  "Statement is true"
  "Statement is false"
```

## Types

are declared with `::` form, and inferenced if not present. Always start with uppercase letter:

```
:: Fn String String String
my-fun = fn (a b) (a + b)
```

Unions

```
union Employee
    Admin (name String) (salary Int)
    FrontLine (salary Int)
``` 

Using types

```
:: Fn String Int
new-admin = fn (name salary) (Admin name salary)
```

With named fields

```
new-admin = fn (adm-name adm-salary)
  Admin
    name adm-name
    salary adm-salary
```

## Callback operator

```
user-name = fn url
  user v= await get-user
  user.name 
```

is equivalent to

```
user-name = fn url
  await get-user (fn user user.name)
```

## Functions

If an identifier is a function and it sits by itself on a line, you need to force calling it by surrounding in parameters (so it turns into a list).
If a function is called with less parameters than it takes, then it is partially applied.

```


```

## Traits

```
Unwrappable = trait T V
    unwrap (Fn T (Fn V V))

BareUnwrappable = trait T V
    bare-unwrap (Fn T (Fn V T))

Result = generic-type (K E)
    Ok
        value K
    Error
        value E

implement Unwrappable Result
    unwrap (value cb)
        case value
	    Ok value
	        Ok (cb value)
	    Err err
	        err

implement BareUnwrappable Result
    bare-unwrap (value cb)
        case value
	    Ok value
                cb value
	    Err err
	        err
```

## Field accessors

Use spaces

```
Reader = type
    Reader
        read (Fn Bytes Result)

r = Reader
(r read 10) unwrap
```

## Symbolic Expressions

* Group of symbols
* The first symbol defines an action:
  * fn: function definition
  * trait: trait definition
  * type: type definition
  * generic-type: generic type definition
  * implement: trait implementation
  * (function): function call
  * variable: value accessor. Any subsequent fields are field accessors.
    * Field accessors may be trait calls 

## Features

* Immutable by default. Mutable with `mut`
* Pass by reference always
* Statically typed with type inference
* Modules: each folder is a module, like go
* Concurrency via channels
* Semi-homoiconic (everything but operators are homoiconic)