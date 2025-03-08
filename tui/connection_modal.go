package tui

import (
	"fmt"
	"strconv"

	"github.com/nao1215/sqluv/config"
	"github.com/rivo/tview"
)

// connectionModal represents the modal for database connection options
type connectionModal struct {
	*tview.Modal
	app          *tview.Application
	onClose      func(conn *config.DBConnection)
	configMgr    *config.DBConfig
	selectedConn *config.DBConnection
}

// newConnectionModal creates a new connection modal
func newConnectionModal(app *tview.Application, onClose func(conn *config.DBConnection)) *connectionModal {
	modal := tview.NewModal().
		SetText("No local files found. Would you like to connect to a database?").
		AddButtons([]string{"New", "List", "Cancel"})

	// Create config manager
	configMgr, err := config.NewDBConfig()
	if err != nil {
		// Handle error gracefully - could show an error modal
		modal.SetText("Error: Could not initialize configuration manager\n" + err.Error())
	}

	cm := &connectionModal{
		Modal:     modal,
		app:       app,
		onClose:   onClose,
		configMgr: configMgr,
	}

	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		switch buttonLabel {
		case "New":
			cm.showNewConnectionForm()
		case "List":
			cm.showConnectionsList()
		case "Cancel":
			if onClose != nil {
				onClose(nil)
			}
			cm.app.Stop()
		}
	})
	return cm
}

// showNewConnectionForm displays the form for adding a new database connection
func (cm *connectionModal) showNewConnectionForm() {
	form := tview.NewForm()

	form.AddInputField("Connection Name", "", 30, nil, nil)

	dbmsTypes := []string{string(config.MySQL), string(config.PostgreSQL)}
	form.AddDropDown("DBMS Type", dbmsTypes, 0, nil)
	form.AddInputField("Host", "127.0.0.1", 50, nil, nil)
	form.AddInputField("Port", "3306", 10, func(_ string, lastChar rune) bool {
		return '0' <= lastChar && lastChar <= '9'
	}, nil)

	form.AddInputField("Username", "", 30, nil, nil)
	form.AddPasswordField("Password", "", 30, '*', nil)
	form.AddInputField("Database Name", "", 30, nil, nil)

	form.AddButton("Save", func() {
		name := form.GetFormItem(0).(*tview.InputField).GetText()
		_, dbmsType := form.GetFormItem(1).(*tview.DropDown).GetCurrentOption()
		host := form.GetFormItem(2).(*tview.InputField).GetText()
		portStr := form.GetFormItem(3).(*tview.InputField).GetText()
		username := form.GetFormItem(4).(*tview.InputField).GetText()
		password := form.GetFormItem(5).(*tview.InputField).GetText()
		database := form.GetFormItem(6).(*tview.InputField).GetText()

		port, _ := strconv.Atoi(portStr) //nolint:errcheck // Error is handled by the form validation

		conn := config.DBConnection{
			Name:     name,
			Type:     config.DBMSType(dbmsType),
			Host:     host,
			Port:     port,
			User:     username,
			Password: password,
			Database: database,
		}

		if err := cm.configMgr.SaveConnection(conn); err != nil {
			cm.showError(err.Error())
			return
		}

		// Set selected connection and call onClose with it
		cm.selectedConn = &conn
		if cm.onClose != nil {
			cm.onClose(cm.selectedConn)
		}
	})

	form.AddButton("Cancel", func() {
		cm.app.SetRoot(cm.Modal, true)
	})

	form.SetBorder(true).SetTitle("New Database Connection")
	cm.app.SetRoot(form, true)
}

// showConnectionsList displays the list of saved connections
func (cm *connectionModal) showConnectionsList() {
	// Load connections from config file
	connections, err := cm.configMgr.LoadConnections()
	if err != nil {
		cm.showError(fmt.Sprintf("Failed to load connections: %v", err))
		return
	}

	if len(connections) == 0 {
		cm.showError("No saved connections found")
		return
	}

	// Create a list to display connections
	list := tview.NewList()
	list.SetTitle("Saved Connections").SetBorder(true)

	// Add connections to the list
	for i, conn := range connections {
		list.AddItem(
			conn.Name,
			fmt.Sprintf("%s %s:%d", conn.Type, conn.Host, conn.Port),
			rune('1'+i),
			func() {
				cm.connectToDatabase(conn)
			})
	}

	// Add a "Back" button at the bottom
	list.AddItem("Back", "Return to connection options", '*', func() {
		cm.app.SetRoot(cm.Modal, true)
	})

	// Show the list
	cm.app.SetRoot(list, true)
}

// connectToDatabase connects to the selected database
func (cm *connectionModal) connectToDatabase(conn config.DBConnection) {
	// Store the selected connection
	cm.selectedConn = &conn

	// Return the selected connection to TUI.Run
	if cm.onClose != nil {
		cm.onClose(cm.selectedConn)
	}
}

// showError displays an error message
func (cm *connectionModal) showError(message string) {
	errorModal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			cm.app.SetRoot(cm.Modal, true)
		})

	cm.app.SetRoot(errorModal, true)
}
