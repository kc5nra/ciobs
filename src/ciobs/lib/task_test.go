package ciobs

import (
	"io/ioutil"
	"os"
	"testing"
)

func newDir() (string, error) {
	f, e := ioutil.TempDir("", "testLoadTask")
	if e != nil {
		return "", e
	}
	return f, nil
}

func remDir(d string) {
	os.RemoveAll(d)
}

func TestNew(t *testing.T) {
	d, e := newDir()
	if e != nil {
		panic(e)
	}
	defer remDir(d)

	task, e := NewTask(d, "task")
	if e != nil {
		panic(e)
	}

	if rcnt := len(task.upTriggers); rcnt != 0 {
		t.Errorf("Task.NewTask() expected 0 up triggers but found %d", rcnt)
	}
	if tcnt := len(task.downTriggers); tcnt != 0 {
		t.Errorf("Task.NewTask() expected 0 down triggers but found %d", tcnt)
	}
}

func TestAddUpTrigger(t *testing.T) {
	d, e := newDir()
	if e != nil {
		panic(e)
	}

	defer remDir(d)
	task, e := NewTask(d, "task")
	if e != nil {
		panic(e)
	}

	task.AddUpTrigger("t1")
	if rcnt := len(task.upTriggers); rcnt != 1 {
		t.Errorf("len(task.upTriggers) expected 1, but found %d", rcnt)
	}
	if !task.upTriggers["t1"] {
		t.Errorf("task.upTriggers['t1'] not found")
	}
}

func TestAddDownTrigger(t *testing.T) {
	d, e := newDir()
	if e != nil {
		panic(e)
	}
	defer remDir(d)

	task, e := NewTask(d, "task")
	if e != nil {
		panic(e)
	}

	task.AddDownTrigger("t1")
	if rcnt := len(task.downTriggers); rcnt != 1 {
		t.Errorf("len(task.downTriggers) expected 1, but found %d", rcnt)
	}
	if !task.downTriggers["t1"] {
		t.Errorf("task.downTriggers['t1'] not found")
	}
}

func TestRemoveUpTrigger(t *testing.T) {
	d, e := newDir()
	if e != nil {
		panic(e)
	}
	defer remDir(d)

	task, e := NewTask(d, "task")
	if e != nil {
		panic(e)
	}

	task.AddUpTrigger("t1")
	task.RemoveUpTrigger("t1")
	if rcnt := len(task.upTriggers); rcnt != 0 {
		t.Errorf("len(task.upTriggers) expected 0, but found %d", rcnt)
	}
	if task.upTriggers["t1"] {
		t.Errorf("task.upTriggers['t1'] expected to be false, but true")
	}
}
