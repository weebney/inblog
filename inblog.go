/*
	(c) weebney 2024
	See `license` for details
*/

package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/weebney/inblog/internal/inblog"
	"github.com/weebney/inblog/internal/wizard"
)

const Version = "1.0"

var flags struct {
	rebuild bool
	refresh bool
	skip    bool
	wizard  bool
	mizard  bool
	help    bool
	version bool

	outputDir  string
	contentDir string
}

func init() {
	flag.BoolVarP(&flags.version, "version", "v", false, "print version info and quit")
	flag.BoolVarP(&flags.rebuild, "rebuild", "r", false, "force a rebuild of the HTML")
	flag.BoolVarP(&flags.skip, "skip", "s", false, "skip fetching emails")
	flag.BoolVarP(&flags.refresh, "refresh", "R", false, "refresh all content, overwriting edits")
	flag.BoolVarP(&flags.wizard, "wizard", "w", false, "force the setup wizard")
	flag.BoolVarP(&flags.mizard, "no-wizard", "m", false, "disable the setup wizard")
	flag.StringVarP(&flags.outputDir, "output-dir", "o", "public", "set the output directory")
	flag.StringVarP(&flags.contentDir, "content-dir", "c", "content", "set the content directory")
	flag.BoolVarP(&flags.help, "help", "h", false, "show this help message")

	flag.Parse()
	if flags.help || flags.version {
		if flags.help {
			fmt.Println("Usage: inblog [OPTIONS...]")
			flag.PrintDefaults()
		}
		fmt.Println("inblog v"+Version, "(c) weebney 2024, BSD-2-Clause")
		os.Exit(0)
	}
}

func main() {
	config, needsWizard := inblog.LoadConfig()
	if (needsWizard || flags.wizard) && !flags.mizard {
		wizard.Wizard()
		config, _ = inblog.LoadConfig()
	}

	inblog.EnsureContent(flags.contentDir, flags.outputDir)

	lastUID := inblog.ReadCache()
	if flags.refresh {
		lastUID = 0
		//log.Debug("read cache", "lastUID", 0)
	}

	if !flags.skip {
		newMessages := inblog.FetchEmails(config, lastUID)
		inblog.ProcessEmails(newMessages, flags.contentDir)
	} else {
		//log.Debug("skipping fetch")
	}

	inblog.GenerateHTML(config, flags.contentDir, flags.outputDir)
}
