package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type Prompter struct {
	reader *bufio.Reader
}

func NewPrompter() *Prompter {
	return &Prompter{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (p *Prompter) PromptUsernameAndPassword() (username, password string, err error) {
	if _, err = fmt.Fprint(os.Stdout, "Username: "); err != nil {
		return "", "", err
	}
	username, err = p.reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	username = strings.TrimSpace(username)

	if _, err = fmt.Fprint(os.Stdout, "Password: "); err != nil {
		return "", "", err
	}
	var passwordBytes []byte
	passwordBytes, err = term.ReadPassword(int(os.Stdin.Fd()))
	if _, werr := fmt.Fprintln(os.Stdout); werr != nil {
		return "", "", werr
	}
	if err != nil {
		return "", "", err
	}
	password = string(passwordBytes)

	return username, password, nil
}
