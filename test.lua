myTable = {}

myTable[1] = "a"
myTable[2] = "b"
myTable[3] = "c"
myTable[4] = "d"
myTable["beans"] = "wtf"
myTable[69] = "420"

-- pairs() iterates over all key-value pairs in a table, in an arbitrary order
print("pairs for-loop:")
for key, value in pairs(myTable) do
    print(key, value)
end

-- ipairs() iterates over the integer keys of a table in ascending order.
-- It stops as soon as it encounters a nil value.
print("\nipairs for-loop:")
for key, value in ipairs(myTable) do
    print(key, value)
end