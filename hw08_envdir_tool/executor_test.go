package main

import (
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name        string
		cmd         []string
		env         Environment
		expectedEnv map[string]string
		expectedRC  int
	}{
		{
			name: "Successful command",
			cmd:  []string{"echo", "Test"},
			env: Environment{
				"FOO": EnvValue{Value: "123"},
				"BAR": EnvValue{Value: "value"},
			},
			expectedEnv: map[string]string{
				"FOO": "123",
				"BAR": "value",
			},
			expectedRC: 0,
		},
		{
			name: "Command with environment removal",
			cmd:  []string{"echo", "Test!"},
			env: Environment{
				"FOO": EnvValue{Value: "123"},
				"BAR": EnvValue{NeedRemove: true},
			},
			expectedEnv: map[string]string{
				"FOO": "123",
			},
			expectedRC: 0,
		},
		{
			name: "Command with error",
			cmd:  []string{"false"}, // Err
			env: Environment{
				"FOO": EnvValue{Value: "123"},
			},
			expectedEnv: map[string]string{
				"FOO": "123",
			},
			expectedRC: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnCode := RunCmd(tt.cmd, tt.env)

			if returnCode != tt.expectedRC {
				t.Errorf("Expected return code %d, got %d", tt.expectedRC, returnCode)
			}

			for key, expectedValue := range tt.expectedEnv {
				actualValue, exists := os.LookupEnv(key)
				if !exists {
					t.Errorf("Expected environment variable %s to be set", key)
				} else if actualValue != expectedValue {
					t.Errorf("For environment variable %s, expected %s, got %s", key, expectedValue, actualValue)
				}
			}

			for key, envValue := range tt.env {
				if envValue.NeedRemove {
					if _, exists := os.LookupEnv(key); exists {
						t.Errorf("Expected environment variable %s to be unset", key)
					}
				}
			}
		})
	}
}
