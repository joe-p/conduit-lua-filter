-- Provides helper methods that can be used in filter.lua

-- Convert a []byte to a string
-- Probably better to do this in Go, but forming the new struct seems like a lot of work... TODO with codegen
function bytesToString(byteSlice)
    local chars = {}
    for i in byteSlice() do
        currentByte = byteSlice[i]
        local utf8byte = currentByte < 0 and (0xff + currentByte + 1) or currentByte
        table.insert(chars, string.char(utf8byte))
    end
    return table.concat(chars)
end
