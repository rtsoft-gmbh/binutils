package decode

import "testing"

func TestReader(t *testing.T) {
	t.Run("positive scenario: read bytes", func(t *testing.T) {
		r := reader{
			payload: []byte{1, 2, 3, 4, 5},
		}
		if v, e := r.getBytes(1); e != nil || len(v) != 1 || v[0] != 1 {
			t.Error("result not expected")
		}
		if v, e := r.getBytes(2); e != nil || len(v) != 2 || v[0] != 2 || v[1] != 3 {
			t.Error("result not expected")
		}
		if v, e := r.getBytes(2); e != nil || len(v) != 2 || v[0] != 4 || v[1] != 5 {
			t.Error("result not expected")
		}
	})
	t.Run("negative scenario: read bytes more then reader contains. Expected ErrorOutOfIndex", func(t *testing.T) {
		r := reader{
			payload: []byte{1, 2},
		}
		var e error
		if _, e = r.getBytes(3); e == nil {
			t.Error("result not expected")
		}
		if e != ErrorOutOfIndex {
			t.Error("expected ErrorOutOfIndex")
		}
	})
}
