package client

// Header const
const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json; charset=utf-8"
	ApplicationXML  = "application/xml; charset=utf-8"
	TextXML         = "text/xml; charset=utf-8"

	HeaderXRequestID    = "X-Request-ID"
	HeaderAuthorization = "Authorization"
)

// Header var
var (
	HeaderApplicationJSON = Header{ContentType: ApplicationJSON}
	HeaderApplicationXML  = Header{ContentType: ApplicationXML}
	HeaderTextXML         = Header{ContentType: TextXML}
)
