package coremodel

import (
	"gitlab.com/auk-go/core/coredata/corejson"
	"gitlab.com/auk-go/core/coreinstruction"
	"gitlab.com/auk-go/core/coreinterface/enuminf"
)

type BasicEnumSourceDestination struct {
	Source, Destination enuminf.BasicEnumer
}

func (it *BasicEnumSourceDestination) Json() corejson.Result {
	return corejson.New(it)
}

func (it *BasicEnumSourceDestination) JsonPtr() *corejson.Result {
	return corejson.NewPtr(it)
}

func (it *BasicEnumSourceDestination) SafeSourceName() string {
	if it == nil || it.Source == nil {
		return ""
	}

	return it.Source.Name()
}

func (it *BasicEnumSourceDestination) SafeDestinationName() string {
	if it == nil || it.Destination == nil {
		return ""
	}

	return it.Destination.Name()
}

func (it *BasicEnumSourceDestination) IsAnyNull() bool {
	return it == nil || it.Source == nil || it.Destination == nil
}

func (it *BasicEnumSourceDestination) IsEmpty() bool {
	return it == nil || it.Source == nil && it.Destination == nil
}

func (it *BasicEnumSourceDestination) IsAllDefined() bool {
	return it != nil && it.Source != nil && it.Destination != nil
}

func (it *BasicEnumSourceDestination) SourceDestination() *coreinstruction.BaseSourceDestination {
	if it == nil {
		return &coreinstruction.BaseSourceDestination{}
	}

	return &coreinstruction.BaseSourceDestination{
		SourceDestination: coreinstruction.SourceDestination{
			Source:      it.SafeSourceName(),
			Destination: it.SafeDestinationName(),
		},
	}
}

func (it *BasicEnumSourceDestination) FromTo() *coreinstruction.BaseFromTo {
	if it == nil {
		return &coreinstruction.BaseFromTo{}
	}

	return &coreinstruction.BaseFromTo{
		FromTo: coreinstruction.FromTo{
			From: it.SafeSourceName(),
			To:   it.SafeDestinationName(),
		},
	}
}

func (it BasicEnumSourceDestination) AsJsoner() corejson.Jsoner {
	return &it
}
