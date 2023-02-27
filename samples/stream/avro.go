package stream

import (
	"log"

	"github.com/hamba/avro/v2"
	"github.com/reugn/go-streams"
)

var _ streams.Flow = (*Avro[any])(nil)

type Avro[T any] struct {
	schema avro.Schema
	in     chan any
	out    chan any
}

func NewAvroUnmarshal[T any](schema avro.Schema) *Avro[T] {
	a := &Avro[T]{
		schema: schema,
		in:     make(chan any),
		out:    make(chan any),
	}

	go a.unmarshal()
	return a
}

func NewAvroMarshal[T any](schema avro.Schema) *Avro[T] {
	a := &Avro[T]{
		schema: schema,
		in:     make(chan any),
		out:    make(chan any),
	}

	go a.marshal()
	return a
}

func (a *Avro[T]) Via(flow streams.Flow) streams.Flow {
	go doStream(a, flow)
	return flow
}

func (a *Avro[T]) To(sink streams.Sink) {
	doStream(a, sink)
}

func (a *Avro[T]) In() chan<- any {
	return a.in
}

func (a *Avro[T]) Out() <-chan any {
	return a.out
}

func (a *Avro[T]) unmarshal() {
	var value T
	for elem := range a.in {
		err := avro.Unmarshal(a.schema, elem.([]byte), &value)
		if err != nil {
			log.Println(err)
			continue
		}

		a.out <- value
	}

	close(a.out)
}

func (a *Avro[T]) marshal() {
	for elem := range a.in {
		bytes, err := avro.Marshal(a.schema, elem.(T))
		if err != nil {
			continue
		}

		a.out <- bytes
	}

	close(a.out)
}
