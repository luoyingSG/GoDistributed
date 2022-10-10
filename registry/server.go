package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// 注册服务
// 注册服务所占用的端口号/地址
const ServicePort = ":3000"
const ServiesURL = "http://localhost" + ServicePort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.RWMutex // 用于保护 registrations 的互斥量
}

// 添加注册
func (r *registry) add(reg Registration) error {
	r.mutex.Lock() // 给 registrations 上锁
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	// 给正要注册的服务它所依赖的服务
	err := r.sendRequiredServices(reg)
	if nil != err {
		return err
	}
	return nil
}

// 找到待注册服务的所有依赖服务，与当前注册表进行比对
// 找到那些已经注册过的依赖，将这些依赖的服务名、URL返回给当前待注册的服务
func (r registry) sendRequiredServices(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, reqService := range reg.RequiredServices {
		for _, serviceReg := range r.registrations {
			if reqService == serviceReg.ServiceName {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
			}
		}
	}
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if nil != err {
		return err
	}
	return nil
}

func (r registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if nil != err {
		return err
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if nil != err {
		return err
	}
	return nil
}

// 取消注册
func (r *registry) remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceURL == url {
			reg.mutex.Lock()
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			reg.mutex.Unlock()

			return nil
		}
	}

	return fmt.Errorf("service at url %s not found", url)
}

// 一个包级的注册管理变量
var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

// 一个 HTTP 服务
type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")

	switch r.Method {
	case http.MethodPost: // 添加注册
		dec := json.NewDecoder(r.Body) // 建立一个解码器
		var r Registration             // 解码的目标
		err := dec.Decode(&r)          // 进行解码
		if nil != err {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Adding service: %v with URL: %v\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if nil != err {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	case http.MethodDelete: // 取消注册
		payload, err := ioutil.ReadAll(r.Body)
		if nil != err {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		url := string(payload)
		log.Printf("Removing service at URL: %s\n", url)
		err = reg.remove(url)
		if nil != err {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
