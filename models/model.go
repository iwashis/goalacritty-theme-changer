package models

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	cf "goalacritty_themes/config"
	it "goalacritty_themes/theme_tools"
)

const (
	listHeight      = 14
	listWidth       = 30
	sampleTextWidth = 50
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("150"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	frameStyle        = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2).Margin(1).BorderForeground(lipgloss.Color("63"))
	frameTitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).PaddingLeft(2)
)

type item struct {
	title, desc string
}

func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// model for the choice of the theme. Usable only once the theme repository
// is in place
type model struct {
	list          list.Model
	choice        string
	initialTheme  string
	quitting      bool
	config        cf.Config
	previousIndex int
	currentTheme  it.ThemeData
	sampleText    string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width - sampleTextWidth - 4) // Adjust list width for the sample text
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.title
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	// Detect if the highlighted item has changed
	currentIndex := m.list.Index()
	if currentIndex != m.previousIndex {
		i, ok := m.list.SelectedItem().(item)
		if ok {
			themeData := it.ThemeData{
				Name:     i.title,
				FullPath: i.desc,
			}
			it.UpdateAlacrittyConfigFile(m.config, themeData)
		}
		m.previousIndex = currentIndex
	}

	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("Selected theme: %s", m.choice))
	}
	if m.quitting {
		// We first have to reset the config file
		it.UpdateAlacrittyConfigFile(m.config, m.currentTheme)
		return quitTextStyle.Render("Not making a selection? That’s cool.")
	}

	sampleText := m.sampleText //formatSampleText(m.sampleText)
	sampleFrame := frameStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			frameTitleStyle.Render("Sample"),
			sampleText,
		),
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.list.View(),
		sampleFrame,
	)
}


func InitializeMainModel(config cf.Config) model {
	currentTheme, err := it.GetCurrentTheme(config)
	if err != nil {
		fmt.Println("Error getting the current theme:", err)
		os.Exit(1)
	}
	themedataList, err := it.GetThemeDataNames(config)
	if err != nil {
		fmt.Println("Error getting theme data:", err)
		os.Exit(1)
	}

	var items []list.Item
	for _, theme := range themedataList {
		items = append(items, item{title: theme.Name, desc: theme.FullPath})
	}

	l := list.New(items, itemDelegate{}, listWidth, listHeight)
	l.Title = "Select a Theme"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	sampleText :=
		"|039| \033[39mDefault \033[m      |049| \033[49mDefault \033[m      |037| \033[37mLight gray \033[m     |047| \033[47mLight gray \033[m" + "\n" +
			"|030| \033[30mBlack \033[m        |040| \033[40mBlack \033[m        |090| \033[90mDark gray \033[m      |100| \033[100mDark gray \033[m" + "\n" +
			"|031| \033[31mRed \033[m          |041| \033[41mRed \033[m          |091| \033[91mLight red \033[m      |101| \033[101mLight red \033[m" + "\n" +
			"|032| \033[32mGreen \033[m        |042| \033[42mGreen \033[m        |092| \033[92mLight green \033[m    |102| \033[102mLight green \033[m" + "\n" +
			"|033| \033[33mYellow \033[m       |043| \033[43mYellow \033[m       |093| \033[93mLight yellow \033[m   |103| \033[103mLight yellow \033[m" + "\n" +
			"|034| \033[34mBlue \033[m         |044| \033[44mBlue \033[m         |094| \033[94mLight blue \033[m     |104| \033[104mLight blue \033[m" + "\n" +
			"|035| \033[35mMagenta \033[m      |045| \033[45mMagenta \033[m      |095| \033[95mLight magenta \033[m  |105| \033[105mLight magenta \033[m" + "\n" +
			"|036| \033[36mCyan \033[m         |046| \033[46mCyan \033[m         |096| \033[96mLight cyan \033[m     |106| \033[106mLight cyan \033[m"

	return model{
		list:          l,
		config:        config,
		previousIndex: -1, // Initialize to an invalid index
		currentTheme:  *currentTheme,
		sampleText:    sampleText, // "Lorem ipsum dolor sit amet,\nconsectetur adipiscing elit.\nPhasellus imperdiet...",
	}
}
