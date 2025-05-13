package pool

import (
	"testing"
)

func TestNewSlice(t *testing.T) {
	capSave := 10
	capCreate := 20
	def := []int{1, 2, 3}

	slicePool := NewSlice(capSave, capCreate, def)

	if slicePool == nil {
		t.Fatal("NewSlice returned nil")
	}

	if slicePool.cap != capSave {
		t.Errorf("Expected capSave to be %d, got %d", capSave, slicePool.cap)
	}

	if len(slicePool.def) != len(def) {
		t.Errorf("Expected default slice length to be %d, got %d", len(def), len(slicePool.def))
	}

	for i, v := range def {
		if slicePool.def[i] != v {
			t.Errorf("Expected default slice element at index %d to be %d, got %d", i, v, slicePool.def[i])
		}
	}
}

func TestSliceGet(t *testing.T) {
	capSave := 10
	capCreate := 20
	def := []int{1, 2, 3}

	slicePool := NewSlice(capSave, capCreate, def)
	slice := slicePool.Get()

	if slice == nil {
		t.Fatal("Get returned nil")
	}

	if len(slice) != len(def) {
		t.Errorf("Expected slice length to be len(def), got %d", len(slice))
	}

	if cap(slice) != capCreate {
		t.Errorf("Expected slice capacity to be %d, got %d", capCreate, cap(slice))
	}

	for i, v := range def {
		if slice[i] != v {
			t.Errorf("Expected slice element at index %d to be %d, got %d", i, v, slice[i])
		}
	}
}

func TestSliceSave(t *testing.T) {
	capSave := 10
	capCreate := 20
	def := []int{1, 2, 3}

	slicePool := NewSlice(capSave, capCreate, def)
	slice := slicePool.Get()

	// Изменяем slice
	slice = append(slice, 4, 5)
	slice[0] = 6

	// Сохраняем slice обратно в пул
	slicePool.Save(slice)

	// Получаем новый slice из пула
	newSlice := slicePool.Get()

	if len(newSlice) != len(def) {
		t.Errorf("Expected new slice length to be %d, got %d", len(def), len(newSlice))
	}

	for i, v := range def {
		if newSlice[i] != v {
			t.Errorf("Expected new slice element at index %d to be %d, got %d", i, v, newSlice[i])
		}
	}

	// Проверяем, что slice с превышающей длиной не сохраняется
	longSlice := make([]int, capSave+1)
	slicePool.Save(longSlice)

	// Получаем новый slice из пула
	anotherSlice := slicePool.Get()

	if len(anotherSlice) != len(def) {
		t.Errorf("Expected another slice length to be %d, got %d", len(def), len(anotherSlice))
	}

	for i, v := range def {
		if anotherSlice[i] != v {
			t.Errorf("Expected another slice element at index %d to be %d, got %d", i, v, anotherSlice[i])
		}
	}
}

func TestSliceParallel(t *testing.T) {
	capSave := 10
	capCreate := 20
	def := []int{1, 2, 3}

	slicePool := NewSlice(capSave, capCreate, def)

	t.Run("ParallelTest1", func(t *testing.T) {
		t.Parallel()
		slice := slicePool.Get()
		if slice == nil {
			t.Fatal("Get returned nil")
		}
		if len(slice) != len(def) {
			t.Errorf("Expected slice length to be len(def), got %d", len(slice))
		}
		if cap(slice) != capCreate {
			t.Errorf("Expected slice capacity to be %d, got %d", capCreate, cap(slice))
		}
		for i, v := range def {
			if slice[i] != v {
				t.Errorf("Expected slice element at index %d to be %d, got %d", i, v, slice[i])
			}
		}
	})

	t.Run("ParallelTest2", func(t *testing.T) {
		t.Parallel()
		slice := slicePool.Get()
		if slice == nil {
			t.Fatal("Get returned nil")
		}
		if len(slice) != len(def) {
			t.Errorf("Expected slice length to be  len(def), got %d", len(slice))
		}
		if cap(slice) != capCreate {
			t.Errorf("Expected slice capacity to be %d, got %d", capCreate, cap(slice))
		}
		for i, v := range def {
			if slice[i] != v {
				t.Errorf("Expected slice element at index %d to be %d, got %d", i, v, slice[i])
			}
		}
	})
}
