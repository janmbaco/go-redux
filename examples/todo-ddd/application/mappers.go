package application

import "github.com/janmbaco/go-redux/v2/examples/todo-ddd/domain"

// ToEntity converts TodoDTO to domain Todo
func ToEntity(dto TodoDTO) *domain.Todo {
	return &domain.Todo{
		ID:        dto.ID,
		Text:      dto.Text,
		Completed: dto.Completed,
		CreatedAt: dto.CreatedAt,
	}
}

// ToDTO converts domain Todo to TodoDTO
func ToDTO(todo *domain.Todo) TodoDTO {
	return TodoDTO{
		ID:        todo.ID,
		Text:      todo.Text,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt,
	}
}

// ToEntities converts slice of TodoDTO to slice of domain Todo
func ToEntities(dtos []TodoDTO) []*domain.Todo {
	entities := make([]*domain.Todo, len(dtos))
	for i, dto := range dtos {
		entities[i] = ToEntity(dto)
	}
	return entities
}

// ToDTOs converts slice of domain Todo to slice of TodoDTO
func ToDTOs(todos []*domain.Todo) []TodoDTO {
	dtos := make([]TodoDTO, len(todos))
	for i, todo := range todos {
		dtos[i] = ToDTO(todo)
	}
	return dtos
}
