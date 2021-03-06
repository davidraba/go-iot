// Copyright 2015 UBIKWA SL.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	//"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	//"time"

	"github.com/davidraba/go-iot/models"
	"github.com/davidraba/go-iot/util/conversions"
	"github.com/davidraba/go-iot/util/gongiutil"
	"github.com/davidraba/go-iot/util/gongsi"
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
		//		battery, _ := gongsiutil.FindAttribute(&cr, string("battery_c"))
		//		temperature, _ := gongsiutil.FindAttribute(&cr, string("temperature_c"))
		//		capacity, _ := gongsiutil.FindAttribute(&cr, string("available_percentage"))
		distance, _ := gongsiutil.FindAttribute(&cr, string("mean_distance"))
		d, _ := conversions.Float64ForValue(distance)

		data := models.SiloData{Distance: d, Temperature: 35, Humidity: 40}
		h.unicast <- DirectMessage{Id, data}
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
