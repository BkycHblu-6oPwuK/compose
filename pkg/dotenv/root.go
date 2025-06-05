package dotenv

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var (
	loadedEnv map[string]string
)
// loads env from file
func Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	loadedEnv = make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		key, value, err := parse(scanner.Text())
		if err != nil {
			continue
		}
		loadedEnv[key] = value
		os.Setenv(key, value)
	}

	return scanner.Err()
}
// return last loaded environments
func GetLoadedEnv() map[string]string {
	return loadedEnv
}

func parse(line string) (key, value string, err error) {
	line = strings.TrimSpace(line)

	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", errors.New("empty or comment line")
	}

	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid line")
	}

	key = strings.TrimSpace(parts[0])
	value = strings.TrimSpace(parts[1])

	value = stripComment(value)

	runes := []rune(value)
	lenRunes := len(runes)
	if lenRunes > 0 {
		if runes[lenRunes-1] == '"' || runes[lenRunes-1] == '\'' {
			runes = runes[:lenRunes-1]
			lenRunes = lenRunes - 1
		}
		if lenRunes > 0 {
			if runes[0] == '"' || runes[0] == '\'' {
				runes = runes[1:]
			}
		}
		value = string(runes)
	}

	return key, value, nil
}

func stripComment(value string) string {
	var (
		inSingleQuote bool
		inDoubleQuote bool
		builder       strings.Builder
	)
	for _, r := range value {
		switch r {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
		case '#', ';':
			if !inSingleQuote && !inDoubleQuote {
				return strings.TrimSpace(builder.String())
			}
		}
		builder.WriteRune(r)
	}

	return strings.TrimSpace(builder.String())
}
