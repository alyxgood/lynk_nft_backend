package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

func Decode(outputParameters []string, input string) (result []interface{}, err error) {
	if outputParameters == nil || len(outputParameters) == 0 {
		result = append(result, input)
	} else {
		inputIndex := 0
		for _, outputParameter := range outputParameters {
			if outputParameter == "uint256" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}
				var data *big.Int
				data, err = parseUint256(input[inputIndex : inputIndex+64])
				if err != nil {
					break
				}

				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "address" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}
				result = append(result, fmt.Sprintf("0x%s", input[inputIndex+24:inputIndex+64]))
				inputIndex += 64
			} else if outputParameter == "string" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}
				var tmpDataInt *big.Int
				tmpDataInt, err = parseUint256(input[inputIndex : inputIndex+64])
				if err != nil {
					break
				}
				dataOffset := int(tmpDataInt.Uint64() * 2)
				if len(input) < dataOffset+64 {
					err = errors.New("invalid data")
					break
				}
				tmpDataInt, err = parseUint256(input[dataOffset : dataOffset+64])
				if err != nil {
					break
				}
				dataLength := int(tmpDataInt.Uint64() * 2)
				if len(input) < dataOffset+64+dataLength {
					err = errors.New("invalid data")
					break
				}
				var data string
				data, err = parseString(input[dataOffset+64 : dataOffset+64+dataLength])
				if err != nil {
					break
				}
				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "bool" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}
				var b uint64
				b, err = strconv.ParseUint(input[inputIndex:inputIndex+64], 16, 64)
				if err != nil {
					break
				}
				var data bool
				if b == 0 {
					data = false
				} else if b == 1 {
					data = true
				}

				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "uint8" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}
				var data uint64
				data, err = strconv.ParseUint(input[inputIndex:inputIndex+64], 16, 64)
				if err != nil {
					break
				}

				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "address[]" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}

				var tmpDataInt *big.Int
				tmpDataInt, err = parseUint256(input[inputIndex : inputIndex+64])
				if err != nil {
					break
				}
				dataOffset := int(tmpDataInt.Uint64() * 2)
				if len(input) < dataOffset+64 {
					err = errors.New("invalid data")
					break
				}
				tmpDataInt, err = parseUint256(input[dataOffset : dataOffset+64])
				if err != nil {
					break
				}
				dataLength := int(tmpDataInt.Uint64())
				if len(input) < dataOffset+64+(dataLength*64) {
					err = errors.New("invalid data")
					break
				}
				var data []string
				for index := 0; index < dataLength; index++ {
					data = append(data, fmt.Sprintf("0x%s", input[dataOffset+64+(index*64)+24:dataOffset+64+(index*64)+64]))
				}
				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "uint256[]" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}

				var tmpDataInt *big.Int
				tmpDataInt, err = parseUint256(input[inputIndex : inputIndex+64])
				if err != nil {
					break
				}
				dataOffset := int(tmpDataInt.Uint64() * 2)
				if len(input) < dataOffset+64 {
					err = errors.New("invalid data")
					break
				}
				tmpDataInt, err = parseUint256(input[dataOffset : dataOffset+64])
				if err != nil {
					break
				}
				dataLength := int(tmpDataInt.Uint64())
				if len(input) < dataOffset+64+(dataLength*64) {
					err = errors.New("invalid data")
					break
				}
				var data []*big.Int
				for index := 0; index < dataLength; index++ {
					var datum *big.Int
					datum, err = parseUint256(input[dataOffset+64+(index*64) : dataOffset+64+(index*64)+64])
					if err != nil {
						break
					}
					data = append(data, datum)
				}
				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "string[]" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}

				var tmpDataInt *big.Int
				tmpDataInt, err = parseUint256(input[inputIndex : inputIndex+64])
				if err != nil {
					break
				}
				dataOffset := int(tmpDataInt.Uint64() * 2)
				if len(input) < dataOffset+64 {
					err = errors.New("invalid data")
					break
				}
				tmpDataInt, err = parseUint256(input[dataOffset : dataOffset+64])
				if err != nil {
					break
				}
				dataLength := int(tmpDataInt.Uint64())
				if len(input) < dataOffset+64+(dataLength*64) {
					err = errors.New("invalid data")
					break
				}
				var data []string
				for index := 0; index < dataLength; index++ {
					var offset, length *big.Int
					offset, err = parseUint256(input[dataOffset+64+(index*64) : dataOffset+64+(index*64)+64])
					length, err = parseUint256(input[dataOffset+64+int(offset.Int64()*2) : dataOffset+64+int(offset.Int64()*2)+64])
					var datum string
					datum, err = parseString(input[dataOffset+64+int(offset.Int64()*2)+64 : dataOffset+64+int(offset.Int64()*2)+64+int(length.Int64()*2)])
					data = append(data, datum)
				}
				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "bytes[]" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}

				var tmpDataInt *big.Int
				tmpDataInt, err = parseUint256(input[inputIndex : inputIndex+64])
				if err != nil {
					break
				}
				dataOffset := int(tmpDataInt.Uint64() * 2)
				if len(input) < dataOffset+64 {
					err = errors.New("invalid data")
					break
				}
				tmpDataInt, err = parseUint256(input[dataOffset : dataOffset+64])
				if err != nil {
					break
				}
				dataLength := int(tmpDataInt.Uint64())
				if len(input) < dataOffset+64+(dataLength*64) {
					err = errors.New("invalid data")
					break
				}
				var data []string
				for index := 0; index < dataLength; index++ {
					var offset, length *big.Int
					offset, err = parseUint256(input[dataOffset+64+(index*64) : dataOffset+64+(index*64)+64])
					length, err = parseUint256(input[dataOffset+64+int(offset.Int64()*2) : dataOffset+64+int(offset.Int64()*2)+64])
					var datum string
					datum = input[dataOffset+64+int(offset.Int64()*2)+64 : dataOffset+64+int(offset.Int64()*2)+64+int(length.Int64()*2)]
					data = append(data, fmt.Sprintf("0x%s", datum))
				}
				result = append(result, data)
				inputIndex += 64
			} else {
				err = errors.New(fmt.Sprintf("unsupported %s data", outputParameter))
				break
			}
		}
	}

	return result, err
}

func parseUint256(input string) (*big.Int, error) {
	bytes, err := hex.DecodeString(input)
	if err != nil {
		return nil, err
	}

	return big.NewInt(0).SetBytes(bytes), nil
}

func parseString(input string) (string, error) {
	bytes, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}
	result := string(bytes)

	return result, nil
}
