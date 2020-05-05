package lru

import lru "github.com/hashicorp/golang-lru"

type Cache interface {
	Purge()
	Len() int
	Contains(k interface{}) bool
	Get(k interface{}) interface{}
	Add(k, v interface{})
	Remove(k interface{})
}

type Simple lru.Cache

var _ Cache = (*Simple)(nil)

func NewSimple(size int) (Cache, error) {
	c, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return (*Simple)(c), nil
}

func (c *Simple) Purge()                        { (*lru.Cache)(c).Purge() }
func (c *Simple) Len() int                      { return (*lru.Cache)(c).Len() }
func (c *Simple) Contains(k interface{}) bool   { return (*lru.Cache)(c).Contains(k) }
func (c *Simple) Get(k interface{}) interface{} { v, _ := (*lru.Cache)(c).Get(k); return v }
func (c *Simple) Add(k, v interface{})          { (*lru.Cache)(c).Add(k, v) }
func (c *Simple) Remove(k interface{})          { (*lru.Cache)(c).Remove(k) }

type Two lru.TwoQueueCache

var _ Cache = (*Two)(nil)

func NewTwoQueue(size int) (Cache, error) {
	c, err := lru.New2Q(size)
	if err != nil {
		return nil, err
	}
	return (*Two)(c), nil
}

func (c *Two) Purge()                        { (*lru.TwoQueueCache)(c).Purge() }
func (c *Two) Len() int                      { return (*lru.TwoQueueCache)(c).Len() }
func (c *Two) Contains(k interface{}) bool   { return (*lru.TwoQueueCache)(c).Contains(k) }
func (c *Two) Get(k interface{}) interface{} { v, _ := (*lru.TwoQueueCache)(c).Get(k); return v }
func (c *Two) Add(k, v interface{})          { (*lru.TwoQueueCache)(c).Add(k, v) }
func (c *Two) Remove(k interface{})          { (*lru.TwoQueueCache)(c).Remove(k) }
