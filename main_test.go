package main

import "testing"

var mongoString = "{\"_id\":{\"$oid\":\"1\"},\"date_value\":{\"$date\":\"2019-06-11T12:23:14.496Z\"},\"decimal_value\":{\"$numberDecimal\":\"42.42\"},\"float_value\":42.42,\"int_value\":42,\"long_value\":{\"$numberLong\":\"42424242424242424242424242\"},\"str_value\":\"test str\",\"x\":{\"a\":1},\"y\":[{\"$oid\":\"2\"},3,4],\"z\":{\"_id\":{\"$oid\":\"3\"}}}"

var bqString = "{\"_id\":\"1\",\"date_value\":\"2019-06-11T12:23:14.496Z\",\"decimal_value\":\"42.42\",\"float_value\":42.42,\"int_value\":42,\"long_value\":\"42424242424242424242424242\",\"str_value\":\"test str\",\"x\":{\"a\":1},\"y\":[\"2\",3,4],\"z\":{\"_id\":\"3\"}}"

func TestConvert(t *testing.T) {
	convertedString := convert([]byte(mongoString))
	if string(convertedString) != bqString {
		t.Error("expected", bqString, "got", convertedString)
	}
}

func BenchmarkConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		convert([]byte(mongoString))
	}
}
