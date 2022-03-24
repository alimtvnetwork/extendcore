package mappedfunc

type ErrorWrapOptions struct {
	IsExcludeUrl,
	IsExcludeHintUrl,
	IsExcludeRootName,
	IsExcludeErrorUrl,
	IsExcludeErrorAdditionalWrap,
	IsSecureText bool // indicates secure text, invert means log payload, plain text. it will not log payload
	IsExcludeDescription bool
}

func (it *ErrorWrapOptions) IsIncludeUrl() bool {
	return it == nil || !it.IsExcludeUrl
}

func (it *ErrorWrapOptions) IsIncludeHintUrl() bool {
	return it == nil || !it.IsExcludeHintUrl
}

func (it *ErrorWrapOptions) IsIncludeRootName() bool {
	return it == nil || !it.IsExcludeRootName
}

func (it *ErrorWrapOptions) IsIncludeErrorUrl() bool {
	return it == nil || !it.IsExcludeErrorUrl
}

func (it *ErrorWrapOptions) IsIncludeErrorAdditionalWrap() bool {
	return it == nil || !it.IsExcludeErrorAdditionalWrap
}

func (it *ErrorWrapOptions) IsIncludeDescription() bool {
	return it == nil || !it.IsExcludeDescription
}

func (it *ErrorWrapOptions) IsIncludePayloads() bool {
	return it == nil || !it.IsSecureText
}
