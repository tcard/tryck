// Package tryck provides something similar to exceptions as found in other languages.
//
// Usually, with 'aggregate' functions like DoThisAndThat(), the only error handling you want
// is returning any errors found, perhaps logging or something simple like that along the way. This can populate
// the code with repetitive 'if err != nil' checks. tryck is an attempt to reproduce try/catch using
// core language features (closures, and panic/recover) and without sacrificing fine grained control; you can
// still do the in-place error checking 'if' when needed inside a tryck block.
//
// This is an experiment which may come handy in some situations. Please don't use this on code that
// the general public may need to hack on. It's not idiomatic Go and adds unnecesary cognitive load to
// your code.
package tryck

import (
	"fmt"
	"runtime/debug"
)

// Wraps an error, indicating in Nth in which call to try it was encountered. It's an error type itself.
type TryError struct {
	Err error
	Nth int
}

func (err TryError) Error() string {
	return err.Err.Error()
}

// Wraps a panic that happens inside a TryCatch. We need this
type TryPanic struct {
	Panic interface{}
	Stack []byte
}

func (p TryPanic) String() string {
	return fmt.Sprintf("TryCatch panic: %v\n%s", p.Panic, p.Stack)
}

// This is the type of the 'try' function that will be passed to a tryBlock. It's last argument
// should be an error; in other case, it has no effect. See TryCatch.
type TryFunc func(...interface{})

// TryCatch executes the tryBlock function passed as argument. Each call to the tryBlocks's 'try' argument
// will pass it's last argument, if it is a non-nil error, to the catch function. The catch function
// determines if the tryBlock execution should be terminated at that point. TryCatch returns the last
// error for which catch returned false (ie. the error that stopped the execution), wrapped in a TryError.
//
// The catch argument can be nil; in that case, a standard 'if err != nil { return false }; return true' will be used.
//
// TryCatch blocks can be nested.
//
// Any panic encountered inside a TryCatch will be recovered and re-panicked. So that the original stack trace doesn't
// get lost in the process, the panic gets wrapped in a TryPanic which saves it.
//
// A limitation is that 'try' is not concurrent-safe. You shouldn't call 'try' from a goroutine other than the one in which 'TryCatch'
// was called.
func TryCatch(tryBlock func(try TryFunc), catch func(error) bool) error {
	err := error(nil)
	p := make(chan interface{})
	sig := p // A panic signal we're sure 'try' won't recover by accident.
	if catch == nil {
		catch = func(e error) bool {
			if e != nil {
				return false
			}
			return true
		}
	}
	go func() {
		defer func() {
			defer close(p)
			if r := recover(); r != nil && r != sig {
				p <- TryPanic{r, debug.Stack()}
			}
			p <- nil
		}()
		tryCounter := 0
		tryBlock(func(v ...interface{}) {
			tryCounter++
			if len(v) < 1 {
				return
			}
			if e, ok := v[len(v)-1].(error); !ok || catch(e) {
				return
			} else {
				err = TryError{e, tryCounter}
				panic(sig)
			}
		})
	}()
	if v := <-p; v != nil {
		panic(v)
	}
	return err
}
