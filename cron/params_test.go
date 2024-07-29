package cron

import (
	"reflect"
	"testing"
)

func TestParseJsonParams(t *testing.T) {
	type testStruct struct {
		UserID   string   `json:"user_id"`
		Contents []string `json:"contents"`
	}
	testCases := []struct {
		name         string
		paramStr     string
		paramStruct  *testStruct
		expectStruct *testStruct
	}{
		{
			name:        "default",
			paramStr:    `{"user_id":"u1","contents":["c1","c2","c3"]}`,
			paramStruct: &testStruct{},
			expectStruct: &testStruct{
				UserID:   "u1",
				Contents: []string{"c1", "c2", "c3"},
			},
		},
		{
			name:         "case_empty",
			paramStr:     ``,
			paramStruct:  &testStruct{},
			expectStruct: &testStruct{},
		},
	}

	for _, tcase := range testCases {
		ParseJsonParams(tcase.paramStr, tcase.paramStruct)
		if !reflect.DeepEqual(tcase.expectStruct, tcase.paramStruct) {
			t.Errorf("parse_params case %s error. expect_struct:%+v, actually:%+v",
				tcase.name, tcase.expectStruct, tcase.paramStruct)
		}
	}
}

func TestParseUrlParams(t *testing.T) {
	type testStruct struct {
		UserID   string   `json:"user_id"`
		Contents []string `json:"contents"`
	}

	testCases := []struct {
		name         string
		paramStr     string
		paramStruct  *testStruct
		expectStruct *testStruct
	}{
		{
			name:        "default",
			paramStr:    `contents=c1,c2,c3&user_id=u1`,
			paramStruct: &testStruct{},
			expectStruct: &testStruct{
				UserID:   "u1",
				Contents: []string{"c1", "c2", "c3"},
			},
		},
	}

	for _, tcase := range testCases {
		ParseUrlParams(tcase.paramStr, tcase.paramStruct)
		if !reflect.DeepEqual(tcase.expectStruct, tcase.paramStruct) {
			t.Errorf("parse_params case %s error. expect_struct:%+v, actually:%+v",
				tcase.name, tcase.expectStruct, tcase.paramStruct)
		}
	}
}
