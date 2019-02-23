package tripoli

import (
	"reflect"
	"sync"
	"sync/atomic"
)

type Jobs struct {
	job interface{}
	arg []interface{}
}

type Results struct {
	arg    []interface{}
	result []interface{}
	error  error
}

type Tripoli struct {
	data    []interface{}
	worker  interface{}
	jobs    chan Jobs
	results chan Results
	done    chan bool
}

func (t *Tripoli) PrepareJobs() {
	for key, item := range t.data {
		t.jobs <- Jobs{job: key, arg: []interface{}{item}}
	}
	close(t.jobs)
}

func (t *Tripoli) ResultsHarvester(storage *atomic.Value, wg *sync.WaitGroup) {
	var stor []interface{}
	for r := range t.results {
		stor = append(stor, reflect.ValueOf(r.result).Interface())
		storage.Store(stor)
	}
	wg.Done()
}

func (t *Tripoli) StartPool(amount int) {
	var wg sync.WaitGroup

	for i := 0; i < amount; i++ {
		wg.Add(1)
		go t.Worker(&wg)
	}

	wg.Wait()
	close(t.results)
}

func (t *Tripoli) Worker(wg *sync.WaitGroup) {
	for job := range t.jobs {
		v := reflect.ValueOf(t.worker)
		p := make([]reflect.Value, v.Type().NumIn())
		p[0] = reflect.ValueOf(job.arg[0])

		var result []interface{}
		for _, v := range v.Call(p) {
			result = append(result, v.Interface())
		}

		t.results <- Results{
			result: result,
			error:  nil,
		}
	}

	wg.Done()
}

func (t *Tripoli) Exec(amount int) []interface{} {
	var storage atomic.Value
	t.jobs = make(chan Jobs, len(t.data))
	t.results = make(chan Results, len(t.data))

	var wg sync.WaitGroup
	wg.Add(1)

	go t.PrepareJobs()
	go t.ResultsHarvester(&storage, &wg)

	t.StartPool(amount)
	wg.Wait()

	return storage.Load().([]interface{})
}

func Run(fn interface{}, amount int, data []interface{}) []interface{} {
	var tp = &Tripoli{}
	tp.data = data
	tp.worker = fn
	return tp.Exec(amount)
}
