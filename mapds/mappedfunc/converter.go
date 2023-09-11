package mappedfunc

import (
	"gitlab.com/auk-go/core/coredata/corejson"
	"gitlab.com/auk-go/core/coredata/corepayload"
	"gitlab.com/auk-go/core/coreinterface"
	"gitlab.com/auk-go/core/errcore"
	"gitlab.com/auk-go/core/isany"
	"gitlab.com/auk-go/errorwrapper"
	"gitlab.com/auk-go/errorwrapper/errdata/errbyte"
	"gitlab.com/auk-go/errorwrapper/errdata/errstr"
	"gitlab.com/auk-go/errorwrapper/errnew"
	"gitlab.com/auk-go/errorwrapper/errtype"
)

type converter struct{}

func (it converter) AnyToPayloadWrapper(
	payloadWrapperAsAny interface{},
) (payloadWrapper *corepayload.PayloadWrapper, parsingErr *errorwrapper.Wrapper) {
	if isany.Null(payloadWrapperAsAny) {
		return nil, errnew.Null.Simple(payloadWrapperAsAny)
	}

	switch casted := payloadWrapperAsAny.(type) {
	case corepayload.PayloadWrapper:
		return &casted, nil
	case *corepayload.PayloadWrapper:
		return casted, nil
	case corejson.Result:
		payload, err := corepayload.
			New.
			PayloadWrapper.
			DeserializeUsingJsonResult(
				casted.Ptr())

		return payload, errnew.Unmarshal.Error(err)
	case *corejson.Result:
		payload, err := corepayload.
			New.
			PayloadWrapper.
			DeserializeUsingJsonResult(
				casted)

		return payload, errnew.Unmarshal.Error(err)
	case []byte:
		payload, err := corepayload.
			New.
			PayloadWrapper.
			Deserialize(
				casted)

		return payload, errnew.Unmarshal.Error(err)
	}

	return nil, errnew.SrcDst.Create(
		errtype.CastingFailed,
		payloadWrapperAsAny,
		&corepayload.PayloadWrapper{},
	)
}

func (it converter) AnyToSerializedString(
	anyItem interface{},
) *errstr.Result {
	result := it.AnyToBytesOrSerializedBytes(anyItem)

	return result.ErrStr()
}

func (it converter) AnyToBytesOrSerializedBytes(
	anyItem interface{},
) *errbyte.Results {
	if isany.Null(anyItem) {
		return errbyte.Empty.ResultsWithError(
			errnew.Null.Simple(anyItem))
	}

	switch casted := anyItem.(type) {
	case []byte:
		return errbyte.New.Results.ValuesOnly(
			casted)
	case string:
		return errbyte.New.Results.ValuesOnly(
			[]byte(casted))
	case corejson.Jsoner:
		jsonResult := casted.Json()
		if jsonResult.HasError() {
			return errbyte.New.Results.Error(
				errtype.Unmarshalling,
				jsonResult.MeaningfulError())
		}

		return errbyte.New.Results.ValuesOnly(
			jsonResult.SafeBytes())
	case corejson.Result:
		if casted.HasError() {
			return errbyte.New.Results.Error(
				errtype.Unmarshalling,
				casted.MeaningfulError())
		}

		return errbyte.New.Results.ValuesOnly(
			casted.SafeBytes())
	case *corejson.Result:
		if casted.HasError() {
			return errbyte.New.Results.Error(
				errtype.Unmarshalling,
				casted.MeaningfulError())
		}

		return errbyte.New.Results.ValuesOnly(
			casted.SafeBytes())
	case coreinterface.Serializer:
		allBytes, err := casted.Serialize()

		if err != nil {
			return errbyte.New.Results.Error(
				errtype.Unmarshalling,
				err)
		}

		return errbyte.New.Results.ValuesOnly(
			allBytes)
	case error:
		return errbyte.New.Results.ValuesOnly(
			[]byte(errcore.ToString(casted)))
	}

	serializeJsonResult := corejson.Serialize.Apply(
		anyItem)

	if serializeJsonResult.IsEmptyError() {
		return errbyte.New.Results.ValuesOnly(
			serializeJsonResult.SafeBytes())
	}

	return errbyte.Empty.ResultsWithError(errnew.SrcDst.Create(
		errtype.CastingFailed,
		anyItem,
		[]byte{},
	))
}
