local vshard = require('vshard')

box.once('messages', function()
    box.schema.create_space('messages', {
        format = {{
            name = 'id',
            type = 'unsigned'
        }, {
            name = 'bucket_id',
            type = 'unsigned'
        }, {
            name = 'sender_id',
            type = 'string'
        }, {
            name = 'receiver_id',
            type = 'string'
        }, {
            name = 'text',
            type = 'string'
        }, {
            name = 'dialog_hash',
            type = 'string'
        }, {
            name = 'sent_at',
            type = 'number'
        }}
    })
    box.space.messages:create_index('primary_index', {
        parts = {{
            field = 1,
            type = 'unsigned'
        }},
        sequence = true
    })
    box.space.messages:create_index('bucket_id', {
        parts = {{
            field = 2,
            type = 'unsigned'
        }},
        unique = false
    })
    box.space.messages:create_index('dialog_hash', {
        parts = {'dialog_hash', 'sent_at'},
        unique = false
    })
end)

function save_message(bucket_id, sender_id, receiver_id, text, dialog_hash, sent_at)
    box.space.messages:insert({box.NULL, bucket_id, sender_id, receiver_id, text, dialog_hash, sent_at})
end

function get_messages_by_dialog_hash(dialog_hash)
    local res = {}
    for _, message in box.space.messages.index.dialog_hash:pairs({dialog_hash}, {iterator = 'REQ'}) do
        table.insert(res, message:tomap({names_only = true}))
    end
    return res
end
