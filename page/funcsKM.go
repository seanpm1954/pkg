package page

import (
	"html/template"
	"strings"

	"github.com/biz/templates"
)

func funcsKM() {
	templates.AddFunc("RadioFieldKM", RadioFieldKM)
	templates.AddFunc("RadioInputClassKM", RadioInputClassKM)
	templates.AddFunc("BasicRadioInputKM", BasicRadioInputKM)
}

func RadioFieldKM(p *Page, field interface{}) template.HTML {

	fo := convert(field)
	value := escapeField(p, fo.Key)

	checked := ""
	if fo.InputValue == value {
		checked = "checked"
	}

	return template.HTML(`
		<label class="mdl-radio mdl-js-radio mdl-js-ripple-effect" for="` + fo.CssID + `">
			<input type="radio" id="` + fo.CssID + `"  name="` + fo.Name + `" value="` + fo.InputValue + `"` + checked + ` />
			<span class="mdl-radio__label">` + fo.Label + `</span>
		</label>
		`)
}

func RadioInputClassKM(label, name, value, class string) *FieldOptions {
	f := BasicRadioInput(label, name, value)
	f.CssClass = class
	return f
}

func BasicRadioInputKM(label, name, value string) *FieldOptions {
	return &FieldOptions{
		Label:      label,
		Name:       name,
		InputValue: value,
		CssClass:   name,
		CssID:      strings.Replace(strings.ToLower(name+"-"+value), " ", "-", -1),
		Key:        name,
	}
}
