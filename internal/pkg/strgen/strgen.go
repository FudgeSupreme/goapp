package strgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"
)

type StringGenerator struct {
	strChan     chan<- string  // String output channel.
	quitChannel chan struct{}  // Quit.
	running     sync.WaitGroup // Running.
}

func New(strChan chan<- string) *StringGenerator {
	s := StringGenerator{}
	s.strChan = strChan
	s.quitChannel = make(chan struct{})
	s.running = sync.WaitGroup{}
	return &s
}

// Start string generator. Stop() must be called at the end.
func (s *StringGenerator) Start() error {
	s.running.Add(1)
	go s.mainLoop()

	return nil
}

func (s *StringGenerator) Stop() {
	close(s.quitChannel)
	s.running.Wait()
}

func (s *StringGenerator) mainLoop() {
	defer s.running.Done()

	for {
		select {
		case <-s.quitChannel:
			return
		default:
			hexValue, err := s.generateHexValues(5)
			if err != nil {
				log.Printf("%v\n", err)
			} else {
				s.strChan <- hexValue
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// Generate a hex value using the crypto/rand package
func (s *StringGenerator) generateHexValues(valueLen int) (string, error) {
	bytes := make([]byte, valueLen)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("could not generate hex value because of %v", err)
	}
	return hex.EncodeToString(bytes), nil
}
