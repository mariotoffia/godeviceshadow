package managermodel

// Manager is a "full" manager interface that includes all the other manager interfaces.
type Manager interface {
	Reportable // Report functions
	Desireable // Desire functions
	Lister     // Query functions
	Receiver   // Read functions
	Remover    // Delete functions
}
