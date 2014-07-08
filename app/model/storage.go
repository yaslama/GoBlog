package model

import (
	"encoding/json"
	"github.com/yaslama/GoBlog/app/utils"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

var (
	appVersion int
	// global data storage instance
	Storage *jsonStorage
	// global tmp data storage instance. Temp data are generated for special usages, will not backup.
	TmpStorage *jsonStorage
)

type jsonStorage struct {
	dir string
}

func (jss *jsonStorage) Init(dir string) {
	jss.dir = dir
}

func (jss *jsonStorage) Has(key string) bool {
	file := path.Join(jss.dir, key+".json")
	_, e := os.Stat(file)
	return e == nil
}

func (jss *jsonStorage) Get(key string, v interface{}) {
	file := path.Join(jss.dir, key+".json")
	bytes, e := ioutil.ReadFile(file)
	if e != nil {
		println("read storage '" + key + "' error")
		return
	}
	e = json.Unmarshal(bytes, v)
	if e != nil {
		println("json decode '" + key + "' error")
	}
}

func (jss *jsonStorage) Set(key string, v interface{}) {
	locker.Lock()
	defer locker.Unlock()

	bytes, e := json.Marshal(v)
	if e != nil {
		println("json encode '" + key + "' error")
		return
	}
	file := path.Join(jss.dir, key+".json")
	e = ioutil.WriteFile(file, bytes, 0777)
	if e != nil {
		println("write storage '" + key + "' error")
	}
}

func (jss *jsonStorage) Dir(name string) {
	os.MkdirAll(path.Join(jss.dir, name), os.ModePerm)
}

func writeDefaultData() {
	// write user
	u := new(User)
	u.Id = Storage.TimeInc(10)
	u.Name = "admin"
	u.Password = utils.Sha1("adminxxxxx")
	u.Nick = "Administrator"
	u.Email = "admin@example.com"
	u.Url = "http://example.com/"
	u.CreateTime = utils.Now()
	u.Bio = "This is site administrator，You can add some personal introduction，Support line breaks do not supportmarkdown"
	u.LastLoginTime = u.CreateTime
	u.Role = "ADMIN"
	Storage.Set("users", []*User{u})

	// write token
	Storage.Set("tokens", map[string]*Token{})

	// write contents
	a := new(Content)
	a.Id = Storage.TimeInc(9)
	a.Title = "Welcome to Fxh.Go"
	a.Slug = "welcome-fxh-go"
	a.Text = "If you have read this article,Indicates that you blog Has been successfully installed."
	a.Tags = []string{"Fxh.Go"}
	a.CreateTime = utils.Now()
	a.EditTime = a.CreateTime
	a.UpdateTime = a.CreateTime
	a.IsComment = true
	a.IsLinked = false
	a.AuthorId = u.Id
	a.Type = "article"
	a.Status = "publish"
	a.Format = "markdown"
	a.Template = "blog.html"
	a.Hits = 1
	// write comments
	co := new(Comment)
	co.Author = u.Nick
	co.Email = u.Email
	co.Url = u.Url
	co.Content = "Welcome to use Fxh.Go"
	co.Avatar = utils.Gravatar(co.Email, "50")
	co.Pid = 0
	co.Ip = "127.0.0.1"
	co.UserAgent = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.17 (KHTML, like Gecko) Chrome/24.0.1312.57 Safari/537.17"
	co.IsAdmin = true
	co.Id = Storage.TimeInc(7)
	co.CreateTime = utils.Now()
	co.Status = "approved"
	co.Cid = a.Id
	a.Comments = []*Comment{co}
	Storage.Set("content/article-"+strconv.Itoa(a.Id), a)

	// write pages
	p := new(Content)
	p.Id = a.Id + Storage.TimeInc(6)
	p.Title = "With regard to"
	p.Slug = "about-me"
	p.Text = "This page by Fxh.Go Create, This is just a test page."
	p.Tags = []string{}
	p.CreateTime = utils.Now()
	p.EditTime = p.CreateTime
	p.UpdateTime = p.UpdateTime
	p.IsComment = true
	p.IsLinked = true
	p.AuthorId = u.Id
	p.Type = "page"
	p.Status = "publish"
	p.Format = "markdown"
	p.Comments = make([]*Comment, 0)
	p.Template = "page.html"
	p.Hits = 1
	Storage.Set("content/page-"+strconv.Itoa(p.Id), p)
	p2 := new(Content)
	p2.Id = p.Id + Storage.TimeInc(6)
	p2.Title = "Buddy"
	p2.Slug = "friends"
	p2.Text = "This page by Fxh.Go Create, This is just a test page."
	p2.Tags = []string{}
	p2.CreateTime = utils.Now()
	p2.EditTime = p2.CreateTime
	p2.UpdateTime = p2.UpdateTime
	p2.IsComment = true
	p2.IsLinked = true
	p2.AuthorId = u.Id
	p2.Type = "page"
	p2.Status = "publish"
	p2.Format = "markdown"
	p2.Comments = make([]*Comment, 0)
	p2.Template = "page.html"
	p2.Hits = 1
	Storage.Set("content/page-"+strconv.Itoa(p2.Id), p2)

	// write new reader
	Storage.Set("readers", map[string]*Reader{})

	// write version
	v := new(version)
	v.Name = "Fxh.Go"
	v.BuildTime = utils.Now()
	v.Version = appVersion
	Storage.Set("version", v)

	// write settings
	s := map[string]string{
		"site_title":         "Fxh.Go",
		"site_sub_title":     "GoDevelopment of simple blog",
		"site_keywords":      "Fxh.Go,Golang,Blog",
		"site_description":   "GoSimple blog language development program",
		"site_url":           "http://localhost/",
		"article_size":       "4",
		"site_theme":         "default",
		"enable_go_markdown": "false",
		"c_footer_weibo":     "#",
		"c_footer_github":    "#",
		"c_footer_email":     "#",
		"c_home_avatar":      "/static/img/site.png",
		"c_footer_ga":        "<!-- google analytics or other -->",
	}
	Storage.Set("settings", s)

	// write files
	Storage.Set("files", []*File{})

	// write message
	Storage.Set("messages", []*Message{})

	// write navigators
	n := new(navItem)
	n.Order = 1
	n.Text = "Articles"
	n.Title = "Articles"
	n.Link = "/"
	n2 := new(navItem)
	n2.Order = 2
	n2.Text = "With regard to"
	n2.Title = "With regard to"
	n2.Link = "/about-me.html"
	n3 := new(navItem)
	n3.Order = 3
	n3.Text = "Buddy"
	n3.Title = "Buddy"
	n3.Link = "/friends.html"
	Storage.Set("navigators", []*navItem{n, n2, n3})

	// write default tmp data
	writeDefaultTmpData()
}

func writeDefaultTmpData() {
	TmpStorage.Set("contents", make(map[string][]int))
}

func loadAllData() {
	loadVersion()
	LoadSettings()
	LoadNavigators()
	LoadUsers()
	LoadTokens()
	LoadContents()
	LoadMessages()
	LoadReaders()
	LoadComments()
	LoadFiles()
}

// TimeInc returns time step value devided by d int with time unix stamp.
func (jss *jsonStorage) TimeInc(d int) int {
	return int(utils.Now())%d + 1
}

// Init does model initialization.
// If first run, write default data.
// v means app.Version number. It's needed for version data.
func Init(v int) {
	appVersion = v
	Storage = new(jsonStorage)
	Storage.Init("data")
	TmpStorage = new(jsonStorage)
	TmpStorage.dir = "tmp/data"
	if !Storage.Has("version") {
		os.Mkdir(Storage.dir, os.ModePerm)
		os.Mkdir(path.Join(Storage.dir, "content"), os.ModePerm)
		os.Mkdir(path.Join(Storage.dir, "plugin"), os.ModePerm)
		writeDefaultData()
	}
}

// All loads all data from storage to memory.
// Start timers for content, comment and message.
func All() {
	loadAllData()
	// generate indexes
	SyncIndexes()
	// start model timer, do all timer stuffs
	StartModelTimer()
}

func SyncIndexes() {
	// generate indexes
	generatePublishArticleIndex()
	generateContentTmpIndexes()
}

// SyncAll writes all current memory data to storage files.
func SyncAll() {
	SyncContents()
	SyncMessages()
	SyncFiles()
	SyncReaders()
	SyncSettings()
	SyncNavigators()
	SyncTokens()
	SyncUsers()
	SyncVersion()
}
