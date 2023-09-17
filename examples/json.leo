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
    

:: generic T (Fn Json (Result T))
decode-json = fn (json)
    # Because this function is generic, T becomes
    # available within scope for reflection.
    map 
    	fn (field)
	    	
        fields T  

    
