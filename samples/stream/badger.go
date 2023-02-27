package stream

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/dgraph-io/badger/v3"
	"github.com/reugn/go-streams"
)

var _ streams.Source = (*BadgerSource)(nil)

type BadgerSource struct {
	ctx    context.Context
	db     *badger.DB
	prefix []byte
	out    chan any
}

func NewBadgerSource(ctx context.Context, options badger.Options, prefix []byte) (*BadgerSource, error) {
	var err error
	bs := &BadgerSource{
		ctx:    ctx,
		prefix: prefix,
		out:    make(chan any),
	}

	bs.db, err = badger.Open(options)
	if err != nil {
		return nil, err
	}

	go bs.init()
	return bs, nil
}

func (bs *BadgerSource) init() {
	var err error
	var last []byte
	err = bs.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(bs.prefix); it.ValidForPrefix(bs.prefix); it.Next() {
			last, err = it.Item().ValueCopy(nil)
			if err != nil {
				return err
			}

			select {
			case <-bs.ctx.Done():
				return nil
			case bs.out <- last:
			}
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}

	bs.db.Close()
	close(bs.out)
}

// Via streams data through the given flow
func (bs *BadgerSource) Via(_flow streams.Flow) streams.Flow {
	go doStream(bs, _flow)
	return _flow
}

// Out returns an output channel for sending data
func (bs *BadgerSource) Out() <-chan interface{} {
	return bs.out
}

var _ streams.Sink = (*BadgerSink)(nil)

type BadgerSink struct {
	db     *badger.DB
	seq    sync.Map
	prefix BadgerPrefixExtractor
	done   chan struct{}
	in     chan any
}

type BadgerPrefixExtractor func([]byte) []byte

func NewBadgerSink(ctx context.Context, options badger.Options, prefix []byte) (*BadgerSink, error) {
	return NewBadgerSinkWithPrefixExtractor(options, func(_ []byte) []byte {
		return prefix
	})
}

func NewBadgerSinkWithPrefixExtractor(options badger.Options, prefix BadgerPrefixExtractor) (*BadgerSink, error) {
	var err error
	bs := &BadgerSink{
		prefix: prefix,
		seq:    sync.Map{},
		done:   make(chan struct{}),
		in:     make(chan any),
	}

	bs.db, err = badger.Open(options)
	if err != nil {
		return nil, err
	}

	go bs.init()
	return bs, nil
}

func (bs *BadgerSink) init() {
	for value := range bs.in {
		prefix := bs.prefix(value.([]byte))
		if len(prefix) == 0 {
			continue
		}
		seq, err := bs.sequenceNext(prefix)
		if err != nil {
			fmt.Println(err)
			continue
		}
		key := append(prefix, itob(seq)...)

		err = bs.db.Update(func(txn *badger.Txn) error {
			return txn.Set(key, value.([]byte))
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	bs.Close()
}

func (bs *BadgerSink) Close() error {
	var err error
	err = bs.db.Sync()
	if err != nil {
		return err
	}

	bs.seq.Range(func(_ any, value any) bool {
		err = value.(*badger.Sequence).Release()
		return err == nil
	})
	if err != nil {
		return err
	}

	return bs.db.Close()
}

func (bs *BadgerSink) sequenceNext(key []byte) (uint64, error) {
	seq, ok := bs.seq.Load(string(key))
	if !ok {
		newSeq, err := bs.db.GetSequence(key, 65_000)
		if err != nil {
			return 0, err
		}
		bs.seq.Store(string(key), newSeq)
		seq = newSeq
	}

	return seq.(*badger.Sequence).Next()
}

func (bs *BadgerSink) In() chan<- interface{} {
	return bs.in
}
