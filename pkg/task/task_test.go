package task

import (
	"testing"

	"github.com/robfig/cron/v3"
)

func TestTaskLifecycleAndURL(t *testing.T) {
	c := cron.New()
	defer c.Stop()

	// Reset tasks
	setListTasks(nil)

	executed := false
	opt := Option{
		Cron:     c,
		Spec:     "@every 1h",
		TaskName: "test-task",
		URL:      "https://store.steampowered.com/feeds/news/app/730/",
		Task: func() {
			executed = true
		},
	}

	err := CreateTask(opt)
	if err != nil {
		t.Fatalf("unexpected error creating task: %v", err)
	}

	if !executed {
		t.Errorf("expected initial task execution on CreateTask")
	}

	tasks := GetListTasks()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task in list, got %d", len(tasks))
	}

	created := tasks[0]
	if created.TaskName != "test-task" {
		t.Errorf("expected task name 'test-task', got %s", created.TaskName)
	}
	if created.URL != "https://store.steampowered.com/feeds/news/app/730/" {
		t.Errorf("expected task URL to match, got %s", created.URL)
	}
	if !created.Valid() {
		t.Errorf("expected created task to be valid")
	}

	retrieved := GetTask(created.ID)
	if retrieved.ID != created.ID || retrieved.URL != created.URL {
		t.Errorf("GetTask failed to retrieve matching task")
	}

	// Delete task
	DeleteTask(c, created.ID)
	tasksAfterDelete := GetListTasks()
	if len(tasksAfterDelete) != 0 {
		t.Errorf("expected 0 tasks after delete, got %d", len(tasksAfterDelete))
	}

	nonExistent := GetTask(9999)
	if nonExistent.Valid() {
		t.Errorf("expected non-existent task to be invalid")
	}
}
