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


The problem: if a function with no arguments returns a function, how do we know to execute it, or pass it as a high order function?

Example:

```
# Is this a function call? Or does it return a function?
call-callback = fn (cb) (cb)
```

Solution 1: always apply by default, unless using `ref` <- this is the more readable and prevents point-free abuse.

```
# This is definitelly a function call
call-callback = fn (cb) cb

# Returns a function referenced by cb
call-callback = fn (cb) (ref cb)
```

in this case, `ref` is a trait of any type.

Solution 2: use an operator, like `&cb`
Solution 3: make function application explicit by changing the syntax - only apply it when it's a member of a list.
Solution 4: always call empty functions with (), since they are Void -> t

Solution 4:

```
# This is definitelly a function call
call-callback = fn (cb) (cb ())

# Returns a function referenced by cb
call-callback = fn (cb) cb

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

Use spaces to separate fields. Wrap it in a list to invoke the function.
Continuing after the function will drill down into subsequent fields / traits.

```
Reader = type
    Reader
        read (Fn Bytes Result)

r = Reader
r (read 10) unwrap
```

## Macros

A compile time function that takes in a SymbolExpr and returns a Result SymbolExpr E.

## Package management

* store a package.lock style file with checksums of all packages
* the package.lock also stores an allow list of all system calls for security

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

## Memory model

* idea: everything is immutable. Pass it as references (move semantics), until that reference has to refer to 
  more than one location. When that happens, create a persistent data structure.
  
## Strings

* Use the `raw-str` form for raw strings
* Use the `template-str` form for template strings

## Modules

* import other modules using urls.
* if alias is not provided, the last part of the path will be the module name (sans extension)
* if folder is imported, then all modules in the folder will be available under the folder alias
* relative folders that lie outside the root where the compiler command is being called cannot be included.

```
imports
    "git://github.com/somepackage"
    "../models/base.neon"
    "std://io"
    
```

* use the `public` form to export items from modules:

```

public
    some-fun = fn (a b) (a + b)
    some-other-fun = fn (b c) c

```

* use `submodules` 

## Features

* Immutable by default. Mutable with `mut`
* Pass by reference always
* Statically typed with type inference
* Modules: each folder is a module, like go
* Concurrency via channels
* Semi-homoiconic (everything but operators are homoiconic)
