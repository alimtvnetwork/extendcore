package mappedfunc

import (
	"gitlab.com/evatix-go/errorwrapper"
	"gitlab.com/evatix-go/errorwrapper/errfunc"
	"gitlab.com/evatix-go/errorwrapper/errnew"
	"gitlab.com/evatix-go/errorwrapper/errtype"
)

type PayloadFuncExecutorWrapper struct {
	BaseExecutorWrapper
	Executor errfunc.PayloadWrapperExecutorFunc
}

func (it *PayloadFuncExecutorWrapper) IsNull() bool {
	return it == nil
}

func (it *PayloadFuncExecutorWrapper) IsAnyNull() bool {
	return it == nil ||
		it.BaseExecutorWrapper.IsNull() ||
		it.Executor == nil
}

func (it *PayloadFuncExecutorWrapper) IsEmpty() bool {
	return it == nil ||
		it.BaseExecutorWrapper.IsEmpty() ||
		it.Executor == nil
}

func (it *PayloadFuncExecutorWrapper) IsEmptyExecutor() bool {
	return it == nil ||
		it.Executor == nil
}

func (it *PayloadFuncExecutorWrapper) HasExecutor() bool {
	return it != nil &&
		it.Executor != nil
}

func (it *PayloadFuncExecutorWrapper) Exec() *errorwrapper.Wrapper {
	if it.HasAnyError() {
		return it.CompiledErrorWrapper()
	}

	return it.Executor(it.PayloadWrapper)
}

func (it *PayloadFuncExecutorWrapper) ExecIf(
	isExecute bool,
) *errorwrapper.Wrapper {
	if !isExecute {
		return nil
	}

	return it.Exec()
}

func (it *PayloadFuncExecutorWrapper) ExecByPayloadBytes() *errorwrapper.Wrapper {
	if it.HasAnyError() {
		return it.CompiledErrorWrapper()
	}

	newPayload, err := it.PayloadWrapper.DeserializePayloadsToPayloadWrapper()

	if err != nil {
		return errnew.Error.Default(errtype.Deserialize, err)
	}

	return it.Executor(newPayload)
}

func (it *PayloadFuncExecutorWrapper) ExecSkipExistingError() *errorwrapper.Wrapper {
	return it.Executor(it.PayloadWrapper)
}
