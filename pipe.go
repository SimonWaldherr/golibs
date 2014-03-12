package pipe

func cmd(cmds ...*exec.Cmd) ([]byte, error) {
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
