package provider

import "testing"

func Test_RunSpeedTest(t *testing.T) {
	netflix, err := NewNetflix()
	if err != nil {
		t.Errorf("create netflix provider: %s", err)
	}

	ookla, err := NewOokla()
	if err != nil {
		t.Errorf("create ookla provider: %s", err)
	}

	tests := []struct {
		name     string
		provider Provider
	}{
		{
			name:     "netflix happy path",
			provider: netflix,
		},
		{
			name:     "ookla happy path",
			provider: ookla,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.provider.RunSpeedTest()
			if err != nil {
				t.Errorf("speed test went wrong: %s", err)
			}
			if res == nil {
				t.Error("Non-nil result expected.")
			}
			if res.Upload == 0 || res.Download == 0 {
				t.Error("Unperformed speed test.")
			}
		})
	}
}
