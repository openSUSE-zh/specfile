package specfile

// Specfile specfile struct
type Specfile struct {
	Subpackage []Tag
	Tags       []Tag
	Macros     Macros
	Parts      []Tag
}
