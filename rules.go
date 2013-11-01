package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/howeyc/fsnotify"
)

type Rule struct {
	eventName string
	match     *regexp.Regexp
	cmd       string
}

func (r Rule) MatchedCommand(ev *fsnotify.FileEvent) (string, bool) {
	if r.eventName == "MODIFY" && !ev.IsModify() {
		return "", false
	}
	if r.eventName == "CREATE" && !ev.IsCreate() {
		return "", false
	}
	if r.eventName == "DELETE" && !ev.IsDelete() {
		return "", false
	}
	match := r.match.FindStringSubmatch(ev.Name)
	if match == nil {
		return "", false
	}
	if len(match) == 1 {
		return r.cmd, true
	}
	return r.match.ReplaceAllString(ev.Name, r.cmd), true
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
		fragments := strings.SplitN(line, "\t", 3)
		if len(fragments) < 3 {
			return nil, fmt.Errorf("Invalid line %s, expected event<TAB>regexp<TAB>command", line)
		}
		exp, err := regexp.Compile(fragments[1])
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{eventName: fragments[0], match: exp, cmd: fragments[2]})
	}
	return rules, nil
}
