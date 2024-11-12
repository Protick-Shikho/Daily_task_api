package tasks

type TaskRepository interface {

    SetupDatabase()
    // NewTaskRepository()
    Create(task *Task) (*Task, error)
    ShowTasks() ([]Task, error)
    UpdateTask(id int) (*Task, error)
    GetTaskByID (id int) (*Task, error)
    DeleteTask(id int) (*Task, error)
    Close() error
}
