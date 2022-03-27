package mappedfunc

import (
	"sort"

	"gitlab.com/evatix-go/core/corecsv"
	"gitlab.com/evatix-go/core/coredata/coredynamic"
	"gitlab.com/evatix-go/core/coredata/corepayload"
	"gitlab.com/evatix-go/core/coreinterface"
	"gitlab.com/evatix-go/core/coreinterface/enuminf"
	"gitlab.com/evatix-go/errorwrapper"
	"gitlab.com/evatix-go/errorwrapper/errnew"
	"gitlab.com/evatix-go/errorwrapper/errwrappers"
)

type BytesFuncMap struct {
	BaseExecutorInfo
	FunctionsMap map[string]BytesExecutorFunc
}

func (it *BytesFuncMap) Append(
	isSkipOnExist bool,
	executorsMap *BytesFuncMap,
) *BytesFuncMap {
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

func (it *BytesFuncMap) ConcatNew(
	isSkipOnExist bool,
	executorsMap *BytesFuncMap,
) *BytesFuncMap {
	if executorsMap.IsEmpty() {
		return it.ClonePtr()
	}

	cloned := it.ClonePtr()

	return cloned.Append(isSkipOnExist, executorsMap)
}

func (it BytesFuncMap) Items() map[string]BytesExecutorFunc {
	return it.FunctionsMap
}

func (it BytesFuncMap) Count() int {
	return it.Length()
}

func (it BytesFuncMap) LastIndex() int {
	return it.Length() - 1
}

func (it BytesFuncMap) HasIndex(index int) bool {
	return it.LastIndex() >= index
}

func (it BytesFuncMap) RemoveAt(index int) (isSuccess bool) {
	panic("not supported for map")
}

func (it BytesFuncMap) Set(
	name string,
	executor BytesExecutorFunc,
) (isAddedNewly bool) {
	isAddedNewly = it.IsMissingFunc(name)
	it.FunctionsMap[name] = executor

	return isAddedNewly
}

func (it BytesFuncMap) SetIf(
	isAdd bool,
	name string,
	executor BytesExecutorFunc,
) (isAddedNewly bool) {
	if !isAdd {
		return it.IsMissingFunc(name)
	}

	return it.Set(name, executor)
}

func (it BytesFuncMap) AddOrUpdate(
	name string,
	executor BytesExecutorFunc,
) (isAddedNewly bool) {
	return it.Set(name, executor)
}

func (it *BytesFuncMap) SetChain(
	name string,
	executor BytesExecutorFunc,
) *BytesFuncMap {
	it.Set(name, executor)

	return it
}

func (it *BytesFuncMap) SetChainIf(
	isSet bool,
	name string,
	executor BytesExecutorFunc,
) *BytesFuncMap {
	if !isSet {
		return it
	}

	it.Set(name, executor)

	return it
}

func (it *BytesFuncMap) AddOrUpdateChain(
	name string,
	executor BytesExecutorFunc,
) *BytesFuncMap {
	it.Set(name, executor)

	return it
}

func (it BytesFuncMap) ListStrings() []string {
	return it.AllNames()
}

func (it BytesFuncMap) AsStandardSlicerContractsBinder() coreinterface.StandardSlicerContractsBinder {
	return &it
}

func (it *BytesFuncMap) Length() int {
	if it == nil {
		return 0
	}

	return len(it.FunctionsMap)
}

func (it BytesFuncMap) IsEmpty() bool {
	return it.Length() == 0
}

func (it BytesFuncMap) HasAnyItem() bool {
	return it.Length() > 0
}

func (it BytesFuncMap) AllKeys() []string {
	return it.AllNames()
}

func (it BytesFuncMap) SortedNamesCsv() string {
	return corecsv.DefaultCsv(it.AllNamesSorted()...)
}

func (it BytesFuncMap) AllKeysSorted() []string {
	return it.AllNamesSorted()
}

func (it BytesFuncMap) AllNames() []string {
	names := make([]string, 0, it.Length())

	for s := range it.FunctionsMap {
		names = append(names, s)
	}

	return names
}

func (it BytesFuncMap) AllNamesSorted() []string {
	names := it.AllNames()
	sort.Strings(names)

	return names
}

func (it BytesFuncMap) Get(name string) BytesExecutorFunc {
	return it.FunctionsMap[name]
}

// GetFuncPayloadWrapperCategoryByRawPayloads
//
//  deserializes payload to payload wrapper and
//  then usages' payload category name to
//  get to the executor
func (it BytesFuncMap) GetFuncPayloadWrapperCategoryByRawPayloads(
	payloads []byte,
) (executorWrapper *BytesFuncExecutorWrapper) {
	return it.GetFuncPayloadWrapperByRawPayloads(
		payloads,
		payloadCategoryNameToExecutor)
}

// GetFuncPayloadWrapperNameByRawPayloads
//
//  deserializes payload to payload wrapper and
//  then usages' payload name to
//  get to the executor
func (it BytesFuncMap) GetFuncPayloadWrapperNameByRawPayloads(
	payloads []byte,
) (executorWrapper *BytesFuncExecutorWrapper) {
	return it.GetFuncPayloadWrapperByRawPayloads(
		payloads,
		payloadNameToExecutor)
}

// GetFuncPayloadWrapperByRawPayloads
//
//  deserializes payload to payload wrapper and
//  then usages' the func executorNameGetter to get the executor name
//  using payloadWrapper then get to the executor and
//  returns the final wrapper
func (it BytesFuncMap) GetFuncPayloadWrapperByRawPayloads(
	payloads []byte,
	executorNameGetter PayloadWrapperToFuncNameGetterFunc,
) (executorWrapper *BytesFuncExecutorWrapper) {
	payloadWrapper, err := corepayload.New.PayloadWrapper.Deserialize(
		payloads)

	if err != nil {
		invalidBaseExecutorWrapper := InvalidBaseExecutorWrapper(
			it.BaseExecutorInfo.ToPtr(),
			payloads,
			err)

		return &BytesFuncExecutorWrapper{
			BaseExecutorWrapper: invalidBaseExecutorWrapper.ToNonPtr(),
		}
	}

	executorName := executorNameGetter(payloadWrapper)
	executor, notFoundErrWp := it.GetFuncOrErrorWrapper(
		executorName)
	validBaseExecutorWrapper := NewBaseExecutorWrapper(
		it.BaseExecutorInfo.ToPtr(),
		payloads,
		notFoundErrWp,
		payloadWrapper)

	return &BytesFuncExecutorWrapper{
		BaseExecutorWrapper: validBaseExecutorWrapper.ToNonPtr(),
		Executor:            executor,
	}
}

func (it BytesFuncMap) GetWithStat(name string) (BytesExecutorFunc, bool) {
	currFunc, has := it.FunctionsMap[name]

	return currFunc, has
}

func (it BytesFuncMap) GetWithStatByNamer(name enuminf.Namer) (BytesExecutorFunc, bool) {
	currFunc, has := it.FunctionsMap[name.Name()]

	return currFunc, has
}

func (it BytesFuncMap) HasFunc(name string) bool {
	_, has := it.FunctionsMap[name]

	return has
}

func (it BytesFuncMap) IsMissingFunc(name string) bool {
	_, has := it.FunctionsMap[name]

	return !has
}

func (it BytesFuncMap) ExecAllLock(
	samePayloads []byte,
) *errwrappers.Collection {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	return it.ExecAll(samePayloads)
}

func (it BytesFuncMap) ExecAll(
	samePayloads []byte,
) *errwrappers.Collection {
	if it.IsEmpty() {
		return errwrappers.Empty()
	}

	if len(it.OrderedNames) > 0 {
		return it.executeByOrderedNamesAll(samePayloads)
	}

	info := it.BaseExecutorInfo.SafeInfo()
	emptyCollector := errwrappers.Empty()

	for _, errFunc := range it.FunctionsMap {
		// found, exec
		errWrap := errFunc(samePayloads)
		wrappedError := wrapErrorWithDetailsInstance.
			wrapErrorWrapper(
				errWrap,
				info,
			)

		emptyCollector.AddWrapperPtr(wrappedError)
	}

	return emptyCollector
}

func (it BytesFuncMap) executeByOrderedNamesAll(
	samePayloads []byte,
) *errwrappers.Collection {
	emptyCollector := errwrappers.Empty()

	for _, name := range it.OrderedNames {
		wrappedErr := it.execInternal(
			name,
			samePayloads)
		emptyCollector.AddWrapperPtr(wrappedErr)
	}

	return emptyCollector
}

func (it BytesFuncMap) ExecLock(
	name string, payloads []byte,
) *errorwrapper.Wrapper {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	return it.Exec(name, payloads)
}

func (it BytesFuncMap) Exec(
	name string,
	payloads []byte,
) *errorwrapper.Wrapper {
	return it.execInternal(
		name,
		payloads)
}

func (it BytesFuncMap) ExecNamer(
	namer enuminf.Namer,
	payloads []byte,
) *errorwrapper.Wrapper {
	return it.execInternal(
		namer.Name(),
		payloads)
}

func (it BytesFuncMap) ExecNamerLock(
	namer enuminf.Namer,
	payloads []byte,
) *errorwrapper.Wrapper {
	return it.ExecLock(
		namer.Name(),
		payloads)
}

func (it BytesFuncMap) NotFoundError(
	name string,
) *errorwrapper.Wrapper {
	_, hasFunc := it.FunctionsMap[name]

	if hasFunc {
		return nil
	}

	// not found
	return wrapErrorWithDetailsInstance.notFoundError(
		name,
		it.SafeInfo(),
	)
}

func (it BytesFuncMap) GetFuncOrErrorWrapper(
	name string,
) (
	executor BytesExecutorFunc,
	notFoundErr *errorwrapper.Wrapper,
) {
	currentFunc, hasFunc := it.FunctionsMap[name]

	if hasFunc {
		return currentFunc, nil
	}

	// not found
	return nil, wrapErrorWithDetailsInstance.notFoundError(
		name,
		it.SafeInfo())
}

func (it BytesFuncMap) execInternal(
	name string,
	payloads []byte,
) *errorwrapper.Wrapper {
	errFunc, hasFunc := it.FunctionsMap[name]

	if !hasFunc {
		return wrapErrorWithDetailsInstance.notFoundError(
			name,
			it.SafeInfo(),
		)
	}

	// found, executor
	errorWrapper := errFunc(payloads)
	if errorWrapper.IsEmpty() {
		return nil
	}

	return wrapErrorWithDetailsInstance.wrapErrorWrapperPayloads(
		errorWrapper,
		it.SafeInfo(),
		payloads)
}

func (it BytesFuncMap) ExecByNamesLock(
	samePayloads []byte,
	names ...string,
) *errwrappers.Collection {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	return it.ExecByNames(
		samePayloads,
		names...)
}

func (it BytesFuncMap) ExecByNames(
	samePayloads []byte,
	names ...string,
) *errwrappers.Collection {
	emptyCollector := errwrappers.Empty()

	if len(names) == 0 {
		return emptyCollector
	}

	for _, name := range names {
		wrappedErr := it.execInternal(
			name,
			samePayloads)

		emptyCollector.AddWrapperPtr(
			wrappedErr)
	}

	return emptyCollector
}

func (it BytesFuncMap) ExecByNamesIf(
	isExecute bool,
	samePayloads []byte,
	names ...string,
) *errwrappers.Collection {
	if !isExecute {
		return errwrappers.Empty()
	}

	return it.ExecByNames(
		samePayloads,
		names...)
}

func (it BytesFuncMap) ExecByNamesLockIf(
	isExecute bool,
	samePayloads []byte,
	names ...string,
) *errwrappers.Collection {
	if !isExecute {
		return errwrappers.Empty()
	}

	return it.ExecByNamesLock(
		samePayloads,
		names...)
}

func (it *BytesFuncMap) ReflectSetTo(toPtr interface{}) error {
	return coredynamic.ReflectSetFromTo(it, toPtr)
}

func (it *BytesFuncMap) ReflectSetToErrWrap(toPtr interface{}) *errorwrapper.Wrapper {
	return errnew.Reflect.SetFromTo(it, toPtr)
}

func (it BytesFuncMap) ClonePtr() *BytesFuncMap {
	cloned := it.Clone()

	return &cloned
}

func (it BytesFuncMap) Clone() BytesFuncMap {
	if it.IsEmpty() {
		return BytesFuncMap{
			BaseExecutorInfo: it.BaseExecutorInfo.Clone(),
			FunctionsMap:     map[string]BytesExecutorFunc{},
		}
	}

	newMap := make(
		map[string]BytesExecutorFunc,
		it.Length())
	for name, executor := range it.FunctionsMap {
		newMap[name] = executor
	}

	return BytesFuncMap{
		BaseExecutorInfo: it.BaseExecutorInfo.Clone(),
		FunctionsMap:     newMap,
	}
}
