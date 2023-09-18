# Leo Lang LISP version

```
#idea #language
```

A high-performance type-safe language with focus on readability, concurrency, and full stack development.

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

## Pritive types

The same as go lang, with maybe the exception of runes.

## Assignments

```
myfun = fn (a b)
  a + b
```

## IF statements

Single line:

```
if (a > b) a else b
```

Multi line:

```
if (a > b)
    a
else
    b
```


## Types

are declared with `::` form, and inferenced if not present. Always start with uppercase letter:

```
:: Fn String String String
my-fun = fn (a b) (a + b)
```

Enums

```
Employee = enum Employee
    Admin String Int
    FrontLine Int
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

## Wrap operator

```
user-name = fn url
  user v= get-user unwrap
  user name 
```

is equivalent to

```
user-name = fn url
  (get-user unwrap) (fn user user.name)
```

## High Order Functions

Functions will be passed by value or by applying depending on the type of the
receptor.

```


```

## Traits

If the `self` keyword is used in the function arg, then the struct will be passed to the the function.

```
Unwrappable = trait T V
    :: Fn T (Fn V V)
    unwrap

BareUnwrappable = trait T V
    :: Fn T (Fn V T)
    bare-unwrap

Result = generic-type (K E)
    Ok
        value K
    Error
        value E

implement Unwrappable Result
    unwrap = fn (value cb)
        case value
	    Ok value
	        Ok (cb value)
	    Err err
	        err

implement BareUnwrappable Result
    bare-unwrap = fn (value cb)
        case value
	    Ok value
                cb value
	    Err err
	        err
```

## Field accessors

Use spaces.

```
Reader = type
    Reader
        read (Fn Bytes Result)

r = Reader
((r read) 10) unwrap
```

## Generic reflection

Generics are available within functions to create type safe versions.

## Symbolic Expressions

* Group of symbols
* The first symbol defines an action:
  * fn: function definition
  * trait: trait definition
  * inherit
  * type: type definition
  * generic-type: generic type definition
  * implement: trait implementation
  * (function): function call
  * variable: value accessor. Any subsequent fields are field accessors.
    * Field accessors may be trait calls 
  * if / else

## Features

* Immutable by default. Mutable with `mut`
* Pass by reference always
* Statically typed with type inference
* Modules: each folder is a module, like go
* Concurrency via channels
* Semi-homoiconic (everything but operators are homoiconic)
