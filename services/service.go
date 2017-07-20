package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type service struct {
	// Имя сервиса
	Name string `json:"-"`
	// Аргументы, которые будут переданы в сервис как параметры командной строки
	Args                  []string         `json:"args"`
	currentServiceChannel chan interface{} `json:"-"`
	currentExternalCmd    *exec.Cmd        `json:"-"`
	// Отметка о старте
	StartTime time.Time `json:"-"`
	IsRunning bool      `json:"-"`
	// *.go файл, который запускает сервис
	Target string `json:"target"`
	// Не даёт запустить новый сервис, если старый еще не закончил работу
	syncMutex sync.Mutex `json:"-"`
}

var allServices map[string]*service

const (
	OK_SIGNAL   = "ok"
	KILL_SIGNAL = "kill"
)

func loadServices() (error, bool) {
	raw, err := ioutil.ReadFile("services.json")
	if err != nil {
		return err, false
	}
	allServices = make(map[string]*service)
	err = json.Unmarshal(raw, &allServices)
	if err != nil {
		return err, false
	}
	for key, val := range allServices {
		if key == "all" || key == "env" {
			return fmt.Errorf("name `%v` is not allowed", key), true
		}
		val.Name = key
	}
	return nil, false
}

func init() {
	err, fatal := loadServices()
	if err != nil {
		fmt.Printf("[!] Can't load services because %v.\n", err)
		if fatal {
			os.Exit(1)
		}
	}
}

func ReloadServices(args ...string) {
	StopAll()
	err, fatal := loadServices()
	if err != nil {
		fmt.Printf("[!] Can't load services because %v.\n", err)
		if fatal {
			os.Exit(1)
		}
	}
	StartAll(args...)
}

func ReloadService(svcName string) {
	StopService(svcName)
	svc := NewService(svcName)
	if svc != nil {
		svc.Start()
	}
}

// Returns true if process should be stopped
func (svc *service) stringSwitch(text string) bool {
	switch text {
	case OK_SIGNAL:
		return true
	case KILL_SIGNAL:
		err := svc.currentExternalCmd.Process.Kill()
		if err != nil {
			fmt.Printf("[!] Service %v can't be killed because %v.\n", svc.Name, err)
		}
		return true
	}
	return false
}

func (svc *service) wait() {
	svc.IsRunning = true
	for {
		select {
		case signal := <-svc.currentServiceChannel:
			switch typedSignal := signal.(type) {
			case string:
				if svc.stringSwitch(typedSignal) {
					return
				}
			case error:
				fmt.Printf("[!] Error with service %v:\n%v\n", svc.Name, typedSignal)
				return
			}
		}
	}
}

func (svc *service) Start() {
	go func() {
		svc.StartTime = time.Now()
		svc.currentExternalCmd.Stdout = os.Stdout
		svc.currentExternalCmd.Stderr = os.Stderr
		err := svc.currentExternalCmd.Start()
		if err != nil {
			fmt.Printf("[!] Can't start service %v because %v.\n", svc.Name, err)
			return
		}
		svc.syncMutex.Lock()
		go func() {
			err := svc.currentExternalCmd.Wait()
			if err != nil {
				svc.currentServiceChannel <- err
			} else {
				svc.currentServiceChannel <- OK_SIGNAL
			}
		}()
		fmt.Printf("[ ] Start service %v.\n", svc.Name)
		svc.wait()
		svc.cleanService()
		svc.StartTime = time.Time{}
		fmt.Printf("[ ] Service %v was stopped\n", svc.Name)
	}()
}

func StopService(svcName string) {
	svc, exist := allServices[svcName]
	if !exist {
		fmt.Printf("[!] Can't find service %v.\n", svcName)
	} else {
		svc.currentServiceChannel <- KILL_SIGNAL
		svc.syncMutex.Lock()
		svc.syncMutex.Unlock()
	}
}

func (svc *service) cleanService() {
	close(svc.currentServiceChannel)
	svc.currentServiceChannel = nil
	svc.currentExternalCmd = nil
	svc.IsRunning = false
	svc.syncMutex.Unlock()
}

func NewService(svcName string, args ...string) *service {
	svc, exist := allServices[svcName]
	if !exist {
		fmt.Printf("[!] Can't find service %v.\n", svcName)
		return nil
	} else {
		if svc.currentExternalCmd != nil || svc.currentServiceChannel != nil {
			fmt.Printf("[?] Service %v already in use. Please stop or restart it.\n", svc.Name)
			return nil
		} else {
			svc.currentServiceChannel = make(chan interface{})
			svc.Args = append(svc.Args, args...)
			runArgs := []string{"run", svc.Target}
			runArgs = append(runArgs, svc.Args...)
			cmd := exec.Command("go", runArgs...)
			svc.currentExternalCmd = cmd
			return svc
		}
	}
}

func ListServices() string {
	now := time.Now()
	var svcStrs []string
	for key, val := range allServices {
		isRunningStr := "Down"
		if val.IsRunning {
			isRunningStr = fmt.Sprintf("Up for %v", now.Sub(val.StartTime))
		}
		svcStrs = append(svcStrs, fmt.Sprintf("%v %v %v", key, val.Args, isRunningStr))
	}
	return strings.Join(svcStrs, "\n")
}

func StartAll(args ...string) {
	for key, val := range allServices {
		if val.IsRunning {
			fmt.Printf("[?] %v already started.\n", key)
		} else {
			svc := NewService(key, args...)
			if svc != nil {
				svc.Start()
			}
		}
	}
}

func StopAll() {
	for key, val := range allServices {
		if val.IsRunning {
			StopService(key)
		}
	}
}