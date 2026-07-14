package repo

import (
	"context"
	"fmt"

	"app/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GovernanceRepo struct {
	db *pgxpool.Pool
}

func NewGovernanceRepo(db *pgxpool.Pool) *GovernanceRepo {
	return &GovernanceRepo{db: db}
}

// Organization CRUD
func (r *GovernanceRepo) CreateOrganization(ctx context.Context, name, code, notes string) (*models.Organization, error) {
	row := r.db.QueryRow(ctx, `INSERT INTO organizations (name, code, notes) VALUES ($1,$2,$3) RETURNING uid, name, code, notes, created_at, updated_at`, name, code, notes)
	org := &models.Organization{}
	var uid pgtype.UUID
	var createdAt pgtype.Timestamptz
	var updatedAt pgtype.Timestamptz
	if err := row.Scan(&uid, &org.Name, &org.Code, &org.Notes, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	org.UID = uid.String()
	if createdAt.Valid {
		t := createdAt.Time
		org.CreatedAt = &t
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		org.UpdatedAt = &t
	}
	return org, nil
}

func (r *GovernanceRepo) ListOrganizations(ctx context.Context) ([]models.Organization, error) {
	rows, err := r.db.Query(ctx, `SELECT uid, name, code, notes, created_at, updated_at FROM organizations ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Organization
	for rows.Next() {
		var org models.Organization
		var uid pgtype.UUID
		var createdAt pgtype.Timestamptz
		var updatedAt pgtype.Timestamptz
		if err := rows.Scan(&uid, &org.Name, &org.Code, &org.Notes, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		org.UID = uid.String()
		if createdAt.Valid {
			t := createdAt.Time
			org.CreatedAt = &t
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			org.UpdatedAt = &t
		}
		out = append(out, org)
	}
	return out, nil
}

func (r *GovernanceRepo) GetOrganizationByUID(ctx context.Context, uid string) (*models.Organization, error) {
	row := r.db.QueryRow(ctx, `SELECT uid, name, code, notes, created_at, updated_at FROM organizations WHERE uid = $1`, uid)
	org := &models.Organization{}
	var pu pgtype.UUID
	var createdAt pgtype.Timestamptz
	var updatedAt pgtype.Timestamptz
	if err := row.Scan(&pu, &org.Name, &org.Code, &org.Notes, &createdAt, &updatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	org.UID = pu.String()
	if createdAt.Valid {
		t := createdAt.Time
		org.CreatedAt = &t
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		org.UpdatedAt = &t
	}
	return org, nil
}

// Sites
func (r *GovernanceRepo) CreateSite(ctx context.Context, organizationUID, name, code, location, notes string) (*models.Site, error) {
	// resolve organization id
	var orgID int
	if err := r.db.QueryRow(ctx, `SELECT id FROM organizations WHERE uid=$1`, organizationUID).Scan(&orgID); err != nil {
		return nil, fmt.Errorf("invalid organization: %w", err)
	}
	row := r.db.QueryRow(ctx, `INSERT INTO sites (organization_id, name, code, location, notes) VALUES ($1,$2,$3,$4,$5) RETURNING uid, name, code, location, notes, created_at, updated_at`, orgID, name, code, location, notes)
	s := &models.Site{}
	var uid pgtype.UUID
	var createdAt pgtype.Timestamptz
	var updatedAt pgtype.Timestamptz
	if err := row.Scan(&uid, &s.Name, &s.Code, &s.Location, &s.Notes, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	s.UID = uid.String()
	s.OrganizationUID = organizationUID
	if createdAt.Valid {
		t := createdAt.Time
		s.CreatedAt = &t
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		s.UpdatedAt = &t
	}
	return s, nil
}

func (r *GovernanceRepo) ListSitesByOrganization(ctx context.Context, organizationUID string) ([]models.Site, error) {
	rows, err := r.db.Query(ctx, `SELECT s.uid, s.name, s.code, s.location, s.notes, s.created_at, s.updated_at FROM sites s JOIN organizations o ON s.organization_id = o.id WHERE o.uid = $1 ORDER BY s.name`, organizationUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Site
	for rows.Next() {
		var s models.Site
		var uid pgtype.UUID
		var createdAt pgtype.Timestamptz
		var updatedAt pgtype.Timestamptz
		if err := rows.Scan(&uid, &s.Name, &s.Code, &s.Location, &s.Notes, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		s.UID = uid.String()
		if createdAt.Valid {
			t := createdAt.Time
			s.CreatedAt = &t
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			s.UpdatedAt = &t
		}
		out = append(out, s)
	}
	return out, nil
}

func (r *GovernanceRepo) GetSiteByUID(ctx context.Context, uid string) (*models.Site, error) {
	row := r.db.QueryRow(ctx, `SELECT s.uid, o.uid as org_uid, s.name, s.code, s.location, s.notes, s.created_at, s.updated_at FROM sites s JOIN organizations o ON s.organization_id = o.id WHERE s.uid = $1`, uid)
	var s models.Site
	var pu pgtype.UUID
	var orgUID pgtype.UUID
	var createdAt pgtype.Timestamptz
	var updatedAt pgtype.Timestamptz
	if err := row.Scan(&pu, &orgUID, &s.Name, &s.Code, &s.Location, &s.Notes, &createdAt, &updatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	s.UID = pu.String()
	s.OrganizationUID = orgUID.String()
	if createdAt.Valid {
		t := createdAt.Time
		s.CreatedAt = &t
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		s.UpdatedAt = &t
	}
	return &s, nil
}

// Roles and assignments
func (r *GovernanceRepo) CreateRole(ctx context.Context, name, notes string) (*models.Role, error) {
	row := r.db.QueryRow(ctx, `INSERT INTO roles (name, notes) VALUES ($1,$2) RETURNING uid, name, notes, created_at, updated_at`, name, notes)
	role := &models.Role{}
	var uid pgtype.UUID
	var createdAt pgtype.Timestamptz
	var updatedAt pgtype.Timestamptz
	if err := row.Scan(&uid, &role.Name, &role.Notes, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	role.UID = uid.String()
	if createdAt.Valid {
		t := createdAt.Time
		role.CreatedAt = &t
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		role.UpdatedAt = &t
	}
	return role, nil
}

func (r *GovernanceRepo) ListRoles(ctx context.Context) ([]models.Role, error) {
	rows, err := r.db.Query(ctx, `SELECT uid, name, notes, created_at, updated_at FROM roles ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Role
	for rows.Next() {
		var role models.Role
		var uid pgtype.UUID
		var createdAt pgtype.Timestamptz
		var updatedAt pgtype.Timestamptz
		if err := rows.Scan(&uid, &role.Name, &role.Notes, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		role.UID = uid.String()
		if createdAt.Valid {
			t := createdAt.Time
			role.CreatedAt = &t
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			role.UpdatedAt = &t
		}
		out = append(out, role)
	}
	return out, nil
}

func (r *GovernanceRepo) AssignRoleToUser(ctx context.Context, roleUID, userUID string, organizationUID, siteUID *string) (*models.RoleAssignment, error) {
	// resolve role id
	var roleID int
	if err := r.db.QueryRow(ctx, `SELECT id FROM roles WHERE uid=$1`, roleUID).Scan(&roleID); err != nil {
		return nil, fmt.Errorf("invalid role: %w", err)
	}
	// optionally resolve org/site ids
	var orgID *int
	var siteID *int
	if organizationUID != nil {
		var oid int
		if err := r.db.QueryRow(ctx, `SELECT id FROM organizations WHERE uid=$1`, *organizationUID).Scan(&oid); err != nil {
			return nil, fmt.Errorf("invalid organization: %w", err)
		}
		orgID = &oid
	}
	if siteUID != nil {
		var sid int
		if err := r.db.QueryRow(ctx, `SELECT id FROM sites WHERE uid=$1`, *siteUID).Scan(&sid); err != nil {
			return nil, fmt.Errorf("invalid site: %w", err)
		}
		siteID = &sid
	}

	row := r.db.QueryRow(ctx, `INSERT INTO role_assignments (role_id, user_uid, organization_id, site_id) VALUES ($1,$2,$3,$4) RETURNING uid, created_at`, roleID, userUID, orgID, siteID)
	var uid pgtype.UUID
	var createdAt pgtype.Timestamptz
	if err := row.Scan(&uid, &createdAt); err != nil {
		return nil, err
	}
	ra := &models.RoleAssignment{UID: uid.String(), RoleUID: roleUID, UserUID: userUID}
	if organizationUID != nil {
		ra.OrganizationUID = organizationUID
	}
	if siteUID != nil {
		ra.SiteUID = siteUID
	}
	if createdAt.Valid {
		t := createdAt.Time
		ra.CreatedAt = &t
	}
	return ra, nil
}

func (r *GovernanceRepo) ListRoleAssignmentsForUser(ctx context.Context, userUID string) ([]models.RoleAssignment, error) {
	rows, err := r.db.Query(ctx, `SELECT ra.uid, r.uid as role_uid, ra.user_uid, o.uid as org_uid, s.uid as site_uid, ra.created_at
	FROM role_assignments ra
	JOIN roles r ON ra.role_id = r.id
	LEFT JOIN organizations o ON ra.organization_id = o.id
	LEFT JOIN sites s ON ra.site_id = s.id
	WHERE ra.user_uid = $1`, userUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.RoleAssignment
	for rows.Next() {
		var ra models.RoleAssignment
		var uid pgtype.UUID
		var roleUID pgtype.UUID
		var createdAt pgtype.Timestamptz
		// note: org/site may be null — use pgtype.Null
		var orgUIDNull pgtype.UUID
		var siteUIDNull pgtype.UUID
		if err := rows.Scan(&uid, &roleUID, &ra.UserUID, &orgUIDNull, &siteUIDNull, &createdAt); err != nil {
			return nil, err
		}
		ra.UID = uid.String()
		ra.RoleUID = roleUID.String()
		if orgUIDNull.Valid {
			v := orgUIDNull.String()
			ra.OrganizationUID = &v
		}
		if siteUIDNull.Valid {
			v := siteUIDNull.String()
			ra.SiteUID = &v
		}
		if createdAt.Valid {
			t := createdAt.Time
			ra.CreatedAt = &t
		}
		out = append(out, ra)
	}
	return out, nil
}
