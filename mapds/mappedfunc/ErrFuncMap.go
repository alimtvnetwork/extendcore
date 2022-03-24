package mappedfunc

import (
	"sort"

	"gitlab.com/evatix-go/core/corecsv"
	"gitlab.com/evatix-go/core/coredata/coredynamic"
	"gitlab.com/evatix-go/core/coreinterface"
	"gitlab.com/evatix-go/core/coreinterface/enuminf"
	"gitlab.com/evatix-go/errorwrapper"
	"gitlab.com/evatix-go/errorwrapper/errfunc"
	"gitlab.com/evatix-go/errorwrapper/errnew"
	"gitlab.com/evatix-go/errorwrapper/errwrappers"
)

type ErrFuncMap struct {
	BaseExecutorInfo
	FunctionsMap map[string]errfunc.WrapperFunc
}

func (it *ErrFuncMap) Append(
	isSkipOnExist bool,
	executorsMap *ErrFuncMap,
) *ErrFuncMap {
	if executorsMap.IsEmpty() {
		return it
	}

	for name, executor := range executorsMap.FunctionsMap {
		if isSkipOnExist && it.HasFunc(name) {
			continue
		}

		it.Set(name, executor)
	}

	return it
}

func (it *ErrFuncMap) ConcatNew(
	isSkipOnExist bool,
	executorsMap *ErrFuncMap,
) *ErrFuncMap {
	if executorsMap.IsEmpty() {
		return it.ClonePtr()
	}

	cloned := it.ClonePtr()

	return cloned.Append(isSkipOnExist, executorsMap)
}

func (it ErrFuncMap) Items() map[string]errfunc.WrapperFunc {
	return it.FunctionsMap
}

func (it ErrFuncMap) Count() int {
	return it.Length()
}

func (it ErrFuncMap) LastIndex() int {
	return it.Length() - 1
}

func (it ErrFuncMap) HasIndex(index int) bool {
	return it.LastIndex() >= index
}

func (it ErrFuncMap) RemoveAt(index int) (isSuccess bool) {
	panic("not supported for map")
}

func (it ErrFuncMap) Set(
	name string,
	executor errfunc.WrapperFunc,
) (isAddedNewly bool) {
	isAddedNewly = it.IsMissingFunc(name)

	it.FunctionsMap[name] = executor

	return isAddedNewly
}

func (it ErrFuncMap) SetIf(
	isAdd bool,
	name string,
	executor errfunc.WrapperFunc,
) (isAddedNewly bool) {
	if !isAdd {
		return it.IsMissingFunc(name)
	}

	return it.Set(name, executor)
}

func (it ErrFuncMap) AddOrUpdate(
	name string,
	executor errfunc.WrapperFunc,
) (isAddedNewly bool) {
	return it.Set(name, executor)
}

func (it *ErrFuncMap) SetChain(
	name string,
	executor errfunc.WrapperFunc,
) *ErrFuncMap {
	it.Set(name, executor)

	return it
}

func (it *ErrFuncMap) SetChainIf(
	isSet bool,
	name string,
	executor errfunc.WrapperFunc,
) *ErrFuncMap {
	if !isSet {
		return it
	}

	it.Set(name, executor)

	return it
}

func (it *ErrFuncMap) AddOrUpdateChain(
	name string,
	executor errfunc.WrapperFunc,
) *ErrFuncMap {
	it.Set(name, executor)

	return it
}

func (it ErrFuncMap) ListStrings() []string {
	return it.AllNames()
}

func (it *ErrFuncMap) Length() int {
	if it == nil {
		return 0
	}

	return len(it.FunctionsMap)
}

func (it ErrFuncMap) IsEmpty() bool {
	return it.Length() == 0
}

func (it ErrFuncMap) HasAnyItem() bool {
	return it.Length() > 0
}

func (it ErrFuncMap) AllKeys() []string {
	return it.AllNames()
}

func (it ErrFuncMap) SortedNamesCsv() string {
	return corecsv.DefaultCsv(it.AllNamesSorted()...)
}

func (it ErrFuncMap) AllKeysSorted() []string {
	return it.AllNamesSorted()
}

func (it ErrFuncMap) AllNames() []string {
	names := make([]string, 0, it.Length())

	for s := range it.FunctionsMap {
		names = append(names, s)
	}

	return names
}

func (it ErrFuncMap) AllNamesSorted() []string {
	names := it.AllNames()
	sort.Strings(names)

	return names
}

func (it ErrFuncMap) Get(name string) errfunc.WrapperFunc {
	return it.FunctionsMap[name]
}

func (it ErrFuncMap) GetWithStat(name string) (errfunc.WrapperFunc, bool) {
	currFunc, has := it.FunctionsMap[name]

	return currFunc, has
}

func (it ErrFuncMap) HasFunc(name string) bool {
	_, has := it.FunctionsMap[name]

	return has
}

func (it ErrFuncMap) IsMissingFunc(name string) bool {
	_, has := it.FunctionsMap[name]

	return !has
}

func (it ErrFuncMap) ExecAllLock() *errwrappers.Collection {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	return it.ExecAll()
}

func (it ErrFuncMap) ExecAll() *errwrappers.Collection {
	if it.IsEmpty() {
		return errwrappers.Empty()
	}

	if len(it.OrderedNames) > 0 {
		return it.executeByOrderedNamesAll()
	}

	displayModel := it.BaseExecutorInfo.ToPtr()
	emptyCollector := errwrappers.Empty()

	for _, errFunc := range it.FunctionsMap {
		// found, exec
		errorWrapper := errFunc()
		wrappedError := wrapErrorWithDetailsInstance.wrapErrorWrapper(
			displayModel,
			it.ErrorWrapOptions,
			errorWrapper,
		)

		emptyCollector.AddWrapperPtr(wrappedError)
	}

	return emptyCollector
}

func (it ErrFuncMap) executeByOrderedNamesAll() *errwrappers.Collection {
	displayModel := it.BaseExecutorInfo.ToPtr()
	emptyCollector := errwrappers.Empty()

	for _, name := range it.OrderedNames {
		wrappedErr := it.execInternal(displayModel, name)
		emptyCollector.AddWrapperPtr(wrappedErr)
	}

	return emptyCollector
}

func (it ErrFuncMap) ExecLock(name string) *errorwrapper.Wrapper {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	return it.Exec(name)
}

func (it ErrFuncMap) Exec(name string) *errorwrapper.Wrapper {
	displayModel := it.BaseExecutorInfo.ToPtr()

	return it.execInternal(displayModel, name)
}

func (it ErrFuncMap) ExecNamer(namer enuminf.Namer) *errorwrapper.Wrapper {
	displayModel := it.BaseExecutorInfo.ToPtr()

	return it.execInternal(displayModel, namer.Name())
}

func (it ErrFuncMap) execInternal(
	displayModel *BaseExecutorInfo,
	name string,
) *errorwrapper.Wrapper {
	errFunc, hasFunc := it.FunctionsMap[name]

	if !hasFunc {
		return wrapErrorWithDetailsInstance.notFoundError(
			displayModel,
			it.ErrorWrapOptions,
			name)
	}

	// found, executor
	errorWrapper := errFunc()
	if errorWrapper.IsEmpty() {
		return nil
	}

	return wrapErrorWithDetailsInstance.wrapErrorWrapper(
		displayModel,
		it.ErrorWrapOptions,
		errorWrapper,
	)
}

func (it ErrFuncMap) ExecByNamesLock(
	names ...string,
) *errwrappers.Collection {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	return it.ExecByNames(names...)
}

func (it ErrFuncMap) ExecByNames(
	names ...string,
) *errwrappers.Collection {
	emptyCollector := errwrappers.Empty()

	if len(names) == 0 {
		return emptyCollector
	}

	displayModel := it.BaseExecutorInfo.ToPtr()

	for _, name := range names {
		wrappedErr := it.execInternal(displayModel, name)
		emptyCollector.AddWrapperPtr(wrappedErr)
	}

	return emptyCollector
}

func (it ErrFuncMap) ExecByNamesIf(
	isExecute bool,
	names ...string,
) *errwrappers.Collection {
	if !isExecute {
		return errwrappers.Empty()
	}

	return it.ExecByNames(names...)
}

func (it ErrFuncMap) ExecByNamesLockIf(
	isExecute bool,
	names ...string,
) *errwrappers.Collection {
	if !isExecute {
		return errwrappers.Empty()
	}

	return it.ExecByNamesLock(names...)
}

func (it *ErrFuncMap) ReflectSetTo(toPtr interface{}) error {
	return coredynamic.ReflectSetFromTo(it, toPtr)
}

func (it *ErrFuncMap) ReflectSetToErrWrap(toPtr interface{}) *errorwrapper.Wrapper {
	return errnew.Reflect.SetFromTo(it, toPtr)
}

func (it ErrFuncMap) ClonePtr() *ErrFuncMap {
	cloned := it.Clone()

	return &cloned
}

func (it ErrFuncMap) Clone() ErrFuncMap {
	if it.IsEmpty() {
		return ErrFuncMap{
			BaseExecutorInfo: it.BaseExecutorInfo.Clone(),
			FunctionsMap:     map[string]errfunc.WrapperFunc{},
		}
	}

	newMap := make(
		map[string]errfunc.WrapperFunc,
		it.Length())
	for name, executor := range it.FunctionsMap {
		newMap[name] = executor
	}

	return ErrFuncMap{
		BaseExecutorInfo: it.BaseExecutorInfo.Clone(),
		FunctionsMap:     newMap,
	}
}

func (it ErrFuncMap) AsStandardSlicerContractsBinder() coreinterface.StandardSlicerContractsBinder {
	return &it
}
