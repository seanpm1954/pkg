package page

import "net/http"
import "strings"

// DefaultMenu is used to store the menu state
var DefaultMenu = Menu{Types: map[string][]Link{}}

// Menu represents the menu
type Menu struct {
	Types map[string][]Link
}

// AddLink adds a links to a menu type
func (m *Menu) AddLink(typ string, links ...Link) {
	if _, ok := m.Types[typ]; !ok {
		m.Types[typ] = links
		return
	}

	m.Types[typ] = append(m.Types[typ], links...)
}

// Get gets menu items by section and sets the active menu item
func (m *Menu) Get(typ string, req *http.Request) []Link {
	links, ok := m.Types[typ]
	if !ok {
		return nil
	}

	return setActiveItems(req, links)
}

func setActiveItems(req *http.Request, links []Link) []Link {
	nl := make([]Link, len(links))

	// TODO: check sub links
	pt := req.URL.Path
	for i := 0; i < len(links); i++ {
		nl[i] = links[i]
		if nl[i].Href == pt {
			nl[i].IsActive = true
		} else if strings.Contains(pt, nl[i].Href) {
			nl[i].IsPartialActive = true
		}
	}

	return nl
}

// AddMenuLink adds links to the DefaultMenu
func AddMenuLink(typ string, links ...Link) {
	DefaultMenu.AddLink(typ, links...)
}

// GetMenu gets the links from the DefaultMenu
func GetMenu(typs string, req *http.Request) []Link {
	return DefaultMenu.Get(typs, req)
}
