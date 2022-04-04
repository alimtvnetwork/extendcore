package mappedfunc

type PayloadWrapperMiddleware struct {
	Name           string
	MiddlewareFunc PayloadWrapperMiddlewareFunc
}
