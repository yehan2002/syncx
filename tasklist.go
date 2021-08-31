package syncx

import (
	"fmt"
	"sync"
)

//TaskList a  list of tasks
type TaskList struct {
	tasks []func()

	once sync.Once
	mux  sync.Mutex
}

//Add adds a task to the task list.
func (t *TaskList) Add(task func()) {
	t.mux.Lock()
	t.tasks = append(t.tasks, task)
	t.mux.Unlock()
}

//Run run the task list.
func (t *TaskList) Run() (errors []error) {
	t.mux.Lock()
	t.once.Do(func() {
		for i := len(t.tasks) - 1; i >= 0; i-- {
			if err := t.run(t.tasks[i]); err != nil {
				errors = append(errors, err)
			}
		}
	})
	t.mux.Unlock()
	return errors
}

func (t *TaskList) run(r func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("panic: %s", r)
			}
		}
	}()
	r()
	return
}
