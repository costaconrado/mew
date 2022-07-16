package mew

import "reflect"

func isPointer(obj any) bool { return reflect.ValueOf(obj).Kind() == reflect.Ptr }
