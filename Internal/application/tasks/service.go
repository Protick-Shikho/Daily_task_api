package tasks

type TaskService struct {
    repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

func (s *TaskService) Create(task *Task) (*Task, error) {
    return s.repo.Create(task)
}

func (s *TaskService) ShowTasks() ([]Task, error) {
	return s.repo.ShowTasks()
}

func (s *TaskService) UpdateTask(id int) (*Task, error) {
	return s.repo.UpdateTask(id)
}

func (s *TaskService) DeleteTask(id int) (*Task, error) {
    return s.repo.DeleteTask(id)
}
