package main

import "reflect"

// Contains checks if the element exists in slice
func Contains(slice interface{}, element interface{}) bool {
	s := reflect.ValueOf(slice)
	
	if s.Kind() != reflect.Slice {
		panic("Function Contains() given a non-slice type")
	}
	// map keys should be comparable
	set := make(map[interface{}]struct{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		set[s.Index(i).Interface()] = struct{}{}
	}
	// 2 interface{} is equal if they 
	// have identical dynamic types and equal dynamic values or if both have value nil.
	_, exists := set[element]
	return exists
}