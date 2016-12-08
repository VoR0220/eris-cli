package abi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/eris-ltd/eris-cli/interpret"
	//"github.com/eris-ltd/eris-cli/log"
)

/*func ReadAbiFormulateCall(abiLocation string, funcName string, args []string, do *definitions.Do) ([]byte, error) {
	abiSpecBytes, err := util.ReadAbi(do.ABIPath, abiLocation)
	if err != nil {
		return []byte{}, err
	}
	log.WithField("=>", string(abiSpecBytes)).Debug("ABI Specification (Formulate)")
	log.WithFields(log.Fields{
		"function":  funcName,
		"arguments": fmt.Sprintf("%v", args),
	}).Debug("Packing Call via ABI")

	return Packer(abiSpecBytes, funcName, args...)
}

func ReadAndDecodeContractReturn(abiLocation, funcName string, resultRaw []byte, do *definitions.Do) ([]*definitions.Variable, error) {
	abiSpecBytes, err := util.ReadAbi(do.ABIPath, abiLocation)
	if err != nil {
		return nil, err
	}
	log.WithField("=>", abiSpecBytes).Debug("ABI Specification (Decode)")

	// Unpack the result
	return Unpacker(abiSpecBytes, funcName, resultRaw)
}*/

func MakeAbi(abiData string) (ABI, error) {
	if len(abiData) == 0 {
		return ABI{}, nil
	}

	abiSpec, err := JSON(strings.NewReader(abiData))
	if err != nil {
		return ABI{}, err
	}

	return abiSpec, nil
}

func ConvertSlice(from []interface{}, to Type) (interface{}, err error){
	if !to.IsSlice {
		return nil, fmt.Errorf("Attempting to convert to non slice type")
	} else if to.SliceSize != -1 && len(from) != to.SliceSize {
		return nil, fmt.Errorf("Length of array does not match, expected %v got %v", to.SliceSize, len(from))
	}
	for i, typ := range from {
		from[i], err = ConvertToPackingType(typ, *to.Elem)
		if err != nil {
			return nil, err
		}
	}
	return from, nil
}

func ConvertToPackingType(from interface{}, to Type) (interface{}, error) {
	if to.IsSlice || to.IsArray && to.T != BytesTy && to.T != FixedBytesTy {
		if typ, ok := from.([]interface{}); !ok {
			return nil, fmt.Errorf("Unexpected non slice type during type conversion, please reformat your run file to use an array/slice.")
		} else {
			return ConvertSlice(typ, to)
		}
	} else {
		switch to.T {
		case IntTy, UintTy:
			var signed bool = to.T == IntTy
			if typ, ok := from.(int); !ok {
				return nil, fmt.Errorf("Unexpected non integer type during type conversion, please reformat your run file to use an integer.")
			} else {
				switch to.Size {
				case 8:
					if signed {
						return int8(typ), nil
					}
					return uint8(typ), nil
				case 16:
					if signed {
						return int16(typ), nil
					}
					return uint16(typ), nil
				case 32:
					if signed {
						return int32(typ), nil
					}
					return uint32(typ), nil
				case 64:
					if signed {
						return int64(typ), nil
					}
					return uint64(typ), nil
				default:
					big := interpret.Big0
					if signed {
						return big.SetInt64(int64(typ)), nil
					}
					return big.SetUint64(uint64(typ)), nil
				}
			}
		case BoolTy:
			if typ, ok := from.(bool); !ok {
				return nil, fmt.Errorf("Unexpected non bool type during type conversion, please reformat your run file to use a bool.")
			} else {
				return typ, nil
			}
		case StringTy:
			if typ, ok := from.(string); !ok {
				return nil, fmt.Errorf("Unexpected non string type during type conversion, please reformat your run file to use a string.")
			} else {
				return typ, nil
			}
		case AddressTy:
			if typ, ok := from.(string); !ok {
				return nil, fmt.Errorf("Unexpected non string type during type conversion, please reformat your run file to use a string.")
			} else {
				return interpret.HexToAddress(typ), nil
			}
		case BytesTy:
			if typ, ok := from.(string); !ok {
				return nil, fmt.Errorf("Unexpected non string type during type conversion, please reformat your run file to use a string.")
			} else {
				return interpret.HexToBytes(typ), nil
			}
		default:
			return nil, fmt.Errorf("Invalid type during type conversion.")
		}
	}
}

func ConvertUnpackedToEpmTypes(from interface{}, reference Type) (string, interface{}, error) {
	if reference.IsSlice || reference.IsArray && reference.T != FixedBytesTy && reference.T != BytesTy {
		var normalSliceString = func(i interface{}) string {
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(i)
			return fmt.Sprintf(buf.String())
		}
		// convert to yaml createable types, ignoring string and bool because those are accounted for already
		sliceVal := reflect.ValueOf(from)
		var stored []interface{}
		for i := 0; i < sliceVal.Len(); i++{
			_, typ, err := ConvertUnpackedToEpmTypes(sliceVal.Index(i).Interface(), *reference.Elem)
			stored := append(stored, typ)
		}
		return normalSliceString(stored), stored, nil
	} else {
		switch reference.T {
		case UintTy, IntTy:
			switch typ := from.(type) {
			case int8, int16, int32, int64:
				return fmt.Sprintf("%v", from), int(typ.(int)), nil
			case uint8, uint16, uint32, uint64:
				return fmt.Sprintf("%v", from), int(typ.(uint)), nil
			case *big.Int:
				val := typ.Int64()
				if val == 0 {
					val := typ.Uint64()
					return typ.String(), int(val), nil
				}
				return typ.String(), int(val), nil
			}
		case StringTy:
			return from.(string), from.(string), nil
		case BoolTy:
			return fmt.Sprintf("%v", from), from.(bool), nil
		case AddressTy:
			return from.(interpret.Address).Str(), from.(interpret.Address).Str(), nil
		case BytesTy, FixedBytesTy:
			return string(bytes.Trim(from.([]byte), "\x00")[:]), string(bytes.Trim(from.([]byte), "\x00")[:]), nil
		default:
			return "", nil, fmt.Errorf("Could not find type to convert.")
		}
	}

}