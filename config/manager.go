package config

type Manager struct {
	source Source
}

func (m *Manager) SetSource(s Source) {
	m.source = s
}

func (m *Manager) GetConfig() (*Config, error) {
	return m.source.GetConfig()
}

func (m *Manager) CloseSource() {
	m.source.Close()
}
