package views // import "github.com/jenkins-x/octant-jx/pkg/plugin/views"

import (
	"fmt"
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/jenkins-x/octant-jx/pkg/admin"
	"github.com/jenkins-x/octant-jx/pkg/admin/workspaces"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

var (
	sortByIndex = true
)

type WorkspacesViewConfig struct {
	TeamFilter        component.TableFilter
	EnvironmentFilter component.TableFilter
}

func (f *WorkspacesViewConfig) TableFilters() []*component.TableFilter {
	return []*component.TableFilter{&f.TeamFilter, &f.EnvironmentFilter}
}

func BuildWorkspacesView(request service.Request, ws []workspaces.WorkspaceOctant) (component.Component, error) {
	//Todo: request argument is not being used, but removing it requires a bit more work...
	log.Logger().Debugf("got list of Workspaces %d\n", len(ws))

	header := viewhelpers.NewMarkdownText(viewhelpers.ToBreadcrumbMarkdown(admin.RootBreadcrumb, "Workspaces"))

	config := &WorkspacesViewConfig{}

	table := component.NewTableWithRows(
		"Workspaces", "There are no Workspaces!",
		component.NewTableCols("Name", "Source", "Team", "Environment"),
		[]component.TableRow{})

	for k := range ws {
		v := &ws[k]
		tr, err := toWorkspaceTableRow(v, k, config)
		if err != nil {
			log.Logger().Infof("failed to create Table Row: %s", err.Error())
			continue
		}
		if tr != nil {
			table.Add(*tr)
		}
	}

	table.Sort("Sort")

	viewhelpers.InitTableFilters(config.TableFilters())

	table.AddFilter("Team", config.TeamFilter)
	table.AddFilter("Environment", config.EnvironmentFilter)

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: table},
	})

	return flexLayout, nil
}

func toWorkspaceTableRow(r *workspaces.WorkspaceOctant, idx int, config *WorkspacesViewConfig) (*component.TableRow, error) {
	w := r.Workspace
	team := r.Team
	env := r.Environment
	viewhelpers.AddFilterValue(&config.TeamFilter, team)
	viewhelpers.AddFilterValue(&config.EnvironmentFilter, env)
	return &component.TableRow{
		"Name":        ToWorkspaceName(r),
		"Team":        component.NewText(team),
		"Environment": component.NewText(env),
		"Source":      ToWorkspaceSource(&w),
		"Sort":        ToWorkspaceSort(&w, idx),
	}, nil
}

func ToWorkspaceName(r *workspaces.WorkspaceOctant) component.Component {
	u := r.URL
	if u == "" {
		return component.NewText(r.Name)
	}
	if r.Default {
		return viewhelpers.NewMarkdownText(viewhelpers.ToMarkdownLink(r.Name, u))
	}
	return viewhelpers.NewMarkdownText(viewhelpers.ToMarkdownExternalLink(r.Name, r.Name, u))
}

func ToWorkspaceSort(r *workspaces.Workspace, idx int) component.Component {
	if sortByIndex {
		return component.NewText(fmt.Sprintf("%019d", idx))
	}
	names := []string{r.Team, r.Name, r.Environment}
	return component.NewText(strings.Join(names, "-"))
}

func ToWorkspaceSource(r *workspaces.Workspace) component.Component {
	return viewhelpers.NewMarkdownText(viewhelpers.ToGitLinkMarkdown(r.GitURL))
}
