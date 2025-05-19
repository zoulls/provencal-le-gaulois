package task

import (
	"sync"

	"github.com/robfig/cron/v3"
)

// ListTasks is the list of tasks
var ListTasks []List

// ListTasksOnce is the once for the task list
var ListTasksOnce sync.Once

// List is the list of task jobs
type List struct {
	TaskName string
	ID       int
}

// Option is the option for task
type Option struct {
	Cron     *cron.Cron
	Spec     string
	TaskName string
	Task     func()
}

// GetListTasks get the task list
func GetListTasks() []List {
	ListTasksOnce.Do(func() {
		ListTasks = make([]List, 0)
	})
	return ListTasks
}

// setListTasks set the task list
func setListTasks(taskList []List) {
	if taskList == nil {
		ListTasksOnce.Do(func() {
			ListTasks = make([]List, 0)
		})
	}
	ListTasks = taskList
}

// CreateTask create a task
func CreateTask(option Option) error {
	// First exec of the job
	option.Task()

	// Add job to task
	id, err := option.Cron.AddFunc(option.Spec, option.Task)
	if err != nil {
		return err
	}

	// Add job to list
	ListTasks = GetListTasks()
	ListTasks = append(ListTasks, List{
		TaskName: option.TaskName,
		ID:       int(id),
	})
	setListTasks(ListTasks)

	// Start task
	option.Cron.Start()

	return nil
}

// DeleteTask delete a task
func DeleteTask(cronParam *cron.Cron, id int) {
	// Get task list
	taskList := GetListTasks()

	// Remove task from cron
	cronParam.Remove(cron.EntryID(id))

	// Delete task from list
	for i, v := range taskList {
		if v.ID == id {
			taskList = append(taskList[:i], taskList[i+1:]...)
			break
		}
	}

	// Set task list
	setListTasks(taskList)
}
