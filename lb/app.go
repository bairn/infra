package lb

import (
	"fmt"
	"github.com/tietang/go-eureka-client/eureka"
	"strings"
)

type Apps struct {
	Client *eureka.Client
}

func (a *Apps) Get(appName string) *App {
	var app *eureka.Application
	for _, a := range a.Client.Applications.Applications {
		if a.Name == strings.ToUpper(appName) {
			app = &a
			break
		}
	}

	if app == nil {
		return nil
	}

	na := &App {
		Name : app.Name,
		Instances: make([]*ServerInstance, 0),
		lb : &RoundRobinBalancer{},
	}

	for _, ins := range app.Instances {
		var port int
		if ins.SecurePort.Enabled {
			port = ins.SecurePort.Port
		} else {
			port = ins.Port.Port
		}

		si := &ServerInstance{
			InstanceId: ins.InstanceId,
			AppName:    appName,
			Address:    fmt.Sprintf("%s:%d", ins.IpAddr, port),
			Status:     Status(ins.Status),
		}
		na.Instances = append(na.Instances, si)
	}

	return na
}

type App struct {
	Name string
	Instances []*ServerInstance
	lb Balancer
}

func (a *App) Get(key string) *ServerInstance {
	ins := a.lb.Next(key, a.Instances)
	return ins
}


type Status string

const (
	StatusEnabled Status = "enabled"
	StatusDisabled Status = "disabled"
)

type ServerInstance struct {
	InstanceId string
	AppName string
	Address string
	Status Status
}