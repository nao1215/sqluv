package tui

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
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
	theme        *Theme
}

// newConnectionModal creates a new connection modal
func newConnectionModal(app *tview.Application, theme *Theme, onClose func(conn *config.DBConnection)) *connectionModal {
	modal := tview.NewModal().
		SetText("No local files found. Would you like to connect to a database?").
		AddButtons([]string{"New", "List", "Cancel"})

	colors := theme.GetColors()
	modal.SetBackgroundColor(colors.Background)
	modal.SetTextColor(colors.Foreground)
	modal.SetButtonBackgroundColor(colors.Button)
	modal.SetButtonTextColor(colors.ButtonText)

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
		theme:     theme,
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
	cm.applyTheme(theme)
	return cm
}

// showNewConnectionForm displays the form for adding a new database connection
func (cm *connectionModal) showNewConnectionForm() {
	form := tview.NewForm()

	form.AddInputField("Connection Name", "", 0, nil, nil)
	dbmsTypes := []string{string(config.MySQL), string(config.PostgreSQL), string(config.SQLite3), string(config.SQLServer)}

	// First add all form fields without the callback
	form.AddDropDown("DBMS Type", dbmsTypes, 0, nil)
	form.AddInputField("Host", "127.0.0.1", 0, nil, nil)
	form.AddInputField("Port", "3306", 0, func(_ string, lastChar rune) bool {
		return '0' <= lastChar && lastChar <= '9'
	}, nil)
	form.AddInputField("Username", "", 0, nil, nil)
	form.AddPasswordField("Password", "", 0, '*', nil)
	form.AddInputField("Database Name", "", 0, nil, nil)
	form.AddInputField("Database File Path (SQLite3 only)", "", 0, nil, nil)

	// Now that all fields exist, we can set up the callback for the dropdown
	dbmsDropdown := form.GetFormItem(1).(*tview.DropDown)
	hostField := form.GetFormItem(2).(*tview.InputField)
	portField := form.GetFormItem(3).(*tview.InputField)
	usernameField := form.GetFormItem(4).(*tview.InputField)
	passwordField := form.GetFormItem(5).(*tview.InputField)
	databaseField := form.GetFormItem(6).(*tview.InputField)
	filePathField := form.GetFormItem(7).(*tview.InputField)

	// Initially disable the SQLite3 file path field
	filePathField.SetDisabled(true)

	// Set the change handler for the dropdown
	dbmsDropdown.SetSelectedFunc(func(option string, _ int) {
		if option == string(config.SQLite3) {
			// For SQLite3, disable network and authentication fields
			hostField.SetLabel("Host (N/A for SQLite3)")
			hostField.SetDisabled(true)
			hostField.SetText("N/A")

			portField.SetLabel("Port (N/A for SQLite3)")
			portField.SetDisabled(true)
			portField.SetText("0")

			usernameField.SetLabel("Username (N/A for SQLite3)")
			usernameField.SetDisabled(true)
			usernameField.SetText("N/A")

			passwordField.SetLabel("Password (N/A for SQLite3)")
			passwordField.SetDisabled(true)
			passwordField.SetText("N/A")

			databaseField.SetLabel("Database Name (N/A for SQLite3)")
			databaseField.SetDisabled(true)
			databaseField.SetText("N/A")

			filePathField.SetLabel("Database File Path")
			filePathField.SetDisabled(false)
		} else {
			// For other DBMSes, enable network and authentication fields
			hostField.SetLabel("Host")
			hostField.SetDisabled(false)
			hostField.SetText("127.0.0.1")

			portField.SetLabel("Port")
			portField.SetDisabled(false)
			if option == string(config.MySQL) {
				portField.SetText("3306")
			} else if option == string(config.PostgreSQL) {
				portField.SetText("5432")
			} else if option == string(config.SQLServer) {
				portField.SetText("1433")
			}

			usernameField.SetLabel("Username")
			usernameField.SetDisabled(false)
			usernameField.SetText("")

			passwordField.SetLabel("Password")
			passwordField.SetDisabled(false)
			passwordField.SetText("")

			databaseField.SetLabel("Database Name")
			databaseField.SetDisabled(false)
			databaseField.SetText("")

			filePathField.SetLabel("Database File Path (SQLite3 only)")
			filePathField.SetDisabled(true)
			filePathField.SetText("")
		}
	})

	form.AddButton("Save", func() {
		name := form.GetFormItem(0).(*tview.InputField).GetText()
		dbmsTypeIndex, _ := form.GetFormItem(1).(*tview.DropDown).GetCurrentOption()
		dbmsType := dbmsTypes[dbmsTypeIndex]
		host := form.GetFormItem(2).(*tview.InputField).GetText()
		portStr := form.GetFormItem(3).(*tview.InputField).GetText()
		username := form.GetFormItem(4).(*tview.InputField).GetText()
		password := form.GetFormItem(5).(*tview.InputField).GetText()
		database := form.GetFormItem(6).(*tview.InputField).GetText()
		filePath := form.GetFormItem(7).(*tview.InputField).GetText()

		port, _ := strconv.Atoi(portStr) //nolint:errcheck // Error is handled by the form validation

		conn := config.DBConnection{
			Name: name,
			Type: config.DBMSType(dbmsType),
		}

		// Set appropriate fields based on DBMS type
		if conn.Type == config.SQLite3 {
			// For SQLite3, only the file path is relevant
			conn.Database = filePath
			// Leave other fields with empty/default values
			conn.Host = "N/A"
			conn.Port = 0
			conn.User = "N/A"
			conn.Password = ""
		} else {
			// For other DBMS types, set all connection properties
			conn.Host = host
			conn.Port = port
			conn.User = username
			conn.Password = password
			conn.Database = database
		}

		// Save connection in configuration
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

	// Retrieve colors from the current theme (via config.ColorConfig)
	colors := cm.theme.GetColors()
	form.SetBackgroundColor(colors.Background)
	form.SetBorderColor(colors.Border)
	form.SetTitleColor(colors.Header)
	form.SetFieldTextColor(colors.Foreground)
	form.SetFieldBackgroundColor(colors.Background)
	form.SetLabelColor(colors.Foreground)
	form.SetBorderStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.BorderFocus))
	form.SetButtonStyle(tcell.StyleDefault.
		Background(colors.Button).
		Foreground(colors.ButtonText))
	form.SetFieldStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.Foreground))

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
	colors := cm.theme.GetColors()
	list := tview.NewList()
	list.SetTitle("Saved Connections").SetBackgroundColor(colors.Background)

	// Add connections to the list
	for i, conn := range connections {
		if conn.Type == config.SQLite3 {
			list.AddItem(
				conn.Name,
				fmt.Sprintf("%s database=%s", conn.Type, conn.Database),
				rune('1'+i),
				func() {
					cm.connectToDatabase(conn)
				})
			continue
		}
		list.AddItem(
			conn.Name,
			fmt.Sprintf("%s %s:%d database=%s", conn.Type, conn.Host, conn.Port, conn.Database),
			rune('1'+i),
			func() {
				cm.connectToDatabase(conn)
			})
	}

	// Add a "Back" button at the bottom
	list.AddItem("Back", "Return to connection options", '*', func() {
		cm.app.SetRoot(cm.Modal, true)
	})

	list.SetBorder(true)
	list.SetTitleColor(colors.Header)
	list.SetBorderStyle(tcell.StyleDefault.
		Background(colors.Border).
		Foreground(colors.BorderFocus))
	list.SetBackgroundColor(colors.Background)
	list.SetMainTextStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.Foreground))
	list.SetSecondaryTextStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.ButtonText))
	list.SetSelectedBackgroundColor(colors.Selection)
	list.SetShortcutStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.Foreground))

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

func (cm *connectionModal) applyTheme(theme *Theme) {
	colors := theme.GetColors()

	cm.Modal.SetBorder(true).SetBorderColor(colors.Border)
	cm.Modal.SetBackgroundColor(colors.Background)
	cm.Modal.SetTextColor(colors.Foreground)
	cm.Modal.SetButtonBackgroundColor(colors.Button)
	cm.Modal.SetButtonTextColor(colors.ButtonText)
	cm.Modal.SetBorderStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.BorderFocus))
	cm.Modal.SetButtonActivatedStyle(tcell.StyleDefault.
		Background(colors.ButtonFocus).
		Foreground(colors.ButtonTextFocus))
	cm.Modal.SetButtonStyle(tcell.StyleDefault.
		Background(colors.Button).
		Foreground(colors.ButtonText))

}
