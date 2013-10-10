tryck
=====

Something similar to exceptions for Go as found in other languages.

Usually, with 'aggregate' functions like `DoThisAndThat()`, the only error handling you want is returning any errors found, perhaps logging or something simple like that along the way. This can populate the code with repetitive `if err != nil` checks. tryck is an attempt to reproduce try/catch using core language features (closures, and panic/recover) and without sacrificing fine grained control; you can still do the in-place error checking `if` when needed inside a tryck block.

This is an experiment which may come handy in some situations. Please don't use this on code that the general public may need to hack on. It's not idiomatic Go and adds unnecesary cognitive load to your code.

Take a look at [the documentation](http://godoc.org/github.com/tcard/tryck).

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