package ciobs

import (
	"fmt"
	"os"
	"path"
)

type Task struct {
	name         string
	path         string
	upTriggers   map[string]bool
	downTriggers map[string]bool
}

func createDir(parentDir string, name string) error {
	p := path.Join(parentDir, name)
	_, e := os.Stat(p)
	if e != nil {
		if os.IsNotExist(e) {
			// replacing e is okay since it was
			// just a marker
			if e := os.MkdirAll(p, 0700); e != nil {
				return e
			}
			return nil
		} else {
			// unspecified error
			// unlikely that we can succeed
			return e
		}
	} else {
		// directory already exists
		return nil
	}
}

func NewTask(tasksDir string, name string) (*Task, error) {

	if e := createDir(tasksDir, name); e != nil {
		return nil, fmt.Errorf("task: unable to create task dir '%s', %s", name, e)
	}

	return &Task{
		name:         name,
		path:         tasksDir,
		upTriggers:   make(map[string]bool),
		downTriggers: make(map[string]bool),
	}, nil
}

func (t *Task) Name() string {
	return t.name
}

func (t *Task) AddUpTrigger(taskName string) {
	t.upTriggers[taskName] = true
}

func (t *Task) AddDownTrigger(taskName string) {
	t.downTriggers[taskName] = true
}

func (t *Task) RemoveUpTrigger(taskName string) {
	delete(t.upTriggers, taskName)
}

func (t *Task) RemoveDownTrigger(taskName string) {
	delete(t.downTriggers, taskName)
}
