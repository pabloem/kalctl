package base

import (
	"fmt"

	"github.com/pabloem/kalctl/auth"
	"github.com/pabloem/kalctl/output"
	"github.com/pabloem/kalctl/reqs"
	"github.com/rs/zerolog/log"
)

type Element interface {
	Name() string
	Description() string
}

type Command interface {
	Element
	Run(args CommandArgs) error
	Arguments() []Argument
}

type Namespace interface {
	Element
	Children() []Element
}

type namespaceImpl struct {
	name     string
	desc     string
	children []Element
}

func (n *namespaceImpl) Name() string {
	return n.name
}

func (n *namespaceImpl) Description() string {
	return n.desc
}

func (n *namespaceImpl) Children() []Element {
	return n.children
}

func NewNamespace(name, desc string, children ...Element) Namespace {
	return &namespaceImpl{
		name:     name,
		desc:     desc,
		children: children,
	}
}

type Argument struct {
	Name     string
	Desc     string
	Position int
	Required bool
}

type httpRequestCommand struct {
	name     string
	desc     string
	template reqs.HttpRequestTemplate
	args     []Argument
}

func (c *httpRequestCommand) Name() string {
	return c.name
}

func (c *httpRequestCommand) Description() string {
	return c.desc
}

func (c *httpRequestCommand) Arguments() []Argument {
	return c.args
}

func (c *httpRequestCommand) Run(args CommandArgs) error {
	log.Trace().Msgf("Running command %s", c.name)
	authToken, err := auth.GetToken()
	if err != nil {
		return fmt.Errorf("unable to get auth token. run 'kalctl auth login' to authenticate", err)
	}

	res, err := reqs.KalshiRequest(c.template, authToken, "")
	if err != nil {
		return err
	}
	outputFormatter := output.GetFormatter()
	fmt.Println(outputFormatter.CommandResult(res))
	return nil
}

func NewCommand(name, desc string, requestTemplate reqs.HttpRequestTemplate, args ...Argument) Command {
	return &httpRequestCommand{
		name:     name,
		desc:     desc,
		args:     args,
		template: requestTemplate,
	}
}
