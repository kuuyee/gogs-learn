// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

import (
	"fmt"
	"time"

	"github.com/go-macaron/session"
	"gopkg.in/ini.v1"
)

type Scheme string

const (
	HTTP  Scheme = "http"
	HTTPS Scheme = "https"
	FCGI  Scheme = "fcgi"
)

type LandingPage string

const (
	LANDING_PAGE_HOME    LandingPage = "/"
	LANDING_PAGE_EXPLORE LandingPage = "/explore"
)

var (
	// Build information
	BuildTime    string
	BuildGitHash string

	// App settings
	AppVer      string
	AppName     string
	AppUrl      string
	AppSubUrl   string
	AppPath     string
	AppDataPath = "data"

	// Server settings
	Protocol           Scheme
	Domain             string
	HttpAddr, HttpPort string
	LocalURL           string
	DisableSSH         bool
	StartSSHServer     bool
	SSHDomain          string
	SSHPort            int
	SSHRootPath        string
	OfflineMode        bool
	DisableRouterLog   bool
	CertFile, KeyFile  string
	StaticRootPath     string
	EnableGzip         bool
	LandingPageUrl     LandingPage

	// Security settings
	InstallLock          bool
	SecretKey            string
	LogInRememberDays    int
	CookieUserName       string
	CookieRememberName   string
	ReverseProxyAuthUser string

	// Database settings
	UseSQLite3    bool
	UseMySQL      bool
	UsePostgreSQL bool
	UseTiDB       bool

	// Webhook settings
	Webhook struct {
		QueueLength    int
		DeliverTimeout int
		SkipTLSVerify  bool
		Types          []string
		PagingNum      int
	}

	// Repository settings
	Repository struct {
		AnsiCharset            string
		ForcePrivate           bool
		MaxCreationLimit       int
		PullRequestQueueLength int
	}
	RepoRootPath string
	ScriptType   string

	// UI settings
	ExplorePagingNum     int
	IssuePagingNum       int
	FeedMaxCommitNum     int
	AdminUserPagingNum   int
	AdminRepoPagingNum   int
	AdminNoticePagingNum int
	AdminOrgPagingNum    int

	// Markdown sttings
	Markdown struct {
		EnableHardLineBreak bool
	}

	// Picture settings
	PictureService   string
	AvatarUploadPath string
	GravatarSource   string
	DisableGravatar  bool

	// Log settings
	LogRootPath string
	LogModes    []string
	LogConfigs  []string

	// Attachment settings
	AttachmentPath         string
	AttachmentAllowedTypes string
	AttachmentMaxSize      int64
	AttachmentMaxFiles     int
	AttachmentEnabled      bool

	// Time settings
	TimeFormat string

	// Cache settings
	CacheAdapter  string
	CacheInternal int
	CacheConn     string

	// Session settings
	SessionConfig session.Options

	// Git settings
	Git struct {
		MaxGitDiffLines int
		GcArgs          []string `delim:" "`
	}

	// Cron tasks
	Cron struct {
		UpdateMirror struct {
			Enabled    bool
			RunAtStart bool
			Schedule   string
		} `ini:"cron.update_mirrors"`
		RepoHealthCheck struct {
			Enabled    bool
			RunAtStart bool
			Schedule   string
			Timeout    time.Duration
			Args       []string `delim:" "`
		} `ini:"cron.repo_health_check"`
		CheckRepoStats struct {
			Enabled    bool
			RunAtStart bool
			Schedule   string
		} `ini:"cron.check_repo_stats"`
	}

	// I18n settings
	Langs, Names []string
	dateLangs    map[string]string

	// Highlight settings are loaded in modules/template/hightlight.go

	// Other settings
	ShowFooterBranding    bool
	ShowFooterVersion     bool
	SupportMiniWinService bool

	// Global setting objects
	Cfg          *ini.File
	CustomPath   string // Custom directory path
	CustomConf   string
	ProdMode     bool
	RunUser      string
	IsWindows    bool
	HasRobotsTxt bool
)

// NewContext initializes configuration context.
// NOTE: do not print any log except error.
func NewContext() {
	fmt.Println("Gogs-Learn Runing...")
}
