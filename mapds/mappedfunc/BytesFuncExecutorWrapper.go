package mappedfunc

import "gitlab.com/evatix-go/errorwrapper"

type BytesFuncExecutorWrapper struct {
	BaseExecutorWrapper
	Executor BytesExecutorFunc
}

func (it *BytesFuncExecutorWrapper) IsNull() bool {
	return it == nil
}

func (it *BytesFuncExecutorWrapper) IsAnyNull() bool {
	return it == nil ||
		it.BaseExecutorWrapper.IsNull() ||
		it.Executor == nil
}

func (it *BytesFuncExecutorWrapper) IsEmpty() bool {
	return it == nil ||
		it.BaseExecutorWrapper.IsEmpty() ||
		it.Executor == nil
}

func (it *BytesFuncExecutorWrapper) IsEmptyExecutor() bool {
	return it == nil ||
		it.Executor == nil
}

func (it *BytesFuncExecutorWrapper) HasExecutor() bool {
	return it != nil &&
		it.Executor != nil
}

func (it *BytesFuncExecutorWrapper) Exec() *errorwrapper.Wrapper {
	if it.HasAnyError() {
		return it.CompiledErrorWrapper()
	}

	return it.Executor(it.RawPayloads)
}

func (it *BytesFuncExecutorWrapper) ExecIf(
	isExecute bool,
) *errorwrapper.Wrapper {
	if !isExecute {
		return nil
	}

	return it.Exec()
}

func (it *BytesFuncExecutorWrapper) ExecByPayloadBytes() *errorwrapper.Wrapper {
	if it.HasAnyError() {
		return it.CompiledErrorWrapper()
	}

	return it.Executor(
		it.PayloadWrapper.Payloads)
}

func (it *BytesFuncExecutorWrapper) ExecSkipExistingError() *errorwrapper.Wrapper {
	return it.Executor(it.RawPayloads)
}
