package notify

import "github.com/mariotoffia/godeviceshadow/model/notifiermodel"

type MainBuilder struct {
	targets []notifiermodel.SelectionTargetImpl
	err     error
}

func NewBuilder() *MainBuilder {
	return &MainBuilder{}
}

func (b *MainBuilder) AddSelectionTarget(target notifiermodel.SelectionTargetImpl) *MainBuilder {
	b.targets = append(b.targets, target)
	return b
}

func (b *MainBuilder) TargetBuilder(target notifiermodel.NotificationTarget) *TargetBuilder {
	return &TargetBuilder{
		main:   b,
		target: target,
	}
}

func (b *MainBuilder) Build() *NotifierImpl {
	return &NotifierImpl{
		Targets: b.targets,
	}
}

type TargetBuilder struct {
	main      *MainBuilder
	selection notifiermodel.Selection
	target    notifiermodel.NotificationTarget
	err       error
}

func (b *TargetBuilder) WithSelection(selection notifiermodel.Selection) *TargetBuilder {
	b.selection = selection
	return b
}

func (b *TargetBuilder) WithSelectionBuilder(sb *notifiermodel.SelectionBuilder) *TargetBuilder {
	if res, err := sb.Build(); err != nil {
		b.err = err
		return b
	} else {
		b.selection = res
	}

	return b
}

func (b *TargetBuilder) Build() *MainBuilder {
	b.main.targets = append(b.main.targets, notifiermodel.SelectionTargetImpl{
		Selection: b.selection,
		Target:    b.target,
	})
	return b.main
}
