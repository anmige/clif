package cli

import (
	"os"
	"reflect"
)

// Cli is a command line interface object
type Cli struct {
	Name        string
	Version     string
	Description string
	Commands    map[string]*Command
	Registry    *Registry
}

func New(name, version, desc string) *Cli {
	this := &Cli{
		Name:        name,
		Version:     version,
		Description: desc,
		Commands:    make(map[string]*Command),
		Registry:    NewRegistry(),
	}
	this.Add(NewHelpCommand())
	out := NewFancyOutput(os.Stdout)
	this.Register(this).Register(out).Register(NewDefaultInput(os.Stdin, out))
	return this
}

// Add is a builder method for adding a new command
func (this *Cli) Add(cmd *Command) *Cli {
	this.Commands[cmd.Name] = cmd.SetCli(this)
	return this
}

// New creates and adds a new command
func (this *Cli) New(name, usage string, call CallMethod) *Cli {
	return this.Add(NewCommand(name, usage, call).SetCli(this))
}

// RegisterAs is builder method and registers object in registry
func (this *Cli) Register(v interface{}) *Cli {
	this.Registry.Register(v)
	return this
}

// RegisterAs is builder method and registers object under alias in registry
func (this *Cli) RegisterAs(n string, v interface{}) *Cli {
	this.Registry.Alias(n, v)
	return this
}

// Run with arguments
func (this *Cli) Run() {
	this.RunWith(os.Args[1:])
}

// Run the cli and be happy
func (this *Cli) RunWith(args []string) {
	if len(args) < 1 {
		Die(DescribeCli(this))
	} else if c, ok := this.Commands[args[0]]; ok {
		this.Register(c)
		method := c.Call.Type()
		input := make([]reflect.Value, method.NumIn())
		for i := 0; i < method.NumIn(); i++ {
			t := method.In(i)
			s := t.String()
			if this.Registry.Has(s) {
				input[i] = this.Registry.Get(s)
			} else {
				Die("Missing parameter %s", s)
			}
		}

		if err := c.Parse(args[1:]); err != nil {
			Die("Parse error: %s", err)
		}

		res := c.Call.Call(input)

		if len(res) > 0 && res[0].Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) && !res[0].IsNil() {
			Die("Failure in execution: %s", res[0].Interface().(error))
		}
	} else {
		Die("Command \"%s\" unknown", args[0])
	}
}

// SetOutput is builder method and replaces current input
func (this *Cli) SetInput(in Input) *Cli {
	t := reflect.TypeOf((*Input)(nil)).Elem()
	this.Registry.Alias(t.String(), in)
	return this
}

// SetOutput is builder method and replaces current output
func (this *Cli) SetOutput(out Output) *Cli {
	t := reflect.TypeOf((*Output)(nil)).Elem()
	this.Registry.Alias(t.String(), out)
	return this
}