package conways

import(
	"math/rand"
	
	mgl "github.com/go-gl/mathgl/mgl32"
)

func GetPositions(num_x, num_y int) []mgl.Vec2 {
	translations := []mgl.Vec2{}
	xOffset := 1.0 / float32(num_x)
	yOffset := 1.0 / float32(num_y)
	for y := -num_y; y < num_y; y += 2 {
		for x := -num_x; x < num_x; x += 2 {
			translations = append(translations,
				mgl.Vec2{float32(x) / float32(num_x) + xOffset,
					float32(y) / float32(num_y) + yOffset})
		}
	}

	return translations
}

func CreateBoard(seed int64, num_x, num_y int) []float32 {
	rand.Seed(seed)
	board := []float32{}

	for i := 0; i < num_x * num_y; i++ {
		if rand.Float32() < 0.5 {
			board = append(board, 1)
		} else {
			board = append(board, 0)
		}
	}

	return board
}

func UpdateBoard(board []float32, num_x, num_y int) {
	helperBoard := [][]float32{}
	for i := 0; i < num_y; i++ {
		row := []float32{}
		for j := 0; j < num_x; j++ {
			row = append(row, board[i * num_y + j])
		}
		helperBoard = append(helperBoard, row)
	}

	for i := 0; i < num_y; i++ {
		for j := 0; j < num_x; j++ {
			// Using x and y wrap around
			neighborCount := helperBoard[i][(num_x + j - 1) % num_x] +
				helperBoard[i][(j + 1) % num_x] +
				helperBoard[(num_y + i - 1) % num_y][j] +
				helperBoard[(i + 1) % num_y][j] +
				helperBoard[(num_y + i - 1) % num_y][(num_x + j - 1) % num_x] +
				helperBoard[(num_y + i - 1) % num_y][(j + 1) % num_x] +
				helperBoard[(i + 1) % num_y][(num_x + j - 1) % num_x] +
				helperBoard[(i + 1) % num_y][(j + 1) % num_x]
			if helperBoard[i][j] == 1.0 {
				if neighborCount < 2 || neighborCount > 3 {
					board[i * num_y + j] = 0.0
				} 
			} else {
				if neighborCount == 3 {
					board[i * num_y + j] = 1.0
				}
			}			
		}
	}
}
