package lrustate

import (
	"sync"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/state"
	"github.com/diamondburned/arikawa/state/lrustate/lru"
)

type GuildStore struct {
	guildIDmu sync.RWMutex
	guildIDs  map[discord.Snowflake]struct{}
	guilds    lru.Cache // discord.Snowflake : *discord.Guild
}

func (s *GuildStore) Reset() {
	s.guilds.Purge()

	s.guildIDmu.Lock()
	s.guildIDs = map[discord.Snowflake]struct{}{}
	s.guildIDmu.Unlock()
}

func (s *GuildStore) Guild(id discord.Snowflake) (*discord.Guild, error) {
	if v := s.guilds.Get(id); v != nil {
		return v.(*discord.Guild), nil
	}
	return nil, state.ErrStoreNotFound
}

func (s *GuildStore) Guilds() ([]discord.Guild, error) {
	s.guildIDmu.RLock()
	defer s.guildIDmu.RUnlock()

	if len(s.guildIDs) > s.guilds.Len() {
		return nil, state.ErrStoreNotFound
	}

	guilds := make([]discord.Guild, 0, len(s.guildIDs))

	for id := range s.guildIDs {
		if v := s.guilds.Get(id); v != nil {
			guilds = append(guilds, *(v.(*discord.Guild)))
		}
	}

	return guilds, nil
}

func (s *GuildStore) GuildSet(g *discord.Guild) error {
	s.guilds.Add(g.ID, g)

	s.guildIDmu.Lock()
	s.guildIDs[g.ID] = struct{}{}
	s.guildIDmu.Unlock()

	return nil
}

func (s *GuildStore) GuildRemove(id discord.Snowflake) error {
	s.guilds.Remove(id)

	s.guildIDmu.Lock()
	delete(s.guildIDs, id)
	s.guildIDmu.Unlock()

	return nil
}
