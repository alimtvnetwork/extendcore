package mappedfunc

import "gitlab.com/auk-go/core/coredata/corepayload"

func payloadCategoryNameToExecutor(
	payloadWrapper *corepayload.PayloadWrapper,
) (executorName string) {
	return payloadWrapper.CategoryName
}

func payloadNameToExecutor(
	payloadWrapper *corepayload.PayloadWrapper,
) (executorName string) {
	return payloadWrapper.PayloadName()
}
