package models

import (
	"math"
)

type (

	// User represents the structure of our resource
	SiloMeasureSimpleRequest struct {
		Distancia            float64 `json:"distancia"`
		SiloHeightCm         int     `json:"silo_height_cm"`
		SiloDiameterCm       int     `json:"silo_diameter_cm"`
		SiloHeightConeCm     int     `json:"silo_height_cone_cm"`
		SiloOffsetDistanceCm int     `json:"silo_offset_distance_cm"` // Distance from sensor to cylinder
		ContentDensityKgm3   float64 `json:"content_density_kgm3"`
	}

	SiloMeasureResponse struct {
		SiloCapacityM3            float64 `json:"silo_capacity_m3"`
		SiloCapacityKg            float64 `json:"silo_capacity_kg"`
		ContentDistanceFromSensor float64 `json:"content_distance_from_sensor"`
		ContentLevelCm            float64 `json:"content_level_cm"`
		ContentVolumeM3           float64 `json:"content_volume_m3"`
		ContentPerc               float64 `json:"content_perc"`
		ContentWeightKg           float64 `json:"content_weight_kg"`
	}
)

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func (req *SiloMeasureSimpleRequest) EvalDistance() (response SiloMeasureResponse) {

	/*
		(altura_con+altura_cil+distancia_cil_sen)-distancia_pienso
	*/
	response.ContentLevelCm = (float64(req.SiloHeightConeCm) + float64(req.SiloHeightCm) + float64(req.SiloOffsetDistanceCm)) - req.Distancia

	/*
		REDONDEAR(
			SI(nivel_pienso>altura_con;
				((PI()/3*(diametro_cil/2)^2*altura_con)+
				 (PI()*(diametro_cil/2)^2*(nivel_pienso-altura_con)))/1000000;
			ELSE
				PI()/3*((diametro_cil*nivel_pienso/altura_con)/2)^2*nivel_pienso/1000000);
		2)
	*/
	if response.ContentLevelCm != -1 {
		if float64(response.ContentLevelCm) > float64(req.SiloHeightConeCm) {
			response.ContentVolumeM3 = Round(
				(math.Pi/3*(float64(req.SiloDiameterCm)/2)*(float64(req.SiloDiameterCm)/2)*float64(req.SiloHeightConeCm)+
					(math.Pi*(float64(req.SiloDiameterCm)/2)*(float64(req.SiloDiameterCm)/2)*(response.ContentLevelCm-float64(req.SiloHeightConeCm))))/1E6,
				.5, 2)
		} else {
			if req.SiloHeightConeCm != 0 {
				response.ContentVolumeM3 = Round(
					math.Pi/3*(math.Pow((float64(req.SiloDiameterCm)*response.ContentLevelCm/float64(req.SiloHeightConeCm))/2, 2)*response.ContentLevelCm/1E6),
					.5, 2)
			} else {
				response.ContentVolumeM3 = -1
			}
		}
	} else {
		response.ContentVolumeM3 = -1
	}

	/*
		volumen_pienso*densidad
	*/
	if response.ContentVolumeM3 != -1 && req.ContentDensityKgm3 != -1 {
		response.ContentWeightKg = Round(response.ContentVolumeM3*req.ContentDensityKgm3, .5, 1)
	} else {
		response.ContentWeightKg = -1
	}

	//REDONDEAR(
	//    (
	//     (PI()/3*(diametro_cil/2)^2*altura_con) +
	//     (PI()*(diametro_cil/2)^2*altura_cil)
	//    ) /1000000;2)
	tmp := (float64(req.SiloDiameterCm) / 2) * (float64(req.SiloDiameterCm) / 2)
	response.SiloCapacityM3 = Round(((math.Pi/3*tmp*float64(req.SiloHeightConeCm))+
		(math.Pi*(float64(req.SiloDiameterCm)/2*float64(req.SiloDiameterCm)/2)*float64(req.SiloHeightCm)))/1E6, .5, 2)

	/*
		REDONDEAR(volumen_pienso/capacidad_vol*100;0)
	*/
	if response.ContentVolumeM3 != -1 && response.SiloCapacityM3 != 0 {
		response.ContentPerc = Round(response.ContentVolumeM3/response.SiloCapacityM3*100, .5, 0)
	} else {
		response.ContentPerc = -1
	}

	/*
		REDONDEAR(densidad*capacidad_vol;0)
	*/
	if response.SiloCapacityM3 != -1 {
		response.SiloCapacityKg = Round(req.ContentDensityKgm3*response.SiloCapacityM3, .5, 0)
	} else {
		response.SiloCapacityKg = -1
	}
	return
}
