package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error {
	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if nil != err {
		return err
	}
	http.Handle(serviceUpdateURL.Path, new(serviceUpdateHandler))

	buff := new(bytes.Buffer)
	enc := json.NewEncoder(buff)
	err = enc.Encode(r)
	if nil != err {
		return err
	}

	res, err := http.Post(ServiesURL, "application/json", buff)
	if nil != err {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service. registry service responsed with code %v", res.StatusCode)
	}

	return nil
}

func Shutdown(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServiesURL, bytes.NewBuffer([]byte(url)))
	if nil != err {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")

	res, err := http.DefaultClient.Do(req)
	if nil != err {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to remove service. registry service responsed with code %v", res.StatusCode)
	}

	return nil
}

type serviceUpdateHandler struct{}

func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	prov.Update(p)
}

// 当当前服务被注册时，会请求到一些被它依赖的服务
// 这里记录当前服务所请求得到的依赖
type providers struct {
	services map[ServiceName][]string // 依赖的服务名 -> 依赖的 URL
	mutex    *sync.RWMutex            // services 的互斥量
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name],
			patchEntry.URL)
	}

	for _, patchEntry := range pat.Removed {
		if providerURLs, ok := p.services[patchEntry.Name]; ok {
			for i := range providerURLs {
				if providerURLs[i] == patchEntry.URL {
					p.services[patchEntry.Name] = append(providerURLs[:i],
						providerURLs[i+1:]...)
				}
			}
		}
	}
}

// 通过服务名称获得服务的 URL
func (p providers) get(name ServiceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("no provider available for service %v", name)
	}

	idx := int(rand.Float32() * float32(len(providers)))
	return providers[idx], nil
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}
