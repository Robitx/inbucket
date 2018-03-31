package config

import (
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	prefix      = "inbucket"
	tableFormat = `Inbucket is configured via the environment. The following environment variables
can be used:

KEY	DEFAULT	DESCRIPTION
{{range .}}{{usage_key .}}	{{usage_default .}}	{{usage_description .}}
{{end}}`
)

var (
	// Version of this build, set by main
	Version = ""

	// BuildDate for this build, set by main
	BuildDate = ""
)

// Root wraps all other configurations.
type Root struct {
	LogLevel string `required:"true" default:"INFO" desc:"DEBUG, INFO, WARN, or ERROR"`
	SMTP     SMTP
	POP3     POP3
	Web      Web
	Storage  Storage
}

// SMTP contains the SMTP server configuration.
type SMTP struct {
	Addr            string        `required:"true" default:"0.0.0.0:2500" desc:"SMTP server IP4 host:port"`
	Domain          string        `required:"true" default:"inbucket" desc:"HELO domain"`
	DomainNoStore   string        `desc:"Load testing domain"`
	MaxRecipients   int           `required:"true" default:"200" desc:"Maximum RCPT TO per message"`
	MaxMessageBytes int           `required:"true" default:"10240000" desc:"Maximum message size"`
	StoreMessages   bool          `required:"true" default:"true" desc:"Store incoming mail?"`
	Timeout         time.Duration `required:"true" default:"300s" desc:"Idle network timeout"`
	Debug           bool          `ignored:"true"`
}

// POP3 contains the POP3 server configuration.
type POP3 struct {
	Addr    string        `required:"true" default:"0.0.0.0:1100" desc:"POP3 server IP4 host:port"`
	Domain  string        `required:"true" default:"inbucket" desc:"HELLO domain"`
	Timeout time.Duration `required:"true" default:"600s" desc:"Idle network timeout"`
	Debug   bool          `ignored:"true"`
}

// Web contains the HTTP server configuration.
type Web struct {
	Addr           string `required:"true" default:"0.0.0.0:9000" desc:"Web server IP4 host:port"`
	UIDir          string `required:"true" default:"ui" desc:"User interface dir"`
	GreetingFile   string `required:"true" default:"ui/greeting.html" desc:"Home page greeting HTML"`
	TemplateCache  bool   `required:"true" default:"true" desc:"Cache templates after first use?"`
	MailboxPrompt  string `required:"true" default:"@inbucket" desc:"Prompt next to mailbox input"`
	CookieAuthKey  string `desc:"Session cipher key (text)"`
	MonitorVisible bool   `required:"true" default:"true" desc:"Show monitor tab in UI?"`
	MonitorHistory int    `required:"true" default:"30" desc:"Monitor remembered messages"`
}

// Storage contains the mail store configuration.
type Storage struct {
	Type            string            `required:"true" default:"memory" desc:"Storage impl: file or memory"`
	Params          map[string]string `desc:"Storage impl parameters, see docs."`
	RetentionPeriod time.Duration     `required:"true" default:"24h" desc:"Duration to retain messages"`
	RetentionSleep  time.Duration     `required:"true" default:"50ms" desc:"Duration to sleep between mailboxes"`
	MailboxMsgCap   int               `required:"true" default:"500" desc:"Maximum messages per mailbox"`
}

// Process loads and parses configuration from the environment.
func Process() (*Root, error) {
	c := &Root{}
	err := envconfig.Process(prefix, c)
	c.SMTP.DomainNoStore = strings.ToLower(c.SMTP.DomainNoStore)
	return c, err
}

// Usage prints out the envconfig usage to Stderr.
func Usage() {
	tabs := tabwriter.NewWriter(os.Stderr, 1, 0, 4, ' ', 0)
	if err := envconfig.Usagef(prefix, &Root{}, tabs, tableFormat); err != nil {
		log.Fatalf("Unable to parse env config: %v", err)
	}
	tabs.Flush()
}
