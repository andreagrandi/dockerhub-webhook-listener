package listener

import "log"

type Handler interface {
	Call(HubMessage)
}

type Logger struct{}

func (l *Logger) Call(msg HubMessage) {
	log.Print(msg)
}

type Registry struct {
	entries []func(HubMessage)
}

func (r *Registry) Add(h func(msg HubMessage)) {
	r.entries = append(r.entries, h)
	return
}

func (r *Registry) Call(msg HubMessage) {
	for _, h := range r.entries {
		go h(msg)
	}
}

func reloadHandler(msg HubMessage) {
	log.Println("received message to reload ...")
	out, err := exec.Command("reload-docker.sh").Output()

	if err != nil {
		log.Println("ERROR EXECUTING COMMAND IN RELOAD HANDLER!!")
		log.Println(err)
		return
	}

	log.Println("output of reload-docker.sh is", string(out))
}

func MsgHandlers() Registry {
	var handlers Registry

	handlers.Add((&Logger{}).Call)
	handlers.Add(&reloadHandler)

	return handlers
}
