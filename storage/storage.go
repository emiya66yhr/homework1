package storage

import (
	"homework1/models"
	"sync"
)

var StudentData = make(map[string]*models.Student)
var Mu sync.RWMutex
