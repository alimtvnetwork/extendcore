package mappedfunc

import "gitlab.com/evatix-go/core/coredata/corepayload"

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
