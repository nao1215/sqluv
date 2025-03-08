package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

// DBMSType represents the type of DBMS
type DBMSType string

const (
	// MySQL represents the MySQL DBMS
	MySQL DBMSType = "MySQL"
	// PostgreSQL represents the PostgreSQL DBMS
	PostgreSQL DBMSType = "PostgreSQL"
)

// DBConnection represents a database connection configuration as a value object
type DBConnection struct {
	Name     string   `yaml:"name"`
	Type     DBMSType `yaml:"type"`
	Host     string   `yaml:"host"`
	Port     int      `yaml:"port"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	Database string   `yaml:"database"`
}

// DBConfigFile represents the structure of the dbms.yml file
type DBConfigFile struct {
	Connections []DBConnection `yaml:"connections"`
}

// DBConfig manages database configurations
type DBConfig struct {
	configPath string
}

// NewDBConfig creates a new database configuration manager
func NewDBConfig() (*DBConfig, error) {
	configDir, err := getConfigDir()
	if err != nil {
		configDir, err = ensureConfigDir()
		if err != nil {
			return nil, err
		}
	}
	return &DBConfig{
		configPath: filepath.Join(configDir, "dbms.yml"),
	}, nil
}

// getConfigDir returns the configuration directory path
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".config", "sqluv")
	if _, err := os.Stat(configDir); err != nil {
		return "", err
	}
	return configDir, nil
}

// ensureConfigDir ensures the config directory exists and returns its path
func ensureConfigDir() (string, error) {
	configDir := filepath.Join(xdg.ConfigHome, ".config", "sqluv")
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return "", err
	}
	return configDir, nil
}

// SaveConnection saves a database connection to the config file
func (cm *DBConfig) SaveConnection(conn DBConnection) error {
	// Encrypt the password before saving
	encryptedPassword, err := EncryptPassword(conn.Password)
	if err != nil {
		return err
	}

	// Create a copy with encrypted password
	encryptedConn := conn
	encryptedConn.Password = encryptedPassword

	config, err := cm.loadConfigFile()
	if err != nil {
		config = &DBConfigFile{Connections: []DBConnection{}}
	}

	// Check if connection with the same name already exists
	for i, existingConn := range config.Connections {
		if existingConn.Name == encryptedConn.Name {
			// Replace existing connection
			config.Connections[i] = encryptedConn
			return cm.saveConfigFile(config)
		}
	}

	// Add new connection
	config.Connections = append(config.Connections, encryptedConn)
	return cm.saveConfigFile(config)
}

// LoadConnections loads all database connections from the config file
func (cm *DBConfig) LoadConnections() ([]DBConnection, error) {
	config, err := cm.loadConfigFile()
	if err != nil {
		return nil, err
	}

	// Decrypt passwords for all connections
	for i, conn := range config.Connections {
		if IsEncrypted(conn.Password) {
			decryptedPassword, err := DecryptPassword(conn.Password)
			if err != nil {
				// If decryption fails, keep the encrypted password
				continue
			}
			config.Connections[i].Password = decryptedPassword
		}
	}

	return config.Connections, nil
}

// loadConfigFile loads the configuration file
func (cm *DBConfig) loadConfigFile() (*DBConfigFile, error) {
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	}

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, err
	}

	config := &DBConfigFile{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

// saveConfigFile saves the configuration to file
func (cm *DBConfig) saveConfigFile(config *DBConfigFile) error {
	// Convert to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(cm.configPath, data, 0644)
}

// GetConnectionByName retrieves a specific connection by name
func (cm *DBConfig) GetConnectionByName(name string) (DBConnection, error) {
	connections, err := cm.LoadConnections()
	if err != nil {
		return DBConnection{}, err
	}
	for _, conn := range connections {
		if conn.Name == name {
			return conn, nil
		}
	}
	return DBConnection{}, fmt.Errorf("connection '%s' not found", name)
}

// RemoveConnection removes a database connection from the config file by name
func (cm *DBConfig) RemoveConnection(name string) error {
	config, err := cm.loadConfigFile()
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	// Find the connection and remove it
	newConnections := make([]DBConnection, 0, len(config.Connections))
	found := false
	for _, conn := range config.Connections {
		if conn.Name != name {
			newConnections = append(newConnections, conn)
		} else {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("connection '%s' not found in configuration", name)
	}

	config.Connections = newConnections
	if err := cm.saveConfigFile(config); err != nil {
		return fmt.Errorf("failed to save updated config file: %w", err)
	}
	return nil
}
