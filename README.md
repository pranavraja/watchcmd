
A file watcher that runs commands. A work in progress.

# Setup

	go get github.com/pranavraja/watchcmd

# Running

Assuming `$GOPATH` is in your `$PATH`:

	watchcmd

...will watch the current directory for changes and run the commands in
`watchcmd.rules` accordingly.

# Usage

See the output of `watchcmd -h`:

	Usage of ./watchcmd:
	  -directory=".": directory to watch (e.g. --directory src). Defaults to
	   the current dir
	  -rules="watchcmd.rules": file containing rules of the form
	   event<TAB>regexp<TAB>command (default filename is watchcmd.rules)
	  -batchUpdate=1: to prevent unnecessary runs, if multiple files tend to be
	   updated in a batch, the typical duration (in milliseconds) to wait for
	   that batch

# Rules

The rules file contains the patterns to watch (as regular expressions), and the
commands to execute. An example rule:

	MODIFY	\.js$	grunt browserify

The above rule executes `grunt browserify` whenever any file (or set of files)
with extension `.js` is modified. Note that if multiple files matching the
pattern are updated at the same time (see Batch updates below), the command
will only be run once.

Regexp substitution is also supported, for example:

	MODIFY	src\/(.+)\.jade$	jade < src/$1.jade > dist/$1.html

This runs the `jade` executable for each modified `.jade` file under the path
`src/`, replacing a html file of the same name under the path `dist/`.
Directory structure is preserved. This rule can be useful when developing
static sites, for example.

You can also listen for `CREATE` and `DELETE` events, for example:

    CREATE	src/(.+)	mkdir -p dist/`dirname $1`

This will copy directory structures for files created under `src/`.

    REMOVE	src/(.+)	rm dist/$1

This will remove the corresponding file from `dist/` when it is deleted from `src/`.

# Batch updates

One of the annoyances with watchers is that if multiple files have changed
within a short timeframe, the tasks can be run multiple times unnecessarily.
For example, if you have a single build task for your entire project, you don't
want to be triggering it 10 times just because you changed 10 files. If the end
result is going to be the same, that's 9 wasted builds. 

`watchcmd` supports passing in a `--batchUpdate` parameter, which represents
the time between the first and last update of a batch. This can be useful in a
bunch of scenarios, for example dependent build tasks, or when you use vim and
a single `:w` ends up writing the file like 6 times.

