package main

import (
	"github.com/weebney/inblog/internal/inblog"
	"github.com/weebney/inblog/internal/wizard"
)

func main() {
	config, needsWizard := inblog.LoadConfig()
	if needsWizard {
		wizard.Wizard()
		config, _ = inblog.LoadConfig()
	}
	inblog.EnsureDirectories()
	inblog.EnsureTemplates()

	lastUID := inblog.ReadCache()
	newMessages := inblog.FetchEmails(config, lastUID)
	inblog.ProcessEmails(newMessages)
	inblog.GenerateHTML(config)
}
