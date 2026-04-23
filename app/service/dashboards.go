package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pgquerynarrative/pgquerynarrative/api/gen/dashboards"
	queriesapi "github.com/pgquerynarrative/pgquerynarrative/api/gen/queries"
	reportsapi "github.com/pgquerynarrative/pgquerynarrative/api/gen/reports"
)

type DashboardsService struct {
	appPool    *pgxpool.Pool
	reportsSvc *ReportsService
	queriesSvc *QueriesService
}

func NewDashboardsService(appPool *pgxpool.Pool, reportsSvc *ReportsService, queriesSvc *QueriesService) *DashboardsService {
	return &DashboardsService{appPool: appPool, reportsSvc: reportsSvc, queriesSvc: queriesSvc}
}

func (s *DashboardsService) List(ctx context.Context) (*dashboards.DashboardListResult, error) {
	rows, err := s.appPool.Query(ctx, `SELECT id FROM app.dashboards ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := &dashboards.DashboardListResult{Items: []*dashboards.Dashboard{}}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		item, err := s.Get(ctx, &dashboards.GetPayload{ID: id})
		if err != nil {
			return nil, err
		}
		out.Items = append(out.Items, item)
	}
	return out, rows.Err()
}

func (s *DashboardsService) Create(ctx context.Context, payload *dashboards.CreatePayload) (*dashboards.Dashboard, error) {
	var id string
	if err := s.appPool.QueryRow(ctx, `INSERT INTO app.dashboards (name) VALUES ($1) RETURNING id`, payload.Name).Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(ctx, &dashboards.GetPayload{ID: id})
}

func (s *DashboardsService) Get(ctx context.Context, payload *dashboards.GetPayload) (*dashboards.Dashboard, error) {
	row := s.appPool.QueryRow(ctx, `SELECT id, name, created_at, updated_at FROM app.dashboards WHERE id = $1`, payload.ID)
	var d dashboards.Dashboard
	var createdAt, updatedAt time.Time
	if err := row.Scan(&d.ID, &d.Name, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &dashboards.NotFoundError{Name: "not_found", Message: "dashboard not found", Code: strPtr("NOT_FOUND")}
		}
		return nil, err
	}
	d.CreatedAt = createdAt.Format(time.RFC3339)
	d.UpdatedAt = updatedAt.Format(time.RFC3339)
	widgets, err := s.getWidgets(ctx, d.ID)
	if err != nil {
		return nil, err
	}
	d.Widgets = widgets
	return &d, nil
}

func (s *DashboardsService) Update(ctx context.Context, payload *dashboards.UpdatePayload) (*dashboards.Dashboard, error) {
	tag, err := s.appPool.Exec(ctx, `UPDATE app.dashboards SET name = $2, updated_at = NOW() WHERE id = $1`, payload.ID, payload.Name)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, &dashboards.NotFoundError{Name: "not_found", Message: "dashboard not found", Code: strPtr("NOT_FOUND")}
	}
	if payload.Widgets != nil {
		if _, err := s.appPool.Exec(ctx, `DELETE FROM app.dashboard_widgets WHERE dashboard_id = $1`, payload.ID); err != nil {
			return nil, err
		}
		for i, w := range payload.Widgets {
			if w == nil {
				continue
			}
			refresh := int32(300)
			if w.RefreshSeconds != nil && *w.RefreshSeconds > 0 {
				refresh = *w.RefreshSeconds
			}
			position := int32(i)
			if w.Position != nil {
				position = *w.Position
			}
			_, err := s.appPool.Exec(ctx, `
				INSERT INTO app.dashboard_widgets (dashboard_id, widget_type, title, report_id, saved_query_id, refresh_seconds, position)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, payload.ID, w.WidgetType, w.Title, w.ReportID, w.SavedQueryID, refresh, position)
			if err != nil {
				return nil, err
			}
		}
	}
	return s.Get(ctx, &dashboards.GetPayload{ID: payload.ID})
}

func (s *DashboardsService) Delete(ctx context.Context, payload *dashboards.DeletePayload) error {
	_, err := s.appPool.Exec(ctx, `DELETE FROM app.dashboards WHERE id = $1`, payload.ID)
	return err
}

func (s *DashboardsService) Resolve(ctx context.Context, payload *dashboards.ResolvePayload) (*dashboards.DashboardResolved, error) {
	d, err := s.Get(ctx, &dashboards.GetPayload{ID: payload.ID})
	if err != nil {
		return nil, err
	}
	res := &dashboards.DashboardResolved{ID: d.ID, Name: d.Name, Widgets: []*dashboards.DashboardResolvedWidget{}}
	for _, w := range d.Widgets {
		if w == nil {
			continue
		}
		item := &dashboards.DashboardResolvedWidget{
			ID:             w.ID,
			WidgetType:     w.WidgetType,
			Title:          w.Title,
			RefreshSeconds: w.RefreshSeconds,
			Position:       w.Position,
		}
		if w.ReportID != nil {
			r, err := s.reportsSvc.Get(ctx, &reportsapi.GetPayload{ID: *w.ReportID})
			if err == nil {
				item.Report = r
			}
		}
		if w.SavedQueryID != nil {
			q, err := s.queriesSvc.GetSaved(ctx, &queriesapi.GetSavedPayload{ID: *w.SavedQueryID})
			if err == nil {
				item.SavedQuery = q
			}
		}
		res.Widgets = append(res.Widgets, item)
	}
	return res, nil
}

func (s *DashboardsService) getWidgets(ctx context.Context, dashboardID string) ([]*dashboards.DashboardWidget, error) {
	rows, err := s.appPool.Query(ctx, `
		SELECT id, widget_type, title, report_id, saved_query_id, refresh_seconds, position
		FROM app.dashboard_widgets
		WHERE dashboard_id = $1
		ORDER BY position ASC, created_at ASC
	`, dashboardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*dashboards.DashboardWidget{}
	for rows.Next() {
		var w dashboards.DashboardWidget
		var title sql.NullString
		var reportID sql.NullString
		var savedID sql.NullString
		if err := rows.Scan(&w.ID, &w.WidgetType, &title, &reportID, &savedID, &w.RefreshSeconds, &w.Position); err != nil {
			return nil, err
		}
		if title.Valid {
			w.Title = &title.String
		}
		if reportID.Valid {
			w.ReportID = &reportID.String
		}
		if savedID.Valid {
			w.SavedQueryID = &savedID.String
		}
		out = append(out, &w)
	}
	return out, rows.Err()
}
