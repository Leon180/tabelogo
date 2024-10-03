package enum

import "github.com/samber/lo"

const MiddleWareEventIDKey = "event_id_key"

type ContextType string

const (
	ContextTypeJSON               ContextType = "application/json"
	ContextTypeFormData           ContextType = "application/form-data"
	ContextTypeXWWWFormUrlencoded ContextType = "application/x-www-form-urlencoded"
	ContextTypeXML                ContextType = "application/xml"
)

func (c ContextType) ToString() string {
	return string(c)
}

type ContextTypeSlice []ContextType

func (c ContextTypeSlice) ToStringSlice() []string {
	return lo.Map(c, func(v ContextType, _ int) string {
		return v.ToString()
	})
}

type ContextTypeGroup string

const (
	ContextTypeGroupDefault ContextTypeGroup = "default"
)

func (c ContextTypeGroup) GetSlice() ContextTypeSlice {
	switch c {
	case ContextTypeGroupDefault:
		return ContextTypeSlice{ContextTypeJSON, ContextTypeFormData, ContextTypeXWWWFormUrlencoded}
	default:
		return ContextTypeSlice{}
	}
}

type RequestMethod string

const (
	RequestMethodGET     RequestMethod = "GET"
	RequestMethodPOST    RequestMethod = "POST"
	RequestMethodPUT     RequestMethod = "PUT"
	RequestMethodDELETE  RequestMethod = "DELETE"
	RequestMethodOPTIONS RequestMethod = "OPTIONS"
	RequestMethodPATCH   RequestMethod = "PATCH"
)

func (r RequestMethod) ToString() string {
	return string(r)
}

type RequestMethodSlice []RequestMethod

func (r RequestMethodSlice) ToStringSlice() []string {
	return lo.Map(r, func(v RequestMethod, _ int) string {
		return v.ToString()
	})
}

type RequestMethodGroup string

const (
	RequestMethodGroupDefault RequestMethodGroup = "default"
)

func (r RequestMethodGroup) GetSlice() RequestMethodSlice {
	switch r {
	case RequestMethodGroupDefault:
		return RequestMethodSlice{RequestMethodGET, RequestMethodPOST, RequestMethodPUT, RequestMethodDELETE, RequestMethodOPTIONS}
	default:
		return RequestMethodSlice{}
	}
}

type RequestHeader string

const (
	RequestHeaderXPINGOTHER     RequestHeader = "X-PINGOTHER"
	RequestHeaderAccept         RequestHeader = "Accept"
	RequestHeaderAuthorization  RequestHeader = "Authorization"
	RequestHeaderContentType    RequestHeader = "Content-Type"
	RequestHeaderXCSRFToken     RequestHeader = "X-CSRF-Token"
	RequestHeaderUpgrade        RequestHeader = "Upgrade"
	RequestHeaderOrigin         RequestHeader = "Origin"
	RequestHeaderConnection     RequestHeader = "Connection"
	RequestHeaderAcceptEncoding RequestHeader = "Accept-Encoding"
	RequestHeaderAcceptLanguage RequestHeader = "Accept-Language"
	RequestHeaderHost           RequestHeader = "Host"
	RequestHeaderAccessControl  RequestHeader = "Access-Control-Request-Method"
	RequestHeaderAccessHeaders  RequestHeader = "Access-Control-Request-Headers"
	RequestHeaderXGoogAPIKey    RequestHeader = "X-Goog-Api-Key"
	RequestHeaderXGoogFieldMask RequestHeader = "X-Goog-FieldMask"
)

func (r RequestHeader) ToString() string {
	return string(r)
}

type RequestHeaderSlice []RequestHeader

func (r RequestHeaderSlice) ToStringSlice() []string {
	return lo.Map(r, func(v RequestHeader, _ int) string {
		return v.ToString()
	})
}

type RequestHeaderGroup string

const (
	RequestHeaderGroupDefault RequestHeaderGroup = "default"
)

func (r RequestHeaderGroup) GetSlice() RequestHeaderSlice {
	switch r {
	case RequestHeaderGroupDefault:
		return RequestHeaderSlice{RequestHeaderXPINGOTHER, RequestHeaderAccept, RequestHeaderAuthorization, RequestHeaderContentType, RequestHeaderXCSRFToken, RequestHeaderUpgrade, RequestHeaderOrigin, RequestHeaderConnection, RequestHeaderAcceptEncoding, RequestHeaderAcceptLanguage, RequestHeaderHost, RequestHeaderAccessControl, RequestHeaderAccessHeaders}
	default:
		return RequestHeaderSlice{}
	}
}

type ExposeHeader string

const (
	ExposeHeaderContentLength ExposeHeader = "Content-Length"
	ExposeHeaderLink          ExposeHeader = "Link"
)

func (e ExposeHeader) ToString() string {
	return string(e)
}

type ExposeHeaderSlice []ExposeHeader

func (e ExposeHeaderSlice) ToStringSlice() []string {
	return lo.Map(e, func(v ExposeHeader, _ int) string {
		return v.ToString()
	})
}

type ExposeHeaderGroup string

const (
	ExposeHeaderGroupDefault ExposeHeaderGroup = "default"
)

func (e ExposeHeaderGroup) GetSlice() ExposeHeaderSlice {
	switch e {
	case ExposeHeaderGroupDefault:
		return ExposeHeaderSlice{ExposeHeaderContentLength, ExposeHeaderLink}
	default:
		return ExposeHeaderSlice{}
	}
}
