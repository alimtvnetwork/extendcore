package mappedfunc

import (
	"gitlab.com/evatix-go/core/coredata/corepayload"
	"gitlab.com/evatix-go/errorwrapper"
)

type (
	PayloadWrapperToFuncNameGetterFunc func(payloadWrapper *corepayload.PayloadWrapper) (executorName string)
	PayloadWrapperMiddlewareFunc       func(
		middlewareName string,
		self *PayloadFuncMap,
		input interface{},
	) (output interface{}, errWrap *errorwrapper.Wrapper)
)
