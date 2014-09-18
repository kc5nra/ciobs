package ciobs

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const taskPath = "tasks"

var baseDir = flag.String("base-dir", "", "base directory")

func init() {
	flag.StringVar(baseDir, "d", "", "base directory")
}

type TaskMan struct {
	tasks []*Task
}

func loadTask(taskDir os.FileInfo) (t *Task, err error) {
	return nil, nil
}

func loadTasks() ([]*Task, error) {
	p := path.Join(*baseDir, taskPath)
	files, e := ioutil.ReadDir(p)
	if e != nil {
		// no tasks, directory will be create when task is created
		return nil, nil
	}

	var tasks []*Task

	for _, d := range files {
		if !d.IsDir() {
			log.Printf("found noise in task dir: '%s'", d)
		} else {
			log.Printf("found task: '%s'", d)
			t, e := loadTask(d)
			if e != nil {
				return nil, e
			}
			tasks = append(tasks, t)
		}
	}

	return nil, nil

}

func (t *TaskMan) AddTask(name string) error {
	for _, t := range t.tasks {
		if (t.name == name) {
			return fmt.Errorf("taskman: task with name '%s' already exists")
		}
	}
	t.tasks = append(t.tasks, NewTask(name))
	return nil
}

func NewTaskMan() (*TaskMan, error) {

	tasks, e := loadTasks()
	if e != nil {
		return nil, e
	}
	tm := &TaskMan{tasks}
	return tm, nil
}
