package buildin_test

import (
	"testing"

	"github.com/vaihdass/webber/buildin"
)

func TestPtr(t *testing.T) {
	t.Parallel()

	t.Run("string", testPtrString)
	t.Run("int", testPtrInt)
	t.Run("struct", testPtrStruct)
	t.Run("slice", testPtrSlice)
	t.Run("zero_value", testPtrZeroValue)
}

func testPtrString(t *testing.T) {
	t.Parallel()
	value := "test"
	ptr := buildin.Ptr(value)

	if ptr == nil {
		t.Fatal("expected non-nil pointer")
	}

	if *ptr != value {
		t.Errorf("expected %q, got %q", value, *ptr)
	}

	// Verify it's a copy, not the same address
	if &value == ptr {
		t.Error("expected different memory addresses")
	}
}

func testPtrInt(t *testing.T) {
	t.Parallel()
	value := 42
	ptr := buildin.Ptr(value)

	if ptr == nil {
		t.Fatal("expected non-nil pointer")
	}

	if *ptr != value {
		t.Errorf("expected %d, got %d", value, *ptr)
	}
}

func testPtrStruct(t *testing.T) {
	t.Parallel()
	type TestStruct struct {
		Name string
		Age  int
	}

	value := TestStruct{Name: "John", Age: 30}
	ptr := buildin.Ptr(value)

	if ptr == nil {
		t.Fatal("expected non-nil pointer")
	}

	if ptr.Name != value.Name || ptr.Age != value.Age {
		t.Errorf("expected %+v, got %+v", value, *ptr)
	}
}

func testPtrSlice(t *testing.T) {
	t.Parallel()
	value := []int{1, 2, 3}
	ptr := buildin.Ptr(value)

	if ptr == nil {
		t.Fatal("expected non-nil pointer")
	}

	if len(*ptr) != len(value) {
		t.Errorf("expected length %d, got %d", len(value), len(*ptr))
	}

	for i, v := range value {
		if (*ptr)[i] != v {
			t.Errorf("expected %d at index %d, got %d", v, i, (*ptr)[i])
		}
	}
}

func testPtrZeroValue(t *testing.T) {
	t.Parallel()
	var value int
	ptr := buildin.Ptr(value)

	if ptr == nil {
		t.Fatal("expected non-nil pointer")
	}

	if *ptr != 0 {
		t.Errorf("expected 0, got %d", *ptr)
	}
}

func TestFromPtr(t *testing.T) {
	t.Parallel()
	t.Run("non_nil_string", testFromPtrNonNilString)
	t.Run("non_nil_int", testFromPtrNonNilInt)
	t.Run("non_nil_struct", testFromPtrNonNilStruct)
	t.Run("nil_string", testFromPtrNilString)
	t.Run("nil_int", testFromPtrNilInt)
	t.Run("nil_struct", testFromPtrNilStruct)
	t.Run("nil_slice", testFromPtrNilSlice)
}

func testFromPtrNonNilString(t *testing.T) {
	t.Parallel()
	value := "test"
	ptr := &value

	result, ok := buildin.FromPtr(ptr)

	if !ok {
		t.Error("expected ok to be true")
	}

	if result != value {
		t.Errorf("expected %q, got %q", value, result)
	}
}

func testFromPtrNonNilInt(t *testing.T) {
	t.Parallel()
	value := 42
	ptr := &value

	result, ok := buildin.FromPtr(ptr)

	if !ok {
		t.Error("expected ok to be true")
	}

	if result != value {
		t.Errorf("expected %d, got %d", value, result)
	}
}

func testFromPtrNonNilStruct(t *testing.T) {
	t.Parallel()
	type TestStruct struct {
		Name string
		Age  int
	}

	value := TestStruct{Name: "Jane", Age: 25}
	ptr := &value

	result, ok := buildin.FromPtr(ptr)

	if !ok {
		t.Error("expected ok to be true")
	}

	if result.Name != value.Name || result.Age != value.Age {
		t.Errorf("expected %+v, got %+v", value, result)
	}
}

func testFromPtrNilString(t *testing.T) {
	t.Parallel()
	var ptr *string

	result, ok := buildin.FromPtr(ptr)

	if ok {
		t.Error("expected ok to be false")
	}

	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func testFromPtrNilInt(t *testing.T) {
	t.Parallel()
	var ptr *int

	result, ok := buildin.FromPtr(ptr)

	if ok {
		t.Error("expected ok to be false")
	}

	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}

func testFromPtrNilStruct(t *testing.T) {
	t.Parallel()
	type TestStruct struct {
		Name string
		Age  int
	}

	var ptr *TestStruct

	result, ok := buildin.FromPtr(ptr)

	if ok {
		t.Error("expected ok to be false")
	}

	var zero TestStruct
	if result != zero {
		t.Errorf("expected zero value %+v, got %+v", zero, result)
	}
}

func testFromPtrNilSlice(t *testing.T) {
	t.Parallel()
	var ptr *[]int

	result, ok := buildin.FromPtr(ptr)

	if ok {
		t.Error("expected ok to be false")
	}

	if result != nil {
		t.Errorf("expected nil slice, got %+v", result)
	}
}

func TestPtrFromPtrRoundtrip(t *testing.T) {
	t.Parallel()
	t.Run("string_roundtrip", func(t *testing.T) {
		t.Parallel()
		original := "hello world"

		ptr := buildin.Ptr(original)
		result, ok := buildin.FromPtr(ptr)

		if !ok {
			t.Error("expected ok to be true")
		}

		if result != original {
			t.Errorf("expected %q, got %q", original, result)
		}
	})

	t.Run("int_roundtrip", func(t *testing.T) {
		t.Parallel()
		original := 123

		ptr := buildin.Ptr(original)
		result, ok := buildin.FromPtr(ptr)

		if !ok {
			t.Error("expected ok to be true")
		}

		if result != original {
			t.Errorf("expected %d, got %d", original, result)
		}
	})

	t.Run("struct_roundtrip", func(t *testing.T) {
		t.Parallel()
		type Person struct {
			Name string
			Age  int
		}

		original := Person{Name: "Alice", Age: 30}

		ptr := buildin.Ptr(original)
		result, ok := buildin.FromPtr(ptr)

		if !ok {
			t.Error("expected ok to be true")
		}

		if result != original {
			t.Errorf("expected %+v, got %+v", original, result)
		}
	})
}
