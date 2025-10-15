package rates

type XMLEcbRss struct {
	Channel XMLEcbChannel `xml:"channel"`
}

type XMLEcbChannel struct {
	Items []XMLEcbItem `xml:"item"`
}

type XMLEcbItem struct {
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}
