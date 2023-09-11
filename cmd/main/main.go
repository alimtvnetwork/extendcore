package main

import (
	"fmt"

	"gitlab.com/auk-go/core/coredata/corepayload"
	"gitlab.com/auk-go/core/coretaskinfo"
	"gitlab.com/auk-go/core/errcore"
	"gitlab.com/auk-go/errorwrapper"
	"gitlab.com/auk-go/errorwrapper/errfunc"
	"gitlab.com/auk-go/errorwrapper/errnew"
	"gitlab.com/auk-go/extendcore/mapds/mappedfunc"
)

func main() {
	poayloadExecutorMap := mappedfunc.PayloadFuncMap{
		BaseExecutorInfo: mappedfunc.BaseExecutorInfo{
			Info: coretaskinfo.Info{
				RootName:    "some root name",
				Description: "desc",
				Url:         "url",
				HintUrl:     "hint",
				ErrorUrl:    "error url",
				ExcludeOptions: &coretaskinfo.ExcludingOptions{
					IsExcludeUrl:                 false,
					IsExcludeHintUrl:             false,
					IsExcludeRootName:            false,
					IsExcludeErrorUrl:            false,
					IsExcludeAdditionalErrorWrap: false,
					IsSecureText:                 false,
					IsExcludeDescription:         false,
				},
			},
			OrderedNames: nil,
		},
		AnyValidatorFunctions: []errfunc.AnyItemValidatorFunc{
			func(anyInput interface{}) *errorwrapper.Wrapper {
				fmt.Println("any validator called")

				return nil
			},
		},
		ValidatorFunctions: []errfunc.PayloadValidatorFunc{
			func(payload *corepayload.PayloadWrapper) *errorwrapper.Wrapper {
				print("validator called")

				return errnew.NotFound.KeyMessage("payload key issues")
			},
		},
		FunctionsMap: map[string]errfunc.PayloadWrapperExecutorFunc{
			"task1": func(payloadWrapper *corepayload.PayloadWrapper) *errorwrapper.Wrapper {
				fmt.Println("task 1")
				fmt.Println(payloadWrapper.PrettyJsonString())
				fmt.Println(payloadWrapper.PayloadEntityType())

				return nil
			},
			"task2": func(payloadWrapper *corepayload.PayloadWrapper) *errorwrapper.Wrapper {
				fmt.Println("task 2")

				return nil
			},
			"task3": func(payloadWrapper *corepayload.PayloadWrapper) *errorwrapper.Wrapper {
				fmt.Println("task 3")

				return nil
			},
		},
	}

	payload, err := corepayload.New.PayloadWrapper.NameIdRecord(
		"payload name",
		"payload id 1",
		"some record")

	errcore.HandleErr(err)

	poayloadExecutorMap.Exec("task1", payload)
	errWp := poayloadExecutorMap.Exec("task2", nil)

	errWp.Log()

	validationErr := poayloadExecutorMap.ValidateExec("task3", payload)
	validationErr.Log()
}
