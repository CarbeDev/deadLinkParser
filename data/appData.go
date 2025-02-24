package data

type AppData struct {
	InitialUrl string
	FoundLinks []FoundLink
}

type FoundLink struct {
	Link    string
	Visited bool
	Alive   bool
}

func InitialiseAppData(baseUrl string) AppData {
	//FIXME add link verification
	return AppData{
		InitialUrl: baseUrl,
		FoundLinks: []FoundLink{},
	}
}

func (data AppData) hasLink(link string) bool {
	for _, foundLink := range data.FoundLinks {
		if link == foundLink.Link {
			return true
		}
	}

	return false
}

func (data AppData) HasUncheckedLink() bool { //TODO check if useful
	for _, link := range data.FoundLinks {
		if !link.Visited {
			return true
		}
	}

	return false
}

func AddLinkFound(link string, appData *AppData) {
	if !appData.hasLink(link) {
		*appData = AppData{
			InitialUrl: appData.InitialUrl,
			FoundLinks: append(appData.FoundLinks, FoundLink{Link: link}),
		}
	}
}
