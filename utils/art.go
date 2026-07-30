package utils

import (
	"fmt"
	"math"
	"os/user"
	"strings"
	"time"

	figure "github.com/common-nighthawk/go-figure"
)

func AnimateWelcome() {
	// get the username
	user, err := user.Current()
	username := "ENGINEEER"

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	username = strings.ToUpper(user.Username)

	// generate the dynamic ASCII art
	titleTxt := fmt.Sprintf("ZSH: %s", username)
	myFigure := figure.NewFigure(titleTxt,"starwars",true)
	bannerLines := myFigure.Slicify()

	// clear the screen and start the interactive prompt loop
	fmt.Print("\033[H\033[2J")

	// run the continous strobe
	runStrobeIntro(bannerLines)
}

func runStrobeIntro(lines []string) {
	//hide the cursor during the animation
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	// run the strobe for a short phase
	for frame := 0; frame < 40; frame++ {
		if frame > 0 {
			fmt.Printf("\033[%dA", len(lines)) // move the cursor back up to redraw
		}

		for i, line := range lines {
			frequency := 0.3

			r := int(math.Sin(frequency*float64(frame)+float64(i))*127+150)
			g := int(math.Sin(frequency*float64(frame)+float64(i)+2)*127+150)
			b := int(math.Sin(frequency*float64(frame)+float64(i)+4)*127+150)

			// print the current line with the calculated strobe color
			fmt.Printf("\033[K\033[38;2;%d;%d;%dm%s\033[0m\n",r,g,b,line)

		}
		time.Sleep(40 *time.Millisecond)
	}
	fmt.Println("\033[1;32m>> Custom Shell Initialization Complete.\033[0m")
}
