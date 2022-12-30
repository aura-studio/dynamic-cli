package pusher

import "log"

type Pusher struct {
	*TaskList
}

func New(taskList *TaskList) *Pusher {
	return &Pusher{
		TaskList: taskList,
	}
}

func (p *Pusher) Push() {
	for remote, files := range p.tasks {
		switch remote {
		case "s3":
			p.pushS3(files)
		default:
			log.Println("Unknown remote:", remote)
		}
	}
}

func (*Pusher) pushS3(files []string) {

}
