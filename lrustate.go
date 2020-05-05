// Package lrustate implements a state storage with LRU capabilities.
package lrustate

// Store implements the state storage using a special LRU cache. This store
// pertains all IDs, but purges structs as usual.
//
// In some cases, functions that return a single struct would be a cache-hit,
// but functions that return a slice would be a cache-miss. This is to prevent
// incomplete slices.
type Store struct {
	GuildStore
	ChannelStore
}
