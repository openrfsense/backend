package stream

import (
	"encoding/json"
	"log"

	"github.com/reugn/go-streams"
)

var _ streams.Flow = (*Avro[any])(nil)

type JSON[T any] struct {
	in  chan any
	out chan any
}

func NewJsonUnmarshal[T any]() *JSON[T] {
	j := &JSON[T]{
		in:  make(chan any),
		out: make(chan any),
	}

	go j.unmarshal()
	return j
}

func NewJsonMarshal[T any]() *JSON[T] {
	j := &JSON[T]{
		in:  make(chan any),
		out: make(chan any),
	}

	go j.marshal()
	return j
}

func (j *JSON[T]) Via(flow streams.Flow) streams.Flow {
	go doStream(j, flow)
	return flow
}

func (j *JSON[T]) To(sink streams.Sink) {
	doStream(j, sink)
}

func (j *JSON[T]) In() chan<- any {
	return j.in
}

func (j *JSON[T]) Out() <-chan any {
	return j.out
}

func (j *JSON[T]) unmarshal() {
	var value T
	for elem := range j.in {
		err := json.Unmarshal(elem.([]byte), &value)
		if err != nil {
			log.Println(err)
			continue
		}

		j.out <- value
	}

	close(j.out)
}

func (j *JSON[T]) marshal() {
	for elem := range j.in {
		bytes, err := json.Marshal(elem.(T))
		if err != nil {
			continue
		}

		j.out <- bytes
	}

	close(j.out)
}
