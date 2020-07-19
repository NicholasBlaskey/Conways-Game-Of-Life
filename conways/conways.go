package conways

import (
	"math/rand"
	"sync"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func GetPositions(numX, numY int) []mgl.Vec2 {
	translations := []mgl.Vec2{}
	xOffset := 1.0 / float32(numX)
	yOffset := 1.0 / float32(numY)
	for y := -numY; y < numY; y += 2 {
		for x := -numX; x < numX; x += 2 {
			translations = append(translations,
				mgl.Vec2{float32(x)/float32(numX) + xOffset,
					float32(y)/float32(numY) + yOffset})
		}
	}
	return translations
}

func CreateBoard(seed int64, numX, numY int) []float32 {
	rand.Seed(seed)
	board := []float32{}

	for i := 0; i < numX*numY; i++ {
		if rand.Float32() < 0.5 {
			board = append(board, 1)
		} else {
			board = append(board, 0)
		}
	}

	return board
}

func UpdateBoard(board []float32, numX, numY int) {
	helperBoard := [][]float32{}
	for i := 0; i < numY; i++ {
		row := []float32{}
		for j := 0; j < numX; j++ {
			row = append(row, board[i*numY+j])
		}
		helperBoard = append(helperBoard, row)
	}

	var wg sync.WaitGroup
	go handleEdgeCases(helperBoard, board, numX, numY, &wg)
	wg.Add(1)
	for i := 1; i < numY-1; i++ {
		wg.Add(1)
		go handleRow(helperBoard, board, numX, numY, i, &wg)
	}
	wg.Wait()
}

func handleRow(helperBoard [][]float32, board []float32,
	numX, numY, i int, wg *sync.WaitGroup) {

	defer wg.Done()
	for j := 1; j < numX-1; j++ {
		neighborCount := helperBoard[i][j-1] +
			helperBoard[i][j+1] +
			helperBoard[i-1][j] +
			helperBoard[i+1][j] +
			helperBoard[i-1][j-1] +
			helperBoard[i-1][j+1] +
			helperBoard[i+1][j-1] +
			helperBoard[i+1][j+1]
		if helperBoard[i][j] == 1.0 {
			if neighborCount < 2 || neighborCount > 3 {
				board[i*numY+j] = 0.0
			}
		} else {
			if neighborCount == 3 {
				board[i*numY+j] = 1.0
			}
		}
	}
}

func handleEdgeCases(helperBoard [][]float32, board []float32,
	numX, numY int, wg *sync.WaitGroup) {

	defer wg.Done()
	// Handle top row and bottom row
	for j := 0; j < numX; j += 1 {
		updateCellEdgeCase(helperBoard, board, numX, numY, 0, j)
		updateCellEdgeCase(helperBoard, board, numX, numY, numY-1, j)
	}

	// Handle left and right columns
	for i := 1; i < numY-1; i += 1 {
		updateCellEdgeCase(helperBoard, board, numX, numY, i, 0)
		updateCellEdgeCase(helperBoard, board, numX, numY, i, numX-1)
	}
}

func updateCellEdgeCase(helperBoard [][]float32, board []float32,
	numX, numY, i, j int) {

	// At edge case we wrap around
	neighborCount := helperBoard[i][(numX+j-1)%numX] +
		helperBoard[i][(j+1)%numX] +
		helperBoard[(numY+i-1)%numY][j] +
		helperBoard[(i+1)%numY][j] +
		helperBoard[(numY+i-1)%numY][(numX+j-1)%numX] +
		helperBoard[(numY+i-1)%numY][(j+1)%numX] +
		helperBoard[(i+1)%numY][(numX+j-1)%numX] +
		helperBoard[(i+1)%numY][(j+1)%numX]

	if helperBoard[i][j] == 1.0 {
		if neighborCount < 2 || neighborCount > 3 {
			board[i*numY+j] = 0.0
		}
	} else {
		if neighborCount == 3 {
			board[i*numY+j] = 1.0
		}
	}
}
