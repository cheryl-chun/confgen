package codegen

type Options struct {
	PackageName string 
	AddComments bool   
}

func DefaultOptions() Options {
	return Options{
		PackageName: "main",
		AddComments: false,
	}
}