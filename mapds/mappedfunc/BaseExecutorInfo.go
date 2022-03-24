package mappedfunc

import (
	"gitlab.com/evatix-go/core/coredata/corejson"
	"gitlab.com/evatix-go/errorwrapper"
)

type BaseExecutorInfo struct {
	RootName          string
	Description, Url  string
	HintUrl, ErrorUrl string
	ErrorWrapOptions  ErrorWrapOptions
	OrderedNames      []string
}

func (it BaseExecutorInfo) IsSecure() bool {
	return it.ErrorWrapOptions.IsSecureText
}

func (it BaseExecutorInfo) IsPlainText() bool {
	return it.ErrorWrapOptions.IsIncludePayloads()
}

func (it BaseExecutorInfo) Name() string {
	return it.RootName
}

func (it BaseExecutorInfo) GetDescription() string {
	return it.Description
}

func (it BaseExecutorInfo) GetUrl() string {
	return it.Url
}

func (it BaseExecutorInfo) GetErrorUrl() string {
	return it.ErrorUrl
}

func (it BaseExecutorInfo) GetErrorWrapOptions() ErrorWrapOptions {
	return it.ErrorWrapOptions
}

func (it BaseExecutorInfo) GetOrderedNames() []string {
	return it.OrderedNames
}

func (it BaseExecutorInfo) Count() int {
	return it.Length()
}

func (it *BaseExecutorInfo) Length() int {
	if it == nil {
		return 0
	}

	return len(it.OrderedNames)
}

func (it *BaseExecutorInfo) WrapErrorWithDetails(
	existingErrWrapper *errorwrapper.Wrapper,
) *errorwrapper.Wrapper {
	if existingErrWrapper.IsEmpty() {
		return nil
	}

	return wrapErrorWithDetailsInstance.wrapErrorWrapper(
		it,
		it.ErrorWrapOptions,
		existingErrWrapper)
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
		it,
		it.ErrorWrapOptions,
		existingErrWrapper,
		result.SafeValues())
}

func (it BaseExecutorInfo) IsEmpty() bool {
	return it.Length() == 0
}

func (it BaseExecutorInfo) HasAnyItem() bool {
	return it.Length() > 0
}

func (it BaseExecutorInfo) Json() corejson.Result {
	return corejson.New(it)
}

func (it BaseExecutorInfo) JsonPtr() *corejson.Result {
	return corejson.NewPtr(it)
}

func (it *BaseExecutorInfo) JsonParseSelfInject(jsonResult *corejson.Result) error {
	return jsonResult.Deserialize(it)
}

func (it BaseExecutorInfo) AsJsonContractsBinder() corejson.JsonContractsBinder {
	return &it
}

func (it BaseExecutorInfo) JsonString() string {
	return it.JsonPtr().PrettyJsonString()
}

func (it BaseExecutorInfo) JsonStringMust() string {
	jsonResult := it.Json()
	jsonResult.MustBeSafe()

	return jsonResult.JsonString()
}

func (it BaseExecutorInfo) String() string {
	return it.Json().PrettyJsonString()
}

func (it BaseExecutorInfo) ToPtr() *BaseExecutorInfo {
	return &it
}

func (it BaseExecutorInfo) ToNonPtr() BaseExecutorInfo {
	return it
}

func (it BaseExecutorInfo) Clone() BaseExecutorInfo {
	return BaseExecutorInfo{
		RootName:         it.RootName,
		Description:      it.Description,
		Url:              it.Url,
		HintUrl:          it.HintUrl,
		ErrorUrl:         it.ErrorUrl,
		ErrorWrapOptions: it.ErrorWrapOptions,
		OrderedNames:     it.OrderedNames,
	}
}

func (it BaseExecutorInfo) ClonePtr() *BaseExecutorInfo {
	return &BaseExecutorInfo{
		RootName:         it.RootName,
		Description:      it.Description,
		Url:              it.Url,
		HintUrl:          it.HintUrl,
		ErrorUrl:         it.ErrorUrl,
		ErrorWrapOptions: it.ErrorWrapOptions,
		OrderedNames:     it.OrderedNames,
	}
}
