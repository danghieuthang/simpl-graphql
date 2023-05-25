package audit

import "sync"

type IChangeTracker interface {
	// Tracking a data by key from current user instance
	Set(key string, data interface{})
	// Get data tracking by key from current user instance
	Get(key string) (interface{}, bool)
	// Clear change tracker data
	Clear()
}

type ChangeTracker struct {
	currentUser string
	data        sync.Map
}

func NewChangeTracker(currentUser string) IChangeTracker {
	return &ChangeTracker{
		currentUser: currentUser,
	}
}
func (c *ChangeTracker) Set(key string, data interface{}) {
	c.data.Store(key, data)

}
func (c *ChangeTracker) Get(key string) (interface{}, bool) {
	data, ok := c.data.Load(key)
	if ok {
		return data, true
	}
	return nil, false
}

func (c *ChangeTracker) Clear() {
	clearSyncMap(c.data)
}

func clearSyncMap(m sync.Map) {
	m.Range(func(key interface{}, value interface{}) bool {
		m.Delete(key)
		return true
	})
}
