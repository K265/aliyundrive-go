package drive

import (
	lru "github.com/hashicorp/golang-lru"
)

type FolderCache interface {
	Get(string) (*Node, bool)
	Put(string, *Node)
	Clear()
}

type CacheImpl struct {
	cache *lru.Cache
}

func NewCache(size int) (FolderCache, error) {
	c, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &CacheImpl{cache: c}, nil
}

func (c *CacheImpl) Get(key string) (*Node, bool) {
	value, ok := c.cache.Get(key)
	if !ok {
		return nil, ok
	}

	node, ok := value.(*Node)
	return node, ok
}

func (c *CacheImpl) Put(key string, node *Node) {
	c.cache.Add(key, node)
}

func (c *CacheImpl) Clear() {
	c.cache.Purge()
}
