Person = {}
MetaPerson = {}
MetaPerson.__index = Person

function Person.new(name)
    local instance = setmetatable({}, MetaPerson)
    instance.name = name
    return instance
end

function Person:printName()
    print(self.name)
end

person1 = Person.new("Bob")
person2 = Person.new("John")
person3 = Person.new("Fred")

person1:printName()
person2:printName()
person3:printName()