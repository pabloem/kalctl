package base

import "strings"

type CommandArgs struct {
	Args   []string
	KwArgs map[string]string
}

func ParseArgs(args []string) CommandArgs {
	var parsed CommandArgs
	for _, elm := range args {
		if strings.HasPrefix(elm, "--") && strings.Contains(elm, "=") {
			parts := strings.SplitN(elm[2:], "=", 2)
			if len(parts) == 2 {
				if parsed.KwArgs == nil {
					parsed.KwArgs = make(map[string]string)
				}
				parsed.KwArgs[parts[0]] = parts[1]
			}
		} else if strings.HasPrefix(elm, "--") {
			if parsed.KwArgs == nil {
				parsed.KwArgs = make(map[string]string)
			}
			parsed.KwArgs[elm[2:]] = ""
		} else {
			parsed.Args = append(parsed.Args, elm)
		}
	}
	return parsed
}
