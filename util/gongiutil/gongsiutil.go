package gongsiutil

import (
	"fmt"
	"github.com/davidraba/go-iot/util/gongsi"
	"time"
)

// Defined NGSI variable types that can be received
const (
	TYPE_IMAGE = iota
	TYPE_BATTERY
	TYPE_TEMPERATURE
	TYPE_DISTANCE
)

type Attribute struct {
	Attr string
	Type string
}

type ContextError struct {
	When time.Time
	What string
}

func (e ContextError) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

/*
Receive an UpdateContext that contains image and sizes associated with distance estimation with a JSON described as:
{
    "contextElement": {
        "id": "SILO001",
        "type": "SILO",
        "isPattern": false,
        "attributes": [
            {
                "name": "distance",
                "type": "float32",
                "value": "11.01"
            },
            {
                "name": "image",
                "type": "string",
                "value": "11086 11060 11154 11134 11276 11026 10876 10998 10828 11158 11428 11267 11002 11277 10850 10807 11014 11124 3747 11061 10994 11251 11121 11335 10945 11353 11095 11121 11263 11261 11159 11070 11364 11162 11133 11194 11063 10830 11265 11058 10892 10788 11342 10938 11136 11964 11050 11567 11249 10928 11248 11167 10919 11179 10844 10881 11318 11227 11286 11039 11061 11179 10915 11085"
            },
            {
                "name": "numcols",
                "type": "int32",
                "value": "8.00"
            },
            {
                "name": "numrows",
                "type": "int32",
                "value": "8.00"
            }
        ]
    }
}
*/
func FindAttribute(cr *gongsi.NotifyContextRequest, attribute string) (string, error) {

	var found bool
	var ret string
	var e error
	/*	cr.ContextResponses[0].ContextElement.Attributes[0].Name
		cr.ContextResponses[0].ContextElement.Attributes[0].Value
		cr.ContextResponses[0].ContextElement.Attributes[0].Type
	*/
	found = false
	idx := 0
	for _, attrib := range cr.ContextResponses[0].ContextElement.Attributes {
		if attrib.Name == attribute {
			ret = attrib.Value
			found = true
		}
		idx++
	}
	if !found {
		e = ContextError{time.Now(), "Parameter not available in provided context"}
	}
	return ret, e
}

func FindTypedAttribute(cr *gongsi.NotifyContextRequest, attribute string) (ret Attribute, e error) {

	var found bool
	/*	cr.ContextResponses[0].ContextElement.Attributes[0].Name
		cr.ContextResponses[0].ContextElement.Attributes[0].Value
		cr.ContextResponses[0].ContextElement.Attributes[0].Type
	*/
	found = false
	idx := 0
	for _, attrib := range cr.ContextResponses[0].ContextElement.Attributes {
		if attrib.Name == attribute {
			ret.Attr = attrib.Value
			ret.Type = attrib.Type
			found = true
		}
		idx++
	}
	if !found {
		e = ContextError{time.Now(), "Parameter not available in provided context"}
	}
	return ret, e
}

/*
	Get Id from NGSI ContextElement

	type ContextElement struct {
		Id                  string             `json:"id"`
		Type                string             `json:"type,omitempty"`
		IsPattern           bool               `json:"isPattern"`
		AttributeDomainName string             `json:"attributeDomainName,omitempty"`
		Attributes          []ContextAttribute `json:"attributes,omitempty"`
		Metadatas           []ContextMetadata  `json:"metadatas,omitempty"`
	}
*/
func GetId(cr *gongsi.NotifyContextRequest) (out string, err error) {
	out = cr.ContextResponses[0].ContextElement.Id
	return out, nil
}

func GetType(cr *gongsi.NotifyContextRequest) (out string, err error) {
	out = cr.ContextResponses[0].ContextElement.Type
	return out, nil
}

func isIncluded(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetContextType(cr *gongsi.NotifyContextRequest) (contextType int) {
	if _, e := FindTypedAttribute(cr, "battery"); e == nil {
		return TYPE_BATTERY
	} else if _, e := FindTypedAttribute(cr, "temperature"); e == nil {
		return TYPE_TEMPERATURE
	} else if _, e := FindTypedAttribute(cr, "pixel"); e == nil {
		return TYPE_IMAGE
	}
	return -1
}
