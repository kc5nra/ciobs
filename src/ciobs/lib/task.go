package ciobs

type Task struct {
	name         string
	upTriggers   map[string]bool
	downTriggers map[string]bool
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

func NewTask(name string) *Task {
	return &Task{
		name:         name,
		upTriggers:   make(map[string]bool),
		downTriggers: make(map[string]bool),
	}
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
