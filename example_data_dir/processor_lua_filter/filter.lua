counter = 0

function filterTxn(x)
    print("Processing transaction " .. counter)
    counter = counter + 1

    govPattern = "^af/gov1:j"
    return bytesToString(x.Note):find(govPattern) ~= nil
end
