package main

import (
	"flag"
	"log"
	"time"
)

var watchDir string
var rulesFile string
var batchUpdateDuration int64

func init() {
	flag.StringVar(&watchDir, "directory", ".", "directory to watch (e.g. --directory src). Defaults to the current dir")
	flag.StringVar(&rulesFile, "rules", "watchcmd.rules", "file containing rules of the form event<TAB>regexp<TAB>command (default filename is watchcmd.rules)")
	flag.Int64Var(&batchUpdateDuration, "batchUpdate", 1, "to prevent unnecessary runs, if multiple files tend to be updated in a batch, the typical duration (in milliseconds) to wait for that batch")
	flag.Parse()
}

func main() {
	rules, err := LoadRules(rulesFile)
	if err != nil {
		log.Fatal(err)
	}
	watcher, err := NewRecursiveWatcher(watchDir)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	for {
		events, err := watcher.NextBatch(time.Duration(batchUpdateDuration) * time.Millisecond)
		if err != nil {
			log.Fatal(err)
		}
		commands := make(map[string]struct{})
		for _, ev := range events {
			for _, rule := range rules {
				if cmd, ok := rule.MatchedCommand(ev); ok {
					log.Println(rule.eventName + " " + ev.Name)
					commands[cmd] = struct{}{}
					break
				}
			}
		}
		for cmd, _ := range commands {
			err := runCommand(cmd)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
