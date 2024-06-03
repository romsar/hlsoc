local vshard = require('vshard')
local log = require('log')

while true do
    local ok, err = vshard.router.bootstrap({
        if_not_bootstrapped = true,
    })
    if ok then
        break
    end
    log.info(('Router bootstrap error: %s'):format(err))
end

local function get_dialog_hash(a, b)
    if a < b then
        return a .. "_" .. b
    else
        return b .. "_" .. a
    end
end

function send_message(request)
    local sender_id = request:param('sender_id')
    local receiver_id = request:param('receiver_id')
    local text = request:param('text')

    local dialog_hash = get_dialog_hash(sender_id, receiver_id)
    local bucket_id = vshard.router.bucket_id_mpcrc32({receiver_id})
    vshard.router.callrw(bucket_id, 'save_message', {bucket_id, sender_id, receiver_id, text, dialog_hash, os.time()})

    return request:render({json = {success = true}})
end

function get_dialog(request)
    local sender_id = request:param('sender_id')
    local receiver_id = request:param('receiver_id')

    local dialog_hash = get_dialog_hash(sender_id, receiver_id)
    local bucket_id = vshard.router.bucket_id_mpcrc32({dialog_hash})
    local dialog = vshard.router.callro(bucket_id, 'get_messages_by_dialog_hash', {dialog_hash})

    return request:render({json = dialog})
end

local server = require('http.server').new(nil, 8083)
server:route({ path = '/api/send_message', method = 'POST' }, send_message)
server:route({ path = '/api/get_dialog', method = 'GET' }, get_dialog)
server:start()