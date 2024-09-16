package main

import (
    "testing"
    "strings"
)


func TestStringsContains(t *testing.T) {
    str := "tarantool: tarantool error when sfdlkjadslfkjaslfkjalkfjlk jdflkjsaflkjas"
    substr := "tarantool error"
    
    res := strings.Contains(str, substr)
    if res != true {
        t.Errorf("must return true, actual %v", res)
    }
}

