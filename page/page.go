package page

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/biz/templates"
	"github.com/edataforms/pkg/session"
	"github.com/edataforms/pkg/util/utilstrings"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// session keys
var (
	InfoMessage  = "InfoMessage"
	ErrorMessage = "ErrorMessage" // key used to hold an error message that should be displayed to a user
	FormErrors   = "FormErrors"   // key used to hold form errors to be displayed to a user
	FormValues   = "FormValues"
	GroupValues  = "GroupValues"
)

// Render
func Render(ctx *gin.Context, view string, data interface{}) {
	templates.MustExecute(ctx.Writer, Layout.Suffix("wrapper"), view, data)
}

/*
type FormField struct {
	Name       string
	Label      string
	Value      string              // Input values
	Error      string              // Input error
	InputAttrs map[string]string   // Attributes
	InputProps map[string]struct{} // Properties
	WrapAttrs  map[string]string   // Attributes that should be assigned to the wrapper div
	WrapProps  map[string]string   // Properties that should be assigned to the wrapper div

	Group string // for inputs that belong to an array/group

	// can be used to create the html for the input
	Func func(FormField) template.HTML
}

// Attributes
func (f *FormField) Attributes() template.HTMLAttr {
	str := ""
	for k, v := range f.Attr {
		str += fmt.Sprintf(`%s="%v" `)
	}
	for k := range f.Props {
		str += k + " "
	}
	if len(str) > 0 {
		str = " " + str
	}

	return str
}
*/

// Page represents the basic elements of an html page
type Page struct {
	Header                   *Header
	Nav                      []Link
	SubNav                   []Link
	Title                    string
	Links                    []string
	Scripts                  []string
	ScriptsNoBust            []string
	ScriptsNoBustPostScripts []string
	InfoMessage              string
	ErrorMessage             string
	FormErrors               map[string]string
	FormValues               map[string]string
	GroupValues              map[string][]string
	BodyClass                string
	//	FormFields  map[string]FormField

	FaviconHTML  template.HTML
	CollapseMenu bool
	GoBack       bool
	BreadCrumbs  []BreadCrumb
}

//BreadCrumb is used to add a navigational link to the top of the content
type BreadCrumb struct {
	Label string
	Link  string
}

// AddBreadCrumb adds a BreadCrumb to p
func (p *Page) AddBreadCrumb(label, link string) *Page {
	p.BreadCrumbs = append(p.BreadCrumbs, BreadCrumb{
		Label: label,
		Link:  link,
	})
	return p
}

// ExistingValues is used to check if the page has pre-existing form values or form errors
func (p *Page) ExistingValues() bool {
	return len(p.FormErrors) > 0 || len(p.FormValues) > 0
}

// AddScript adds a script to the page if it does not already exist
func (p *Page) AddScript(scripts ...string) {
	for i := 0; i < len(scripts); i++ {
		p.Scripts = utilstrings.AddUnique(p.Scripts, scripts[i])
	}
}

// AddScriptNoBust adds a script to the page but does not add the cache buster
func (p *Page) AddScriptNoBust(scripts ...string) {
	for i := 0; i < len(scripts); i++ {
		p.ScriptsNoBust = utilstrings.AddUnique(p.ScriptsNoBust, scripts[i])
	}
}

// AddScriptNoBustPostScripts adds a script to the page but does not add the cache buster
func (p *Page) AddScriptNoBustPostScripts(scripts ...string) {
	for i := 0; i < len(scripts); i++ {
		p.ScriptsNoBustPostScripts = utilstrings.AddUnique(p.ScriptsNoBustPostScripts, scripts[i])
	}
}

// AddLink adds a link to the Page if it does not already exists
func (p *Page) AddLink(links ...string) {
	for i := 0; i < len(links); i++ {
		p.Links = utilstrings.AddUnique(p.Links, links[i])
	}
}

func (p *Page) HydrateFromSession(s *session.Session) {
	p.FormErrors = merge(p.FormErrors, GetFormErrors(s))
	p.FormValues = merge(p.FormValues, GetFormValues(s))
	if len(p.ErrorMessage) == 0 {
		p.ErrorMessage = GetErrorMessage(s)
	}
	if len(p.InfoMessage) == 0 {
		p.InfoMessage = GetInfoMessage(s)
	}
	if p.FormValues == nil {
		p.FormValues = map[string]string{}
	}

	if p.GroupValues == nil {
		p.GroupValues = map[string][]string{}
	}
	for k, v := range GetGroupValues(s) {
		p.GroupValues[k] = v
	}
}

type HeaderLogo struct {
	Img    template.URL
	Height string
	Width  string
}

type Header struct {
	Logo  HeaderLogo
	Title string
	Nav   []Link
}

type Link struct {
	Links           []Link
	Href            string
	Name            string
	IsActive        bool
	IsPartialActive bool
	Title           string
	Attrs           map[string]string
}

// SetActive sets the active link
func SetActive(r *http.Request, links []Link) []Link {
	nl := make([]Link, len(links))
	for i := 0; i < len(links); i++ {
		nl[i] = links[i]
		if nl[i].Href == r.URL.Path {
			nl[i].IsActive = true
		}
	}

	return nl
}

// FormError represents an error that has occured for a form field
type FormError struct {
	Field   string
	Message string
}

// SetErrors sets an error message and form errors on to the user's session
func SetErrors(s *session.Session, message string, errs map[string]string) {
	s.Data[ErrorMessage] = message
	SetFormErrors(s, errs)
	s.ShouldSave = true
}

// SetInfoMessage adds an info message to the user's session
func SetInfoMessage(s *session.Session, message string) {
	s.Data[InfoMessage] = message
	s.ShouldSave = true
}

func SetFormError(s *session.Session, key, value string) {
	SetFormErrors(s, map[string]string{key: value})
}

// setFormErrors adds form errors to the user's session
func SetFormErrors(s *session.Session, errs map[string]string) {
	defer func() {
		s.ShouldSave = true
	}()

	fv, ok := s.Data[FormErrors]
	if !ok {
		s.Data[FormErrors] = errs
		return
	}

	m, ok := fv.(map[string]interface{})
	if !ok {
		s.Data[FormErrors] = errs
		return
	}

	for k, v := range errs {
		m[k] = v
	}

	s.Data[FormErrors] = m
}

// SetFormValues adds form values to the session
func SetFormValues(s *session.Session, values map[string]string) {
	fv, ok := s.Data[FormValues]
	if !ok {
		s.Data[FormValues] = values
		return
	}

	m, ok := fv.(map[string]string)
	if !ok {
		s.Data[FormValues] = values
		return
	}

	for k, v := range values {
		m[k] = v
	}

	s.Data[FormValues] = m

	s.ShouldSave = true
}

// SetFormValue sets a single key/value form value
func SetFormValue(s *session.Session, key, value interface{}) {
	SetFormValues(s, map[string]string{fmt.Sprint(key): fmt.Sprint(value)})
}

// SetFormArrayValue sets a single key/value form value
func SetFormArrayValue(s *session.Session, key, value, id interface{}) {
	SetFormValues(s, map[string]string{fmt.Sprintf("%v:%v", key, id): fmt.Sprint(value)})
}

// SetFormArrayError sets a single key/value form error
func SetFormArrayError(s *session.Session, key, value, id interface{}) {
	SetFormErrors(s, map[string]string{fmt.Sprintf("%v:%v", key, id): fmt.Sprint(value)})
}

// SetGroupValue is used in conjunction with the FieldGroup template function to group
// related fields in an array
func SetGroupValue(s *session.Session, key string, id interface{}) {
	str := fmt.Sprint(id)
	SetFormValue(s, key+":"+str, str)
}

// SetGroup is used to save an ordered list of keys that can be looped to look up other keys belonging
// to the same group
func SetGroup(s *session.Session, key string, id interface{}) {
	s.ShouldSave = true

	gv, ok := s.Data[GroupValues]
	if !ok {
		s.Data[GroupValues] = map[string][]string{
			key: []string{fmt.Sprint(id)},
		}
		return
	}
	m, ok := gv.(map[string][]string)
	if !ok {
		s.Data[GroupValues] = map[string][]string{
			key: []string{fmt.Sprint(id)},
		}
		return
	}

	m[key] = append(m[key], fmt.Sprint(id))

	s.Data[GroupValues] = m
}

// SetErrorMessage adds an error message to the user's session
func SetErrorMessage(s *session.Session, message string) {
	s.Data[ErrorMessage] = message
	s.ShouldSave = true
}

// GetInfoMessage gets the InfoMessage from the session. If it's not found an empty string is returned.
// If it is found it is returned and the InfoMessage is removed from the session
func GetInfoMessage(s *session.Session) string {
	v, ok := s.Data[InfoMessage]
	if !ok {
		return ""
	}

	s.ShouldSave = true
	delete(s.Data, InfoMessage)

	str, ok := v.(string)
	if !ok {
		logrus.WithField("type", fmt.Sprintf("%T %v", v, v)).Error("invalid InfoMessage stored in session")
		return ""
	}

	return str
}

// GetErrorMessage gets an error message from the user's session if it exists it is removed from the users session
func GetErrorMessage(s *session.Session) string {
	v, ok := s.Data[ErrorMessage]
	if !ok {
		return ""
	}

	s.ShouldSave = true
	delete(s.Data, ErrorMessage)

	str, ok := v.(string)
	if !ok {
		logrus.WithField("type", fmt.Sprintf("%T %v", v, v)).Error("invalid ErrorMessage stored in session")
		return ""
	}

	return str
}

func GetGroupValues(s *session.Session) map[string][]string {
	v, ok := s.Data[GroupValues]
	if !ok {
		return nil
	}

	// remove form errors from session data
	delete(s.Data, GroupValues)
	s.ShouldSave = true

	values, ok := v.(map[string]interface{})
	if !ok {
		if m, ok := v.(map[string][]string); ok {
			return m
		}
		logrus.WithFields(logrus.Fields{
			"type": fmt.Sprintf("%T", v),
		}).Error("page: invalid GroupValues stored in session")
		return nil
	}

	m := map[string][]string{}

	for k, v := range values {
		i, ok := v.([]interface{})
		if !ok {
			logrus.WithFields(logrus.Fields{
				"type": fmt.Sprintf("%T", v),
			}).Error("page: invalid GroupValues")
			return nil
		}
		for _, s := range i {
			m[k] = append(m[k], fmt.Sprint(s))
		}
	}

	return m
}

// GetFormValues gets the form values stored in the user's session
func GetFormValues(s *session.Session) map[string]string {
	v, ok := s.Data[FormValues]
	if !ok {
		return nil
	}

	// remove form errors from session data
	delete(s.Data, FormValues)
	s.ShouldSave = true

	values, ok := v.(map[string]interface{})
	if !ok {
		if m, ok := v.(map[string]string); ok {
			return m
		}

		logrus.WithFields(logrus.Fields{
			"type": fmt.Sprintf("%T", v),
		}).Error("page: invalid FormValues stored in session")
		return nil
	}

	return msiTomss(values)
}

// GetFormErrors gets the user's form errors from the session, if any exist they are removed from the user's session
func GetFormErrors(s *session.Session) map[string]string {
	v, ok := s.Data[FormErrors]
	if !ok {
		return nil
	}

	// remove form errors from session data
	delete(s.Data, FormErrors)
	s.ShouldSave = true

	fe, ok := v.(map[string]interface{})
	if !ok {
		if m, ok := v.(map[string]string); ok {
			return m
		}
		logrus.WithFields(logrus.Fields{
			"type": fmt.Sprintf("%T", v),
		}).Error("page: invalid FormErrors stored in session")
		return nil
	}

	return msiTomss(fe)
}

// merge adds values from m2 to m1
func merge(m1, m2 map[string]string) map[string]string {
	if m1 == nil {
		m1 = map[string]string{}
	}
	for k, v := range m2 {
		m1[k] = v
	}

	return m1
}

func msiTomss(msi map[string]interface{}) map[string]string {
	mss := map[string]string{}
	for k, v := range msi {
		if s, ok := v.(string); ok {
			mss[k] = s
		} else {
			logrus.WithFields(logrus.Fields{
				"type":  fmt.Sprintf("%T", v),
				"value": v,
			}).Error("expected value to be a string")
		}
	}
	return mss
}
