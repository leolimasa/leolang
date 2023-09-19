# This describes how to use generic reflection

Json = enum
    JsonObj    (Map Text Json)
    JsonStr    Text
    JsonNumber Float
    JsonArray  (List Json)

ReadJsonError = inherit JsonError 
    enum
        JsonParseError String

read-json = fn (file)
    contents v= open-file file
    
# --------------------------------
# Without macros
# --------------------------------

Person = struct
    name Text
    age  Int

JsonError = enum
    NotAString
    NotANumber

decode-json-str = fn (json)
    switch json
        case JsonStr s
            Ok s
        else
            Error NotAString
           
decode-json-number = fn (json)
    switch json
        case JsonNumber s
            Ok s
        default
            Error NotANumber

DecodePersonError = inherit (DictError JsonError)
    enum JsonNotAnObject

decode-person = fn (json)
    switch json
        case (JsonObj obj)
            name-json v= (dict-key "name" obj) unwrap-bare
            name v= (decode-json-str name-json) unwrap-bare
            age-json v= (dict-key "age" obj) unwrap-bare
            age v= (decode-json-int age-json) unwrap
           
            Person
                name name
                age age

        default
            Error JsonNotAnObject

# --------------------------------
# With macros
# --------------------------------

SymbolExpr =
    Atom Text
    SymbolExprList (List SymbolExpr)

DecodeJsonErr = enum
    StructRequired

GetAtomError = enum AtomRequired

GetListError = enum ListRequired

get-atom = fn expr 
    switch expr 
        case (Atom a) a
        default Error AtomRequired

get-list = fn expr 
    switch expr 
        case (SymbolExprList a) a
        default Error AtomRequired

# expr is a struct VALUE
decode-json-struct = macro expr
    var-name v= (get-atom expr) unwrap
    struct-info v= (resolve-type var-name) unwrap

    # Require a struct for the type
    if struct-info type != StructType
        Error StructRequired

    # field assignments
    field-assns = map
        fn field
            SymbolExprList 
                list 
                    Atom (field name) 
                    SymbolExprList 
                        list
                            Atom "decode-json" 
                            var-name (field-name)
        struct-info fields
        
    SymbolExpr ((struct-info type) + field-assns)

# Example
some-fun = fn (a)
    decode-json-struct Person person-json
