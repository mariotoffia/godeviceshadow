package dynamodbnotifier

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/notify/dynamodbnotifier/stream"
)

type ProcessorBuilder struct {
	stream         *stream.DynamoDBStream
	tr             model.TypeRegistry
	eventModel     ProcessorEventModel
	sb             *stream.StreamBuilder
	pcCallback     bool
	pcProcessImage ProcessImagesFunc
	startDone      ProcessorDoneFunc
}

func NewProcessorBuilder() *ProcessorBuilder {
	return &ProcessorBuilder{}
}

func (b *ProcessorBuilder) WithStreamBuilder(sb *stream.StreamBuilder) *ProcessorBuilder {
	b.sb = sb
	return b
}

func (b *ProcessorBuilder) WithCallbackProcessorCallback() *ProcessorBuilder {
	b.pcCallback = true
	return b
}

func (b *ProcessorBuilder) WithTypeRegistry(tr model.TypeRegistry) *ProcessorBuilder {
	b.tr = tr
	return b
}

func (b *ProcessorBuilder) WithEventModel(eventModel ProcessorEventModel) *ProcessorBuilder {
	b.eventModel = eventModel
	return b
}

func (b *ProcessorBuilder) WithProcessImageCallback(cb ProcessImagesFunc) *ProcessorBuilder {
	b.pcProcessImage = cb
	return b
}

func (b *ProcessorBuilder) WithStartDoneFunc(cb ProcessorDoneFunc) *ProcessorBuilder {
	b.startDone = cb
	return b
}

func (b *ProcessorBuilder) Build() (*Processor, error) {
	p := &Processor{
		stream:         b.stream,
		tr:             b.tr,
		eventModel:     b.eventModel,
		cbProcessImage: b.pcProcessImage,
		startDoneFunc:  b.startDone,
	}

	if b.sb != nil {
		if b.startDone != nil {
			b.sb.WithStartDone(func(ctx context.Context, err error) {
				b.startDone(ctx, err)
			})
		}

		if b.pcCallback {
			b.sb.WithCallback(p.HandleRequest)
		}

		if stream, err := b.sb.Build(); err != nil {
			return nil, err
		} else {
			p.stream = stream
		}
	}

	return p, nil
}
