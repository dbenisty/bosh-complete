package main

import (
	"fmt"
	"os"
	"strings"
)

type compContext struct {
	CurrentToken string
	Command      string
	Args         []string
	//Long flag string of the current flag that is being completed
	CurrentFlag string
	//Long flag string to value(s)
	Flags    map[string][]string
	Switches []string
}

func (c compContext) Complete() ([]string, error) {
	var compFn compFunc

	log.Write("Current Token: %s", c.CurrentToken)

	//determine what we're completing
	if strings.HasPrefix(c.CurrentToken, "-") {
		log.Write("Completing flag names")
		compFn = compFlagNames
	} else if c.CurrentFlag != "" {
		log.Write("Checking current flag: %s", c.CurrentFlag)
		flag, found := flags[c.CurrentFlag]
		if found {
			log.Write("Completing flag value for flag %s", c.CurrentFlag)
			compFn = flag.Complete
		}
	} else if c.Command == "" {
		log.Write("Completing command names")
		compFn = compCommandNames
	} else {
		if cmd, found := commands.Find(c.Command); found {
			if len(c.Args) < len(cmd.ArgComps) {
				position := len(c.Args)
				log.Write("Completing positional arg for command `%s' position %d (0 indexed)", c.Command, position)
				compFn = cmd.ArgComps[position]
			}
		}
	}

	if compFn == nil {
		log.Write("No completion registered")
		compFn = compNoop
	}

	candidates, err := compFn(c)
	if err != nil {
		return nil, err
	}

	log.Write("Completion candidates: \n---START---\n%s\n---END---\n", strings.Join(candidates, "\n"))

	ret := []string{}
	for _, val := range candidates {
		if strings.HasPrefix(val, c.CurrentToken) {
			ret = append(ret, val)
		}
	}
	log.Write("Completion return: \n---START---\n%s\n---END---\n", strings.Join(ret, "\n"))
	return ret, nil
}

type compFunc func(compContext) ([]string, error)

func doComplete(boshArgs []string) {
	log.Write("in complete")
	argsString := ""
	if len(boshArgs) > 0 {
		argsString = fmt.Sprintf(`'%s'`, strings.Join(boshArgs, `', '`))
	}
	log.Write("Bosh args: [%s]", argsString)

	insertGlobalFlags()
	commands.Populate()

	compContext := parseContext(boshArgs)
	results, err := compContext.Complete()
	if err != nil {
		log.Write("Completion error: %s", err.Error())
		return
	}

	response := strings.Join(results, "\n")
	log.Write(response)
	fmt.Print(response)
}

func parseContext(args []string) compContext {
	if len(args) < 2 {
		os.Exit(0)
	}

	ret := compContext{
		CurrentToken: args[len(args)-1],
		Flags:        map[string][]string{},
	}

	ret.CurrentFlag = ""

	//loop over all but last token - the last one is the token
	// we're suggesting changes to.
	for i := 0; i < len(args)-1; i++ {
		token := args[i]
		if strings.HasPrefix(token, "-") {
			//Check if value or not
			f := flags[token]
			ret.CurrentFlag = "--" + f.Long
			if !f.TakesValue {
				log.Write("current flag is switch: %s", ret.CurrentFlag)
				ret.Switches = append(ret.Switches)
				ret.CurrentFlag = ""
			}
		} else {
			if ret.CurrentFlag != "" {
				//This is the value to a flag
				ret.Flags[ret.CurrentFlag] = append(ret.Flags[ret.CurrentFlag], token)
				ret.CurrentFlag = ""
			} else if ret.Command == "" {
				//This is a command name
				if cmd, found := commands.Find(token); found {
					ret.Command = token
					cmd.InsertFlags()
				}
			} else {
				//This is a positional argument
				ret.Args = append(ret.Args, token)
			}
		}
	}

	return ret
}