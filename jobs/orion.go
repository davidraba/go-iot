package jobs

import (
	"ubikwa/backend/microservice/ubk-silo/main/controllers"
)

type OrionJob struct {
	Stream   []byte
	Orion    controllers.OrionController
	OrionUrl string
}

func NewOrionJob(stream []byte, orionUrl string) *OrionJob {
	//--- API Controllers Orion ------------------------------------------
	orion_ctl := controllers.OrionController{}
	orion_ctl.Init(nil)
	orion_ctl.Set(orionUrl)
	return &OrionJob{Stream: stream, Orion: orion_ctl, OrionUrl: orionUrl}
}

func (p *OrionJob) Do() {
	p.Orion.UpdateContextStream(p.Stream)
}
