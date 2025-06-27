package types

type ParseResult struct {
	Data		interface{}
	FilePath	string
	Size		int64
	Type		string
	KeyCount	int
	ArrayLength	int
}

type ValidationResult struct {
	FilePath	string
	IsValid		bool
	Errors		[]string
	Warnings	[]string
}

type OutputOptions struct {
	Format		string
	Indent		int
	Colors		bool
	ShowKeys	bool
	ShowTypes	bool
	MaxDepth	int
	SortKeys	bool
	UseTabs		bool
}