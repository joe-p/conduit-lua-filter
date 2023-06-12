counter = 0

function filterTxn(signedTxn)
    txn = signedTxn.Txn

    -- TODO define logger
    print("Processing transaction " .. counter)
    counter = counter + 1

    if (txn.Type ~= 'pay') then
        return false
    end

    govPattern = "^af/gov1:j"
    return bytesToString(txn.Note):find(govPattern) ~= nil
end
