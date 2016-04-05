// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

import (
	"fmt"
	"time"

	"github.com/kuuyee/gogs-learn/modules/log"
	"github.com/kuuyee/gogs-learn/modules/bindata"
	"github.com/go-macaron/session"
	"github.com/Unknwon/com"
	"gopkg.in/ini.v1"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"runtime"
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
	StaticRootPath     string //模板文件和静态文件的上级目录，默认为应用二进制所在的位置
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

// execPath returns the executable path.
func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func init()  {
	// 判断是否是window系统
	IsWindows = runtime.GOOS == "windows"
	log.NewLogger(0,"console",`{"level":0}`)

	var err error
	if AppPath,err = execPath(); err != nil {
		log.Fatal(4, "fail to get app path: %v\n", err)
	}
	log.Info("AppPath = %s",AppPath)

	// Note: we don't use path.Dir here because it does not handle case
	//	which path starts with two "/" in Windows: "//psf/Home/..."
	AppPath = strings.Replace(AppPath, "\\", "/", -1)
}

func WorkDir()(string,error)  {
	wd := os.Getenv("GOGS_WORK_DIR") //检查是否存在环境变量GOGS_WORK_DIR
	if len(wd) >0{
		return wd,nil
	}

	// 找不到GOGS_WORK_DIR，就取AppPath得值
	i:=strings.LastIndex(AppPath,"/")
	if i ==-1{
		return AppPath,nil
	}

	return AppPath[:i],nil
}

func forcePathSeparator(path string) {
	if strings.Contains(path, "\\") {
		log.Fatal(4, "Do not use '\\' or '\\\\' in paths, instead, please use '/' in all places")
	}
}

// NewContext 初始化配置文件上下文.
// NOTE: do not print any log except error.
func NewContext() {
	fmt.Println("Gogs-Learn Runing...")
	workDir,err := WorkDir() //取得目录
	if err!=nil{
		log.Fatal(4, "Fail to get work directory: %v", err)
	}

	Cfg,err = ini.Load(bindata.MustAsset("conf/app.ini"))
	if err != nil {
		log.Fatal(4, "Fail to parse 'conf/app.ini': %v", err)
	}

	CustomPath = os.Getenv("GOGS_CUSTOM")
	if len(CustomPath) == 0 {
		CustomPath = workDir + "/custom"
	}

	if len(CustomConf) == 0 {
		CustomConf = CustomPath + "/conf/app.ini"
	}

	if com.IsFile(CustomConf) {
		if err = Cfg.Append(CustomConf); err != nil {
			log.Fatal(4, "Fail to load custom conf '%s': %v", CustomConf, err)
		}
	} else {
		log.Warn("Custom config '%s' not found, ignore this if you're running first time", CustomConf)
	}
	Cfg.NameMapper = ini.AllCapsUnderscore

	homeDir, err := com.HomeDir()
	if err != nil {
		log.Fatal(4, "Fail to get home directory: %v", err)
	}
	homeDir = strings.Replace(homeDir, "\\", "/", -1)

	LogRootPath = Cfg.Section("log").Key("ROOT_PATH").MustString(path.Join(workDir, "log"))
	forcePathSeparator(LogRootPath)

	sec := Cfg.Section("server")
	StaticRootPath = sec.Key("STATIC_ROOT_PATH").MustString(workDir)
}


