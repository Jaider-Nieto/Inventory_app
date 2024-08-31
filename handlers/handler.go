package handlers

import interfaces "github.com/jaider-nieto/ecommerce-go/interfaces"


type Handler struct {
	userRepository interfaces.UserRepositoryInterface
	taskRepository interfaces.TaskRepositoryInterface
}
