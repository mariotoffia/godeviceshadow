package dynamodbnotifier

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/loggers/desirelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/notify/dynamodbnotifier/loggers"
	"github.com/mariotoffia/godeviceshadow/notify/dynamodbnotifier/stream"
	"github.com/mariotoffia/godeviceshadow/utils/loggerutils"
)

type ProcessorBuilder struct {
	stream                          *stream.DynamoDBStream
	tr                              model.TypeRegistry
	eventModel                      ProcessorEventModel
	sb                              *stream.StreamBuilder
	pcCallback                      bool
	reportedLoggers, desiredLoggers []model.CreatableMergeLogger
	pcProcessImage                  ProcessImagesFunc
	startDone                       ProcessorDoneFunc
	notificationManager             notifiermodel.Notifier
	cbNotifyProcess                 notifiermodel.ProcessFunc
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

func (b *ProcessorBuilder) WithNotificationManager(nm notifiermodel.Notifier) *ProcessorBuilder {
	b.notificationManager = nm
	return b
}

func (b *ProcessorBuilder) WithNotifyProcessCallback(cb notifiermodel.ProcessFunc) *ProcessorBuilder {
	b.cbNotifyProcess = cb
	return b
}

// WithReportMergeLogger adds a logger that logs the reported values.
//
// It will automatically add the `changelogger.ChangeMergeLogger` to detect the merge of reports.
func (b *ProcessorBuilder) WithReportMergeLogger(logger ...model.CreatableMergeLogger) *ProcessorBuilder {
	b.reportedLoggers = append(b.reportedLoggers, logger...)
	return b
}

// WithDesireMergeLogger adds a logger that logs the desired values.
//
// It will automatically add the `changelogger.ChangeMergeLogger` and `loggers.DynamoDbDesireLogger` to detect
// when a desire has been acknowledged. The former will log a plain desire document merge and the latter will
// produce a log when a desire has been acknowledged.
func (b *ProcessorBuilder) WithDesireMergeLogger(logger ...model.CreatableMergeLogger) *ProcessorBuilder {
	b.desiredLoggers = append(b.desiredLoggers, logger...)
	return b
}

func (b *ProcessorBuilder) Build() (*Processor, error) {
	// Add default loggers if not already added.
	if loggerutils.FindCreatableMerge[*changelogger.ChangeMergeLogger](b.reportedLoggers) == nil {
		b.reportedLoggers = append(b.reportedLoggers, changelogger.New())
	}

	if loggerutils.FindCreatableMerge[*loggers.DynamoDbDesireLogger](b.desiredLoggers) == nil {
		b.desiredLoggers = append(b.desiredLoggers, changelogger.New())
		b.desiredLoggers = append(b.desiredLoggers, &loggers.DynamoDbDesireLogger{DesireLogger: desirelogger.New()})
	}

	p := &Processor{
		stream:              b.stream,
		tr:                  b.tr,
		eventModel:          b.eventModel,
		cbProcessImage:      b.pcProcessImage,
		startDoneFunc:       b.startDone,
		reportedLoggers:     b.reportedLoggers,
		desiredLoggers:      b.desiredLoggers,
		notificationManager: b.notificationManager,
		cbNotifyProcess:     b.cbNotifyProcess,
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
