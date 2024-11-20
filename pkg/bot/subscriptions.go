package bot

import (
	"sync"
)

// SubscriptionData stores information about a pair's subscribers and its last processed blockNumber.
type SubscriptionData struct {
	TokenAttributes *TokenAttributes
	ChatIDs         []int64
	BlockNumber     uint64
}

// Subscriptions struct to group subscription data by pair
type Subscriptions struct {
	pairs map[string]SubscriptionData
	mu    sync.RWMutex
}

func (s *Subscriptions) AddSubscription(chatID int64, pair string, tokenTokenAttributes *TokenAttributes) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the pair already exists
	if data, ok := s.pairs[pair]; ok {

		// Modify the existing SubscriptionData directly in the map
		data.ChatIDs = append(data.ChatIDs, chatID)
		data.TokenAttributes = tokenTokenAttributes
		s.pairs[pair] = data // Store the modified data back into the map
	} else {
		// If the pair doesn't exist, create a new entry
		s.pairs[pair] = SubscriptionData{
			BlockNumber:     0,
			ChatIDs:         []int64{chatID},
			TokenAttributes: tokenTokenAttributes,
		}
	}
}

func (s *Subscriptions) RemoveSubscription(chatID int64, pair string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data, ok := s.pairs[pair]; ok {
		// Filter out the chatID
		newChatIDs := make([]int64, 0, len(data.ChatIDs))
		for _, id := range data.ChatIDs {
			if id != chatID {
				newChatIDs = append(newChatIDs, id)
			}
		}
		if len(newChatIDs) > 0 {
			data.ChatIDs = newChatIDs
		} else {
			delete(s.pairs, pair)
		}
	}
}

func (s *Subscriptions) GetPairs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	pairs := make([]string, 0, len(s.pairs))
	for pair := range s.pairs {
		pairs = append(pairs, pair)
	}
	return pairs
}

func (s *Subscriptions) GetSubscriptionData(pair string) SubscriptionData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.pairs[pair]
}

func (s *Subscriptions) UpdateBlockNumber(pair string, blockNumber uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data, ok := s.pairs[pair]; ok {
		data.BlockNumber = blockNumber
	}
}

func (s *Subscriptions) SetLastProgressed(pair string, blockNumber uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if subscription, exists := s.pairs[pair]; exists {
		subscription.BlockNumber = blockNumber
		s.pairs[pair] = subscription
	}
}
