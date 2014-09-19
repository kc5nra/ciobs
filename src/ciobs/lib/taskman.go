package ciobs

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path"
)

const taskPath = "tasks"

var baseDir = flag.String("base-dir", "", "base directory")

func init() {
	flag.StringVar(baseDir, "d", "", "base directory")
}

type TaskMan struct {
	tasks []*Task
	path  string
}

func loadTasks(tasksDir string) ([]*Task, error) {
	files, e := ioutil.ReadDir(tasksDir)
	if e != nil {
		// no tasks, directory will be create when task is created
		return nil, nil
	}

	var tasks []*Task

	for _, d := range files {
		if !d.IsDir() {
			log.Printf("found noise in task dir: '%s'", d)
		} else {
			log.Printf("found task: '%s'", d.Name())
			t, e := NewTask(tasksDir, d.Name())
			if e != nil {
				return nil, e
			}
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}

func appendUnique(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

func appendGraph(graph map[string][]string, k string, v string) {
	var deps []string
	if d, ok := graph[k]; ok {
		deps = d
	}
	graph[k] = appendUnique(deps, v)
}

func normalizeDeps(tasks []*Task) map[string][]string {
	graph := make(map[string][]string)

	// Naive, but most likely okay
	for _, task := range tasks {
		for k, _ := range task.upTriggers {
			appendGraph(graph, k, task.Name())
		}
		for k, _ := range task.downTriggers {
			appendGraph(graph, task.Name(), k)
		}
		if _, ok := graph[task.Name()]; !ok {
			graph[task.Name()] = make([]string, 0)
		}
	}

	return graph
}

func topSort(g map[string][]string) (order, cyclic []string) {
	L := make([]string, len(g))
	i := len(L)
	temp := map[string]bool{}
	perm := map[string]bool{}
	var cycleFound bool
	var cycleStart string
	var visit func(string)
	visit = func(n string) {
		switch {
		case temp[n]:
			cycleFound = true
			cycleStart = n
			return
		case perm[n]:
			return
		}
		temp[n] = true
		for _, m := range g[n] {
			visit(m)
			if cycleFound {
				if cycleStart > "" {
					cyclic = append(cyclic, n)
					if n == cycleStart {
						cycleStart = ""
					}
				}
				return
			}
		}
		delete(temp, n)
		perm[n] = true
		i--
		L[i] = n
	}
	for n := range g {
		if perm[n] {
			continue
		}
		visit(n)
		if cycleFound {
			return nil, cyclic
		}
	}
	return L, nil
}

func (tm *TaskMan) AddTask(name string) (*Task, error) {
	for _, t := range tm.tasks {
		if t.name == name {
			return nil, fmt.Errorf("taskman: task with name '%s' already exists")
		}
	}

	t, e := NewTask(*baseDir, name)
	if e != nil {
		return nil, e
	}

	tm.tasks = append(tm.tasks, t)

	return t, nil
}

func NewTaskMan() (*TaskMan, error) {
	return NewTaskManUsing(*baseDir)
}

func NewTaskManUsing(baseDir string) (*TaskMan, error) {
	tasksDir := path.Join(baseDir, taskPath)
	tasks, e := loadTasks(tasksDir)
	if e != nil {
		return nil, e
	}
	tm := &TaskMan{
		tasks: tasks,
		path:  baseDir,
	}
	return tm, nil
}
