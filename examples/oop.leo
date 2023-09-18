
Person = struct
    name Text
    age Int

PersonMethods = trait

    :: Fn Person Text
    to-text = fn (self) 
    	self.name + " " + self.age.to-text

new-person = fn (name age)
    Person
        name name
	age age

main = fn ()
    p = new-person "joe" 32
    print p.to-text

