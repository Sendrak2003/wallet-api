package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type EnvironmentValidator struct {
	requiredVars []string
}

func NewEnvironmentValidator() *EnvironmentValidator {
	return &EnvironmentValidator{
		requiredVars: []string{
			"APP_PORT",
			"DB_HOST", 
			"DB_PORT",
			"DB_USER",
			"DB_PASSWORD",
			"DB_NAME",
		},
	}
}

func (v *EnvironmentValidator) ValidateRequired() error {
	var missingVars []string
	
	for _, varName := range v.requiredVars {
		if value := os.Getenv(varName); value == "" {
			missingVars = append(missingVars, varName)
		}
	}
	
	if len(missingVars) > 0 {
		return fmt.Errorf("отсутствуют обязательные переменные окружения: %s", 
			strings.Join(missingVars, ", "))
	}
	
	return nil
}

func (v *EnvironmentValidator) ValidateFormats() error {
	if err := v.validatePort("APP_PORT"); err != nil {
		return err
	}
	
	if err := v.validatePort("DB_PORT"); err != nil {
		return err
	}
	
	stringVars := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, varName := range stringVars {
		if err := v.validateNonEmptyString(varName); err != nil {
			return err
		}
	}
	
	return nil
}

func (v *EnvironmentValidator) validatePort(varName string) error {
	value := os.Getenv(varName)
	if value == "" {
		return fmt.Errorf("переменная %s не установлена", varName)
	}
	
	port, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("переменная %s должна быть числом, получено: %s", varName, value)
	}
	
	if port < 1 || port > 65535 {
		return fmt.Errorf("переменная %s должна быть в диапазоне 1-65535, получено: %d", varName, port)
	}
	
	return nil
}

func (v *EnvironmentValidator) validateNonEmptyString(varName string) error {
	value := os.Getenv(varName)
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("переменная %s не может быть пустой", varName)
	}
	
	return nil
}

func (v *EnvironmentValidator) ValidateAll() error {
	if err := v.ValidateRequired(); err != nil {
		return err
	}
	
	if err := v.ValidateFormats(); err != nil {
		return err
	}
	
	return nil
}