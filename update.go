// Copyright 2015 UBIKWA SL.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	"ubikwa/backend/microservice/ubk-silo/main/models"
	"ubikwa/backend/util/conversions"
	"ubikwa/backend/util/gongiutil"
	"ubikwa/backend/util/gongsi"
)

type (
	// UpdateController represents the controller for operating on the User resource
	UpdateController struct {
		status int
	}
)

func NewUpdateController() *UpdateController {
	return &UpdateController{0}
}

// Rception POST with NGSI update from Orion
// It computes distance from image and updates:
// - apiweb via POST
// - Orion via POST UpdateContext with mean distance field and distance image
func (uc UpdateController) OnUpdateContext(c *gin.Context) {
	var status int
	var details string
	// This method receive JSON POST into NGSI format that contains update data to be processed..
	cr := gongsi.NotifyContextRequest{}

	//////////////////////////////////////////////////////////////////////////
	// Bind JSON into interval structure NotifyContextRequest
	if err := c.Bind(&cr); err != nil {
		log.Println("Error in request:", err)
		status = http.StatusBadRequest
		details = "JSON parsing error"
	} else {
		// Depending on ID send channel updates....
		Id, _ := gongsiutil.GetId(&cr)
		_ = Id
		battery, _ := gongsiutil.FindAttribute(&cr, string("battery_c"))
		temperature, _ := gongsiutil.FindAttribute(&cr, string("temperature_c"))
		capacity, _ := gongsiutil.FindAttribute(&cr, string("available_percentage"))

		battery_f, _ := conversions.Float64ForValue(battery)
		temperature_f, _ := conversions.Float64ForValue(temperature)
		capacity_f, _ := conversions.Float64ForValue(capacity)

		json_str, _ := json.Marshal(cr)
		fmt.Println(string(json_str))
		t := time.Now()
		data := models.WebsocketData{
			Timestamp: t.Unix() * 1000,
			Analog: models.AnalogData{
				Capacity:    capacity_f,
				Battery:     battery_f,
				Temperature: temperature_f,
			},
		}
		// {"timestamp":1448021916,"analog":{"capacity":25.520000457763672,"battery":95.12000274658203,"temp":31.600000381469727}}
		if b, err := json.Marshal(data); err == nil {
			h.unicast <- DirectMessage{Id, string(b)}
			//h.broadcast <- string(b)
		}

		// TODO: Depending on device SerialNumber send to client
	}

	//////////////////////////////////////////////////////////////////////////
	// Work with data sent
	response := gongsi.NotifyContextResponse{}
	sc := gongsi.StatusCode{}
	sc.Code = status
	sc.Details = details
	response.ResponseCode = sc

	// Marshal provided interface into XML structure
	uj, _ := json.Marshal(response)

	// Write content-type, statuscode, payload
	c.JSON(http.StatusCreated, uj)
}
