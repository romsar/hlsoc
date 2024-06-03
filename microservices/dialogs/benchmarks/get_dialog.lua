function getAlphaChar()
    selection = math.random(1, 3)
    if selection == 1 then return string.char(math.random(65, 90)) end
    if selection == 2 then return string.char(math.random(97, 122)) end
    return string.char(math.random(48, 57))
end

function getRandomString(length)
           length = length or 1
                if length < 1 then return nil end
                local array = {}
                for i = 1, length do
                    array[i] = getAlphaChar()
                end
                return table.concat(array)
end

request = function()
    local path = wrk.path .. "?sender_id=" .. getRandomString(10) .. "&receiver_id=" .. getRandomString(10)
    return wrk.format(wrk.method, path, wrk.headers, wrk.body)
end