package tasks

type TaskRepository interface {
    Create(task *Task) (*Task, error)
    ShowTasks() ([]Task, error)
    UpdateTask(id int) (*Task, error)
    GetTaskByID (id int) (*Task, error)
    DeleteTask(id int) (*Task, error)
}
