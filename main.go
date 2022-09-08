package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

// // define the board and the snake
type game_state_t struct {
	num_rows uint
	width    uint
	height   uint
	board    *[][]string
	over     bool
	snake    *snake
	prevdir  rune
	food     bool
}
type snake struct {
	tail_row uint
	tail_col uint
	body     [][]uint
	head_row uint
	head_col uint
	live     uint
	length   uint
}

func next_square(gs *game_state_t) (s string) {
	/*
	   returns the string that the snake is moving into
	*/
	scol := gs.snake.head_col
	srow := gs.snake.head_row
	nextrow := get_next_row(srow, (*gs.board)[srow][scol])
	nextcol := get_next_col(scol, (*gs.board)[srow][scol])
	return (*gs.board)[nextrow][nextcol]
}

func get_next_col(cur_col uint, c string) uint {
	if c == ">" {
		return cur_col + 1
	} else if c == "<" {
		return cur_col - 1
	}
	return cur_col

}

func get_next_row(cur_row uint, c string) uint {
	if c == "^" {
		return cur_row - 1
	} else if c == "v" {
		return cur_row + 1
	}
	return cur_row
}

func SetUp(game *game_state_t) *game_state_t {
	// The fruit is at row 2, column 9 (zero-indexed).
	// The tail is at row 2, column 2,
	//  and the head is at row 2, column 4.
	// **> o  <>^>
	board := make([][]string, (*game).height, (*game).height)
	sbody := make([][]uint, 1)

	for i := uint(0); uint(i) < (*game).height; i++ {
		board[i] = make([]string, game.width)
	}

	for i := uint(0); uint(i) < 1; i++ {
		sbody[i] = make([]uint, 2)
	}
	for i, row := range board {
		for j := uint(0); j < (*game).width; j++ {
			if i == 0 || i == int(game.height)-1 || j == 0 || j == uint(game.width)-1 {
				row[j] = "#"
			} else {
				row[j] = " "
			}

		}

	}
	sbody[0][0] = 1
	sbody[0][1] = 3
	game.board = &board
	theSnake := &snake{
		tail_row: uint(1),
		tail_col: uint(2),
		body:     sbody,
		head_row: uint(1),
		head_col: uint(4),
	}
	game.food = true
	board[1][2] = "x"
	board[1][3] = "o"
	board[1][4] = ">"
	board[3][5] = "*"

	game.snake = theSnake
	game.board = &board
	return game
}
func Draw(game game_state_t) {
	board := game.board
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	// draw body
	for _, v := range (game).snake.body {
		(*game.board)[v[0]][v[1]] = "o"
	}
	//tail
	(*game.board)[game.snake.tail_row][game.snake.tail_col] = "x"
	for _, row := range *board {
		for _, v := range row {
			fmt.Printf(v)
		}
		fmt.Println()
	}
	// (*game.board)[game.snake.head_row][game.snake.head_col] = ">"
	// (*game.board)[game.snake.tail_row][game.snake.tail_col-1] = " "

}

func Input(input *chan rune) {
	for {
		ru, _, err := keyboard.GetKey()
		if err != nil {
			panic("Getting keys was wrong ")
		}
		*input <- ru
	}

}
func Check_Collion(game *game_state_t) bool {
	check := next_square(game)
	if check == "o" || check == "#" {
		fmt.Println("Collison")
		fmt.Println(check)
		Quit()
	}

	if check == "*" {
		game.food = false
		return true
	} else {
		return false
	}

}

func randomfood(game *game_state_t) {
	rand.Seed(time.Now().UnixNano())
	for game.food == false {
		xmin := uint(1)
		xmax := game.height + uint(1)
		ymin := uint(1)
		ymax := game.width - uint(1)

		y := rand.Intn(int(ymax-ymin)) + int(ymin)
		x := rand.Intn(int(xmax-xmin)) + int(xmin)
		if (*game.board)[x][y] != "o" && (*game.board)[x][y] != ">" && (*game.board)[x][y] != "x" {
			(*game.board)[x][y] = "*"
			game.food = true
		}

	}

}
func Logic(direction rune, game *game_state_t) {
	prevcol := (*game.snake).tail_col
	prevrow := (*game.snake).tail_row
	if !game.food {
		randomfood(game)
	}
	if direction == 119 && game.prevdir == 115 {
		direction = 115
	}
	if direction == 115 && game.prevdir == 119 {
		direction = 119
	}
	if direction == 97 && game.prevdir == 100 {
		direction = 100
	}
	if direction == 100 && game.prevdir == 97 {
		direction = 97
	}
	switch check := direction; check {
	case 119:
		game.prevdir = 119
		fmt.Println("UP")
		(*game.board)[game.snake.head_row][game.snake.head_col] = "^"
		grow := Check_Collion(game)
		if grow {
			game.snake.body = append(game.snake.body, []uint{game.snake.head_row, game.snake.head_col})
			(*game.snake).head_row -= 1
			(*game.board)[game.snake.head_row+1][game.snake.head_col] = "o"
			(*game.board)[game.snake.head_row][game.snake.head_col] = "o"

		} else {
			/// body tail become last place of body
			(*game.snake).tail_row = game.snake.body[0][0]
			(*game.snake).tail_col = game.snake.body[0][1]
			for i, v := range (*game).snake.body {
				if i == len(game.snake.body)-1 {
					(*game).snake.body[i][0] = game.snake.head_row
					(*game).snake.body[i][1] = game.snake.head_col
				} else {
					for j, _ := range v {
						(*game).snake.body[i][j] = (*game).snake.body[i+1][j]
					}
				}
			}
			game.snake.head_row -= 1
			(*game.board)[game.snake.head_row][game.snake.head_col] = "^"
			(*game.board)[prevrow][prevcol] = " "

		}

	case 97:
		game.prevdir = 97

		fmt.Println("Left")
		(*game.board)[game.snake.head_row][game.snake.head_col] = "<"
		grow := Check_Collion(game)
		if grow {
			game.snake.body = append(game.snake.body, []uint{game.snake.head_row, game.snake.head_col})
			(*game.snake).head_col -= 1
			(*game.board)[game.snake.head_row][game.snake.head_col] = "<"

		} else {
			/// body tail become last place of body
			(*game.snake).tail_row = game.snake.body[0][0]
			(*game.snake).tail_col = game.snake.body[0][1]
			for i, v := range (*game).snake.body {
				if i == len(game.snake.body)-1 {
					(*game).snake.body[i][0] = game.snake.head_row
					(*game).snake.body[i][1] = game.snake.head_col
				} else {
					for j, _ := range v {
						(*game).snake.body[i][j] = (*game).snake.body[i+1][j]
					}
				}
			}
			(*game.snake).head_col -= 1
			(*game.board)[game.snake.head_row][game.snake.head_col] = "<"
			(*game.board)[prevrow][prevcol] = " "

		}

	case 100:
		fmt.Println("Right")
		game.prevdir = 100
		(*game.board)[game.snake.head_row][game.snake.head_col] = ">"
		grow := Check_Collion(game)
		if grow {
			game.snake.body = append(game.snake.body, []uint{game.snake.head_row, game.snake.head_col})
			(*game.snake).head_col += 1
			(*game.board)[game.snake.head_row][game.snake.head_col] = ">"

		} else {
			/// body tail become last place of body
			(*game.snake).tail_row = game.snake.body[0][0]
			(*game.snake).tail_col = game.snake.body[0][1]
			for i, v := range (*game).snake.body {
				if i == len(game.snake.body)-1 {
					(*game).snake.body[i][0] = game.snake.head_row
					(*game).snake.body[i][1] = game.snake.head_col
				} else {
					for j, _ := range v {
						(*game).snake.body[i][j] = (*game).snake.body[i+1][j]
					}
				}
			}
			(*game.snake).head_col += 1
			(*game.board)[game.snake.head_row][game.snake.head_col] = ">"
			(*game.board)[prevrow][prevcol] = " "

		}
	case 115:
		game.prevdir = 115
		grow := Check_Collion(game)
		if grow {
			game.snake.body = append(game.snake.body, []uint{game.snake.head_row, game.snake.head_col})
			game.snake.head_row += 1
			(*game.board)[game.snake.head_row][game.snake.head_col] = "v"
		} else {
			/// body tail become last place of body
			(*game.snake).tail_row = game.snake.body[0][0]
			(*game.snake).tail_col = game.snake.body[0][1]
			for i, v := range (*game).snake.body {
				if i == len(game.snake.body)-1 {
					(*game).snake.body[i][0] = game.snake.head_row
					(*game).snake.body[i][1] = game.snake.head_col
				} else {
					for j, _ := range v {
						(*game).snake.body[i][j] = (*game).snake.body[i+1][j]
					}
				}
			}
			game.snake.head_row += 1
			(*game.board)[game.snake.head_row][game.snake.head_col] = "v"
			(*game.board)[prevrow][prevcol] = " "

		}

	case 113: // Q
		Quit()
	}

}

func Quit() {
	fmt.Println("Quit")
	os.Exit(1)
}
func main() {
	game := game_state_t{
		height: 30,
		width:  30,
	}
	// open key board
	keyboard.Open()
	//close keyboard
	defer keyboard.Close()
	// set up chan
	userInput := make(chan rune)
	gs := SetUp(&game)
	// Go routine to poll keys
	go Input(&userInput)
	// time.Sleep(20 * time.Second)
	for !game.over {
		Draw(*gs)
		select {
		case dir, ok := <-userInput:
			if ok {
				Logic(dir, gs)
				time.Sleep(100 * time.Millisecond)
			} else {
			}
		default:
			Logic(gs.prevdir, gs)
			time.Sleep(100 * time.Millisecond)
		}
	}

}
