package parsing

var escapableInDoubleQuotes = map[rune]int{
	'"':  1,
	'\\': 1,
	'$':  1,
	'`':  1,
	'\n': 1,
}

// The first item from the returned array is the command itself
// the rest of the returned array contains arguments for the command
func ParseArgsWithQuotes(command string, index int) []string {
	runes := []rune(command)
	args := make([]string, 0)

	currArg := ""
	for index < len(command) {
		curr := index
		switch command[curr] {
		case '"':
			curr++
			index = curr
			for index < len(command) && command[index] != '"' {
				if command[index] == '\\' {
					if _, ok := escapableInDoubleQuotes[runes[index+1]]; ok {
						index++
					}
				}
				currArg += string(command[index])
				index++
			}
			index++
		case '\'':
			curr++
			index = curr
			for index < len(command) && command[index] != '\'' {
				// if command[index] == '\\' {
				// 	index++
				// }
				currArg += string(command[index])
				index++
			}
			index++
		default:
			for index < len(command) && (command[index] != ' ' && command[index] != '\'' && command[index] != '"') {
				if command[index] == '\\' {
					index++
				}
				currArg += string(command[index])
				index++
			}
		}
		if (index >= len(command) || command[index] == ' ') && len(currArg) > 0 {
			args = append(args, currArg)
			currArg = ""
		}
		for index < len(command) && command[index] == ' ' {
			index++
		}
	}

	return args
}
