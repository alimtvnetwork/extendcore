package mappedfunc

import (
	"gitlab.com/evatix-go/core/coredata/corepayload"
	"gitlab.com/evatix-go/errorwrapper"
)

type (
	BytesExecutorFunc                  func(payloads []byte) *errorwrapper.Wrapper
	PayloadWrapperExecutorFunc         func(payloadWrapper *corepayload.PayloadWrapper) *errorwrapper.Wrapper
	PayloadWrapperToFuncNameGetterFunc func(payloadWrapper *corepayload.PayloadWrapper) (executorName string)
	PayloadValidatorFunc               func(payload *corepayload.PayloadWrapper) *errorwrapper.Wrapper
	AnyItemValidatorFunc               func(anyInput interface{}) *errorwrapper.Wrapper
)
