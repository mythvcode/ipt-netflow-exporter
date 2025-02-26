package statparser

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unicode"
)

const (
	cpuStatFieldsCount  = 12
	sockStatFieldsCount = 10
)

var (
	errNotFoundField   = errors.New("stat field not found")
	errUpdate          = errors.New("struct field cannot be updated")
	errUnsupportedType = errors.New("unsupported type")
)

type statType interface {
	fieldCount() int
	typeName() string
}

// metric with names from stats files
type Statistics struct {
	InBitRate    uint64
	InPacketRate uint64
	InFlows      uint64
	InPackets    uint64
	InBytes      uint64
	HashMetric   float64
	HashMemory   uint64
	HashFlows    uint64
	HashPackets  uint64
	HashBytes    uint64
	DropPackets  uint64
	DropBytes    uint64
	OutByteRate  uint64
	OutFlows     uint64
	OutPackets   uint64
	OutBytes     uint64
	LostFlows    uint64
	LostPackets  uint64
	LostBytes    uint64
	ErrTotal     uint64
	SndbufPeak   uint64
	CPUStatList  []CPUStat
	SockStatList []NFSockEntry
}

// Order of fields must be the same as in pt_netflow_snmp
type CPUStat struct {
	CPU             string
	CPUInPacketRate uint64
	CPUInFlows      uint64
	CPUInPackets    uint64
	CPUInBytes      uint64
	CPUHashMetric   float64
	CPUDropPackets  uint64
	CPUuDropBytes   uint64
	CPUErrTrunc     uint64
	CPUErrFrag      uint64
	CPUErrAlloc     uint64
	CPUErrMaxflows  uint64
}

// Order of fields must be the same as in pt_netflow_snmp
type NFSockEntry struct {
	SockName        string
	SockDestination string
	SockActive      uint32
	SockErrConnect  uint32
	SockErrFull     uint32
	SockErrCberr    uint32
	SockErrOther    uint32
	SockSndbuf      uint32
	SockSndbufFill  uint32
	SockSndbufPeak  uint32
}

func (c *CPUStat) fieldCount() int {
	return cpuStatFieldsCount
}

func (c *CPUStat) typeName() string {
	return "cpu"
}

func (c *NFSockEntry) fieldCount() int {
	return sockStatFieldsCount
}

func (c *NFSockEntry) typeName() string {
	return "socket"
}

func toUpperFirstChar(str string) string {
	if len(str) == 0 {
		return str
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func checkStructField(structField reflect.Value) error {
	if !structField.IsValid() {
		return errNotFoundField
	}

	if !structField.CanSet() {
		return errUpdate
	}

	return nil
}

// parses value according to field type and return parsed value as reflect.Value
func getValueByType(field reflect.Value, value string) (reflect.Value, error) {
	switch field.Kind() { //nolint
	case reflect.Uint64:
		intVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(intVal), nil
	case reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(floatVal), nil
	case reflect.Uint32:
		intVal, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(uint32(intVal)), nil

	case reflect.String:

		return reflect.ValueOf(value), nil

	default:
		return reflect.Value{}, errUnsupportedType
	}
}

func setValueByName(stat *Statistics, metricName, value string) error {
	structVal := reflect.ValueOf(stat).Elem()
	structField := structVal.FieldByName(toUpperFirstChar(metricName))
	if err := checkStructField(structField); err != nil {
		return err
	}
	val, err := getValueByType(structField, value)
	if err != nil {
		if errors.Is(err, errUnsupportedType) {
			return fmt.Errorf("error unsupported field name type: field %s, type: %s", toUpperFirstChar(metricName), structField.Kind().String())
		}

		return err
	}
	structField.Set(val)

	return nil
}

func setValues(stat statType, fields []string) error {
	if len(fields) != stat.fieldCount() {
		return fmt.Errorf("error parse fields count for %s stat: must be %d, actual %d", stat.typeName(), stat.fieldCount(), len(fields))
	}
	structVal := reflect.ValueOf(stat).Elem()
	if !structVal.IsValid() {
		return fmt.Errorf("error value of stat parameter")
	}

	for index, rawValue := range fields {
		structField := structVal.Field(index)
		if err := checkStructField(structField); err != nil {
			return err
		}
		val, err := getValueByType(structField, rawValue)
		if err != nil {
			if errors.Is(err, errUnsupportedType) {
				return fmt.Errorf("error parse field type: position %d, type: %s", index, structField.Kind().String())
			}

			return err
		}
		structField.Set(val)
	}

	return nil
}
