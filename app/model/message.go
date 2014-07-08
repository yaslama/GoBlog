package model

import (
	"github.com/yaslama/GoBlog/app/utils"
	"strings"
)

var (
	messages         []*Message
	messageMaxId     int
	messageGenerator map[string]func(v interface{}) string
)

func init() {
	messageGenerator = make(map[string]func(v interface{}) string)
	messageGenerator["comment"] = generateCommentMessage
	messageGenerator["backup"] = generateBackupMessage
}

type Message struct {
	Id         int
	Type       string
	CreateTime int64
	Data       string
	IsRead     bool
}

func CreateMessage(tp string, data interface{}) *Message {
	m := new(Message)
	m.Type = tp
	m.Data = messageGenerator[tp](data)
	if m.Data == "" {
		println("message generator returns empty")
		return nil
	}
	m.CreateTime = utils.Now()
	m.IsRead = false
	messageMaxId += Storage.TimeInc(3)
	m.Id = messageMaxId
	messages = append([]*Message{m}, messages...)
	SyncMessages()
	return m
}

func SetMessageGenerator(name string, fn func(v interface{}) string) {
	messageGenerator[name] = fn
}

func GetMessage(id int) *Message {
	for _, m := range messages {
		if m.Id == id {
			return m
		}
	}
	return nil
}

func GetUnreadMessages() []*Message {
	ms := make([]*Message, 0)
	for _, m := range messages {
		if m.IsRead {
			continue
		}
		ms = append(ms, m)
	}
	return ms
}

func GetMessages() []*Message {
	return messages
}

func GetTypedMessages(tp string, unread bool) []*Message {
	ms := make([]*Message, 0)
	for _, m := range messages {
		if m.Type == tp {
			if unread {
				if !m.IsRead {
					ms = append(ms, m)
				}
			} else {
				ms = append(ms, m)
			}
		}
	}
	return ms
}

func SaveMessageRead(m *Message) {
	m.IsRead = true
	SyncMessages()
}

func SyncMessages() {
	Storage.Set("messages", messages)
}

func LoadMessages() {
	messages = make([]*Message, 0)
	messageMaxId = 0
	Storage.Get("messages", &messages)
	for _, m := range messages {
		if m.Id > messageMaxId {
			messageMaxId = m.Id
		}
	}
}

func RecycleMessages() {
	for i, m := range messages {
		if m.CreateTime+3600*24*3 < utils.Now() {
			messages = messages[:i]
			return
		}
	}
}

func generateCommentMessage(co interface{}) string {
	c, ok := co.(*Comment)
	if !ok {
		return ""
	}
	cnt := GetContentById(c.Cid)
	s := ""
	if c.Pid < 1 {
		s = "<p>" + c.Author + "Students，In the article《" + cnt.Title + "》Leave a comment："
		s += utils.Html2str(c.Content) + "</p>"
	} else {
		p := GetCommentById(c.Pid)
		s = "<p>" + p.Author + "Students，In the article《" + cnt.Title + "》Comments："
		s += utils.Html2str(p.Content) + "</p>"
		s += "<p>" + c.Author + "Student response："
		s += utils.Html2str(c.Content) + "</p>"
	}
	return s
}

func generateBackupMessage(co interface{}) string {
	str := co.(string)
	if strings.HasPrefix(str, "[0]") {
		return "Backup failed: " + strings.TrimPrefix(str, "[0]") + "."
	}
	return "Backup station " + strings.TrimPrefix(str, "[1]") + " Success."
}

func startMessageTimer() {
	SetTimerFunc("message-sync", 9, func() {
		println("write messages in 1.5 hour timer")
		RecycleMessages()
		SyncMessages()
	})
}
