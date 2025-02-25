package factory

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/OctopusDeploy/cli/pkg/apiclient"
	"github.com/OctopusDeploy/cli/pkg/question"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
)

// wrapper over the underlying spinner so we can mock it
type Spinner interface {
	Start()
	Stop()
}

type factory struct {
	client       apiclient.ClientFactory
	asker        question.AskProvider
	spinner      Spinner
	buildVersion string
}

type Factory interface {
	GetSystemClient(requester apiclient.Requester) (*client.Client, error)
	GetSpacedClient(requester apiclient.Requester) (*client.Client, error)
	GetCurrentSpace() *spaces.Space
	GetCurrentHost() string
	Spinner() Spinner
	IsPromptEnabled() bool
	Ask(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
	BuildVersion() string
}

func New(clientFactory apiclient.ClientFactory, asker question.AskProvider, s Spinner, buildVersion string) Factory {
	return &factory{
		client:       clientFactory,
		asker:        asker,
		spinner:      s,
		buildVersion: buildVersion,
	}
}

func (f *factory) GetSystemClient(requester apiclient.Requester) (*client.Client, error) {
	// Fine to start the spinner here because we guarantee no
	// user-input is triggered by GetSystemClient
	f.spinner.Start()
	defer f.spinner.Stop()
	return f.client.GetSystemClient(requester)
}

func (f *factory) GetSpacedClient(requester apiclient.Requester) (*client.Client, error) {
	// NOTE: Don't start the spinner here because it may prompt for user input,
	// and GetSpacedClient has no access to the spinner in order to stop it.
	//
	// We could link the spinner into GetSpacedClient, but that feels like we're jumping through
	// too many hoops when this is a niche situation and we can just turn the spinner off.
	return f.client.GetSpacedClient(requester)
}

func (f *factory) GetCurrentSpace() *spaces.Space {
	return f.client.GetActiveSpace()
}

func (f *factory) GetCurrentHost() string {
	return f.client.GetHostUrl()
}

func (f *factory) IsPromptEnabled() bool {
	return f.asker.IsInteractive()
}

func (f *factory) Ask(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	return f.asker.Ask(p, response, opts...)
}

func (f *factory) Spinner() Spinner {
	return f.spinner
}

func (f *factory) BuildVersion() string {
	return f.buildVersion
}

// NoSpinner is a static singleton "does nothing" stand-in for spinner if you want to
// call an API that expects a spinner while you're in automation mode.
var NoSpinner Spinner = &noSpinner{}

type noSpinner struct{}

func (f *noSpinner) Start() {}
func (f *noSpinner) Stop()  {}
