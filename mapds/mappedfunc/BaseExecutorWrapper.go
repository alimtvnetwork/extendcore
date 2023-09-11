package mappedfunc

import (
	"gitlab.com/auk-go/core/coredata/corejson"
	"gitlab.com/auk-go/core/coredata/corepayload"
	"gitlab.com/auk-go/core/errcore"
	"gitlab.com/auk-go/errorwrapper"
	"gitlab.com/auk-go/errorwrapper/errnew"
	"gitlab.com/auk-go/errorwrapper/errtype"
)

type BaseExecutorWrapper struct {
	BaseExecutorInfo
	RawPayloads         []byte
	ErrorWrapper        *errorwrapper.Wrapper
	ParsingErrorWrapper *errorwrapper.Wrapper
	PayloadWrapper      *corepayload.PayloadWrapper
	IsExecutorFound     bool
}

func InvalidBaseExecutorWrapper(
	baseExecutorInfo *BaseExecutorInfo,
	rawBytes []byte,
	err error,
) *BaseExecutorWrapper {
	return &BaseExecutorWrapper{
		BaseExecutorInfo:    baseExecutorInfo.ToNonPtr(),
		RawPayloads:         rawBytes,
		ParsingErrorWrapper: errnew.Error.Default(errtype.Deserialize, err),
		IsExecutorFound:     false,
	}
}

func NewBaseExecutorWrapper(
	baseExecutorInfo *BaseExecutorInfo,
	rawPayloads []byte,
	executorErrorWp *errorwrapper.Wrapper,
	payloadWrapper *corepayload.PayloadWrapper,
) *BaseExecutorWrapper {
	payloadErrWp := errnew.
		Payload.
		Create(payloadWrapper).
		ConcatNew().
		Wrapper(executorErrorWp)

	return &BaseExecutorWrapper{
		BaseExecutorInfo: baseExecutorInfo.ToNonPtr(),
		RawPayloads:      rawPayloads,
		ErrorWrapper:     payloadErrWp,
		PayloadWrapper:   payloadWrapper,
		IsExecutorFound:  executorErrorWp.IsEmpty(),
	}
}

func (it *BaseExecutorWrapper) WrapErrorWithDetails(
	existingErrorWrapper *errorwrapper.Wrapper,
) *errorwrapper.Wrapper {
	return it.BaseExecutorInfo.WrapErrorWithDetails(
		existingErrorWrapper)
}

func (it *BaseExecutorWrapper) IsNull() bool {
	return it == nil
}

func (it *BaseExecutorWrapper) IsExecutorInvalid() bool {
	return it == nil || !it.IsExecutorFound
}

func (it *BaseExecutorWrapper) IsEmpty() bool {
	return it == nil ||
		it.RawPayloads == nil &&
			it.ParsingErrorWrapper.IsEmpty() &&
			it.ErrorWrapper.IsEmpty() &&
			it.PayloadWrapper.IsEmpty()
}

func (it *BaseExecutorWrapper) HasAnyError() bool {
	return it != nil &&
		it.ParsingErrorWrapper.HasError() ||
		it.ErrorWrapper.HasError()
}

func (it *BaseExecutorWrapper) HasPayloadError() bool {
	return it != nil &&
		it.ErrorWrapper.HasError()
}

func (it *BaseExecutorWrapper) HasParsingError() bool {
	return it != nil &&
		it.ParsingErrorWrapper.HasError()
}

func (it *BaseExecutorWrapper) IsEmptyPayloadError() bool {
	return it == nil ||
		it.ErrorWrapper.IsEmpty()
}

func (it *BaseExecutorWrapper) IsEmptyParsingError() bool {
	return it == nil ||
		it.ParsingErrorWrapper.IsEmpty()
}

func (it *BaseExecutorWrapper) MustBeSafe() {
	it.CompiledErrorWrapper().HandleError()
}

func (it *BaseExecutorWrapper) MustBeEmptyError() {
	it.CompiledErrorWrapper().HandleError()
}

func (it *BaseExecutorWrapper) HandleCompiledErrorWrapper() {
	it.CompiledErrorWrapper().HandleError()
}

func (it *BaseExecutorWrapper) CompiledErrorWrapper() *errorwrapper.Wrapper {
	if it == nil {
		return nil
	}

	return errnew.Merge.New(
		it.ErrorWrapper,
		it.ParsingErrorWrapper)
}

func (it *BaseExecutorWrapper) CompiledErrorWrapperWrappedDetails() *errorwrapper.Wrapper {
	if it == nil {
		return nil
	}

	return it.WrapErrorWithDetails(it.CompiledErrorWrapper())
}

func (it BaseExecutorWrapper) HasAnyItem() bool {
	return !it.IsEmpty()
}

func (it BaseExecutorWrapper) Json() corejson.Result {
	return corejson.New(it)
}

func (it BaseExecutorWrapper) JsonPtr() *corejson.Result {
	return corejson.NewPtr(it)
}

func (it *BaseExecutorWrapper) JsonParseSelfInject(jsonResult *corejson.Result) error {
	return jsonResult.Deserialize(it)
}

func (it BaseExecutorWrapper) AsJsonContractsBinder() corejson.JsonContractsBinder {
	return &it
}

func (it BaseExecutorWrapper) JsonString() string {
	return it.JsonPtr().PrettyJsonString()
}

func (it BaseExecutorWrapper) JsonStringMust() string {
	jsonResult := it.Json()
	jsonResult.MustBeSafe()

	return jsonResult.JsonString()
}

func (it BaseExecutorWrapper) String() string {
	return it.JsonPtr().PrettyJsonString()
}

func (it BaseExecutorWrapper) ToPtr() *BaseExecutorWrapper {
	return &it
}

func (it BaseExecutorWrapper) ToNonPtr() BaseExecutorWrapper {
	return it
}

func (it BaseExecutorWrapper) Clone(isDeepClone bool) BaseExecutorWrapper {
	clonedPayloadWrapper, err := it.PayloadWrapper.ClonePtr(isDeepClone)
	errcore.MustBeEmpty(err)

	return BaseExecutorWrapper{
		RawPayloads:         it.RawPayloads,
		ErrorWrapper:        it.ErrorWrapper.ClonePtr(),
		ParsingErrorWrapper: it.ParsingErrorWrapper.ClonePtr(),
		PayloadWrapper:      clonedPayloadWrapper,
	}
}

func (it BaseExecutorWrapper) ClonePtr(isDeepClone bool) *BaseExecutorWrapper {
	cloned := it.Clone(isDeepClone)

	return &cloned
}
