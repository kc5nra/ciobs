package ciobs

import (
	"fmt"
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	task := NewTask("task")
	if rcnt := len(task.upTriggers); rcnt != 0 {
		t.Errorf("Task.NewTask() expected 0 up triggers but found %d", rcnt)
	}
	if tcnt := len(task.downTriggers); tcnt != 0 {
		t.Errorf("Task.NewTask() expected 0 down triggers but found %d", tcnt)
	}
}

func TestAddUpTrigger(t *testing.T) {
	task := NewTask("task")
	task.AddUpTrigger("t1")
	if rcnt := len(task.upTriggers); rcnt != 1 {
		t.Errorf("len(task.upTriggers) expected 1, but found %d", rcnt)
	}
	if !task.upTriggers["t1"] {
		t.Errorf("task.upTriggers['t1'] not found")
	}
}

func TestAddDownTrigger(t *testing.T) {
	task := NewTask("task")
	task.AddDownTrigger("t1")
	if rcnt := len(task.downTriggers); rcnt != 1 {
		t.Errorf("len(task.downTriggers) expected 1, but found %d", rcnt)
	}
	if !task.downTriggers["t1"] {
		t.Errorf("task.downTriggers['t1'] not found")
	}
}

func TestRemoveUpTrigger(t *testing.T) {
	task := NewTask("task")
	task.AddUpTrigger("t1")
	task.RemoveUpTrigger("t1")
	if rcnt := len(task.upTriggers); rcnt != 0 {
		t.Errorf("len(task.upTriggers) expected 0, but found %d", rcnt)
	}
	if task.upTriggers["t1"] {
		t.Errorf("task.upTriggers['t1'] expected to be false, but true")
	}
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

func ExampleNormalizeDepsSimple() {
	t1, t2, t3 := NewTask("t1"), NewTask("t2"), NewTask("t3")
	tasks := make([]*Task, 3)
	tasks[0] = t1
	tasks[1] = t2
	tasks[2] = t3

	t2.AddUpTrigger("t1")

	deps := normalizeDeps(tasks)
	printDepsAlpha(deps)

	// Output:
	// t1
	// - t2
	// t2
	// t3
}

func ExampleNormalizeDeps() {
	t1, t2, t3 := NewTask("t1"), NewTask("t2"), NewTask("t3")
	tasks := make([]*Task, 3)
	tasks[0] = t1
	tasks[1] = t2
	tasks[2] = t3

	// Illegal, but for this function creates a circle
	t2.AddUpTrigger("t1")
	t3.AddUpTrigger("t2")
	t3.AddDownTrigger("t1")

	deps := normalizeDeps(tasks)
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
	t1, t2, t3 := NewTask("t1"), NewTask("t2"), NewTask("t3")
	tasks := make([]*Task, 3)
	tasks[0] = t1
	tasks[1] = t2
	tasks[2] = t3

	// Illegal, but for this function creates a circle
	t2.AddUpTrigger("t1")
	t3.AddUpTrigger("t2")
	t3.AddDownTrigger("t1")

	deps := normalizeDeps(tasks)
	order, cycle := topSort(deps)
	sort.Strings(cycle)
	fmt.Printf("%s\n", order)
	fmt.Printf("%s\n", cycle)

	// Output:
	// []
	// [t1 t2 t3]
}

func ExampleTopSortDeps() {
	t1, t2, t3 := NewTask("t1"), NewTask("t2"), NewTask("t3")
	tasks := make([]*Task, 3)
	tasks[0] = t1
	tasks[1] = t2
	tasks[2] = t3

	t1.AddDownTrigger("t3")
	t2.AddUpTrigger("t1")
	t3.AddUpTrigger("t2")

	deps := normalizeDeps(tasks)
	order, cycle := topSort(deps)
	sort.Strings(cycle)
	fmt.Printf("%s\n", order)
	fmt.Printf("%s\n", cycle)

	// Output:
	// [t1 t2 t3]
	// []
}
