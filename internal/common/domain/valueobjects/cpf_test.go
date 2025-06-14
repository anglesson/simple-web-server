package common_vo_test

import (
	"testing"

	common_vo "github.com/anglesson/simple-web-server/internal/common/domain/valueobjects"
)

func TestNewCPF(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid CPF",
			input:   "529.982.247-25",
			wantErr: false,
		},
		{
			name:    "valid CPF without formatting",
			input:   "52998224725",
			wantErr: false,
		},
		{
			name:    "invalid CPF - wrong length",
			input:   "123.456.789-0",
			wantErr: true,
		},
		{
			name:    "invalid CPF - all same digits",
			input:   "111.111.111-11",
			wantErr: true,
		},
		{
			name:    "invalid CPF - wrong check digit",
			input:   "529.982.247-26",
			wantErr: true,
		},
		{
			name:    "invalid CPF - empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpf, err := common_vo.NewCPF(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCPF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && cpf == nil {
				t.Error("NewCPF() returned nil CPF when no error was expected")
			}
		})
	}
}

func TestCPF_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "format valid CPF",
			input:    "52998224725",
			expected: "529.982.247-25",
		},
		{
			name:     "format already formatted CPF",
			input:    "529.982.247-25",
			expected: "529.982.247-25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpf, err := common_vo.NewCPF(tt.input)
			if err != nil {
				t.Fatalf("NewCPF() error = %v", err)
			}
			if got := cpf.String(); got != tt.expected {
				t.Errorf("CPF.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCPF_Equal(t *testing.T) {
	cpf1, _ := common_vo.NewCPF("529.982.247-25")
	cpf2, _ := common_vo.NewCPF("529.982.247-25")
	cpf3, _ := common_vo.NewCPF("123.456.789-09")

	tests := []struct {
		name     string
		cpf1     *common_vo.CPF
		cpf2     *common_vo.CPF
		expected bool
	}{
		{
			name:     "equal CPFs",
			cpf1:     cpf1,
			cpf2:     cpf2,
			expected: true,
		},
		{
			name:     "different CPFs",
			cpf1:     cpf1,
			cpf2:     cpf3,
			expected: false,
		},
		{
			name:     "nil CPF",
			cpf1:     cpf1,
			cpf2:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cpf1.Equal(tt.cpf2); got != tt.expected {
				t.Errorf("CPF.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCPF_Value(t *testing.T) {
	cpf, err := common_vo.NewCPF("529.982.247-25")
	if err != nil {
		t.Fatalf("NewCPF() error = %v", err)
	}

	expected := "52998224725"
	if got := cpf.Value(); got != expected {
		t.Errorf("CPF.Value() = %v, want %v", got, expected)
	}
}
