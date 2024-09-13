package main_test


import (
    "testing"
)


type Caller[T any] interface {
    Call(string) (T, error)
}


type Caller[T any] func(name string) (T, error)


func (f Caller[T any]) Debug(name string) (T any) {
    return T
}

func TestInterfaceCasting(t *testing.T) {
}

