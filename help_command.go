package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)


type CommandsSort []*Command

func (this CommandsSort) Len() int {
	return len(this)
}

func (this CommandsSort) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this CommandsSort) Less(i, j int) bool {
	return this[i].Name < this[j].Name
}

func NewHelpCommand() *Command {
	return NewCommand("help", "Show this help", func(o *Command, out Output) {
		if n := o.Argument("command").String(); n != "" {
			if cmd, ok := o.Cli.Commands[n]; ok {
				out.Printf(DescribeCommand(cmd))
			} else {
				out.Printf(DescribeCli(o.Cli))
				out.Printf("\n\n<error>Unknown command \"%s\"<reset>\n", n)
			}
		} else {
			out.Printf(DescribeCommand(o))
		}
	}).NewArgument("command", "Command to show help for", "", false, false)
}

var DescribeCommand = func(c *Command) string {
	lines := []string{"Command: <headline>" + c.Name + "<reset>"}

	if c.Description != "" {
		lines = append(lines, []string{"<info>" + c.Description + "<reset>", ""}...)
	} else if c.Usage != "" {
		lines = append(lines, []string{"<info>" + c.Usage + "<reset>", ""}...)
	}

	lines = append(lines, "<subline>Usage:<reset>")
	usage := []string{c.Name}
	args := make([][]string, 0)
	argMax := 0
	opts := make([][]string, 0)
	optMax := 0
	for _, p := range c.Arguments {
		var short string
		usg := p.Usage
		if p.Required {
			short = fmt.Sprintf("<%s>", p.Name)
			usg += " (req)"
		} else {
			short = fmt.Sprintf("[%s]", p.Name)
		}
		if p.Multiple {
			short = "(" + short + " ...)"
			usg += " (mult)"
		}
		if p.Default != "" {
			usg += fmt.Sprintf(" (default: \"%s\")", p.Default)
		}
		if l := len(p.Name); l > argMax {
			argMax = l
		}
		usage = append(usage, short)
		args = append(args, []string{p.Name, usg})
	}
	for _, p := range c.Options {
		short := fmt.Sprintf("--%s", p.Name)
		if p.Alias != "" {
			short += "|-" + p.Alias
		}
		if !p.Flag {
			short += " <val>"
		}
		long := short
		usg := p.Usage
		if !p.Required {
			short = "(" + short + ")"
		} else {
			usg += " (req)"
		}
		if p.Multiple {
			short = "(" + short + " ...)"
			usg += " (mult)"
		}
		if p.Default != "" {
			usg += fmt.Sprintf(" (default: \"%s\")", p.Default)
		}
		if l := len(long); l > optMax {
			optMax = l
		}
		usage = append(usage, short)
		opts = append(opts, []string{long, usg})
	}
	lines = append(lines, "\t"+strings.Join(usage, " "))
	lines = append(lines, "")

	if len(args) > 0 {
		lines = append(lines, "<subline>Arguments:<reset>")
		for _, l := range args {
			lines = append(lines, fmt.Sprintf("\t<info>%-"+fmt.Sprintf("%d", argMax)+"s<reset>  %s", l[0], l[1]))
		}
		lines = append(lines, "")
	}

	if len(opts) > 0 {
		lines = append(lines, "<subline>Options:<reset>")
		for _, l := range opts {
			lines = append(lines, fmt.Sprintf("\t<info>%-"+fmt.Sprintf("%d", optMax)+"s<reset>  %s", l[0], l[1]))
		}
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n") + "\n"
}

var DescribeCli = func(c *Cli) string {

	// name + version
	line := "<headline>" + c.Name + "<reset>"
	if c.Version != "" {
		line += " <debug>(" + c.Version + ")<reset>"
	}
	lines := []string{line}

	// description
	if c.Description != "" {
		lines = append(lines, "<info>"+c.Description+"<reset>\n")
	}

	// usage
	prog := filepath.Base(os.Args[0])
	lines = append(lines, fmt.Sprintf("<subline>Usage:<reset>\n\t%s <command> [<arg> ..] [--opt <val> ..]\n", prog))

	// commands
	lines = append(lines, "<subline>Available commands:<reset>")
	max := 0
	commands := make([]*Command, len(c.Commands))
	i := 0
	for _, cmd := range c.Commands {
		commands[i] = cmd
		i++
		if l := len(cmd.Name); l > max {
			max = l
		}
	}
	sort.Sort(CommandsSort(commands))
	for _, cmd := range commands {
		lines = append(lines, fmt.Sprintf("\t<info>%-"+fmt.Sprintf("%d", max)+"s<reset>  %s", cmd.Name, cmd.Usage))
	}

	return strings.Join(lines, "\n") + "\n"
}