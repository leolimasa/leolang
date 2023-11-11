# Example of a sql query builder

Join = struct
    table Text
    alias Text
    conditions (List Condition)

SqlQuery = struct
    fields (List Text)
    table Text
    joins (List Join)

trait QueryBuilder
    :: Fn SqlQuery (List Text) SqlQuery
    select = fn (self fields)
        self fields = fields
	self

    :: Fn SqlQuery  SqlQuery
    join = fn (self table alias conditions)
    	self joins = self joins + 
	    Join
	    	table table
		alias alias
        	condition conditions
	self

    from = fn (self table)

implement QueryBuilder SqlQuery

sql = fn () (SqlQuery () "" ())

# Usage example
sql
    select ("a" "b")
    from "people"
    join "employees" "emp"
    	


