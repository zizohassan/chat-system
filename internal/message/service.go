package message

import (
	"chat-system/internal/utils"
	"encoding/json"
	"errors"
	"time"

	"chat-system/internal/db"
	"chat-system/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
)

// SendMessage sends a new message
func SendMessage(sender, recipient, content string) (models.Message, error) {
	if sender == "" || recipient == "" || content == "" {
		return models.Message{}, errors.New("all fields are required")
	}

	message := models.Message{
		ID:        gocql.TimeUUID(),
		Sender:    sender,
		Recipient: recipient,
		Timestamp: time.Now(),
		Content:   content,
	}

	err := db.Session.Query(`INSERT INTO messages (id, sender, recipient, content, timestamp) VALUES (?, ?, ?, ?, ?)`,
		message.ID, message.Sender, message.Recipient, message.Content, message.Timestamp).Exec()
	if err != nil {
		utils.GetLogger().Error("Error save message database: ", err)
		return models.Message{}, err
	}

	// Cache message in Redis
	messageJSON, err := json.Marshal(message)
	if err != nil {
		utils.GetLogger().Error("Marshal message: ", err)
		return models.Message{}, err
	}
	err = db.RedisClient.LPush(db.Ctx, "messages:"+recipient, messageJSON).Err()
	if err != nil {
		utils.GetLogger().Error("Error lpush send message: ", err)
		return models.Message{}, err
	}

	return message, nil
}

// GetMessages retrieves all messages for a user
func GetMessages(username string) ([]models.Message, error) {
	var messages []models.Message

	// Check Redis cache first
	messageJSONs, err := db.RedisClient.LRange(db.Ctx, "messages:"+username, 0, -1).Result()
	if err == redis.Nil || len(messageJSONs) == 0 {
		// Messages not found in cache, check Cassandra
		var id gocql.UUID
		var sender, recipient, content string
		var timestamp time.Time

		// Query messages where the user is the sender
		iter := db.Session.Query(`SELECT id, sender, recipient, content, timestamp FROM messages WHERE sender = ? ALLOW FILTERING`, username).Iter()
		for iter.Scan(&id, &sender, &recipient, &content, &timestamp) {
			message := models.Message{
				ID:        id,
				Sender:    sender,
				Recipient: recipient,
				Timestamp: timestamp,
				Content:   content,
			}
			messages = append(messages, message)

			// Cache message in Redis
			messageJSON, err := json.Marshal(message)
			if err == nil {
				db.RedisClient.LPush(db.Ctx, "messages:"+username, messageJSON)
			}
		}

		if err := iter.Close(); err != nil {
			return nil, err
		}

		// Query messages where the user is the recipient
		iter = db.Session.Query(`SELECT id, sender, recipient, content, timestamp FROM messages WHERE recipient = ? ALLOW FILTERING`, username).Iter()
		for iter.Scan(&id, &sender, &recipient, &content, &timestamp) {
			message := models.Message{
				ID:        id,
				Sender:    sender,
				Recipient: recipient,
				Timestamp: timestamp,
				Content:   content,
			}
			messages = append(messages, message)

			// Cache message in Redis
			messageJSON, err := json.Marshal(message)
			if err == nil {
				db.RedisClient.LPush(db.Ctx, "messages:"+username, messageJSON)
			}
		}

		if err := iter.Close(); err != nil {
			return nil, err
		}

		if len(messages) == 0 {
			return nil, errors.New("no messages found")
		}
	} else if err != nil {
		return nil, err
	} else {
		for _, messageJSON := range messageJSONs {
			var message models.Message
			err := json.Unmarshal([]byte(messageJSON), &message)
			if err != nil {
				return nil, err
			}
			messages = append(messages, message)
		}
	}

	return messages, nil
}
