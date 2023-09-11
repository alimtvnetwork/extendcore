package mappedfunc

import (
	"errors"

	"gitlab.com/auk-go/core/codestack"
	"gitlab.com/auk-go/core/coredata/corejson"
	"gitlab.com/auk-go/core/coretaskinfo"
	"gitlab.com/auk-go/errorwrapper"
	"gitlab.com/auk-go/errorwrapper/errnew"
	"gitlab.com/auk-go/errorwrapper/errtype"
)

type wrapErrorWithDetails struct{}

func (it wrapErrorWithDetails) wrapErrorWrapper(
	existingErr *errorwrapper.Wrapper,
	info *coretaskinfo.Info,
) *errorwrapper.Wrapper {
	if it.isExcludeErrorWrap(info) || existingErr.IsEmpty() {
		return existingErr
	}

	return errnew.Friendly.WrapErrorWithDetailsNoPayloads(
		existingErr,
		info)
}

func (it wrapErrorWithDetails) isExcludeErrorWrap(info *coretaskinfo.Info) bool {
	return info != nil && info.ExcludeOptions.IsExcludeAdditionalErrorWrap
}

func (it wrapErrorWithDetails) wrapErrorWrapperPayloads(
	existingErr *errorwrapper.Wrapper,
	info *coretaskinfo.Info,
	payloads []byte,
) *errorwrapper.Wrapper {
	if it.isExcludeErrorWrap(info) || existingErr.IsEmpty() {
		return existingErr
	}

	return errnew.Friendly.WrapErrorWithDetailsPayloads(
		existingErr,
		info,
		payloads)
}

func (it wrapErrorWithDetails) wrapErrorWrapperPayloadsAny(
	existingErr *errorwrapper.Wrapper,
	info *coretaskinfo.Info,
	payloadsAny interface{},
) *errorwrapper.Wrapper {
	if it.isExcludeErrorWrap(info) || existingErr.IsEmpty() {
		return existingErr
	}

	return errnew.Friendly.WrapErrorWithDetailsPayloads(
		existingErr,
		info,
		corejson.AnyTo.SerializedJsonResult(payloadsAny).SafeBytes())
}

func (it wrapErrorWithDetails) notFoundError(
	keyName string,
	info *coretaskinfo.Info,
) *errorwrapper.Wrapper {
	references := errnew.Friendly.CompiledReferencesNoPayloads(
		nil,
		info,
	)

	return errorwrapper.NewErrorPlusMsgUsingAllParamsPtr(
		codestack.Skip1,
		errtype.KeyNotFoundInMap,
		true,
		errors.New(keyName+" is not found in the processing map."),
		"",
		references.ToPtr())
}

// notFoundErrorWithPayloads
//
//	can be bytes or any
func (it wrapErrorWithDetails) notFoundErrorWithPayloads(
	keyName string,
	info *coretaskinfo.Info,
	payloads interface{}, // will be cast to bytes to marshalled
) *errorwrapper.Wrapper {
	references := errnew.Friendly.CompiledReferencesPayloads(
		nil,
		info,
		corejson.AnyTo.SerializedJsonResult(payloads).SafeBytes())

	return errorwrapper.NewErrorPlusMsgUsingAllParamsPtr(
		codestack.Skip1,
		errtype.KeyNotFoundInMap,
		true,
		errors.New(keyName+" is not found in the processing map."),
		"",
		references.ToPtr())
}
