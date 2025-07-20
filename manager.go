package mint

type Manager struct {
	config Config
}

func (m *Manager) Init(config Config) {
	m.config = config
}

func (m *Manager) Start() {}

func New() Manager {
	return Manager{}
}
