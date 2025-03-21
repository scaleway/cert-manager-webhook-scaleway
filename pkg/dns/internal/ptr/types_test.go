package ptr

import (
	"testing"
)

func TestDeref(t *testing.T) {
	t.Run("Deref with nil pointer", func(t *testing.T) {
		var nilInt *int
		result := Deref(nilInt)
		var expected int
		if result != expected {
			t.Errorf("expected %v, but got %v", expected, result)
		}
	})

	t.Run("Deref with non-nil pointer", func(t *testing.T) {
		value := 42
		result := Deref(&value)
		if result != value {
			t.Errorf("expected %v, but got %v", value, result)
		}
	})

	t.Run("Deref with nil string pointer", func(t *testing.T) {
		var nilString *string
		result := Deref(nilString)
		var expected string
		if result != expected {
			t.Errorf("expected %v, but got %v", expected, result)
		}
	})

	t.Run("Deref with non-nil string pointer", func(t *testing.T) {
		value := "hello"
		result := Deref(&value)
		if result != value {
			t.Errorf("expected %v, but got %v", value, result)
		}
	})
}

func TestPointer(t *testing.T) {
	t.Run("Pointer with int value", func(t *testing.T) {
		value := 42
		result := Pointer(value)
		if *result != value {
			t.Errorf("expected %v, but got %v", value, *result)
		}
	})

	t.Run("Pointer with string value", func(t *testing.T) {
		value := "hello"
		result := Pointer(value)
		if *result != value {
			t.Errorf("expected %v, but got %v", value, *result)
		}
	})
}
