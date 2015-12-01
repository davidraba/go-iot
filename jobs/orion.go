package jobs

//import (
//	"ubikwa/backend/microservice/ubk-silo/main/controllers"
//)

type OrionJob struct {
	Stream   []byte
	OrionUrl string
}

func NewOrionJob(stream []byte, orionUrl string) *OrionJob {
	return &OrionJob{}
}

func (p *OrionJob) Do() {
}
