package tetris

type MaskType int
const (
	ThreeByTwo
	TwoByTwo

)

type Piece struct {
	Mask     []bool
	Size     int
	Rotation int
}

func SquarePiece() Piece {
	return Piece{
		Size: 2,
		Mask: []bool{
			true, true,
			true, true,
		},
	}
}

func LPiece() Piece {
	return Piece{
		Size: 3,
		Mask: []bool{
			true, true, true,
			true, false, false,
		},
	}
}

func JPiece() Piece {
	return Piece{
		Size: 3,
		Mask: []bool{
			true, false, false,
			true, true, true,
			false, false, false,
		},
	}
}

func ZPiece() Piece {
	return Piece{
		Size: 3,
		Mask: []bool{
			true, true, false,
			false, true, true,
			false, false, false,
		},
	}
}
