package output

import "github.com/charmbracelet/lipgloss"

type OutputFormatter interface {
	Title(title string) string
	Description(desc string) string

	Attribute(name string) string
	AttributeDescription(desc string) string

	CommandResult(res string) string
}

type lipglossOutputFormatter struct {
	titleStyle    lipgloss.Style
	descStyle     lipgloss.Style
	attrStyle     lipgloss.Style
	attrDescStyle lipgloss.Style
	resultStyle   lipgloss.Style
}

var defaultOutputFormatter = lipglossOutputFormatter{
	titleStyle:    lipgloss.NewStyle().Bold(true),
	descStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC66")),
	attrStyle:     lipgloss.NewStyle().Bold(true).PaddingLeft(4).Width(15),
	attrDescStyle: lipgloss.NewStyle().PaddingLeft(8).Width(58),
	resultStyle:   lipgloss.NewStyle(),
}

func GetFormatter() OutputFormatter {
	return &defaultOutputFormatter
}

func (d *lipglossOutputFormatter) Title(title string) string {
	return d.titleStyle.Render(title)
}

func (d *lipglossOutputFormatter) Description(desc string) string {
	return d.descStyle.Render(desc)
}

func (d *lipglossOutputFormatter) Attribute(name string) string {
	return d.attrStyle.Render(name)
}

func (d *lipglossOutputFormatter) AttributeDescription(desc string) string {
	return d.attrDescStyle.Render(desc)
}

func (d *lipglossOutputFormatter) CommandResult(result string) string {
	return d.resultStyle.Render(result)
}
