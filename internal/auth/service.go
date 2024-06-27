package auth

import (
	"chat-system/internal/utils"
	"encoding/json"
	"errors"
	"github.com/gocql/gocql"
	"time"

	"chat-system/internal/db"
	"chat-system/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Register a new user
func Register(username, password string) (models.User, error) {
	// Check if the username already exists
	var existingUsername string
	err := db.Session.Query(`SELECT username FROM users WHERE username = ? LIMIT 1`, username).Scan(&existingUsername)
	if err == nil {
		return models.User{}, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		ID:        gocql.TimeUUID(),
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	err = db.Session.Query(`INSERT INTO users (id, username, password, created_at) VALUES (?, ?, ?, ?)`,
		user.ID, user.Username, user.Password, user.CreatedAt).Exec()
	if err != nil {
		utils.GetLogger().Error("Error save new user on database: ", err)
		return models.User{}, err
	}

	// Cache user data in Redis
	userJSON, err := json.Marshal(user)
	if err != nil {
		utils.GetLogger().Error("Error marshal user object: ", err)
		return models.User{}, err
	}
	err = db.RedisClient.Set(db.Ctx, "user:"+username, userJSON, 0).Err()
	if err != nil {
		utils.GetLogger().Error("Error caching user: ", err)
		return models.User{}, err
	}

	return user, nil
}

// Authenticate a user
func Login(username, password string) (models.User, string, error) {
	var user models.User

	// Always check Cassandra for the latest credentials
	err := db.Session.Query(`SELECT id, username, password, created_at FROM users WHERE username = ? LIMIT 1 ALLOW FILTERING`, username).
		Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == gocql.ErrNotFound {
			utils.GetLogger().Error("User not found in Cassandra: ", err)
			return models.User{}, "", errors.New("invalid username or password")
		}
		utils.GetLogger().Error("Error querying Cassandra when login: ", err)
		return models.User{}, "", err
	}

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, "", errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		utils.GetLogger().Error("Error generate jwt: ", err)
		return models.User{}, "", err
	}

	// Cache user data in Redis after a successful login
	userJSON, err := json.Marshal(user)
	if err != nil {
		utils.GetLogger().Error("Error Marshal user login function : ", err)
		return models.User{}, "", err
	}
	err = db.RedisClient.Set(db.Ctx, "user:"+username, userJSON, 0).Err()
	if err != nil {
		utils.GetLogger().Error("Error save user after login redis : ", err)
		return models.User{}, "", err
	}

	return user, token, nil
}
