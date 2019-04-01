package helper

type HealthCheck struct {
}

func (h *HealthCheck) Check() error {
	tool := Utils{}
	if err := tool.SystemHealth(); err != nil {
		return err
	} else {
		return nil
	}
}