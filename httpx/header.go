package httpx

// Header const
const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json; charset=utf-8"
	ApplicationXML  = "application/xml; charset=utf-8"
	TextXML         = "text/xml; charset=utf-8"

	HeaderXRequestID    = "X-Request-ID"
	HeaderAuthorization = "Authorization"
)

// Header wrap map[string]string for header
type Header map[string]string

func HeaderApplicationJSON() Header {
	return Header{ContentType: ApplicationJSON}
}

func HeaderApplicationXML() Header {
	return Header{ContentType: ApplicationXML}
}

func HeaderTextXML() Header {
	return Header{ContentType: TextXML}
}
