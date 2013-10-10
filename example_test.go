package tryck_test

import (
	"errors"
	"fmt"
	. "github.com/tcard/tryck"
)

func ExampleTryCatch() {
	var fact func(n int) (int, error)
	fact = func(n int) (int, error) {
		if n < 0 {
			return 0, errors.New("n should be >= 0")
		} else if n == 0 {
			return 1, nil
		}
		i, err := fact(n - 1)
		return n * (i - 1), err
	}
	okError := func() error {
		return errors.New("OK")
	}
	err := TryCatch(func(try TryFunc) {
		try(fact(10))
		try(okError())
		try(fact(-1))
		fmt.Println("Shouldn't be reached!")
	}, func(err error) bool {
		fmt.Printf("Entered catch block with error '%v';", err)
		if err == nil || err.Error() == "OK" {
			fmt.Println(" continuing.")
			return true
		}
		fmt.Println(" stopping.")
		return false
	})
	fmt.Printf("Stopped at error '%v' in the %d'th try().", err, err.(TryError).Nth)

	// Output:
	// Entered catch block with error '<nil>'; continuing.
	// Entered catch block with error 'OK'; continuing.
	// Entered catch block with error 'n should be >= 0'; stopping.
	// Stopped at error 'n should be >= 0' in the 3'th try().
}
