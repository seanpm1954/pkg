package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/biz/templates"
	"github.com/edataforms/pkg/config"
	"github.com/edataforms/pkg/defaultassets"
	"github.com/edataforms/pkg/errorpages"
	"github.com/edataforms/pkg/health"
	"github.com/edataforms/pkg/log"
	"github.com/edataforms/pkg/page"
	"github.com/edataforms/pkg/robots"
	"github.com/edataforms/pkg/session"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
)

var (
	SessionCtxKey = session.CtxKey
	PageCtxKey    = "Page" // CtxKey represents where the page will be stored on the request's context
)

// Default adds commen middleware and routes
//
//  Middleware:
//		gzip,
//		panic recovery,
//		session,
//		page,
//
//	Routes:
//		/healthz
//		Not found
//		/robots.txt
func Default(e *gin.Engine) {
	e.Use(
		Logger,
		gzip.Gzip(gzip.DefaultCompression),
		Panic,
		Session,
		Page(DefaultPage),
		TooManySessions,
	)

	health.Routes(e)

	// handle 404 pages
	e.NoRoute(errorpages.NotFoundHandler)

	// Instruct search engines not to index
	robots.DontIndex(e)
}

// DefaultPage is used to setup basic css and javascript files to the page
var DefaultPage = &page.Page{
	Scripts: defaultassets.Scripts,
	Links:   defaultassets.Links,
	Header:  &page.Header{},
}

// Nav adds the nav menu based on the first role found that matches a menu set
func Nav(ctx *gin.Context) {
	if isAsset(ctx) {
		return
	}
	if ctx.Request.Method != "GET" {
		ctx.Next()
		return
	}

	s := SessionFromCtx(ctx)
	p := PageFromCtx(ctx)

	r := s.MenuType
	if r == "super-admin" {
		r = "admin"
	}

	n := page.GetMenu(r, ctx.Request)
	if n == nil {
		return
	}

	p.Header.Nav = page.GetMenu(r, ctx.Request)
	p.Nav = p.Header.Nav
}

// Logger logs the http request and adds a logger to the context with information about the request
func Logger(ctx *gin.Context) {
	if isAsset(ctx) {
		return
	}

	start := time.Now().UTC()

	//s := SessionFromCtx(ctx)
	f := logrus.Fields{}

	// this will be unique per request
	id, err := session.GenerateRandomString(10)
	if err == nil {
		f["req_id"] = id
		ctx.Writer.Header().Set("x-req-id", id)
	}

	f["req_url"] = ctx.Request.URL.String()
	f["req_referer"] = ctx.Request.Referer()
	f["req_ip"] = ctx.Request.RemoteAddr
	f["req_method"] = ctx.Request.Method
	f["req_useragent"] = ctx.Request.UserAgent()
	f["req_handler"] = path.Base(ctx.HandlerName())

	l := logrus.WithFields(f)

	ctx.Set(log.LoggerCtxKey, l)

	ctx.Next()

	l = LoggerFromCtx(ctx)
	l.WithFields(logrus.Fields{
		"response_code":   ctx.Writer.Status(),
		"response_size":   ctx.Writer.Size(),
		"req_duration_ms": time.Since(start).Nanoseconds() / int64(time.Millisecond),
	}).Info("http_request")

	if config.Conf.Env == "dev" {
		fmt.Println("")
	}
}

// RenderMiddleware is a helper function used to render a view. ctx.Keys will be used
// as the data in the template
func Render(baseView, view string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		templates.MustExecute(ctx.Writer, baseView, view, ctx.Keys)
	}
}

func TooManySessions(ctx *gin.Context) {
	if _, ok := ctx.Get("Sessions"); !ok {
		ctx.Next()
		return
	}

	session.TooManySessions(ctx)
	s := SessionFromCtx(ctx)
	s.Save()
}

// Session adds a session to the context
//
// THINK: Not sure if I want to move this func to the session pkg
func Session(ctx *gin.Context) {
	if isAsset(ctx) {
		return
	}

	s := session.New()
	s.Useragent = ctx.Request.UserAgent()
	s.IP = ctx.Request.RemoteAddr

	defer func() {
		prepSession(ctx, s)
	}()

	q, err := session.NewQuery(ctx)
	if err != nil {
		return
	}

	if len(q.SessionID) > 0 {
		s.ID = q.SessionID
	}

	if session.NumberOfSessions != -1 {
		q.Multiple = true
	}

	ss, err := session.GetAll(q)
	if err != nil {
		return
	}

	// find this user's actual session
	remove := -1
	for i := 0; i < len(ss); i++ {
		if ss[i].ID == q.SessionID {
			s = ss[i]
			remove = i
			break
		}
	}

	// if we have unlimited sessions or we are under the limit move on
	if session.NumberOfSessions == -1 || session.NumberOfSessions >= len(ss) {
		return
	}

	// remove the found session
	if remove != -1 {
		ss = append(ss[:remove], ss[remove+1:]...)
	}

	// check if we are deleting a session
	if ctx.Request.URL.Path == "/sessions/delete" {
		return
	}

	ctx.Set("Sessions", ss)
}

func prepSession(ctx *gin.Context, s *session.Session) {
	ctx.Set(SessionCtxKey, s)

	// update logger with session info
	l := LoggerFromCtx(ctx)
	l = l.WithFields(logrus.Fields{
		"session_id":       s.ID,
		"session_user_id":  s.UserID,
		"session_username": s.Username,
	})
	ctx.Set(log.LoggerCtxKey, l)

	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}

	ctx.Next()

	s.Save()
}

func Links(links ...string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != "GET" || isAsset(ctx) {
			return
		}
		p := PageFromCtx(ctx)
		p.AddLink(links...)
	}
}

func Scripts(scripts ...string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != "GET" || isAsset(ctx) {
			ctx.Next()
			return
		}
		p := PageFromCtx(ctx)
		p.AddScript(scripts...)
	}
}

// FormAssets adds appropriate JS and CSS needed to perform common frontend form tasks.
// example, date picker
func FormAssets(ctx *gin.Context) {
	p := PageFromCtx(ctx)
	p.AddLink(
		"/static/lib/pikaday/css/pikaday.css",
	)
	p.AddScript(
		"/static/lib/pikaday/js/pikaday.js",
		"/static/lib/fuse/fuse.min.js",
		"/static/js/forms.js",
	)
}

// NOTE: this is temporary until we move out the trpipe specific code
func FormAssetsV2(ctx *gin.Context) {
	p := PageFromCtx(ctx)
	p.AddLink(
		"/static/lib/pikaday/css/pikaday.css",
	)
	p.AddScript(
		"/static/lib/pikaday/js/pikaday.js",
		"/static/lib/fuse/fuse.min.js",
		"/static/lib/edf/form.js",
	)
}

// Page middleware adds a Page object to the request's context. Form values, errors and messages are also added to the page. The app's menu is also setup in the middelware
func Page(op *page.Page) func(*gin.Context) {
	return func(ctx *gin.Context) {
		if isAsset(ctx) {
			return
		}

		// only setup the Page on get requests
		if ctx.Request.Method != "GET" {
			ctx.Next()
			return
		}

		// copy page - could be a better way
		p := &page.Page{}
		*p = *op
		copy(p.Links, op.Links)
		copy(p.Nav, op.Nav)
		copy(p.Scripts, op.Scripts)
		copy(p.SubNav, op.SubNav)
		*p.Header = *op.Header
		copy(p.Header.Nav, op.Header.Nav)

		s := SessionFromCtx(ctx)
		p.HydrateFromSession(s)

		ctx.Set("Req", ctx.Request)
		ctx.Set("Ctx", ctx)
		ctx.Set(PageCtxKey, p)
		ctx.Next()
	}
}

// PageFromCtx returns the page stored on the context
func PageFromCtx(ctx *gin.Context) *page.Page {
	v, ok := ctx.Get(PageCtxKey)
	if !ok {
		panic(fmt.Sprintf("page not found on context: %s", ctx.Request.URL))
	}

	p, ok := v.(*page.Page)
	if !ok {
		panic("invalid Page stored on context")
	}

	return p
}

// LoggerFromCtx returns the Logger stored on a gin Context
func LoggerFromCtx(ctx *gin.Context) *logrus.Entry {
	v, ok := ctx.Get(log.LoggerCtxKey)
	if !ok {
		panic("logger not found. You are missing middleware")
	}

	l, ok := v.(*logrus.Entry)
	if !ok {
		panic("invalid logger stored on context")
	}

	return l
}

// SessionFromCtx returns the session stored on a gin Context
func SessionFromCtx(ctx *gin.Context) *session.Session {
	v, ok := ctx.Get(SessionCtxKey)
	if !ok {
		panic("session not found. You are missing middleware")
	}

	s, ok := v.(*session.Session)
	if !ok {
		panic("invalid session stored on context")
	}

	return s
}

// Panic middleware catches all panics and serves up an internal server error page
func Panic(ctx *gin.Context) {
	lg := log.New(ctx)
	defer func() {
		if err := recover(); err != nil {
			lg.WithFields(logrus.Fields{
				"stack": string(Stack(2)),
				"error": err,
			}).Error("panic recovery")
			errorpages.InternalServerError(ctx)
		}
	}()
	ctx.Next()
}

// stack returns a nicely formated stack frame, skipping skip frames
func Stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func isAsset(ctx *gin.Context) bool {
	// don't do anything on assets
	if strings.HasPrefix(ctx.Request.URL.Path, "/assets/") {
		ctx.Next()
		return true
	}
	return false
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)
