package daog

import (
	"errors"
	"fmt"
)

type FieldPointExtractor interface {
	// Extract 返回值内存储一个指针，当给定的fieldName不存在时，返回nil
	Extract(fieldName string) any
}

type AfterTransBeginInterceptor func(tx *TransContext) error
type ChangeFieldValueBeforeWriteInterceptor func(valueMap map[string]any, extractor FieldPointExtractor) error
type AddNewModifyFieldBeforeUpdateInterceptor func(valueMap map[string]any, modifier Modifier, existField func (filedName string) bool) error


var TransBegunInterceptor AfterTransBeginInterceptor

var ChangeFieldOfInsBeforeWrite ChangeFieldValueBeforeWriteInterceptor
var AddNewModifyFieldBeforeUpdate AddNewModifyFieldBeforeUpdateInterceptor


func ChangeInt64ByFieldNameCallback(valueMap map[string]any, fieldName string, extractor FieldPointExtractor) error{
	value,ok := valueMap[fieldName]
	if !ok {
		return nil
	}
	id := value.(int64)
	pt := extractor.Extract(fieldName)
	if pt == nil {
		return nil
	}
	targetPt,ok := pt.(*int64)
	if !ok {
		return errors.New(fmt.Sprintf("%s must be int64",fieldName))
	}

	*targetPt = id
	return nil
}

func ChangeModifierByFieldNameCallback(valueMap map[string]any, fieldName string, modifier Modifier, existField func (filedName string) bool) error{
	value,ok := valueMap[fieldName]
	if !ok {
		return nil
	}

	if !existField(fieldName)  {
		return nil
	}

	modifier.Add(fieldName,value)
	return nil
}