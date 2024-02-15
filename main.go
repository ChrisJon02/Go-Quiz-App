package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func problemPuller(fileName string) ([]problem, error) {
	// Read all the problems from the quiz.csv

	//1. Open the file
	if fileObj, err := os.Open(fileName); err == nil {

		//2. Create a new reader
		csvRead := csv.NewReader(fileObj)

		//3. Read the file
		if cLines, err := csvRead.ReadAll(); err == nil {

			//4. Call the parseProblem function
			return parseProblem(cLines), nil

		} else {
			return nil, fmt.Errorf("error in reading data in csv"+"format from %s file; %s", fileName, err.Error())
		}

	} else {
		return nil, fmt.Errorf("error in opening %s file; %s", fileName, err.Error())
	}

}

func main() {

	//1. Input the name of the file
	fName := flag.String("f", "quiz.csv", "path of the csv file")

	//2. Set the duration of the timer
	timer := flag.Int("t", 30, "timer for the quiz")
	flag.Parse()

	//3. Pull the problems from the file ( call the problem puller function)
	problems, err := problemPuller(*fName)

	//4. Handle the error
	if err != nil {
		exit(fmt.Sprintf("Something went wrong : %s", err.Error()))
	}

	//5. Create a variable to count our correct answers
	correctAns := 0

	//6. Using the duration of the timer, we must initialize the timer
	timerInit := time.NewTimer(time.Duration(*timer) * time.Second)
	ansChan := make(chan string)

	//7. Loop through the problems, print the questions, will accept the answers
problemLoop:

	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d : %s = ", i+1, p.ques)

		go func() {
			fmt.Scanf("%s", &answer)
			ansChan <- answer
		}()
		select {
		case <-timerInit.C:
			fmt.Println()
			break problemLoop
		case iAns := <-ansChan:
			if iAns == p.ans {
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansChan)
			}
		}
	}
	//8. We'll calculate the answers and print the result
	fmt.Printf("Your result is %d out of %d\n", correctAns, len(problems))
	fmt.Printf("Press enter to exit")
	<-ansChan
}

func parseProblem(lines [][]string) []problem {
	// Go over the lines and parse them, with problem struct
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{ques: lines[i][0], ans: lines[i][1]}
	}
	return r
}

type problem struct {
	ques string
	ans  string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
