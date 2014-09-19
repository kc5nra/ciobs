package ciobs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"testing"
)

func NewTaskManTemp() *TaskMan {
	f, e := ioutil.TempDir("", "taskManTemp")
	if e != nil {
		panic(e)
	}

	tm, e := NewTaskManUsing(f)
	if e != nil {
		panic(e)
	}

	return tm
}

func (tm *TaskMan) Delete() {
	os.RemoveAll(tm.path)
}

func printDepsAlpha(graph map[string][]string) {
	keys := make([]string, 0, len(graph))
	for k, _ := range graph {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s\n", k)
		for _, d := range graph[k] {
			fmt.Printf("- %s\n", d)
		}
	}
}

func TestLoadTasks(t *testing.T) {
	f, e := ioutil.TempDir("", "testLoadTask")
	if e != nil {
		panic(e)
	}

	p := path.Join(f, taskPath)
	os.Mkdir(p, 0700)

	defer os.RemoveAll(f)

	t1, e := NewTask(p, "t1")
	if e != nil {
		panic(e)
	}

	t2, e := NewTask(p, "t2")
	if e != nil {
		panic(e)
	}

	tasks, e := loadTasks(p)
	if e != nil {
		panic(e)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks but found %d", len(tasks))
		return
	}
	if t1.Name() != tasks[0].Name() {
		t.Errorf("tasks[0].Name() expected %s but found %s", t1.Name(), tasks[0].Name())
	}
	if t2.Name() != tasks[1].Name() {
		t.Errorf("tasks[0].Name() expected %s but found %s", t2.Name(), tasks[1].Name())
	}
}

func TestNewTaskMan(t *testing.T) {
	tm := NewTaskManTemp()
	defer tm.Delete()

	tm.AddTask("t1")
	tm.AddTask("t2")
	tm.AddTask("t3")

	if len(tm.tasks) != 3 {
		t.Errorf("expected 3 tasks but found %d", len(tm.tasks))
	}
}

func ExampleNormalizeDepsSimple() {
	tm := NewTaskManTemp()
	defer tm.Delete()

	tm.AddTask("t1")
	t2, _ := tm.AddTask("t2")
	tm.AddTask("t3")

	t2.AddUpTrigger("t1")

	deps := normalizeDeps(tm.tasks)
	printDepsAlpha(deps)

	// Output:
	// t1
	// - t2
	// t2
	// t3
}

func ExampleNormalizeDeps() {
	tm := NewTaskManTemp()
	defer tm.Delete()

	tm.AddTask("t1")
	t2, _ := tm.AddTask("t2")
	t3, _ := tm.AddTask("t3")

	//// Illegal, but for this function creates a circle
	t2.AddUpTrigger("t1")
	t3.AddUpTrigger("t2")
	t3.AddDownTrigger("t1")

	deps := normalizeDeps(tm.tasks)
	printDepsAlpha(deps)

	// Output:
	// t1
	// - t2
	// t2
	// - t3
	// t3
	// - t1
}

func ExampleTopSortDepsCircle() {
	tm := NewTaskManTemp()
	defer tm.Delete()

	tm.AddTask("t1")
	t2, _ := tm.AddTask("t2")
	t3, _ := tm.AddTask("t3")

	//// Illegal, but for this function creates a circle
	t2.AddUpTrigger("t1")
	t3.AddUpTrigger("t2")
	t3.AddDownTrigger("t1")

	deps := normalizeDeps(tm.tasks)
	order, cycle := topSort(deps)
	sort.Strings(cycle)
	fmt.Printf("%s\n", order)
	fmt.Printf("%s\n", cycle)

	// Output:
	// []
	// [t1 t2 t3]
}

func ExampleTopSortDeps() {
	tm := NewTaskManTemp()
	defer tm.Delete()

	t1, _ := tm.AddTask("t1")
	t2, _ := tm.AddTask("t2")
	t3, _ := tm.AddTask("t3")

	t1.AddDownTrigger("t3")
	t2.AddUpTrigger("t1")
	t3.AddUpTrigger("t2")

	deps := normalizeDeps(tm.tasks)
	order, cycle := topSort(deps)
	sort.Strings(cycle)
	fmt.Printf("%s\n", order)
	fmt.Printf("%s\n", cycle)

	// Output:
	// [t1 t2 t3]
	// []
}
