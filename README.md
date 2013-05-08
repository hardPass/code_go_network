- [ ] a task list item
- [ ] list syntax required
- [ ] normal **formatting**, @mentions, #1234 refs
- [ ] incomplete
- [x] completed
- ri
- fuc
- ss


```package mvc

const (

  // Attribute name for Session ID 
	SESSION_ID = `GSESSIONID`
	// CTX_SESSION = "CTX_SESSION"

	// dir of themes
	THEMES_DIR = `themes`

	// Attribute name for user theme 
	THEME_USER = `theme`

	// dir of tpl files for base theme
	THEME_BASE = `base`
)

//path string, tplName string, autoRender bool
func (m *mvc) RegisterController(path string, controller *Controller) *mvcInfo {
  info := &mvcInfo{}
	m.mvcInfoList = append(m.mvcInfoList, mi)
	info.path = path
	info.controller = controller
	info.concurrent = true
	info.autoRender = true
	info.params = make(map[string]interface{})
}


```




package mvc

const (

  // Attribute name for Session ID 
	SESSION_ID = `GSESSIONID`
	// CTX_SESSION = "CTX_SESSION"

	// dir of themes
	THEMES_DIR = `themes`

	// Attribute name for user theme 
	THEME_USER = `theme`

	// dir of tpl files for base theme
	THEME_BASE = `base`
)

//path string, tplName string, autoRender bool
func (m *mvc) RegisterController(path string, controller *Controller) *mvcInfo {
  info := &mvcInfo{}
	m.mvcInfoList = append(m.mvcInfoList, mi)
	info.path = path
	info.controller = controller
	info.concurrent = true
	info.autoRender = true
	info.params = make(map[string]interface{})
}

