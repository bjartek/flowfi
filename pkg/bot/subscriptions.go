package bot

import (
	"sync"
)

// SubscriptionData stores information about a pair's subscribers and its last processed blockNumber.
type SubscriptionData struct {
	blockNumber uint64  // Last processed block number
	chatIDs     []int64 // List of chat IDs subscribed to the pair
}

// Subscriptions struct to group subscription data by pair
type Subscriptions struct {
	mu    sync.RWMutex
	pairs map[string]*SubscriptionData // Pair -> SubscriptionData
}

func (s *Subscriptions) AddSubscription(chatID int64, pair string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if data, ok := s.pairs[pair]; ok {
		data.chatIDs = append(data.chatIDs, chatID)
	} else {
		s.pairs[pair] = &SubscriptionData{
			blockNumber: 0,
			chatIDs:     []int64{chatID},
		}
	}
}

func (s *Subscriptions) RemoveSubscription(chatID int64, pair string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data, ok := s.pairs[pair]; ok {
		// Filter out the chatID
		newChatIDs := make([]int64, 0, len(data.chatIDs))
		for _, id := range data.chatIDs {
			if id != chatID {
				newChatIDs = append(newChatIDs, id)
			}
		}
		if len(newChatIDs) > 0 {
			data.chatIDs = newChatIDs
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

func (s *Subscriptions) GetSubscriptionData(pair string) *SubscriptionData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.pairs[pair]
}

func (s *Subscriptions) UpdateBlockNumber(pair string, blockNumber uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data, ok := s.pairs[pair]; ok {
		data.blockNumber = blockNumber
	}
}
