package readertools

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func AskYesNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt + " [Yy/Nn]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка ввода. Попробуйте ещё раз.")
			continue
		}

		input = strings.TrimSpace(strings.ToLower(input))

		if input == "y" {
			return true
		} else if input == "n" || input == "" {
			return false
		} else {
			fmt.Println("Пожалуйста, введите 'Yy' или 'Nn'.")
		}
	}
}

func ReadLine(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(input)
}

func ReadPath(prompt string) string {
	path := ReadLine(prompt)
	path = strings.Trim(path, "/\\")
	return path
}

func ChooseFromList(prompt string, options []string) (int, string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		for i, opt := range options {
			fmt.Printf("  %d) %s\n", i+1, opt)
		}

		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > len(options) {
			fmt.Println("Неверный ввод. Попробуйте снова.")
			continue
		}

		return num - 1, options[num-1]
	}
}

func GetOrChoose(prompt, value string, options []string) string {
	if value == "" {
		_, value = ChooseFromList(prompt, options)
	}
	return value
}