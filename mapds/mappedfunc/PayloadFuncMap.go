package mappedfunc

import (
	"sort"

	"gitlab.com/evatix-go/core/codestack"
	"gitlab.com/evatix-go/core/corecsv"
	"gitlab.com/evatix-go/core/coredata/coredynamic"
	"gitlab.com/evatix-go/core/coredata/corejson"
	"gitlab.com/evatix-go/core/coredata/corepayload"
	"gitlab.com/evatix-go/core/coreinterface"
	"gitlab.com/evatix-go/core/coreinterface/enuminf"
	"gitlab.com/evatix-go/errorwrapper"
	"gitlab.com/evatix-go/errorwrapper/errnew"
	"gitlab.com/evatix-go/errorwrapper/errtype"
	"gitlab.com/evatix-go/errorwrapper/errwrappers"
	"gitlab.com/evatix-go/errorwrapper/ref"
)

type PayloadFuncMap struct {
	BaseExecutorInfo
	IsSkipPreCheckNull    bool
	Middlewares           map[string]PayloadWrapperMiddleware
	AnyValidatorFunctions []AnyItemValidatorFunc
	ValidatorFunctions    []PayloadValidatorFunc
	FunctionsMap          map[string]PayloadWrapperExecutorFunc
}

func (it *PayloadFuncMap) AnyValidate(
	anyInput interface{},
) *errorwrapper.Wrapper {
	if len(it.AnyValidatorFunctions) == 0 {
		return nil
	}

	for _, anyValidatorFunc := range it.AnyValidatorFunctions {
		errWp := anyValidatorFunc(anyInput)

		if errWp.HasError() {
			return it.WrapErrorWithDetailsPlusPayload(
				errWp,
				anyInput)
		}
	}

	return nil
}

func (it *PayloadFuncMap) AnyValidateIf(
	isValidate bool,
	anyInput interface{},
) *errorwrapper.Wrapper {
	if !isValidate {
		return nil
	}

	return it.AnyValidate(anyInput)
}

func (it *PayloadFuncMap) Validate(
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	if len(it.ValidatorFunctions) == 0 {
		return nil
	}

	for _, validatorFunc := range it.ValidatorFunctions {
		errWp := validatorFunc(payloadWrapper)

		if errWp.HasError() {
			return it.WrapErrorWithDetailsPlusPayload(
				errWp,
				payloadWrapper)
		}
	}

	return nil
}

func (it *PayloadFuncMap) ValidateIf(
	isValidate bool,
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	if !isValidate {
		return nil
	}

	return it.Validate(payloadWrapper)
}

func (it *PayloadFuncMap) Append(
	isSkipOnExist bool,
	executorsMap *PayloadFuncMap,
) *PayloadFuncMap {
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

func (it *PayloadFuncMap) ConcatNew(
	isSkipOnExist bool,
	executorsMap *PayloadFuncMap,
) *PayloadFuncMap {
	if executorsMap.IsEmpty() {
		return it.ClonePtr()
	}

	cloned := it.ClonePtr()

	return cloned.Append(isSkipOnExist, executorsMap)
}

func (it PayloadFuncMap) Items() map[string]PayloadWrapperExecutorFunc {
	return it.FunctionsMap
}

func (it PayloadFuncMap) Count() int {
	return it.Length()
}

func (it PayloadFuncMap) LastIndex() int {
	return it.Length() - 1
}

func (it PayloadFuncMap) HasIndex(index int) bool {
	return it.LastIndex() >= index
}

func (it *PayloadFuncMap) MiddlewaresLength() int {
	if it == nil || it.Middlewares == nil {
		return 0
	}

	return len(it.Middlewares)
}

func (it *PayloadFuncMap) HasAnyMiddlewares() bool {
	return it.MiddlewaresLength() > 0
}

func (it *PayloadFuncMap) IsEmptyMiddlewares() bool {
	return it.MiddlewaresLength() == 0
}

func (it *PayloadFuncMap) ExecAllMiddlewares(
	isContinueOnError bool,
	input interface{},
) *errwrappers.Collection {
	if it.IsEmptyMiddlewares() {
		return nil
	}
	errCollection := errwrappers.Empty()
	var previousOutput interface{}
	var errWrap *errorwrapper.Wrapper

	index := -1
	for middlewareName, middleWareProcessor := range it.Middlewares {
		index++

		if index == 0 {
			previousOutput, errWrap = middleWareProcessor(
				middlewareName,
				it,
				input)
		} else {
			previousOutput, errWrap = middleWareProcessor(
				middlewareName,
				it,
				previousOutput)
			// todo fix mutating
		}

		if errWrap.HasError() && index == 0 {
			errCollection.AddWrapperPtr(
				it.failedToExecuteMiddlewareErrorWrap(
					middlewareName,
					errWrap,
					input))
		} else if errWrap.HasError() && index == 0 {
			errCollection.AddWrapperPtr(
				it.failedToExecuteMiddlewareErrorWrap(
					middlewareName,
					errWrap,
					previousOutput))
		}

		if !isContinueOnError && errCollection.HasError() {
			return errCollection
		}
	}

	return errCollection
}

func (it PayloadFuncMap) Set(
	name string,
	executor PayloadWrapperExecutorFunc,
) (isAddedNewly bool) {
	isAddedNewly = it.IsMissingFunc(name)
	it.FunctionsMap[name] = executor

	return isAddedNewly
}

func (it PayloadFuncMap) SetIf(
	isAdd bool,
	name string,
	executor PayloadWrapperExecutorFunc,
) (isAddedNewly bool) {
	if !isAdd {
		return it.IsMissingFunc(name)
	}

	return it.Set(name, executor)
}

func (it PayloadFuncMap) AddOrUpdate(
	name string,
	executor PayloadWrapperExecutorFunc,
) (isAddedNewly bool) {
	return it.Set(name, executor)
}

func (it *PayloadFuncMap) SetChain(
	name string,
	executor PayloadWrapperExecutorFunc,
) *PayloadFuncMap {
	it.Set(name, executor)

	return it
}

func (it *PayloadFuncMap) SetChainIf(
	isSet bool,
	name string,
	executor PayloadWrapperExecutorFunc,
) *PayloadFuncMap {
	if !isSet {
		return it
	}

	it.Set(name, executor)

	return it
}

func (it *PayloadFuncMap) AddOrUpdateChain(
	name string,
	executor PayloadWrapperExecutorFunc,
) *PayloadFuncMap {
	it.Set(name, executor)

	return it
}

func (it PayloadFuncMap) ListStrings() []string {
	return it.AllNames()
}

func (it PayloadFuncMap) AsStandardSlicerContractsBinder() coreinterface.StandardSlicerContractsBinder {
	return &it
}

func (it *PayloadFuncMap) Length() int {
	if it == nil {
		return 0
	}

	return len(it.FunctionsMap)
}

func (it PayloadFuncMap) IsEmpty() bool {
	return it.Length() == 0
}

func (it PayloadFuncMap) HasAnyItem() bool {
	return it.Length() > 0
}

func (it PayloadFuncMap) AllKeys() []string {
	return it.AllNames()
}

func (it PayloadFuncMap) SortedNamesCsv() string {
	return corecsv.DefaultCsv(it.AllNamesSorted()...)
}

func (it PayloadFuncMap) AllKeysSorted() []string {
	return it.AllNamesSorted()
}

func (it PayloadFuncMap) AllNames() []string {
	names := make([]string, 0, it.Length())

	for s := range it.FunctionsMap {
		names = append(names, s)
	}

	return names
}

func (it PayloadFuncMap) AllNamesSorted() []string {
	names := it.AllNames()
	sort.Strings(names)

	return names
}

func (it PayloadFuncMap) Get(name string) PayloadWrapperExecutorFunc {
	return it.FunctionsMap[name]
}

// GetFuncPayloadWrapperCategoryByRawPayloads
//
//  deserializes payload to payload wrapper and
//  then usages' payload category name to
//  get to the executor
func (it PayloadFuncMap) GetFuncPayloadWrapperCategoryByRawPayloads(
	payloads []byte,
) (executorWrapper *PayloadFuncExecutorWrapper) {
	return it.GetFuncPayloadWrapperByRawPayloads(
		payloads,
		payloadCategoryNameToExecutor)
}

// GetFuncPayloadWrapperNameByRawPayloads
//
//  deserializes payload to payload wrapper and
//  then usages' payload name to
//  get to the executor
func (it PayloadFuncMap) GetFuncPayloadWrapperNameByRawPayloads(
	payloads []byte,
) (executorWrapper *PayloadFuncExecutorWrapper) {
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
func (it PayloadFuncMap) GetFuncPayloadWrapperByRawPayloads(
	payloads []byte,
	executorNameGetter PayloadWrapperToFuncNameGetterFunc,
) (executorWrapper *PayloadFuncExecutorWrapper) {
	payloadWrapper, err := corepayload.New.PayloadWrapper.Deserialize(
		payloads)

	if err != nil {
		invalidBaseExecutorWrapper := InvalidBaseExecutorWrapper(
			it.BaseExecutorInfo.ToPtr(),
			payloads,
			err)

		return &PayloadFuncExecutorWrapper{
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

	return &PayloadFuncExecutorWrapper{
		BaseExecutorWrapper: validBaseExecutorWrapper.ToNonPtr(),
		Executor:            executor,
	}
}

func (it PayloadFuncMap) GetWithStat(name string) (PayloadWrapperExecutorFunc, bool) {
	currFunc, has := it.FunctionsMap[name]

	return currFunc, has
}

func (it PayloadFuncMap) HasFunc(name string) bool {
	_, has := it.FunctionsMap[name]

	return has
}

func (it PayloadFuncMap) IsMissingFunc(name string) bool {
	_, has := it.FunctionsMap[name]

	return !has
}

func (it PayloadFuncMap) AnyValidateExec(
	name string,
	validationInput interface{},
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	errWp := it.AnyValidate(validationInput)

	if errWp.HasError() {
		return errWp
	}

	return it.Exec(name, payloadWrapper)
}

func (it PayloadFuncMap) ValidateExec(
	name string,
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	errWp := it.Validate(payloadWrapper)

	if errWp.HasError() {
		return errWp
	}

	return it.Exec(name, payloadWrapper)
}

func (it PayloadFuncMap) AnyValidateExecIf(
	isValidate, isExecute bool,
	name string,
	validationInput interface{},
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	errWp := it.AnyValidateIf(
		isValidate,
		validationInput)

	if errWp.HasError() {
		return errWp
	}

	return it.ExecIf(
		isExecute,
		name,
		payloadWrapper)
}

func (it PayloadFuncMap) ValidateExecIf(
	isValidate, isExecute bool,
	name string,
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	errWp := it.ValidateIf(
		isValidate,
		payloadWrapper)

	if errWp.HasError() {
		return errWp
	}

	return it.ExecIf(
		isExecute,
		name,
		payloadWrapper)
}

func (it PayloadFuncMap) ExecAllLock(
	payloadWrapper *corepayload.PayloadWrapper,
) *errwrappers.Collection {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	return it.ExecAll(payloadWrapper)
}

func (it PayloadFuncMap) ExecAll(
	payloadWrapper *corepayload.PayloadWrapper,
) *errwrappers.Collection {
	if it.IsEmpty() {
		return errwrappers.Empty()
	}

	if len(it.OrderedNames) > 0 {
		return it.executeByOrderedNamesAll(payloadWrapper)
	}

	payloadBytes, _ := payloadWrapper.Serialize()
	emptyCollector := errwrappers.Empty()

	for _, errFunc := range it.FunctionsMap {
		// found, exec
		errorWrapper := errFunc(payloadWrapper)
		wrappedError := wrapErrorWithDetailsInstance.
			wrapErrorWrapperPayloads(
				errorWrapper,
				it.SafeInfo(),
				payloadBytes)

		emptyCollector.AddWrapperPtr(wrappedError)
	}

	return emptyCollector
}

func (it PayloadFuncMap) executeByOrderedNamesAll(
	payloadWrapper *corepayload.PayloadWrapper,
) *errwrappers.Collection {
	emptyCollector := errwrappers.Empty()

	for _, name := range it.OrderedNames {
		wrappedErr := it.execInternal(
			name,
			payloadWrapper)
		emptyCollector.AddWrapperPtr(wrappedErr)
	}

	return emptyCollector
}

func (it PayloadFuncMap) ExecLock(
	name string,
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	return it.Exec(name, payloadWrapper)
}

func (it PayloadFuncMap) ExecIf(
	isExecute bool,
	name string,
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	if !isExecute {
		return nil
	}

	return it.Exec(
		name,
		payloadWrapper)
}

func (it PayloadFuncMap) Exec(
	name string,
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	return it.execInternal(
		name,
		payloadWrapper)
}

func (it PayloadFuncMap) ExecUsingAny(
	name string,
	payloadWrapperAsAny interface{},
) *errorwrapper.Wrapper {
	payloadWrapper, parsingErr := Converter.
		AnyToPayloadWrapper(payloadWrapperAsAny)

	if parsingErr.HasError() {
		return it.failedToExecuteExecutorErrorWrap(
			name,
			parsingErr,
			payloadWrapperAsAny)
	}

	return it.Exec(name, payloadWrapper)
}

func (it PayloadFuncMap) ExecUsingBytes(
	name string,
	payloadWrapperAsBytes []byte,
) *errorwrapper.Wrapper {
	payloadWrapper, parsingErr := corepayload.
		New.
		PayloadWrapper.
		Deserialize(payloadWrapperAsBytes)

	if parsingErr != nil {
		return it.failedToExecuteExecutorErrorWrap(
			name,
			errnew.Error.Default(errtype.ParsingFailed, parsingErr),
			payloadWrapperAsBytes)
	}

	return it.Exec(name, payloadWrapper)
}

func (it PayloadFuncMap) ExecNamer(
	namer enuminf.Namer,
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	return it.execInternal(
		namer.Name(),
		payloadWrapper)
}

func (it PayloadFuncMap) ExecNamerIf(
	isExecute bool,
	namer enuminf.Namer,
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	if !isExecute {
		return nil
	}

	return it.execInternal(
		namer.Name(),
		payloadWrapper)
}

func (it PayloadFuncMap) ExecNamerLock(
	namer enuminf.Namer,
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	return it.ExecLock(
		namer.Name(),
		payloadWrapper)
}

func (it PayloadFuncMap) NotFoundError(
	name string,
) *errorwrapper.Wrapper {
	_, hasFunc := it.FunctionsMap[name]

	if hasFunc {
		return nil
	}

	// not found
	return wrapErrorWithDetailsInstance.notFoundError(
		name,
		it.SafeInfo())
}

func (it PayloadFuncMap) GetFuncOrErrorWrapper(
	name string,
) (
	executor PayloadWrapperExecutorFunc,
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

func (it PayloadFuncMap) IsPreCheckNull() bool {
	return !it.IsSkipPreCheckNull
}

func (it PayloadFuncMap) execInternal(
	name string,
	payloadWrapper *corepayload.PayloadWrapper,
) *errorwrapper.Wrapper {
	if it.IsPreCheckNull() && payloadWrapper == nil {
		return it.nullErrorWrapper(
			name,
			payloadWrapper)
	}

	errFunc, hasFunc := it.FunctionsMap[name]

	if !hasFunc {
		return wrapErrorWithDetailsInstance.notFoundErrorWithPayloads(
			name,
			it.SafeInfo(),
			payloadWrapper)
	}

	// found, executor
	errorWrapper := errFunc(payloadWrapper)
	if errorWrapper.IsEmpty() {
		return nil
	}

	return wrapErrorWithDetailsInstance.wrapErrorWrapperPayloadsAny(
		errorWrapper,
		it.SafeInfo(),
		payloadWrapper)
}

func (it PayloadFuncMap) nullErrorWrapper(
	executorName string,
	anyItem interface{},
) *errorwrapper.Wrapper {
	return wrapErrorWithDetailsInstance.wrapErrorWrapper(
		errnew.Null.WithRefs(
			"null received for "+executorName+" - executor",
			anyItem,
			ref.Value{
				Variable: "executor-name",
				Value:    executorName,
			}),
		it.SafeInfo(),
	)
}

func (it PayloadFuncMap) failedToExecuteExecutorErrorWrap(
	executorName string,
	existingErr *errorwrapper.Wrapper,
	payloadAny interface{},
) *errorwrapper.Wrapper {
	if existingErr.IsEmpty() {
		return nil
	}

	wrappedErr := existingErr.ConcatNew().MsgRefOne(
		codestack.Skip1,
		executorName+" executor couldn't execute.",
		"executor-name",
		executorName)

	return wrapErrorWithDetailsInstance.wrapErrorWrapperPayloads(
		wrappedErr,
		it.SafeInfo(),
		corejson.AnyTo.SerializedJsonResult(payloadAny).SafeBytes(),
	)
}

func (it PayloadFuncMap) failedToExecuteMiddlewareErrorWrap(
	middleware string,
	existingErr *errorwrapper.Wrapper,
	payloadAny interface{},
) *errorwrapper.Wrapper {
	if existingErr.IsEmpty() {
		return nil
	}

	wrappedErr := existingErr.ConcatNew().MsgRefOne(
		codestack.Skip1,
		middleware+" middleware failed execute.",
		"middleware-name",
		middleware)

	return wrapErrorWithDetailsInstance.wrapErrorWrapperPayloads(
		wrappedErr,
		it.SafeInfo(),
		corejson.AnyTo.SerializedJsonResult(payloadAny).SafeBytes(),
	)
}

func (it PayloadFuncMap) ExecByNamesLock(
	payloadWrapper *corepayload.PayloadWrapper,
	names ...string,
) *errwrappers.Collection {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	return it.ExecByNames(
		payloadWrapper,
		names...)
}

func (it PayloadFuncMap) ExecByNames(
	payloadWrapper *corepayload.PayloadWrapper,
	names ...string,
) *errwrappers.Collection {
	emptyCollector := errwrappers.Empty()

	if len(names) == 0 {
		return emptyCollector
	}

	for _, name := range names {
		wrappedErr := it.execInternal(
			name,
			payloadWrapper)

		emptyCollector.AddWrapperPtr(wrappedErr)
	}

	return emptyCollector
}

func (it PayloadFuncMap) ExecByNamesChecking(
	payloadWrapper *corepayload.PayloadWrapper,
	isExecuteFunc func(payloadWrapper *corepayload.PayloadWrapper, name string) (isExecute bool),
	names ...string,
) *errwrappers.Collection {
	emptyCollector := errwrappers.Empty()

	if len(names) == 0 {
		return emptyCollector
	}

	for _, name := range names {
		isExec := isExecuteFunc(
			payloadWrapper,
			name)

		if !isExec {
			continue
		}

		wrappedErr := it.execInternal(
			name,
			payloadWrapper)
		emptyCollector.AddWrapperPtr(wrappedErr)
	}

	return emptyCollector
}

func (it PayloadFuncMap) ExecByNamesIf(
	isExecute bool,
	payloadWrapper *corepayload.PayloadWrapper,
	names ...string,
) *errwrappers.Collection {
	if !isExecute {
		return errwrappers.Empty()
	}

	return it.ExecByNames(
		payloadWrapper,
		names...)
}

func (it PayloadFuncMap) ExecByNamesLockIf(
	isExecute bool,
	payloadWrapper *corepayload.PayloadWrapper,
	names ...string,
) *errwrappers.Collection {
	if !isExecute {
		return errwrappers.Empty()
	}

	return it.ExecByNamesLock(
		payloadWrapper,
		names...)
}

func (it *PayloadFuncMap) ReflectSetTo(toPtr interface{}) error {
	return coredynamic.ReflectSetFromTo(it, toPtr)
}

func (it *PayloadFuncMap) ReflectSetToErrWrap(toPtr interface{}) *errorwrapper.Wrapper {
	return errnew.Reflect.SetFromTo(it, toPtr)
}

func (it PayloadFuncMap) ClonePtr() *PayloadFuncMap {
	cloned := it.Clone()

	return &cloned
}

func (it PayloadFuncMap) Clone() PayloadFuncMap {
	if it.IsEmpty() {
		return PayloadFuncMap{
			BaseExecutorInfo: it.BaseExecutorInfo.Clone(),
			FunctionsMap:     map[string]PayloadWrapperExecutorFunc{},
		}
	}

	newMap := make(
		map[string]PayloadWrapperExecutorFunc,
		it.Length())
	for name, executor := range it.FunctionsMap {
		newMap[name] = executor
	}

	return PayloadFuncMap{
		BaseExecutorInfo: it.BaseExecutorInfo.Clone(),
		FunctionsMap:     newMap,
	}
}
