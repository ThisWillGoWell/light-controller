package tetris

// golang imp of https://github.com/LeeYiyuan/tetrisai/blob/gh-pages/js/ai.js
type AI struct {
	heightWeight    float64
	linesWeight     float64
	holesWeight     float64
	bumpinessWeight float64
}

func (*AI) Best(working)
