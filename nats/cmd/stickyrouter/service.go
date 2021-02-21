package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"golang.org/x/sync/errgroup"
)

// NewService returns a new sticky router service.
func NewService(config Config, nats *nats.Conn, redis redis.UniversalClient) (*Service, error) {
	s := &Service{
		nats:   nats,
		redis:  redis,
		config: &config,
	}
	if err := s.config.parseSubjectPattern(); err != nil {
		return nil, err
	}
	return s, nil
}

// Service is the sticky router service.
type Service struct {
	nats   *nats.Conn
	redis  redis.UniversalClient
	config *Config
}

const queueSize = 16

// Run runs the service until the Done() channel of the context is closed.
func (s *Service) Run(ctx context.Context) error {
	var err error
	ch := make(chan *nats.Msg, queueSize)
	log.Printf("Subscribe to %q", s.config.subject)
	sub, err := s.nats.ChanQueueSubscribe(s.config.subject, s.config.Queue, ch)
	if err != nil {
		return err
	}
	errGroup, errGroupCtx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		<-errGroupCtx.Done()
		log.Printf("Start draining subscription to %q", s.config.subject)
		if err := sub.Drain(); err != nil {
			return err
		}
		log.Printf("Drained subscription to %q", s.config.subject)
		close(ch)
		return nil
	})
	for i := 0; i < s.config.Workers; i++ {
		errGroup.Go(func() error {
			return s.handle(ctx, ch)
		})
	}
	return errGroup.Wait()
}

type msg struct {
	duration time.Duration
	hash     string
	data     []byte
	target   string
}

func (s *Service) readMsg(m *nats.Msg) (*msg, error) {
	tokens := strings.Split(m.Subject, ".")
	if len(tokens) < s.config.subjectTokens {
		return nil, fmt.Errorf("insufficient tokens in subject %q", m.Subject)
	}
	duration, err := time.ParseDuration(tokens[s.config.durationTokenIdx])
	if err != nil {
		return nil, fmt.Errorf("invalid duration in subject %q: %w", m.Subject, err)
	}
	hash := tokens[s.config.hashTokenIdx]
	return &msg{
		duration: duration,
		hash:     hash,
		data:     m.Data,
		target:   m.Reply,
	}, nil
}

func (s *Service) handle(ctx context.Context, ch <-chan *nats.Msg) error {
	for natsMsg := range ch {
		msg, err := s.readMsg(natsMsg)
		if err != nil {
			log.Printf("Drop invalid message on subject %q: %v", natsMsg.Subject, err)
			continue
		}
		if err = s.route(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) route(ctx context.Context, msg *msg) error {
	var res *redis.StringCmd
	_, err := s.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SetNX(ctx, msg.hash, msg.target, msg.duration)
		res = pipe.Get(ctx, msg.hash)
		return nil
	})
	if err != nil {
		return err
	}
	return s.nats.Publish(res.Val(), msg.data)
}
