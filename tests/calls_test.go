package tests

import (
	"math"
	"testing"
)

func TestVoid_Void(t *testing.T) {
	var set bool
	Void_Void(func() {
		set = true
	})
	if !set {
		t.Error("Functions was not called")
	}
}

func TestVoid_Int(t *testing.T) {
	var set cint
	const expected = cint(math.MinInt32 + 12)
	Void_Int(func(arg1 cint) {
		set = arg1
	}, expected)
	if set != expected {
		t.Errorf("Bad call. Expected %v, got %v", expected, set)
	}
}

func TestVoid_Uint(t *testing.T) {
	var set cuint
	const expected = cuint(math.MaxUint32 - 134)
	Void_Uint(func(arg1 cuint) {
		set = arg1
	}, expected)
	if set != expected {
		t.Errorf("Bad call. Expected %v, got %v", expected, set)
	}
}

func TestVoid_IntInt(t *testing.T) {
	var set1, set2 cint
	const (
		expected1 = cint(math.MinInt32 + 12)
		expected2 = cint(math.MaxInt32)
	)
	Void_IntInt(func(arg1, arg2 cint) {
		set1 = arg1
		set2 = arg2
	}, expected1, expected2)
	if set1 != expected1 {
		t.Errorf("Bad call. Expected %v, got %v", expected1, set1)
	}
	if set2 != expected2 {
		t.Errorf("Bad call. Expected %v, got %v", expected2, set2)
	}
}

func TestVoid_Float(t *testing.T) {
	var set float32
	const expected float32 = 3.14
	Void_Float(func(arg float32){
		set = arg
	}, expected)
	if set != expected {
		t.Errorf("Bad call. Expected %v, got %v", expected, set)
	}
}

func TestVoid_Double(t *testing.T) {
	var set float64
	const expected float64 = 3.14
	Void_Double(func(arg float64){
		set = arg
	}, expected)
	if set != expected {
		t.Errorf("Bad call. Expected %v, got %v", expected, set)
	}
}
