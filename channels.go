package lrustate

import (
	"sync"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/state"
	"github.com/diamondburned/arikawa/state/lrustate/lru"
)

type ChannelStore struct {
	channelIDmu sync.RWMutex
	// map of guildID -> channelID -> nil
	channelIDs map[discord.Snowflake]map[discord.Snowflake]struct{}
	channels   lru.Cache // discord.Snowflake : *discord.Channel

	privateIDmu sync.RWMutex
	privateIDs  map[discord.Snowflake]struct{}
	privates    lru.Cache // discord.Snowflake : *discord.Channel
}

func (s *ChannelStore) Reset() {
	s.channels.Purge()

	s.channelIDmu.Lock()
	s.channelIDs = map[discord.Snowflake]map[discord.Snowflake]struct{}{}
	s.channelIDmu.Unlock()

	s.privates.Purge()

	s.privateIDmu.Lock()
	s.privateIDs = map[discord.Snowflake]struct{}{}
	s.privateIDmu.Unlock()
}

func (s *ChannelStore) Channel(id discord.Snowflake) (*discord.Channel, error) {
	if v := s.channels.Get(id); v != nil {
		return v.(*discord.Channel), nil
	}
	return nil, state.ErrStoreNotFound
}

func (s *ChannelStore) Channels(guildID discord.Snowflake) ([]discord.Channel, error) {
	s.channelIDmu.RLock()
	defer s.channelIDmu.RUnlock()

	chIDs, ok := s.channelIDs[guildID]
	if !ok {
		return nil, state.ErrStoreNotFound
	}

	// Don't allow returns if the LRU cache is incomplete.
	if len(chIDs) > s.channels.Len() {
		return nil, state.ErrStoreNotFound
	}

	chans := make([]discord.Channel, 0, len(chIDs))

	for id := range chIDs {
		if v := s.channels.Get(id); v != nil {
			chans = append(chans, *(v.(*discord.Channel)))
		}
	}

	return chans, nil
}

func (s *ChannelStore) CreatePrivateChannel(user discord.Snowflake) (*discord.Channel, error) {
	if v := s.privates.Get(user); v != nil {
		return v.(*discord.Channel), nil
	}
	return nil, state.ErrStoreNotFound
}

func (s *ChannelStore) PrivateChannels() ([]discord.Channel, error) {
	s.privateIDmu.RLock()
	defer s.privateIDmu.RUnlock()

	// Don't allow returns if the LRU cache is incomplete.
	if len(s.privateIDs) > s.privates.Len() {
		return nil, state.ErrStoreNotFound
	}

	var chans = make([]discord.Channel, 0, len(s.privateIDs))

	for id := range s.privateIDs {
		if v := s.privates.Get(id); v != nil {
			chans = append(chans, *(v.(*discord.Channel)))
		}
	}

	return chans, nil
}

func (s *ChannelStore) ChannelSet(ch *discord.Channel) error {
	if !ch.GuildID.Valid() {
		s.privates.Add(ch.ID, ch)

		s.privateIDmu.Lock()
		s.privateIDs[ch.ID] = struct{}{}
		s.privateIDmu.Unlock()

		return nil
	}

	s.channels.Add(ch.ID, ch)

	s.channelIDmu.Lock()
	id, ok := s.channelIDs[ch.GuildID]
	if !ok {
		id = make(map[discord.Snowflake]struct{}, 1)
		s.channelIDs[ch.GuildID] = id
	}
	id[ch.ID] = struct{}{}
	s.channelIDmu.Unlock()

	return nil
}

func (s *ChannelStore) ChannelRemove(ch *discord.Channel) error {
	if !ch.GuildID.Valid() {
		s.privates.Remove(ch.ID)

		s.privateIDmu.Lock()
		delete(s.privateIDs, ch.ID)
		s.privateIDmu.Unlock()

		return nil
	}

	s.channels.Remove(ch.ID)

	s.channelIDmu.Lock()
	if m, ok := s.channelIDs[ch.GuildID]; ok {
		delete(m, ch.ID)
	}
	s.channelIDmu.Unlock()

	return nil
}
