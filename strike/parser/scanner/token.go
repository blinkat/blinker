package scanner

type Token struct {
	Type          int    //token type
	Value         string //token value
	Line          int
	Col           int
	Pos           int
	Endpos        int
	Nlb           bool
	CommentsBefor []*Token //comments
}
