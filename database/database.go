package database

import (
	"sync"

	"github.com/yaonkey/gokws/models"
)

var (
	db []*models.User
	mu sync.Mutex
)

func Connect() {
	db = make([]*models.User, 0)
}

func Get() []*models.User {
	return db
}

func Insert(user *models.User) {
	mu.Lock()
	db = append(db, user)
	mu.Unlock()
}
