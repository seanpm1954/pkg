package page

import (
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/biz/templates"
	"github.com/edataforms/pkg/html/htmlselect"
	"github.com/edataforms/pkg/util/utilnumbers"
)

func funcs() {
	templates.AddFunc("IsValid", IsValid)
	templates.AddFunc("FieldError", FieldError)
	templates.AddFunc("ArrayFieldError", ArrayFieldError)
	templates.AddFunc("TextField", TextField)
	templates.AddFunc("RequiredTextField", RequiredTextField)
	templates.AddFunc("HiddenField", HiddenField)
	templates.AddFunc("PositiveNumberField", PositiveNumberField)
	templates.AddFunc("NumberField", NumberField)
	templates.AddFunc("NumberFieldMinMax", NumberFieldMinMax)
	templates.AddFunc("TextAreaField", TextAreaField)
	templates.AddFunc("RequiredTextAreaField", RequiredTextAreaField)
	templates.AddFunc("TextAreaFieldReadOnly", TextAreaFieldReadOnly)
	templates.AddFunc("SelectField", SelectField)
	templates.AddFunc("SelectField4Col", SelectField4Col)
	templates.AddFunc("MultiSelectField", MultiSelectField)
	templates.AddFunc("SelectFieldWithDefault", SelectFieldWithDefault)
	templates.AddFunc("DateField", DateField)
	templates.AddFunc("NativeDateField", NativeDateField)
	templates.AddFunc("BoolCheckBox", BoolCheckBox)
	templates.AddFunc("ArrayCheckBox", ArrayCheckBox)
	templates.AddFunc("ArrayLabelFieldDefault", ArrayLabelFieldDefault)
	templates.AddFunc("SubmitButton", SubmitButton)
	templates.AddFunc("ValueExists", ValueExists)
	templates.AddFunc("FieldGroup", FieldGroup)
	templates.AddFunc("FormGroupValues", FormGroupValues)
	templates.AddFunc("ArrayFieldValue", ArrayFieldValue)
	templates.AddFunc("FieldValue", FieldValue)
	templates.AddFunc("IntSelectField", IntSelectField)
	templates.AddFunc("Options", htmlselect.Options)
	templates.AddFunc("StringSelectField", StringSelectField)
	templates.AddFunc("LabelAndField", LabelAndField)
	templates.AddFunc("LabelArrayField", LabelArrayField)
	templates.AddFunc("Join", Join)
	templates.AddFunc("LabelNameKey", LabelNameKey)
	templates.AddFunc("FieldOrderByDate", FieldOrderByDate)
	templates.AddFunc("RadioField", RadioField)
	templates.AddFunc("BasicRadioInput", BasicRadioInput)
	templates.AddFunc("RadioInputClass", RadioInputClass)
	templates.AddFunc("BoolValue", BoolValue)
	templates.AddFunc("CssClass", CssClass)
	templates.AddFunc("PhoneField", PhoneField)
	templates.AddFunc("KeyArrayID", KeyArrayID)
	templates.AddFunc("KeyNameLabel", KeyNameLabel)
	templates.AddFunc("NameValue", NameValue)
}

func PhoneField(p *Page, field interface{}) template.HTML {
	fo := convert(field)

	value := escapeField(p, fo.Key)

	return template.HTML(`
		<input class="phone-input ` + fo.CssClass + `" type="tel" id="` + fo.CssID + `" name="` + fo.Name + `" value="` + value + `">
	`)
}

func CssClass(typ, class string) string {
	switch typ {
	case "report":
		switch class {
		case "table":
			return " report-tbl "
		case "td-non-numeric":
			fallthrough
		case "th-non-numeric":
			return " report-tbl__text "
		default:
			return ""
		}
	case "default":
		fallthrough
	default:
		switch class {
		case "table":
			return " edf-table mdl-data-table mdl-js-data-table "
		case "td-non-numeric":
			fallthrough
		case "th-non-numeric":
			return " mdl-data-table__cell--non-numeric "
		default:
			return ""
		}
	}
}

func BoolValue(p *Page, field string) bool {
	return escapeField(p, field) == "true"
}

func RadioInputClass(label, name, value, class string) *FieldOptions {
	f := BasicRadioInput(label, name, value)
	f.CssClass = class
	return f
}

// NameValue is usefull in conjuction with setting hidden field values
func NameValue(name string, value interface{}) *FieldOptions {
	return &FieldOptions{
		Name:  name,
		Value: fmt.Sprint(value),
	}
}

func KeyNameLabel(key, name, label string) *FieldOptions {
	return &FieldOptions{
		Label:    label,
		Name:     name,
		CssClass: name,
		CssID:    strings.Replace(strings.ToLower(name+"-"+key), " ", "-", -1),
		Key:      key,
	}
}

func BasicRadioInput(label, name, value string) *FieldOptions {
	return &FieldOptions{
		Label:      label,
		Name:       name,
		InputValue: value,
		CssClass:   name,
		CssID:      strings.Replace(strings.ToLower(name+"-"+value), " ", "-", -1),
		Key:        name,
	}
}

func RadioField(p *Page, field interface{}) template.HTML {
	fo := convert(field)
	value := escapeField(p, fo.Key)

	checked := ""
	if fo.InputValue == value {
		checked = "checked"
	}

	return template.HTML(`
	<label class="mdl-radio mdl-js-radio mdl-js-ripple-effect" for="` + fo.CssID + `">
		<input type="radio" id="` + fo.CssID + `" class="mdl-radio__button ` + fo.CssClass + `" name="` + fo.Name + `" value="` + fo.InputValue + `"` + checked + ` />
		<span class="mdl-radio__label">` + fo.Label + `</span>
	</label>
	`)
}

func KeyArrayID(name string, id interface{}) *FieldOptions {
	key := fmt.Sprintf("%v:%v", name, id)
	return &FieldOptions{
		Name:     name,
		Key:      key,
		CssClass: name,
		CssID:    strings.Replace(key, ":", "-", -1),
	}
}

func HiddenField(p *Page, field interface{}) template.HTML {
	fo := convert(field)

	value := ""
	if len(fo.Value) > 0 {
		value = fo.Value
	} else {
		value = escapeField(p, fo.Key)
	}

	return template.HTML(`
		<input type="hidden" id="` + fo.CssID + `" name="` + fo.Name + `" value="` + value + `">
	`)
}

// Join is used to concat parts together with a given separator
func Join(sep string, parts ...string) string {
	return strings.Join(parts, sep)
}

func LabelNameKey(label, name, key string) *FieldOptions {
	return &FieldOptions{
		Label:    label,
		Name:     name,
		Key:      key,
		CssClass: name,
		CssID:    strings.Replace(key, ":", "-", -1),
	}
}

func LabelArrayField(label, field string, id interface{}) *FieldOptions {
	return LabelAndField(label, fmt.Sprintf("%v:%v", field, id))
}

func ArrayLabelFieldDefault(label, field string, deflt string, id interface{}) *FieldOptions {
	f := LabelAndField(label, fmt.Sprintf("%v:%v", field, id))
	f.Default = deflt
	return f
}

func LabelAndField(label, field string) *FieldOptions {
	l, f := getLookup(field)
	return &FieldOptions{
		Label:    label,
		Name:     f,
		Key:      l,
		CssClass: f,
		CssID:    strings.Replace(l, ":", "-", -1),
	}
}

// TODO: finish implementing everywhere - only using for label and name for now
type ToFieldOptions interface {
	ToFieldOptions() *FieldOptions
}

type FieldOptions struct {
	Label      string
	Name       string
	Key        string // key is used to lookup values from the Page, like errors and values
	CssClass   string
	CssID      string
	InputValue string // used for default values and radio/checkbox values
	Value      string
	Default    string
}

// IntSelectField is used to generate a select box with ints between the start and end parameters
func IntSelectField(p *Page, field interface{}, start, end int) template.HTML {
	options := htmlselect.IntOptions(utilnumbers.IntRange(start, end))
	return SelectField(p, field, options)
}

func StringSelectField(p *Page, field interface{}, strs ...string) template.HTML {
	options := htmlselect.StringOptions(strs)
	return SelectField(p, field, options)
}

// NOTE: this will have to be explained
func FieldGroup(p *Page, group string) []string {
	var flds []string
	for k, v := range p.FormValues {
		if strings.HasPrefix(k, group+":") {
			flds = append(flds, v)
		}
	}

	return flds
}

func FormGroupValues(p *Page, group string) []string {
	return p.GroupValues[group]
}

type dateItem struct {
	time  time.Time
	index string
}

type dateOrder []*dateItem

func (d dateOrder) Len() int           { return len(d) }
func (d dateOrder) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d dateOrder) Less(i, j int) bool { return d[i].time.Before(d[j].time) }

type dateOrderDesc struct{ dateOrder }

func (d dateOrderDesc) Less(i, j int) bool { return d.dateOrder[i].time.After(d.dateOrder[j].time) }

func FieldOrderByDate(p *Page, group, sortField, layout, dir string) []string {
	ids := FieldGroup(p, group)
	if len(ids) == 0 {
		return ids
	}

	dir = strings.ToLower(dir)

	ds := make([]*dateItem, len(ids))

	for i, id := range ids {
		ds[i] = &dateItem{index: id}

		s := escapeField(p, fmt.Sprintf("%s:%s", sortField, id))

		d, err := time.Parse(layout, s)
		if err != nil {
			continue
		}
		ds[i].time = d
	}

	switch dir {
	case "desc":
		sort.Sort(dateOrderDesc{ds})
	case "asc":
		fallthrough
	default:
		sort.Sort(dateOrder(ds))
	}

	sorted := make([]string, len(ds))
	for i, t := range ds {
		sorted[i] = t.index
	}

	return sorted
}

// ValueExists checks if a form field exists and has a value
func ValueExists(p *Page, field interface{}) bool {
	o := convert(field)
	v, ok := p.FormValues[o.Key]
	if !ok {
		return false
	}

	return len(v) > 0
}

// BoolCheckBox is used to create a single checkbox with no value - the server will have to check for "on" or "off"
func BoolCheckBox(p *Page, field interface{}, width string) template.HTML {
	fo := convert(field)
	width = treatWidth(width)

	checked := ""
	if v, ok := p.FormValues[fo.Key]; ok && v != "false" && v != "off" {
		checked = " checked "
	}

	return template.HTML(`
	<label class="mdl-checkbox mdl-js-checkbox mdl-js-ripple-effect" for="` + fo.CssID + `">
		<input ` + checked + ` value="on" type="checkbox" id="` + fo.CssID + `" name="` + fo.Name + `" class="mdl-checkbox__input ` + fo.CssClass + `">
		<span class="mdl-checkbox__label">` + fo.Label + `</span>
	</label>
	`)
}

// ArrayCheckBox is used to create a checkbox without a label and uses the id to lookup the value in Page.FormValues
func ArrayCheckBox(p *Page, field interface{}, idi interface{}) template.HTML {
	fo := convert(field)

	id := fmt.Sprint(idi)
	id = template.HTMLEscapeString(id)
	lookup := fmt.Sprintf("%v:%v", fo.Name, id)

	checked := ""
	if v, ok := p.FormValues[lookup]; ok && v != "false" && v != "off" {
		checked = " checked "
	}

	return template.HTML(`
	<label class="mdl-checkbox mdl-js-checkbox" for="` + fo.Name + `-` + id + `">
		<input ` + checked + ` type="checkbox" id="` + fo.Name + `-` + id + `" name="` + fo.Name + `" value="` + id + `" class="` + fo.CssClass + `-checkbox mdl-checkbox__input">
	</label>
	`)
}

// SubmitButton is a template function used to create a submit button for a form
func SubmitButton(name string) template.HTML {
	return template.HTML(`
	<button type="submit" class="mdl-button mdl-js-button mdl-button--raised mdl-button--accent">
		` + template.HTMLEscapeString(strings.ToUpper(name)) + `
	</button>
	`)
}

// DateField is used to add a date picker to a text field.
func DateField(p *Page, field interface{}, width string) template.HTML {
	fo := convert(field)

	width = treatWidth(width)

	is := IsValid(p, fo.Key)
	fe := string(FieldError(p, fo.Key))
	value := escapeField(p, fo.Key)

	return template.HTML(`
	<div class="` + is + `mdl-textfield mdl-js-textfield mdl-textfield--floating-label ` + width + `">
		<input class="date-picker mdl-textfield__input ` + fo.CssClass + `" type="text" id="` + fo.CssID + `" name="` + fo.Name + `" value="` + value + `">
		<label class="mdl-textfield__label" for="` + fo.CssID + `">` + fo.Label + `</label>
	` + fe + `
	</div>
	`)
}

func NativeDateField(p *Page, field interface{}, width string) template.HTML {
	fo := convert(field)

	width = treatWidth(width)

	is := IsValid(p, fo.Key)
	fe := string(FieldError(p, fo.Key))
	value := escapeField(p, fo.Key)

	// prevents a user from selecting a date from the year 1
	if value == "0001-01-01" {
		value = ""
	}

	return template.HTML(`
	<div class="` + is + `mdl-textfield mdl-js-textfield ` + width + `">
		<input class="mdl-textfield__input ` + fo.CssClass + `" type="date" id="` + fo.CssID + `" name="` + fo.Name + `" value="` + value + `">
	` + fe + `
	</div>
	`)
}

// TextField is a template function that is used to render an HTML text input and label.
// If the input is array the expected fieldName should in "<key>:<value>" format. This format facilitates input arrays
func TextField(p *Page, field interface{}, width string) template.HTML {
	fo := convert(field)

	width = treatWidth(width)

	is := IsValid(p, fo.Key)
	fe := string(FieldError(p, fo.Key))
	value := escapeField(p, fo.Key)

	if len(value) == 0 {
		value = fo.Default
	}

	return template.HTML(`
	<div class="` + is + `mdl-textfield mdl-js-textfield mdl-textfield--floating-label` + width + `">
		<input class="mdl-textfield__input ` + fo.CssClass + `" type="text" id="` + fo.CssID + `" name="` + fo.Name + `" value="` + value + `">
		<label class="mdl-textfield__label" for="` + fo.CssID + `">` + fo.Label + `</label>
	` + fe + `
	</div>
	`)
}

// TextField is a template function that is used to render an HTML text input and label.
// If the input is array the expected fieldName should in "<key>:<value>" format. This format facilitates input arrays
func RequiredTextField(p *Page, field interface{}, width string) template.HTML {
	fo := convert(field)

	width = treatWidth(width)

	is := IsValid(p, fo.Key)
	fe := string(FieldError(p, fo.Key))
	value := escapeField(p, fo.Key)

	return template.HTML(`
	<div class="` + is + `mdl-textfield mdl-js-textfield mdl-textfield--floating-label` + width + `">
		<input required class="mdl-textfield__input ` + fo.CssClass + `" type="text" id="` + fo.CssID + `" name="` + fo.Name + `" value="` + value + `">
		<label class="mdl-textfield__label" for="` + fo.CssID + `">` + fo.Label + `</label>
	` + fe + `
	</div>
	`)
}

func PositiveNumberField(p *Page, field interface{}, width string) template.HTML {
	fo := convert(field)

	width = treatWidth(width)

	is := IsValid(p, fo.Key)
	fe := string(FieldError(p, fo.Key))
	value := escapeField(p, fo.Key)

	return template.HTML(`
	<div class="` + is + `mdl-textfield mdl-js-textfield mdl-textfield--floating-label` + width + `">
		<input class="mdl-textfield__input ` + fo.CssClass + `" type="number" min="1" id="` + fo.CssID + `" name="` + fo.Name + `" value="` + value + `">
		<label class="mdl-textfield__label" for="` + fo.CssID + `">` + fo.Label + `</label>
	` + fe + `
	</div>
	`)
}

func NumberField(p *Page, field interface{}, width string) template.HTML {
	fo := convert(field)

	width = treatWidth(width)

	is := IsValid(p, fo.Key)
	fe := string(FieldError(p, fo.Key))
	value := escapeField(p, fo.Key)

	return template.HTML(`
	<div class="` + is + `mdl-textfield mdl-js-textfield mdl-textfield--floating-label` + width + `">
		<input class="mdl-textfield__input ` + fo.CssClass + `" type="number" min="0" id="` + fo.CssID + `" name="` + fo.Name + `" value="` + value + `">
		<label class="mdl-textfield__label" for="` + fo.CssID + `">` + fo.Label + `</label>
	` + fe + `
	</div>
	`)
}

func NumberFieldMinMax(p *Page, field interface{}, min, max int64, width string) template.HTML {
	fo := convert(field)

	width = treatWidth(width)

	is := IsValid(p, fo.Key)
	fe := string(FieldError(p, fo.Key))
	value := escapeField(p, fo.Key)

	return template.HTML(`
	<div class="` + is + `mdl-textfield mdl-js-textfield mdl-textfield--floating-label` + width + `">
		<input class="mdl-textfield__input ` + fo.CssClass + `" type="number" min="` + fmt.Sprint(min) + `" max="` + fmt.Sprint(max) + `" id="` + fo.CssID + `" name="` + fo.Name + `" value="` + value + `">
		<label class="mdl-textfield__label" for="` + fo.CssID + `">` + fo.Label + `</label>
	` + fe + `
	</div>
	`)
}

// TextAreaField is a template function that is used to render an HTML text input and label.
// If the input is array the expected fieldName should in "<key>:<value>" format. This format facilitates input arrays
func TextAreaField(p *Page, field interface{}, rows string, width string) template.HTML {
	fo := convert(field)

	rows = template.HTMLEscapeString(rows)
	width = treatWidth(width)

	is := IsValid(p, fo.Key)
	fe := string(FieldError(p, fo.Key))
	value := escapeField(p, fo.Key)

	return template.HTML(`
	<div class="` + is + `mdl-textfield mdl-js-textfield mdl-textfield--floating-label` + width + `">
		<textarea class="mdl-textfield__input ` + fo.CssClass + `" type="text" id="` + fo.CssID + `" name="` + fo.Name + `" rows="` + rows + `">` + value + `</textarea>
		<label class="mdl-textfield__label" for="` + fo.CssID + `">` + fo.Label + `</label>
	` + fe + `
	</div>
	`)
}

// TextAreaField is a template function that is used to render an HTML text input and label.
// If the input is array the expected fieldName should in "<key>:<value>" format. This format facilitates input arrays
func RequiredTextAreaField(p *Page, field interface{}, rows string, width string) template.HTML {
	fo := convert(field)

	rows = template.HTMLEscapeString(rows)
	width = treatWidth(width)

	is := IsValid(p, fo.Key)
	fe := string(FieldError(p, fo.Key))
	value := escapeField(p, fo.Key)

	return template.HTML(`
	<div class="` + is + `mdl-textfield mdl-js-textfield mdl-textfield--floating-label` + width + `">
		<textarea required class="mdl-textfield__input ` + fo.CssClass + `" type="text" id="` + fo.CssID + `" name="` + fo.Name + `" rows="` + rows + `">` + value + `</textarea>
		<label class="mdl-textfield__label" for="` + fo.CssID + `">` + fo.Label + `</label>
	` + fe + `
	</div>
	`)
}

// TextAreaFieldReadOnly is a template function that is used to render an HTML text input and label.
// If the input is array the expected fieldName should in "<key>:<value>" format. This format facilitates input arrays
func TextAreaFieldReadOnly(p *Page, field interface{}, rows string, width string) template.HTML {
	fo := convert(field)

	rows = template.HTMLEscapeString(rows)
	width = treatWidth(width)

	is := IsValid(p, fo.Key)
	fe := string(FieldError(p, fo.Key))
	value := escapeField(p, fo.Key)

	return template.HTML(`
	<div class="` + is + `mdl-textfield mdl-js-textfield mdl-textfield--floating-label` + width + `">
		<textarea readonly class="mdl-textfield__input ` + fo.CssClass + `" type="text" id="` + fo.CssID + `" name="` + fo.Name + `" rows="` + rows + `">` + value + `</textarea>
		<label class="mdl-textfield__label" for="` + fo.CssID + `">` + fo.Label + `</label>
	` + fe + `
	</div>
	`)
}

// SelectField is used to create a select field
func SelectField(p *Page, field interface{}, options []htmlselect.Option) template.HTML {
	fo := convert(field)

	id := template.HTMLEscapeString(p.FormValues[fo.Key])

	fe := string(FieldError(p, fo.Key))
	is := IsValid(p, fo.Key)

	slct := `<select id="` + fo.CssID + `" name="` + fo.Name + `" class="mdl-select__input ` + fo.CssClass + `">`
	if len(id) == 0 {
		slct += `<option value="0" selected disabled>` + fo.Label + "</option>"
	} else {
		slct += `<option value="0" disabled>` + fo.Label + "</option>"
	}

	for _, op := range options {
		value, label := op.OptionValue()
		value = template.HTMLEscapeString(value)
		label = template.HTMLEscapeString(label)
		selected := ""
		if value == id {
			selected = " selected"
		}
		slct += `<option value="` + value + `"` + selected + `>` + label + "</option>"
	}
	slct += "</select>"

	return template.HTML(`
			<div class="mdl-select mdl-js-select mdl-select--floating-label ` + is + ` mdl-cell mdl-cell--12-col">
				` + slct + `
				` + fe + `
			</div>
	`)
}

// SelectField4Col is used to create a select field
func SelectField4Col(p *Page, field interface{}, options []htmlselect.Option) template.HTML {
	fo := convert(field)

	id := template.HTMLEscapeString(p.FormValues[fo.Key])

	fe := string(FieldError(p, fo.Key))
	is := IsValid(p, fo.Key)

	slct := `<select id="` + fo.CssID + `" name="` + fo.Name + `" class="mdl-select__input ` + fo.CssClass + `">`
	if len(id) == 0 {
		slct += `<option value="0" selected>` + fo.Label + "</option>"
	} else {
		slct += `<option value="0">` + fo.Label + "</option>"
	}

	for _, op := range options {
		value, label := op.OptionValue()
		value = template.HTMLEscapeString(value)
		label = template.HTMLEscapeString(label)
		selected := ""
		if value == id {
			selected = " selected"
		}
		slct += `<option value="` + value + `"` + selected + `>` + label + "</option>"
	}
	slct += "</select>"

	return template.HTML(`
			<div class="mdl-select mdl-js-select mdl-select--floating-label ` + is + ` mdl-cell mdl-cell--4-col">
				` + slct + `
				` + fe + `
			</div>
	`)
}

// MultiSelectField is used to create a multi-select field
func MultiSelectField(p *Page, field interface{}, options []htmlselect.Option) template.HTML {
	fo := convert(field)

	vals := p.GroupValues[fo.Key]

	fe := string(FieldError(p, fo.Key))
	is := IsValid(p, fo.Key)

	slct := `<select multiple id="` + fo.CssID + `" name="` + fo.Name + `" class="mdl-select__input ` + fo.CssClass + `">`
	if len(vals) == 0 {
		slct += `<option value="0" selected>` + fo.Label + "</option>"
	} else {
		slct += `<option value="0">` + fo.Label + "</option>"
	}

	for _, op := range options {
		value, label := op.OptionValue()
		value = template.HTMLEscapeString(value)
		label = template.HTMLEscapeString(label)
		selected := ""
		for _, v := range vals {
			if value == v {
				selected = " selected"
				break
			}
		}
		slct += `<option value="` + value + `"` + selected + `>` + label + "</option>"
	}
	slct += "</select>"

	return template.HTML(`
		<div class="mdl-select mdl-js-select mdl-select--floating-label ` + is + ` mdl-cell mdl-cell--12-col">
			` + slct + `
			` + fe + `
		</div>
	`)
}

func SelectFieldWithDefault(p *Page, field interface{}, defaultValue, defaultLabel interface{}, options []htmlselect.Option) template.HTML {
	fo := convert(field)

	id := template.HTMLEscapeString(p.FormValues[fo.Key])

	fe := string(FieldError(p, fo.Key))
	is := IsValid(p, fo.Key)

	slct := `<select id="` + fo.CssID + `" name="` + fo.Name + `" class="mdl-select__input ` + fo.CssClass + `">`
	if len(id) == 0 {
		slct += `<option value="` + fmt.Sprint(defaultValue) + `" selected>` + fmt.Sprint(defaultLabel) + `</option>`
	}

	for _, op := range options {
		value, label := op.OptionValue()
		value = template.HTMLEscapeString(value)
		label = template.HTMLEscapeString(label)
		selected := ""
		if value == id {
			selected = " selected"
		}
		slct += `<option value="` + value + `"` + selected + `>` + label + "</option>"
	}
	slct += "</select>"

	return template.HTML(`
			<div class="mdl-select mdl-js-select mdl-select--floating-label ` + is + ` mdl-cell mdl-cell--12-col">
				` + slct + `
				` + fe + `
			</div>
	`)
}

// IsValid returns the is invalid class name if the field has an error
func IsValid(p *Page, key string) string {
	if _, ok := p.FormErrors[key]; ok {
		return "is-invalid "
	}
	return ""
}

// ArrayFieldValue
func ArrayFieldValue(p *Page, fieldName string, id string) string {
	v, ok := p.FormValues[fieldName+":"+id]
	if !ok {
		return ""
	}
	return template.HTMLEscapeString(v)
}

// FieldValue returns the field value if there is one
func FieldValue(p *Page, key string) string {
	v, ok := p.FormValues[key]
	if !ok {
		return ""
	}

	return v
}

// FieldError returns the error message HTML if the field has an error message
func FieldError(p *Page, key string) template.HTML {
	// TODO(james): create a helper type to fetch from a map[string]interface{}
	e, ok := p.FormErrors[key]
	if !ok {
		return ""
	}

	return template.HTML(`
		<span class="mdl-textfield__error">` + template.HTMLEscapeString(e) + `</span>
	`)
}

func ArrayFieldError(p *Page, field string, id interface{}) template.HTML {
	return FieldError(p, fmt.Sprintf("%s:%v", field, id))
}

func title(s string) string {
	return template.HTMLEscapeString(strings.Title(strings.Replace(s, "-", " ", -1)))
}

func treatWidth(width string) string {
	width = template.HTMLEscapeString(width)
	if width == "mdl-cell" {
		return " " + width + " "
	}
	if len(width) > 0 {
		return " mdl-cell mdl-cell--" + width + "-col "
	}

	return ""
}

//escapeField returns an HTML escaped field value, if there is one
func escapeField(p *Page, key string) string {
	v, ok := p.FormValues[key]
	if !ok {
		return ""
	}

	return template.HTMLEscapeString(v)
}

func getLookup(field string) (lookup string, fieldName string) {
	fieldName = template.HTMLEscapeString(field)
	parts := strings.Split(fieldName, ":")
	if len(parts) > 0 {
		lookup = fieldName
		fieldName = parts[0]
	} else {
		lookup = fieldName
	}

	return lookup, fieldName
}

func convert(field interface{}) *FieldOptions {
	switch t := field.(type) {
	case string:
		l, f := getLookup(t)
		return &FieldOptions{
			Key:      l,
			Name:     f,
			Label:    title(f),
			CssClass: f,
			CssID:    strings.Replace(l, ":", "-", -1),
		}
	case *FieldOptions:
		return t
	case ToFieldOptions:
		return t.ToFieldOptions()
	default:
		return &FieldOptions{}
	}
}
