package lib

import (
    "log"
    marathon "github.com/gambol99/go-marathon"
)

type Marathon struct {
    Host string
    Apps []marathon.Application
    _client marathon.Marathon
}

type MarathonApps struct {
    Apps   []*MarathonApp
}

type MarathonApp struct {
    Id              string
    Tasks           []*MarathonTask
}

type MarathonTask struct {
    Id        string
    Container string
    SlaveId   string
    SlaveHost string
    Mesos     bool
    Marathon  bool
}

func (m *Marathon) Client() marathon.Marathon {
    if m._client != nil {
        return m._client
    }

    config := marathon.NewDefaultConfig()
    config.URL = m.Host
    if client, err := marathon.NewClient(config); err != nil {
        log.Fatalf("Failed to create a client for marathon, error: %s", err)
    } else {
        m._client = client
    }
    return m._client
}

func (m *Marathon) LoadApps() error {
    m.Client().ListApplications()
    if applications, err := m.Client().Applications(); err != nil {
        log.Println("Failed to list applications")
        return ErrMesosError
    } else {
        m.Apps = applications.Apps
    }
    return nil
}


func (ma *MarathonApps) GetAppById(appId string) *MarathonApp {
    var app *MarathonApp
    for _, a := range ma.Apps {
        if a.Id == appId {
            app = a
        }
    }
    return app
}


func (a *MarathonApp) GetTaskById(taskId string) *MarathonTask {
    var task *MarathonTask
    for _, t := range a.Tasks {
        if t.Id == taskId {
            task = t
        }
    }
    return task
}


func (a *MarathonApp) AddTask(taskId, slaveId, slaveHost string, mesos, marathon bool) *MarathonTask {
    var task *MarathonTask
    task = a.GetTaskById(taskId)
    if task == nil {
        task = &MarathonTask{ Id: taskId, SlaveHost: slaveHost, SlaveId: slaveId, Mesos: mesos, Marathon: marathon }
        a.Tasks = append(a.Tasks, task)
    }
    return task
}


func (ma *MarathonApps) AddApp(appId string) *MarathonApp {
    var app *MarathonApp
    for _, a := range ma.Apps {
        if a.Id == appId {
            app = a
        }
    }
    if app == nil {
        app = &MarathonApp{ Id: appId }
        ma.Apps = append(ma.Apps, app)
    }
    return app
}


func (ma *MarathonApps) AddTask(taskId, appId, slaveId, slaveHost string, mesos, marathon bool) *MarathonTask {
    var app *MarathonApp
    var task *MarathonTask

    if app = ma.GetAppById(appId);app != nil {
        task = app.GetTaskById(taskId)
    }

    if app == nil {
        app = ma.AddApp(appId)
    }
    if task == nil {
        task = app.AddTask(taskId, slaveId, slaveHost, mesos, marathon)
    }

    if mesos { task.Mesos = mesos }
    if marathon { task.Marathon = marathon }

    return task
}