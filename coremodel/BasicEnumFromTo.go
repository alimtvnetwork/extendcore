package coremodel

import (
	"gitlab.com/auk-go/core/coredata/corejson"
	"gitlab.com/auk-go/core/coreinstruction"
	"gitlab.com/auk-go/core/coreinterface/enuminf"
)

type BasicEnumFromTo struct {
	From, To enuminf.BasicEnumer
}

func (it *BasicEnumFromTo) Json() corejson.Result {
	return corejson.New(it)
}

func (it *BasicEnumFromTo) JsonPtr() *corejson.Result {
	return corejson.NewPtr(it)
}

func (it *BasicEnumFromTo) SafeSourceName() string {
	if it == nil || it.From == nil {
		return ""
	}

	return it.From.Name()
}

func (it *BasicEnumFromTo) SafeDestinationName() string {
	if it == nil || it.To == nil {
		return ""
	}

	return it.To.Name()
}

func (it *BasicEnumFromTo) IsAnyNull() bool {
	return it == nil || it.From == nil || it.To == nil
}

func (it *BasicEnumFromTo) IsEmpty() bool {
	return it == nil || it.From == nil && it.To == nil
}

func (it *BasicEnumFromTo) IsAllDefined() bool {
	return it != nil && it.From != nil && it.To != nil
}

func (it *BasicEnumFromTo) SourceDestination() *coreinstruction.BaseSourceDestination {
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

func (it *BasicEnumFromTo) FromTo() *coreinstruction.BaseFromTo {
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

func (it BasicEnumFromTo) AsJsoner() corejson.Jsoner {
	return &it
}
