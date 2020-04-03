package client

// Header const
const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
	ApplicationXML  = "application/xml"
	TextXML         = "text/xml"

	HeaderXRequestID    = "X-Request-ID"
	HeaderAuthorization = "Authorization"
)

// Header var
var (
	HeaderApplicationJSON = Header{ContentType: ApplicationJSON}
	HeaderApplicationXML  = Header{ContentType: ApplicationXML}
	HeaderTextXML         = Header{ContentType: TextXML}
)
