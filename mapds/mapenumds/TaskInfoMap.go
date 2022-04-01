package mapenumds

import (
	"fmt"

	"gitlab.com/evatix-go/core/converters"
	"gitlab.com/evatix-go/core/coreinterface"
	"gitlab.com/evatix-go/core/coreinterface/enuminf"
	"gitlab.com/evatix-go/core/coretaskinfo"
	"gitlab.com/evatix-go/enum/logtype"
	"gitlab.com/evatix-go/enum/strtype"
	"gitlab.com/evatix-go/errorwrapper"
	"gitlab.com/evatix-go/errorwrapper/errnew"
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

func (it *TaskInfoMap) ErrorWrapperWrap(
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

func (it *TaskInfoMap) ErrorWrapperWrapPayloads(
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
