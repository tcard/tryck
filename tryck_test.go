package tryck

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func fact(n int) int {
	if n == 0 {
		return 1
	}
	return n * fact(n-1)
}

func TestNoErrors(t *testing.T) {
	err := TryCatch(func(try TryFunc) {
		try(fact(5))
	}, func(err error) bool {
		if err != nil {
			t.FailNow()
		}
		return true
	})
	if err != nil {
		t.FailNow()
	}
}

func TestErrorAtThirdTry(t *testing.T) {
	err := TryCatch(func(try TryFunc) {
		try(fact(10))
		try()
		try(errors.New("stop"))
		try(errors.New("shouldn't be reached"))
	}, nil)
	if err == nil {
		t.FailNow()
	}
	if tryErr, ok := err.(TryError); !ok || tryErr.Nth != 3 || tryErr.Error() != "stop" {
		t.FailNow()
	}
}

func TestCustomCatch(t *testing.T) {
	sideEffect := false
	err := TryCatch(func(try TryFunc) {
		try(errors.New("OK"))
		try(errors.New("stop"))
		try(errors.New("shouldn't be reached"))
	}, func(err error) bool {
		sideEffect = true
		if err.Error() == "stop" {
			return false
		}
		return true
	})
	if err.Error() != "stop" || !sideEffect {
		t.FailNow()
	}
}

func TestPanic(t *testing.T) {
	panicMsg := ""
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("recovered panic shouldn't be nil.")
		}
		p, ok := r.(TryPanic)
		if !ok {
			t.Fatal("recovered panic should be a TryPanic")
		}
		if !strings.HasPrefix(p.String(), fmt.Sprintf("TryCatch panic: %v\n", panicMsg)) {
			t.Fatal("incorrect recovered panic:", p)
		}
	}()
	TryCatch(func(try TryFunc) {
		defer func() {
			r := recover()
			panicMsg = fmt.Sprintf("%v", r)
			panic(r)
		}()
		_ = []int{}[1]
	}, nil)
}

func TestNested(t *testing.T) {
	err := TryCatch(func(try TryFunc) {
		err := TryCatch(func(try TryFunc) {
			try()
			try(errors.New("stop"))
			try(errors.New("should't be reached"))
		}, nil)
		if tryErr, ok := err.(TryError); !ok || tryErr.Nth != 2 || tryErr.Error() != "stop" {
			t.FailNow()
		}
	}, func(err error) bool {
		t.FailNow()
		return false
	})
	if err != nil {
		t.FailNow()
	}
}

func BenchmarkFactWithoutTryck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fact(99999)
	}
}

func BenchmarkFactWithNoTryCalls(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TryCatch(func(try TryFunc) {
			fact(99999)
		}, nil)
	}
}

func factWithTry(n int, try TryFunc, nTry int) int {
	for i := 0; i < nTry; i++ {
		try("not an error")
	}
	if n == 0 {
		return 1
	}
	return n * fact(n-1)
}

func BenchmarkFactWith1TryCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TryCatch(func(try TryFunc) {
			factWithTry(99999, try, 1)
		}, nil)
	}
}

func BenchmarkFactWith100TryCalls(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TryCatch(func(try TryFunc) {
			factWithTry(99999, try, 100)
		}, nil)
	}
}
