package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	// Define styled printers
	successSymbol = color.New(color.FgGreen, color.Bold).SprintFunc()("✓")
	errorSymbol   = color.New(color.FgRed, color.Bold).SprintFunc()("✗")
	infoSymbol    = color.New(color.FgBlue, color.Bold).SprintFunc()("ℹ")
	warnSymbol    = color.New(color.FgYellow, color.Bold).SprintFunc()("!")

	successText = color.New(color.FgGreen).SprintfFunc()
	errorText   = color.New(color.FgRed).SprintfFunc()
	infoText    = color.New(color.FgBlue).SprintfFunc()
	warnText    = color.New(color.FgYellow).SprintfFunc()
	dimText     = color.New(color.Faint).SprintfFunc()
	boldText    = color.New(color.Bold).SprintfFunc()
)

// Success prints a success message
func Success(message string) {
	fmt.Printf("%s %s\n", successSymbol, successText(message))
}

// Error prints an error message
func Error(message string) {
	fmt.Printf("%s %s\n", errorSymbol, errorText(message))
}

// Info prints an info message
func Info(message string) {
	fmt.Printf("%s %s\n", infoSymbol, infoText(message))
}

// Warning prints a warning message
func Warning(message string) {
	fmt.Printf("%s %s\n", warnSymbol, warnText(message))
}

// Section prints a section header
func Section(title string) {
	fmt.Printf("\n%s\n%s\n", 
		boldText(title),
		dimText(strings.Repeat("─", len(title))))
}

// CodeBlock prints text in a subtle code block style
func CodeBlock(text string) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	maxLength := 0
	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}

	fmt.Printf("%s\n", dimText("┌"+strings.Repeat("─", maxLength+2)+"┐"))
	for _, line := range lines {
		padding := strings.Repeat(" ", maxLength-len(line))
		fmt.Printf("%s %s%s %s\n",
			dimText("│"),
			line,
			padding,
			dimText("│"))
	}
	fmt.Printf("%s\n", dimText("└"+strings.Repeat("─", maxLength+2)+"┘"))
}

// Prompt asks for user input with styling
func Prompt(question string) string {
	fmt.Print(boldText(question + " "))
	var response string
	fmt.Scanln(&response)
	return response
}

// Bold wraps text in bold styling
func Bold(text string) string {
	return boldText(text)
}