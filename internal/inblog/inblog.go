package inblog

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"mime/quotedprintable"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"github.com/yuin/goldmark"
)

type Config struct {
	Email          string
	Password       string
	ImapServer     string
	Mailbox        string
	ApprovedSender string
	BlogName       string
}

func LoadConfig() (conf Config, needsWizard bool) {
	conf = Config{
		Email:          os.Getenv("INBLOG_EMAIL"),
		Password:       os.Getenv("INBLOG_PASSWORD"),
		ImapServer:     os.Getenv("INBLOG_IMAPSSERVER"),
		Mailbox:        os.Getenv("INBLOG_MAILBOX"),
		ApprovedSender: os.Getenv("INBLOG_APPROVED_SENDER"),
		BlogName:       os.Getenv("INBLOG_NAME"),
	}

	if conf.ApprovedSender == "" || conf.Email == "" || conf.ImapServer == "" || conf.Password == "" || conf.Mailbox == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatalf("Failed to get user's config dir: %v", err)
		}
		inblogDir := filepath.Join(configDir, "inblog")
		configFile := filepath.Join(inblogDir, "config.gob")

		_, err = os.Stat(configFile)
		if !os.IsNotExist(err) {
			// If config file exists, read and decode it
			data, err := os.ReadFile(configFile)
			if err != nil {
				log.Fatalf("Failed to read config file: %v", err)
			}

			dec := gob.NewDecoder(bytes.NewReader(data))
			err = dec.Decode(&conf)
			if err != nil {
				log.Fatalf("Failed to decode config file: %v", err)
			}
		}
	}

	if conf.ApprovedSender == "" || conf.Email == "" || conf.ImapServer == "" || conf.Password == "" {
		return Config{}, true
	}

	return conf, false
}

func EnsureContent(contentDir string, outputDir string) {
	dirs := []string{contentDir, contentDir + "/posts", contentDir + "/templates", outputDir, outputDir + "/posts"}
	for _, dir := range dirs {
		os.MkdirAll(dir, os.ModePerm)
	}

	templates := map[string]string{
		contentDir + "/templates/listitem.template.html": `<li><a href="%HYPERLINK%">%SUBJECT%</a> <span>%DATE%</span>`,
		contentDir + "/templates/index.template.html":    `<!doctype html><html lang="en"><head><meta charset="utf-8"><title>Index - %BLOG_NAME%</title></head><body><h1>Index of %BLOG_NAME%</h1><ul>%LIST%</ul><footer><span>Powered by inblog</span></footer></body></html>`,
		contentDir + "/templates/post.template.html":     `<!doctype html><html lang="en"><head><meta charset="utf8"><title>%SUBJECT% - %BLOG_NAME%</title></head><body><a href="../index.html">back</a><h1 id="subject">%SUBJECT%</h1><p id="date">%DATE%</p><div id="content">%BODY%</div></body></html>`,
	}

	for file, content := range templates {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			os.WriteFile(file, []byte(content), 0o644)
		}
	}
}

func ReadCache() uint32 {
	data, err := os.ReadFile(".cache")
	if err != nil {
		return 0
	}
	var lastUID uint32
	fmt.Sscanf(string(data), "%d", &lastUID)
	return lastUID
}

func FetchEmails(config Config, lastUID uint32) []*imap.Message {
	c, err := client.DialTLS(config.ImapServer, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	if err := c.Login(config.Email, config.Password); err != nil {
		log.Fatal(err)
	}

	mbox, err := c.Select(config.Mailbox, false)
	if err != nil {
		log.Fatal(err)
	}

	criteria := imap.NewSearchCriteria()
	criteria.Uid = new(imap.SeqSet)
	criteria.Uid.AddRange(lastUID+1, mbox.UidNext)

	uids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	if len(uids) == 0 {
		fmt.Printf("No new messages found!\n")
		return nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)

	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchRFC822Text, imap.FetchUid}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	var newMessages []*imap.Message
	for msg := range messages {
		// sender check
		sender := msg.Envelope.From[0].MailboxName + "@" + msg.Envelope.From[0].HostName
		if sender == config.ApprovedSender {
			newMessages = append(newMessages, msg)
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched %d new messages!\n", len(newMessages))

	os.WriteFile(".cache", []byte(fmt.Sprintf("%d", mbox.UidNext-1)), 0o644)

	return newMessages
}

func ProcessEmails(messages []*imap.Message, contentDir string) {
	for _, msg := range messages {

		subject := msg.Envelope.Subject
		if subject == "" {
			subject = "(No subject)"
		}
		date := msg.Envelope.Date.Format("2006-01-02 15:04:05")

		body, err := GetMessageBody(msg)
		if err != nil {
			log.Printf("Error getting message body: %v", err)
			continue
		}

		content := fmt.Sprintf("%s\n%s\n%s", subject, date, body)
		filename := fmt.Sprintf(contentDir+"/posts/%d.md", msg.Uid)
		os.WriteFile(filename, []byte(content), 0o644)
	}
}

// TODO: replace this with a munpack style library
// or at least make it MIME compliant. this seems way too fkn delicate
func GetMessageBody(msg *imap.Message) (string, error) {
	var content string
	for _, part := range msg.Body {
		entity, err := message.New(message.Header{}, part)
		if err != nil {
			log.Printf("Error creating entity: %v", err)
			continue
		}

		mimeType, _, err := entity.Header.ContentType()
		if err != nil {
			log.Printf("Error getting MIME type: %v", err)
			continue
		}
		if strings.ToLower(mimeType) == "text/plain" {
			bodyBytes, err := io.ReadAll(entity.Body)
			if err != nil {
				log.Printf("Error reading plain text body: %v", err)
				continue
			}
			content = string(bodyBytes)
			break
		}
	}

	// everything after first "text/plain"
	content = strings.Split(content, "\nContent-Type: text/plain; charset=\"UTF-8\"")[1]
	// everything before first mime boundary in remaining content
	content = strings.Split(content, "\n--0000")[0]

	// remove this shit
	content = strings.ReplaceAll(content, "Content-Transfer-Encoding: quoted-printable", "")

	reader := quotedprintable.NewReader(strings.NewReader(content))
	var buf bytes.Buffer
	_, err := io.Copy(&buf, reader)
	if err != nil {
		return "", err
	}
	content = buf.String()

	// trim spaces and remove empty lines
	content = regexp.MustCompile(`(?m)^\s*$`).ReplaceAllString(content, "")
	content = strings.TrimSpace(content)

	if content == "" {
		return "", fmt.Errorf("missing content!")
	}

	return content, nil
}

func GenerateHTML(config Config, contentDir string, outputDir string) {
	files, _ := filepath.Glob(contentDir + "/posts/*.md")
	sort.Slice(files, func(i, j int) bool {
		numI, _ := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(files[i], ".md"), contentDir+"/posts/"))
		numJ, _ := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(files[j], ".md"), contentDir+"/posts/"))
		return numI > numJ
	})

	var listItems []string
	for _, file := range files {
		content, _ := os.ReadFile(file)
		lines := strings.SplitN(string(content), "\n", 3)
		subject, date := lines[0], lines[1]

		// generate post HTML
		md := goldmark.New()
		// TODO: Unsafe/html
		var buf strings.Builder
		md.Convert([]byte(lines[2]), &buf)
		md.Renderer()
		postHTML := strings.Replace(readTemplate("post", contentDir), "%SUBJECT%", subject, -1)
		postHTML = strings.Replace(postHTML, "%DATE%", date, -1)
		postHTML = strings.Replace(postHTML, "%BODY%", buf.String(), -1)
		postHTML = strings.Replace(postHTML, "%BLOG_NAME%", config.BlogName, -1)

		postFilename := fmt.Sprintf(outputDir+"/posts/%s%s.html", subject, date)
		os.WriteFile(postFilename, []byte(postHTML), 0o644)

		// generate list item
		listItem := strings.Replace(readTemplate("listitem", contentDir), "%SUBJECT%", subject, -1)
		listItem = strings.Replace(listItem, "%DATE%", date, -1)
		listItem = strings.Replace(listItem, "%HYPERLINK%", fmt.Sprintf("./posts/%s%s.html", subject, date), -1)
		listItems = append(listItems, listItem)
	}

	// generate index
	indexHTML := strings.Replace(readTemplate("index", contentDir), "%BLOG_NAME%", config.BlogName, -1)
	indexHTML = strings.Replace(indexHTML, "%LIST%", strings.Join(listItems, ""), -1)
	os.WriteFile(outputDir+"/index.html", []byte(indexHTML), 0o644)
}

func readTemplate(name string, contentDir string) string {
	content, _ := os.ReadFile(fmt.Sprintf(contentDir+"/templates/%s.template.html", name))
	return string(content)
}
