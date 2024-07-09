package wizard

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emersion/go-imap/client"
	"github.com/weebney/inblog/internal/inblog"
)

type model struct {
	config inblog.Config
	input  textinput.Model
	step   int
	err    error
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(1, 1, 1, 1).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 0, 1, 2)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Padding(0, 0, 0, 2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Padding(0, 0, 1, 2)
)

type step struct {
	title    string
	subtitle string
	field    *string
}

func getSteps(c *inblog.Config) []step {
	return []step{
		{title: "Email Address", subtitle: "Enter the email address you want to send new posts to.\ne.g. myInblogEmail@example.com\n\nYou SHOULD make a new, dedicated email account to power your inblog.\nOutlook.com is the easiest to use with inblog; press enter to go to\nthe registration page. For technical reasons, Gmail is not supported.", field: &c.Email},
		{title: "Password", subtitle: "\n\n\n\nEnter the password for the email account\n", field: &c.Password},
		{title: "IMAP Server", subtitle: "\n\n\nEnter the IMAP server address\n(e.g., imap.example.com:993)\n", field: &c.ImapServer},
		{title: "Mailbox", subtitle: "\n\n\nEnter the mailbox name to use.\nLeave Blank for 'INBOX'\n", field: &c.Mailbox},
		{title: "Approved Sender", subtitle: "\n\nEnter the email address of the approved sender.\nOnly emails sent from this email address will be processed and posted.\nThis should probably be your personal email.\n", field: &c.ApprovedSender},
		{title: "Blog Title", subtitle: "\n\n\nEnter your blog title\n(e.g. My Awesome Blog)\n", field: &c.BlogName},
	}
}

func initialModel() model {
	ti := textinput.New()
	ti.Focus()

	return model{
		input: ti,
		step:  0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.step >= len(getSteps(&m.config)) {
				return m, tea.Quit
			}

			// open outlook.com registration if no email provided
			if getSteps(&m.config)[m.step].title == "Email Address" && m.input.Value() == "" {
				openBrowser("https://signup.live.com/")
				openBrowser("https://support.microsoft.com/en-us/office/pop-imap-and-smtp-settings-for-outlook-com-d088b986-291d-42b8-9564-9c414e2aa040")
				os.Exit(1)
			}

			*getSteps(&m.config)[m.step].field = m.input.Value()
			if getSteps(&m.config)[m.step].title == "Mailbox" && m.input.Value() == "" {
				m.config.Mailbox = "INBOX"
			}
			m.step++
			m.input.SetValue("")
			m.err = nil

			if m.step < len(getSteps(&m.config)) {
				if getSteps(&m.config)[m.step].title == "Password" {
					m.input.EchoMode = textinput.EchoPassword
					m.input.EchoCharacter = '•'
				} else {
					m.input.EchoMode = textinput.EchoNormal
				}

				// validate account info
				if getSteps(&m.config)[m.step].title == "Mailbox" {
					if !strings.Contains(m.config.ImapServer, ":") {
						m.config.ImapServer += ":993"
					}
					c, err := client.DialTLS(m.config.ImapServer, nil)
					if err != nil {
						log.Fatal("Error while connecting to the IMAP server: ", err)
					}
					if err := c.Login(m.config.Email, m.config.Password); err != nil {
						log.Fatal("Error while testing your login credentials: ", err)
					}
					c.Logout()
				}
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			os.Exit(1)
		}

	// Handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	// autofill server for outlook.com
	if m.step < len(getSteps(&m.config)) && getSteps(&m.config)[m.step].title == "IMAP Server" && m.input.Value() == "" && strings.Contains(m.config.Email, "outlook.com") {
		m.input.SetValue("outlook.office365.com:993")
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	steps := getSteps(&m.config)
	s := titleStyle.Render("inblog Setup Wizard 🧙🏻‍♂️")

	if m.step < len(steps) {
		currentStep := steps[m.step]
		s += titleStyle.Render(currentStep.title)
		s += "\n"
		s += subtitleStyle.Render(currentStep.subtitle)
		s += "\n"
		s += inputStyle.Render(m.input.View())
		s += "\n"
	} else {
		s += "\n\n\n"
		s += titleStyle.Render("Setup Complete!")
		s += "\n\n\n"
		s += subtitleStyle.Render("Press any key to continue")
		s += "\n\n"
	}

	if m.err != nil {
		s += errorStyle.Render(m.err.Error())
	}

	return s
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
		fallthrough
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		fallthrough
	case "darwin":
		err = exec.Command("open", url).Start()
		fallthrough
	default:
		fmt.Printf("\n\rURL: %v\n", url)
	}

	if err != nil {
		log.Fatal(err)
	}

}

func Wizard() {
	os.Remove(".cache")

	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		log.Fatalf("Error in wizard: %v", err)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Failed to get user's config dir: %v", err)
	}
	inblogDir := filepath.Join(configDir, "inblog")
	configFile := filepath.Join(inblogDir, "config.gob")
	err = os.MkdirAll(inblogDir, 0o755)
	if err != nil {
		log.Fatalf("Failed to make config dir: %v", err)
	}

	gobBuf := &bytes.Buffer{}
	genc := gob.NewEncoder(gobBuf)
	genc.Encode(m.(model).config)

	err = os.WriteFile(configFile, gobBuf.Bytes(), 0o644)
	if err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}
}
