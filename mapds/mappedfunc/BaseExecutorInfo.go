package mappedfunc

import (
	"gitlab.com/auk-go/core/coredata/corejson"
	"gitlab.com/auk-go/core/coretaskinfo"
	"gitlab.com/auk-go/errorwrapper"
)

type BaseExecutorInfo struct {
	coretaskinfo.Info
	OrderedNames []string
}

func (it *BaseExecutorInfo) GetOrderedNames() []string {
	return it.OrderedNames
}

func (it *BaseExecutorInfo) Count() int {
	return it.Length()
}

func (it *BaseExecutorInfo) Length() int {
	if it == nil {
		return 0
	}

	return len(it.OrderedNames)
}

func (it *BaseExecutorInfo) SafeInfo() *coretaskinfo.Info {
	if it == nil {
		return nil
	}

	return &it.Info
}

func (it *BaseExecutorInfo) WrapErrorWithDetails(
	existingErrWrapper *errorwrapper.Wrapper,
) *errorwrapper.Wrapper {
	if existingErrWrapper.IsEmpty() {
		return nil
	}

	return wrapErrorWithDetailsInstance.wrapErrorWrapper(
		existingErrWrapper,
		it.SafeInfo())
}

// WrapErrorWithDetailsPlusPayload
//
// If IsPlainText then logs payloads
func (it *BaseExecutorInfo) WrapErrorWithDetailsPlusPayload(
	existingErrWrapper *errorwrapper.Wrapper,
	payloads interface{},
) *errorwrapper.Wrapper {
	if existingErrWrapper.IsEmpty() {
		return nil
	}

	if it.IsSecure() {
		return it.WrapErrorWithDetails(existingErrWrapper)
	}

	result := Converter.AnyToBytesOrSerializedBytes(payloads)
	if result.HasError() {
		existingErrWrapper.ConcatNew().Wrapper(result.ErrorWrapper)
	}

	return wrapErrorWithDetailsInstance.wrapErrorWrapperPayloads(
		existingErrWrapper,
		it.SafeInfo(),
		result.SafeValues())
}

func (it *BaseExecutorInfo) IsEmpty() bool {
	return it.Length() == 0
}

func (it *BaseExecutorInfo) HasAnyItem() bool {
	return it.Length() > 0
}

func (it *BaseExecutorInfo) Json() corejson.Result {
	return corejson.New(it)
}

func (it *BaseExecutorInfo) JsonPtr() *corejson.Result {
	return corejson.NewPtr(it)
}

func (it *BaseExecutorInfo) JsonParseSelfInject(jsonResult *corejson.Result) error {
	return jsonResult.Deserialize(it)
}

func (it BaseExecutorInfo) AsJsonContractsBinder() corejson.JsonContractsBinder {
	return &it
}

func (it *BaseExecutorInfo) JsonString() string {
	return it.JsonPtr().PrettyJsonString()
}

func (it *BaseExecutorInfo) JsonStringMust() string {
	jsonResult := it.Json()
	jsonResult.MustBeSafe()

	return jsonResult.JsonString()
}

func (it *BaseExecutorInfo) String() string {
	return it.JsonPtr().PrettyJsonString()
}

func (it BaseExecutorInfo) ToPtr() *BaseExecutorInfo {
	return &it
}

func (it *BaseExecutorInfo) ToNonPtr() BaseExecutorInfo {
	if it == nil {
		return BaseExecutorInfo{}
	}

	return *it
}

func (it *BaseExecutorInfo) Clone() BaseExecutorInfo {
	if it == nil {
		return BaseExecutorInfo{}
	}

	return BaseExecutorInfo{
		Info:         it.Info.Clone(),
		OrderedNames: it.OrderedNames,
	}
}

func (it *BaseExecutorInfo) ClonePtr() *BaseExecutorInfo {
	if it == nil {
		return &BaseExecutorInfo{}
	}

	return &BaseExecutorInfo{
		Info:         it.Info.Clone(),
		OrderedNames: it.OrderedNames,
	}
}
