package mapenumds

import (
	"fmt"

	"gitlab.com/auk-go/core/converters"
	"gitlab.com/auk-go/core/coreinterface"
	"gitlab.com/auk-go/core/coreinterface/enuminf"
	"gitlab.com/auk-go/core/coretaskinfo"
	"gitlab.com/auk-go/enum/logtype"
	"gitlab.com/auk-go/enum/strtype"
	"gitlab.com/auk-go/errorwrapper"
	"gitlab.com/auk-go/errorwrapper/errnew"
)

type TaskInfoMap struct {
	Items map[fmt.Stringer]*coretaskinfo.Info
}

func (it *TaskInfoMap) Length() int {
	if it == nil || it.Items == nil {
		return 0
	}

	return len(it.Items)
}

func (it *TaskInfoMap) Count() int {
	return it.Length()
}

func (it *TaskInfoMap) IsEmpty() bool {
	return it.Length() == 0
}

func (it *TaskInfoMap) HasAnyItem() bool {
	return it.Length() > 0
}

func (it *TaskInfoMap) LastIndex() int {
	return it.Length() - 1
}

func (it *TaskInfoMap) HasIndex(index int) bool {
	return it.LastIndex() >= index
}

func (it *TaskInfoMap) Add(name string, info *coretaskinfo.Info) *TaskInfoMap {
	it.Items[strtype.Variant(name)] = info

	return it
}

func (it *TaskInfoMap) AddStinger(
	stringer fmt.Stringer,
	info *coretaskinfo.Info,
) *TaskInfoMap {
	it.Items[stringer] = info

	return it
}

func (it *TaskInfoMap) AddNamer(
	namer enuminf.Namer,
	info *coretaskinfo.Info,
) *TaskInfoMap {
	if namer == nil {
		return it
	}

	it.Items[strtype.New(namer.Name())] = info

	return it
}

func (it *TaskInfoMap) AddEnum(
	enum enuminf.BaseEnumer,
	info *coretaskinfo.Info,
) *TaskInfoMap {
	if enum == nil {
		return it
	}

	it.Items[enum] = info

	return it
}

func (it *TaskInfoMap) AddSimpleEnum(
	enum enuminf.SimpleEnumer,
	info *coretaskinfo.Info,
) *TaskInfoMap {
	if enum == nil {
		return it
	}

	it.Items[enum] = info

	return it
}

func (it *TaskInfoMap) GetInfo(
	name string,
) *coretaskinfo.Info {
	if it.IsEmpty() {
		return nil
	}

	return it.Items[strtype.New(name)]
}

func (it *TaskInfoMap) GetInfoUsingStringer(
	stringer fmt.Stringer,
) *coretaskinfo.Info {
	if it.IsEmpty() || stringer == nil {
		return nil
	}

	return it.Items[stringer]
}

func (it *TaskInfoMap) GetInfoUsingNamer(
	namer enuminf.Namer,
) *coretaskinfo.Info {
	if it.IsEmpty() || namer == nil {
		return nil
	}

	return it.Items[strtype.New(namer.Name())]
}

func (it *TaskInfoMap) ErrWrap(
	name string,
	existingErrWrap *errorwrapper.Wrapper,
) *errorwrapper.Wrapper {
	if it.IsEmpty() || name == "" {
		return existingErrWrap
	}

	info := it.GetInfo(name)

	return errnew.Ref.WrapExistingErrWrapWithInfo(
		info,
		existingErrWrap)
}

func (it *TaskInfoMap) ErrWrapWithPayloadsAny(
	name string,
	existingErrWrap *errorwrapper.Wrapper,
	payloadsAny interface{},
) *errorwrapper.Wrapper {
	if it.IsEmpty() || name == "" {
		return existingErrWrap
	}

	info := it.GetInfo(name)

	return errnew.Ref.WrapExistingErrWrapWithInfoPayloadsAny(
		info,
		existingErrWrap,
		payloadsAny)
}

func (it *TaskInfoMap) FriendlyErrWrap(
	name string,
	friendlyMessage string,
	existingErrWrap *errorwrapper.Wrapper,
) *errorwrapper.FriendlyError {
	if it.IsEmpty() || name == "" {
		return errnew.Friendly.Create(
			friendlyMessage,
			existingErrWrap)
	}

	if existingErrWrap.IsEmpty() {
		return nil
	}

	info := it.GetInfo(name)

	return errnew.Friendly.All(
		friendlyMessage,
		existingErrWrap,
		logtype.Error,
		info,
		nil)
}

func (it *TaskInfoMap) FriendlyErrWrapWithPayloads(
	name string,
	friendlyMessage string,
	existingErrWrap *errorwrapper.Wrapper,
	payloads []byte,
) *errorwrapper.FriendlyError {
	if it.IsEmpty() || name == "" {
		return errnew.Friendly.Create(
			friendlyMessage,
			existingErrWrap)
	}

	if existingErrWrap.IsEmpty() {
		return nil
	}

	info := it.GetInfo(name)

	return errnew.Friendly.All(
		friendlyMessage,
		existingErrWrap,
		logtype.Error,
		info,
		payloads)
}

func (it *TaskInfoMap) FriendlyErrWrapWithPayloadsAny(
	name string,
	friendlyMessage string,
	existingErrWrap *errorwrapper.Wrapper,
	payloadsAny interface{},
) *errorwrapper.FriendlyError {
	if it.IsEmpty() || name == "" {
		return errnew.Friendly.Create(
			friendlyMessage,
			existingErrWrap)
	}

	if existingErrWrap.IsEmpty() {
		return nil
	}

	info := it.GetInfo(name)

	return errnew.Friendly.PayloadsAnyInfo(
		friendlyMessage,
		existingErrWrap,
		info,
		payloadsAny)
}

func (it *TaskInfoMap) String() string {
	if it == nil || it.IsEmpty() {
		return ""
	}

	return converters.AnyToValueString(
		it.Items)
}

func (it TaskInfoMap) AsBasicMapper() coreinterface.BasicMapper {
	return &it
}

func (it TaskInfoMap) ToPtr() *TaskInfoMap {
	return &it
}

func (it *TaskInfoMap) ToNonPtr() TaskInfoMap {
	if it == nil {
		return TaskInfoMap{}
	}

	return *it
}
