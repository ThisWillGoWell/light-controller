package tetris

type BoardBlock struct {
}

type Board struct {
	Board    [][]BoardBlock
	TopLevel []int
	NumRows  int
	NumCols  int
}

func (b *Board) TopLevel() {

}
