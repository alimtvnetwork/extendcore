package mappedfunc

import (
	"errors"

	"gitlab.com/evatix-go/core/codestack"
	"gitlab.com/evatix-go/errorwrapper"
	"gitlab.com/evatix-go/errorwrapper/errtype"
	"gitlab.com/evatix-go/errorwrapper/refs"
)

type wrapErrorWithDetails struct{}

func (it wrapErrorWithDetails) wrapErrorWrapper(
	displayModel *BaseExecutorInfo,
	options ErrorWrapOptions,
	errorWrapper *errorwrapper.Wrapper,
) *errorwrapper.Wrapper {
	if options.IsExcludeErrorAdditionalWrap || errorWrapper.IsEmpty() {
		return errorWrapper
	}

	references := it.generateReferences(
		displayModel,
		options)

	return errorWrapper.ConcatNew().MsgRefsOnly(
		references.Items()...)
}

func (it wrapErrorWithDetails) wrapErrorWrapperPayloads(
	displayModel *BaseExecutorInfo,
	options ErrorWrapOptions,
	errorWrapper *errorwrapper.Wrapper,
	payloads []byte,
) *errorwrapper.Wrapper {
	if options.IsExcludeErrorAdditionalWrap || errorWrapper.IsEmpty() {
		return errorWrapper
	}

	references := it.generateReferences(
		displayModel,
		options)

	var bytesToString string
	if len(payloads) > 0 {
		bytesToString = string(payloads)
	}

	references.AddIf(
		options.IsIncludePayloads(),
		"Payloads",
		bytesToString,
	)

	return errorWrapper.ConcatNew().MsgRefsOnly(
		references.Items()...)
}

func (it wrapErrorWithDetails) generateReferences(
	displayModel *BaseExecutorInfo,
	options ErrorWrapOptions,
) *refs.Collection {
	references := refs.New4()

	references.AddIf(
		options.IsIncludeRootName(),
		"Executor Mapper Name",
		displayModel.RootName)

	references.AddIf(
		options.IsIncludeDescription(),
		"Description",
		displayModel.Description)

	references.AddIf(
		options.IsIncludeUrl(),
		"Url",
		displayModel.Url)

	references.AddIf(
		options.IsIncludeHintUrl(),
		"Hint Url",
		displayModel.HintUrl)

	references.AddIf(
		options.IsIncludeErrorUrl(),
		"Error Url",
		displayModel.ErrorUrl)

	return references
}

func (it wrapErrorWithDetails) notFoundError(
	displayModel *BaseExecutorInfo,
	options ErrorWrapOptions,
	keyName string,
) *errorwrapper.Wrapper {
	references := it.generateReferences(
		displayModel,
		options)

	return errorwrapper.NewErrorPlusMsgUsingAllParamsPtr(
		codestack.Skip1,
		errtype.KeyNotFoundInMap,
		true,
		errors.New(keyName+" is not found in the processing map."),
		"",
		references)
}
