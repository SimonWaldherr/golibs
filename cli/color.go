package cli

import "fmt"

type Col int

const (
	Black = Col(iota)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func Color(str string, col Col) string {
	return fmt.Sprintf("\033[3%vm%v\033[0m", col, str)
}

func Bold(str string) string {
	return fmt.Sprintf("\033[1m%v\033[0m", str)
}

func Underline(str string) string {
	return fmt.Sprintf("\033[4m%v\033[0m", str)
}

func Cmd(cmds ...*exec.Cmd) ([]byte, error) {
	for i, cmd := range cmds[:len(cmds)-1] {
		out, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		cmd.Start()
		cmds[i+1].Stdin = out
	}

	ret, err := cmds[len(cmds)-1].Output()
	if err != nil {
		return nil, err
	}

	return ret, nil
}
