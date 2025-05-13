package pool

import (
	"testing"
)

func TestNewMap(t *testing.T) {
	capSave := 10
	capCreate := 20
	def := map[string]int{"a": 1, "b": 2, "c": 3}

	mapPool := NewMap(capSave, capCreate, def)

	if mapPool == nil {
		t.Fatal("NewMap returned nil")
	}

	if mapPool.cap != capSave {
		t.Errorf("Expected capSave to be %d, got %d", capSave, mapPool.cap)
	}

	if len(mapPool.def) != len(def) {
		t.Errorf("Expected default map length to be %d, got %d", len(def), len(mapPool.def))
	}

	for k, v := range def {
		if mapPool.def[k] != v {
			t.Errorf("Expected default map element with key %s to be %d, got %d", k, v, mapPool.def[k])
		}
	}
}

func TestMapGet(t *testing.T) {
	capSave := 10
	capCreate := 20
	def := map[string]int{"a": 1, "b": 2, "c": 3}

	mapPool := NewMap(capSave, capCreate, def)
	m := mapPool.Get()

	if m == nil {
		t.Fatal("Get returned nil")
	}

	if len(m) != len(def) {
		t.Errorf("Expected map length to be len(def), got %d", len(m))
	}

	for k, v := range def {
		if m[k] != v {
			t.Errorf("Expected map element with key %s to be %d, got %d", k, v, m[k])
		}
	}
}

func TestMapSave(t *testing.T) {
	capSave := 10
	capCreate := 20
	def := map[string]int{"a": 1, "b": 2, "c": 3}

	mapPool := NewMap(capSave, capCreate, def)
	m := mapPool.Get()

	// Изменяем map
	m["d"] = 4
	m["e"] = 5

	// Сохраняем map обратно в пул
	mapPool.Save(m)

	// Получаем новый map из пула
	newMap := mapPool.Get()

	if len(newMap) != len(def) {
		t.Errorf("Expected new map length to be %d, got %d", len(def), len(newMap))
	}

	for k, v := range def {
		if newMap[k] != v {
			t.Errorf("Expected new map element with key %s to be %d, got %d", k, v, newMap[k])
		}
	}

	// Проверяем, что map с превышающей длиной не сохраняется
	longMap := make(map[string]int, capSave+1)
	mapPool.Save(longMap)

	// Получаем новый map из пула
	anotherMap := mapPool.Get()

	if len(anotherMap) != len(def) {
		t.Errorf("Expected another map length to be %d, got %d", len(def), len(anotherMap))
	}

	for k, v := range def {
		if anotherMap[k] != v {
			t.Errorf("Expected another map element with key %s to be %d, got %d", k, v, anotherMap[k])
		}
	}
}

func TestMapParallel(t *testing.T) {
	capSave := 10
	capCreate := 20
	def := map[string]int{"a": 1, "b": 2, "c": 3}

	mapPool := NewMap(capSave, capCreate, def)

	t.Run("ParallelTest1", func(t *testing.T) {
		t.Parallel()
		m := mapPool.Get()
		if m == nil {
			t.Fatal("Get returned nil")
		}
		if len(m) != len(def) {
			t.Errorf("Expected map length to be len(def), got %d", len(m))
		}
		for k, v := range def {
			if m[k] != v {
				t.Errorf("Expected map element with key %s to be %d, got %d", k, v, m[k])
			}
		}
	})

	t.Run("ParallelTest2", func(t *testing.T) {
		t.Parallel()
		m := mapPool.Get()
		if m == nil {
			t.Fatal("Get returned nil")
		}
		if len(m) != len(def) {
			t.Errorf("Expected map length to be len(def), got %d", len(m))
		}
		for k, v := range def {
			if m[k] != v {
				t.Errorf("Expected map element with key %s to be %d, got %d", k, v, m[k])
			}
		}
	})
}
