package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

type field struct {
	k string
	v any
}
type tag struct {
	k string
	v string
}

func main() {
	ctx := context.Background()
	ctx.Value(context.Background())

	var arr []any

	arr = append(arr, &tag{k: "name", v: "John"})
	arr = append(arr, &field{k: "age", v: 30})
	arr = append(arr, fmt.Errorf("err"))

	for _, a := range arr {
		if v, ok := a.(*tag); ok {
			fmt.Println("tag", v)
		} else {
			fmt.Println("tag ok", ok)
		}
		if v, ok := a.(*field); ok {
			fmt.Println("field", v)
		} else {
			fmt.Println("field ok", ok)
		}
		if v, ok := a.(error); ok {
			fmt.Println("err", v)
		} else {
			fmt.Println("err ok", ok)
		}
	}
	_, _ := zap.NewProduction()
}
