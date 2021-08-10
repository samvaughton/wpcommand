package types

type CommandResult struct {
	Command string
	Output  string
}

type CommandResults struct {
	Items []CommandResult
}

func (cr *CommandResults) Add(result CommandResult) {
	cr.Items = append(cr.Items, result)
}
