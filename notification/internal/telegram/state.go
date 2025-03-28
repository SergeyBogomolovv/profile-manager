package telegram

import "sync"

const (
	stateWaitingToken          = "waiting_token"
	stateWaitingSubTypeEnable  = "waiting_sub_type_enable"
	stateWaitingSubTypeDisable = "waiting_sub_type_disable"
)

type state struct {
	data map[int64]string
	mu   sync.RWMutex
}

func NewState() *state {
	return &state{data: make(map[int64]string)}
}

func (s *state) Set(telegramID int64, data string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[telegramID] = data
}

func (s *state) Get(telegramID int64) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.data[telegramID]
	return data, ok
}

func (s *state) Delete(telegramID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, telegramID)
}
