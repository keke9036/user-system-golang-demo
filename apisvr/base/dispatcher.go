// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/7/12

package base

import (
	"entry-task/util"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	Get  = "GET"
	Post = "POST"
)

var dispatcher *Dispatcher

type HandlerFunc func(*GContext)

type Dispatcher struct {
	Handlers  map[string]map[string]HandlerFunc
	Pages     map[string]string
	Filters   []Filter
	PageFiles map[string]string
}

func NewDispatcher() *Dispatcher {
	dispatcher = &Dispatcher{}
	dispatcher.Handlers = make(map[string]map[string]HandlerFunc)
	dispatcher.Pages = make(map[string]string)
	dispatcher.PageFiles = make(map[string]string)

	return dispatcher
}

func (d *Dispatcher) RegisterHandler(url string, method string, handler HandlerFunc) {
	util.Logger.Infof("Register handler: url %s, method %s, %v", url, method, reflect.TypeOf(handler))
	handlerMap := make(map[string]HandlerFunc)
	handlerMap[method] = handler
	d.Handlers[url] = handlerMap
}

func (d *Dispatcher) RegisterPageDir(path, dir string) {
	util.Logger.Infof("Register page dir: path %s, dir %s", path, dir)
	d.Pages[path] = dir
}

func (d *Dispatcher) RegisterPageFile(path, file string) {
	util.Logger.Infof("Register page file: path %s, dir %s", path, file)
	d.PageFiles[path] = file
}

func (d *Dispatcher) RegisterFilter(filter Filter) {
	filterValue := reflect.ValueOf(filter)
	filterType := reflect.Indirect(filterValue).Type()
	filterName := filterType.Name()
	util.Logger.Infof("Register filter %s", filterName)
	d.Filters = append(d.Filters, filter)
}

// main entry
func (d *Dispatcher) ServeHTTP(objResp http.ResponseWriter, objReq *http.Request) {
	// static file mapping

	url := getUrl(objReq)
	gContext := GContext{
		Request: objReq,
		Writer:  objResp,
		URL:     url,
		Keys:    make(map[string]any),
	}

	pageFile, ok := d.PageFiles[url]
	if ok {
		err := d.RenderPage(objResp, pageFile)
		if err != nil {
			objResp.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	i := strings.Index(url[1:], "/")
	staticPath := url[0 : i+1]
	page, ok := d.Pages[staticPath]
	if ok {
		util.Logger.Infof("Static page matched, url %s, page %s", url, page)

		localPath := page + strings.TrimPrefix(url, staticPath)
		err := d.RenderPage(objResp, localPath)
		if err != nil {
			objResp.WriteHeader(http.StatusNotFound)
			return
		}
	}

	// filter mapping
	for _, filter := range d.Filters {
		if ok := filter.MatchUrl(&gContext); !ok {
			continue
		}

		if ok := filter.DoFilter(&gContext); !ok {
			return
		}
	}

	// handler mapping
	handler, ok := d.Handlers[url][objReq.Method]
	if !ok {
		// return 404
		util.Logger.Warnf("url not found %s", url)
		objResp.WriteHeader(http.StatusNotFound)
		return
	}

	handler(&gContext)
}

func getUrl(req *http.Request) string {
	return req.URL.Path
}

//返回Page
func (d *Dispatcher) RenderPage(objResp http.ResponseWriter, localPath string) error {
	fileBytes, err := ioutil.ReadFile(localPath)
	if err != nil {
		return err
	}
	ext := filepath.Ext(localPath)
	if ext != "" {
		objResp.Header().Set("Content-Type", mime.TypeByExtension(ext))
	}
	objResp.Write(fileBytes)
	return nil
}
