package router

import (
	"reflect"
	"testing"
)

func TestInsertToIndex(t *testing.T) {
    initialArray := []int{1,2,3,4}
    expected1 := []int{99,1,2,3,4}
    expected2 := []int{1,2,99,3,4}
    expected3 := []int{1,2,3,4,99}

    result := insertToIndex(initialArray, 0, 99)
    if !reflect.DeepEqual(expected1, result) {
        t.Errorf("Insert to index: expected %v got %v", expected1, result)
    }

    result = insertToIndex(initialArray, 2, 99)
    if !reflect.DeepEqual(expected2, result) {
        t.Errorf("Insert to index: expected %v got %v", expected2, result)
    }

    result = insertToIndex(initialArray, 4, 99)
    if !reflect.DeepEqual(expected3, result) {
        t.Errorf("Insert to index: expected %v got %v", expected3, result)
    }
}

func TestKeys(t *testing.T) {
    testCase := map[string]int{
        "1" : 2,
        "33": 4,
        "asd": 123,
    }

    result := keys(testCase)
    expected := []string{"1","33","asd"}
    if !reflect.DeepEqual(expected, result) {
        t.Errorf("utils keys: expected %q got %q",expected, result)
    }
}
