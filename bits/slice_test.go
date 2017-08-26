package bits

import (
	"testing"
)

func AssertPanic(t *testing.T) {
	if ok := recover(); ok == nil {
		t.Errorf("AssertPanic failed")
	}
}

func TestNewSlice(t *testing.T) {
	slice := *NewSlice(5, 0x1f)
	if slice.length != 5 || slice.data != 0x1f {
		t.Errorf("NewSlice failed")
	}

	func() {
		defer AssertPanic(t)
		_ = NewSlice(100, 0x0)
	}()
}

func TestSliceAppendSlice(t *testing.T) {
	slice := *NewSlice(4, 0xf)

	slice.AppendSlice(NewSlice(8, 0xff))
	if slice.length != 12 || slice.data != 0xfff {
		t.Errorf("AppendSlice failed: %v", slice)
	}

	slice.AppendSlice(NewSlice(4, 0x0))
	if slice.length != 16 || slice.data != 0xfff0 {
		t.Errorf("AppendSlice failed: %v", slice)
	}

	func() {
		defer AssertPanic(t)
		slice.AppendSlice(NewSlice(50, 0x0))
	}()
}

func TestSliceAppendBit(t *testing.T) {
	slice := NewSlice(0, 0x0)

	slice.AppendBit(true)
	if slice.length != 1 || slice.data != 0x1 {
		t.Errorf("AppendBit failed: %v", slice)
	}

	slice.AppendBit(false)
	if slice.length != 2 || slice.data != 0x2 {
		t.Errorf("AppendBit failed: %v", slice)
	}

	slice.AppendBit(false)
	if slice.length != 3 || slice.data != 0x4 {
		t.Errorf("AppendBit failed: %v", slice)
	}

	func() {
		defer AssertPanic(t)
		slice = NewSlice(64, 0x0)
		slice.AppendBit(true)
	}()

	func() {
		defer AssertPanic(t)
		slice = NewSlice(64, 0x0)
		slice.AppendBit(false)
	}()

}

func TestAppendPadding(t *testing.T) {
	for i := 0; i < 64; i++ {
		data := (uint64(1) << uint(i)) - 1
		slice := NewSlice(i, data)
		slice.AppendPadding()

		padding := 0
		if i%8 != 0 {
			padding = 8 - (i % 8)
		}

		expected := NewSlice(i+padding, data<<uint(padding))

		if *slice != *expected {
			t.Errorf("Expected %v, got %v", *expected, *slice)
		}
	}

}
