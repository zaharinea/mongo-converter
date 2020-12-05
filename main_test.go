package main

import "testing"

var mongoString = `{"_id":{"$oid":"1"},"date_value":{"$date":"2019-06-11T12:23:14.496Z"},"decimal_value":{"$numberDecimal":"42.42"},"float_value":42.42,"int_value":42,"long_value":{"$numberLong":"9223372036854775807"},"str_value":"test str","uuid_field":{"$binary":{"base64":"ZhIYnDxBQsKhG4XbsucYSQ==","subType":"04"}},"x":{"a":1},"y":[{"$oid":"2"},3,4],"z":{"_id":{"$oid":"3"}}}`

var bqString = `{"_id":"1","date_value":"2019-06-11T12:23:14.496Z","decimal_value":"42.42","float_value":42.42,"int_value":42,"long_value":"9223372036854775807","str_value":"test str","uuid_field":"6612189c-3c41-42c2-a11b-85dbb2e71849","x":{"a":1},"y":["2",3,4],"z":{"_id":"3"}}`

var bqStringUnwrapedNumbers = `{"_id":"1","date_value":"2019-06-11T12:23:14.496Z","decimal_value":42.42,"float_value":42.42,"int_value":42,"long_value":9223372036854775807,"str_value":"test str","uuid_field":"6612189c-3c41-42c2-a11b-85dbb2e71849","x":{"a":1},"y":["2",3,4],"z":{"_id":"3"}}`

func TestConvert(t *testing.T) {
	convertedString := convert([]byte(mongoString), false)
	if string(convertedString) != bqString {
		t.Error("expected", bqString, "got", string(convertedString))
	}
}

func TestConvertUnwrapNumbers(t *testing.T) {
	convertedString := convert([]byte(mongoString), true)
	if string(convertedString) != bqStringUnwrapedNumbers {
		t.Error("expected", bqStringUnwrapedNumbers, "got", string(convertedString))
	}
}

func BenchmarkConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		convert([]byte(mongoString), false)
	}
}
