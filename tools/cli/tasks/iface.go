package tasks

type Task interface {
	ID() string
}

type genericTask struct {
	id string
}

func NewGenericTask(id string) Task {
	return &genericTask{
		id: id,
	}
}

func (g *genericTask) ID() string {
	return g.id
}
