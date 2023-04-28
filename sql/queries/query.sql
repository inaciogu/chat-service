-- name: CreateChat :exec
INSERT INTO chats (id, user_id, initial_message_id, status, token_usage, model, model_max_tokens, temperature, top_p, n, stop, max_tokens, presence_penalty, frequency_penalty, created_at, updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);

-- name: AddMessage :exec
INSERT INTO messages (id, chat_id, role, content, tokens, model, erased, order_msg, created_at) VALUES(?,?,?,?,?,?,?,?,?);

-- name: FindMessagesByChatId :many
SELECT * FROM messages WHERE chat_id = ? and erased=0 ORDER BY order_msg ASC;

-- name: FindErasedMessagesByChatId :many
SELECT * FROM messages WHERE chat_id = ? and erased=1 ORDER BY order_msg ASC;

-- name: FindChatById :one
SELECT * FROM chats WHERE id = ?;

-- name: SaveChat :exec
UPDATE chats SET user_id = ?, initial_message_id = ?, status = ?, token_usage = ?, model = ?, model_max_tokens=?, temperature = ?, top_p = ?, n = ?, stop = ?, max_tokens = ?, presence_penalty = ?, frequency_penalty = ?, updated_at = ? WHERE id = ?;

-- name: DeleteChatMessages :exec
DELETE FROM messages WHERE chat_id = ?;

-- name: DeleteErasedChatMessages :exec
DELETE FROM messages WHERE erased=1 and chat_id = ?;