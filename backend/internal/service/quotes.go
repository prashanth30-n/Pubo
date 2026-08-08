package service

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/prashanth30-n/Pubo/internal/model/Quotes"
	"github.com/prashanth30-n/Pubo/internal/server"
)

type QuoteService struct {
	server *server.Server
	client *http.Client
	quotes []Quotes.Quote
	index  int
	mu     sync.Mutex
}

func NewQuoteService(s *server.Server) *QuoteService {
	service := &QuoteService{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		server: s,
	}
	// Populate the existing quote cache once when this service is created.
	go service.FetchQuotes()
	return service
}

var fallbackQuotes = []Quotes.Quote{
	{Text: "Build in public Learn faster.", Author: "Pubo"},
	{Text: "Consistency beats intensity.", Author: "Pubo"},
	{Text: "Start before you are ready.", Author: "Pubo"},
}

func randomFallback() Quotes.Quote {
	return fallbackQuotes[rand.IntN(len(fallbackQuotes))]
}
func (s *QuoteService) FetchQuotes() {
	var collected []Quotes.Quote
	for len(collected) < 50 {
		resp, err := s.client.Get("https://zenquotes.io/api/quotes")
		if err != nil {
			break
		}
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			resp.Body.Close()
			break
		}
		var data []Quotes.ZenQuoteDTO
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			resp.Body.Close()
			break
		}
		resp.Body.Close()
		if len(data) == 0 {
			break
		}
		for _, q := range data {
			collected = append(collected, Quotes.Quote{
				Text:   q.Q,
				Author: q.A,
			})
			if len(collected) >= 50 {
				break
			}

		}
	}
	s.mu.Lock()
	if len(collected) > 0 {
		s.quotes = collected
		s.index = 0
	} else {
		s.quotes = fallbackQuotes
	}
	s.mu.Unlock()

}

func (s *QuoteService) NextQuote() Quotes.Quote {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.quotes) == 0 {
		return randomFallback()
	}
	quote := s.quotes[s.index]
	s.index = (s.index + 1) % len(s.quotes)
	return quote

}
