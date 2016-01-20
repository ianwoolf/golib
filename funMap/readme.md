registered func and call by name

    funcs = NewFuncs(100)

        testcases = map[string]interface{}{
                    "hello":      func() { print("hello") },
                                "foobar":     func(a, b, c int) int { return a + b + c },
                                    }



    for k, v := range testcases {
                err := funcs.Bind(k, v)
                                if err != nil {
                                                t.Error("Bind %s: %s", k, err)
                                                            }
                        }



    if _, err := funcs.Call("foobar"); err == nil {
                log.Error("Call %s: %s", "foobar", "error happen.")
                        }
    if _, err := funcs.Call("foobar", 0, 1, 2); err != nil {
                log.Error("Call %s: %s", "foobar", err)
                        }

