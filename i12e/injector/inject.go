package injector

import (
	"github.com/codegangsta/inject"
)

// Inject 注入函数
func Inject(a interface{}, params ...interface{}) {
	ior := inject.New()
	for _, v := range params {
		ior.Map(v)
	}
	err := ior.Apply(a)
	if err != nil {
		panic(err)
	}
}
