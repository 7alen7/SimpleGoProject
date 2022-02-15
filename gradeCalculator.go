/*	CS424 Spring 2022
	Prog_Asgn1.go
	Dan Lenz
	A program to take in an input file, grade weight percents, grade count, and calculate the average grades for a set of students

	This program was tested in a Windows 11 environment. It was run using the command line found in the Visual Studio Code program.
------------------------------------------------------------------------------------------------*/

package main

import (
	/*
		Inports required to run program
	*/
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
	Student structure, contains information about each student

	firstName: Student's firstname
	lastName: Student's last name
	testAvg: Students average test score
	hwAvg: Student's average HW score
	testSubmitted: number of tests submitted
	hwSubmitted: number of HWs submitted
*/
type Student struct {
	firstName string
	lastName  string
	testAvg   float32
	hwAvg     float32
	totalAvg  float32

	testSubmitted int
	hwSubmitted   int
}

/*
	Global variable declarations. These are used due to some errors that were propogating in their previous implementation

	testPercent: contains the percent that tests will make up of your total grade
	studentCount: number of students in input file
	testMax: The maximum number of tests expected
	hwMax: the maximum number of HW assignments expected
*/
var testPercent float32
var studentCount int
var testMax int
var hwMax int

/*
	Main function, handles function calls and file I/O details
*/
func main() {

	var keyboard *bufio.Scanner //Scanner to check for inputs
	var inputFile string        // input file name
	var inFile *os.File         // input file reference
	var estat error             // error message, if the file cannot open

	var studentSlice []Student // Slice that will contain the students from the input file

	var currentStudent Student // The current student we're working with
	keyboard = bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the gradebook calculator test program. I am going to read students from an input data file. You will tell me the name of your input file\n")

	fmt.Println("Enter the name of your input file: ")
	keyboard.Scan()

	inputFile = keyboard.Text()

	inFile, estat = os.Open(inputFile)

	if estat != nil { // Checks to see if file can be opened
		fmt.Println("This file was not able to be opened. Try again later")
	} else {
		// If the file can be opened, the user is asked to input their information
		fmt.Println("Enter the % amount to weight test in overall avg: ")
		keyboard.Scan()

		testPercentOne, _ := strconv.ParseFloat(keyboard.Text(), 32) // Test percent as a float
		testPercent = float32(testPercentOne)                        // Test percent that is used
		hwPercent := 100 - testPercent                               // HW percent is equal to 100% - testPercent

		s := fmt.Sprintf("\nTests will be weighted %.1f%%, Homework weighted %.1f%%", testPercent, hwPercent)
		fmt.Println(s + "\n")

		fmt.Println("How many homework assignments are there? ") // User is asked for the number of test and hw assignments expected
		keyboard.Scan()
		hwMax, _ = strconv.Atoi(keyboard.Text()) // number of hw assignments expected
		fmt.Println("How many test grades are there? ")
		keyboard.Scan()
		testMax, _ = strconv.Atoi(keyboard.Text()) // number of tests expected
	}
	defer inFile.Close() // will close the file when done

	studentSlice = scanStudents(inFile, currentStudent, studentSlice) // Calls function scanStudents, is used to add students from file into the slice

	// This will sort the student slice based on the lastName variable in the Student structure
	sort.Slice(studentSlice, func(i, j int) bool {
		return studentSlice[i].lastName < studentSlice[j].lastName
	})

	// This function is called to print out the information collected from the input file.
	printData(studentSlice)
	// Program ends here
}

/*
	scanStudents function

	inputs: inputFile reference, Student object, []Student slice
	outputs: altered slice of Student objects that contains their relevant information
*/
func scanStudents(inputFile *os.File, student Student, studentSlice []Student) []Student {
	fileWords := bufio.NewScanner(inputFile) // new scanner created in this function
	for fileWords.Scan() {                   // for all words in the file. stops when end of file is reached
		text := strings.Split(fileWords.Text(), " ") // gets the new line and splits it
		studentCount++                               // increase the number of students in the list
		student.firstName = text[0]
		student.lastName = text[1]

		fileWords.Scan() // scan the next line

		text = strings.Split(fileWords.Text(), " ") // split next line

		/*
			This section will be repeated, but for homeworks, in the outlined section below
		*/
		var runningTotal int = 0          // total test grade
		var number int = 0                // new test score to be added to total
		student.testSubmitted = len(text) // contains the total number of tests submitted
		i := 0
		for i < len(text) { // increases running total
			number, _ = strconv.Atoi(text[i])
			runningTotal = runningTotal + number
			i++
		}
		student.testAvg = float32(runningTotal) / float32(len(text)) // divide running total by total number of tests

		fileWords.Scan() // get next line

		/*
			This section repeats the previous area, but now for homeworks instead of test assignments
		*/
		text = strings.Split(fileWords.Text(), " ")
		student.hwSubmitted = len(text)
		runningTotal = 0
		i = 0
		for i < len(text) {
			number, _ := strconv.Atoi(text[i])
			runningTotal = runningTotal + number
			i++
		}
		student.hwAvg = float32(runningTotal) / float32(len(text)) // updates the hw average for the student

		student.totalAvg = gradeAverage(student)     // updates the total average of the student
		studentSlice = append(studentSlice, student) // adds student to the end of the slice
	}
	return studentSlice // return the newly altered slice
}

/*
	This function calculates the total average for a Student object

	input: Student object
	output: float32 number that is equal to the total average for the student
*/
func gradeAverage(student Student) float32 {
	totalVal := (student.testAvg*(testPercent/100.0) + (student.hwAvg * ((100.0 - testPercent) / 100.0))) // Get total value
	return totalVal                                                                                       // return total value
}

/*
	This function is used to print out the data contained within a []Student object

	input: []Student
	output: only text in the terminal. No return value.
*/
func printData(studentSlice []Student) {
	fmt.Printf("GRADE REPORT --- %d STUDENTS FOUND IN FILE\n", studentCount) // Formatting for the output
	fmt.Printf("TEST WEIGHT: %.1f%%\n", testPercent)

	hwPercent := 100.0 - testPercent // HW percentage is equal to 100% - testPercent

	fmt.Printf("HOMEWORK WEIGHT: %.1f%%\n", hwPercent)
	classAvg := classAverage(studentSlice) // calls function to compute the class average.
	fmt.Printf("OVERALL AVERAGE is %.1f\n\n", classAvg)

	// More formatting
	fmt.Println("\tSTUDENT NAME\t:\tTESTS\t\tHOMEWORKS\tAVG")
	fmt.Println("---------------------------------------------------------------------------")

	/*
		For the length of studentSlice, print out the students information. If the student is missing homework or tests, create a notice next to their name.
	*/
	i := 0
	for i < len(studentSlice) {
		firstLast := studentSlice[i].lastName + ", " + studentSlice[i].firstName
		fmt.Printf("%23v :\t%5.1f (%d)\t%5.1f (%d)\t%4.1f", firstLast, studentSlice[i].testAvg, studentSlice[i].testSubmitted, studentSlice[i].hwAvg, studentSlice[i].hwSubmitted, studentSlice[i].totalAvg)
		if (studentSlice[i].hwSubmitted < hwMax) && (studentSlice[i].testSubmitted < testMax) {
			fmt.Print("\t** May be missing a homework and test **\n")
		} else if studentSlice[i].hwSubmitted < hwMax {
			fmt.Print("\t** May be missing a homework **\n")
		} else if studentSlice[i].testSubmitted < testMax {
			fmt.Print("\t** May be missing a test **\n")
		} else {
			fmt.Print("\n")
		}
		i++ // Move on to the next student
	}

}

/*
	This function calculates the average score for the class

	input: []Student object
	output: float32 value equal to the class' average grade
*/
func classAverage(studentSlice []Student) float32 {
	var runningTotal float32 // running total will be increased with each student examined. Contains the combination of all their grades.

	/*
		For each student, add their total grade to the running total
	*/
	i := 0
	for i < len(studentSlice) {
		runningTotal = runningTotal + studentSlice[i].totalAvg
		i++
	}
	classAvg := runningTotal / float32(len(studentSlice)) // Divide the running total by the number of students in the class
	return classAvg                                       // return the class average (float32)
}
