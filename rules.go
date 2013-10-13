package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Rule struct {
	match *regexp.Regexp
	cmd   string
}

func (r Rule) MatchedCommand(filename string) (string, bool) {
	match := r.match.FindStringSubmatch(filename)
	if match == nil {
		return "", false
	}
	if len(match) == 1 {
		return r.cmd, true
	}
	return r.match.ReplaceAllString(filename, r.cmd), true
}

func LoadRules(filename string) ([]Rule, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var rules []Rule
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fragments := strings.SplitN(line, "\t", 2)
		if len(fragments) < 2 {
			return nil, fmt.Errorf("Invalid line %s, expected regexp<TAB>command", line)
		}
		exp, err := regexp.Compile(fragments[0])
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{match: exp, cmd: fragments[1]})
	}
	return rules, nil
}
