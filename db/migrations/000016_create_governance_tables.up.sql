-- Organizations, Sites, Roles, Role assignments (simple RBAC scaffolding)

CREATE TABLE IF NOT EXISTS organizations (
  id SERIAL PRIMARY KEY,
  uid UUID NOT NULL DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  code TEXT,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_organizations_uid ON organizations(uid);

CREATE TABLE IF NOT EXISTS sites (
  id SERIAL PRIMARY KEY,
  uid UUID NOT NULL DEFAULT gen_random_uuid(),
  organization_id INTEGER REFERENCES organizations(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  code TEXT,
  location TEXT,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sites_uid ON sites(uid);
CREATE INDEX IF NOT EXISTS idx_sites_organization_id ON sites(organization_id);

CREATE TABLE IF NOT EXISTS roles (
  id SERIAL PRIMARY KEY,
  uid UUID NOT NULL DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_roles_uid ON roles(uid);

CREATE TABLE IF NOT EXISTS role_assignments (
  id SERIAL PRIMARY KEY,
  uid UUID NOT NULL DEFAULT gen_random_uuid(),
  role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
  user_uid TEXT NOT NULL,
  organization_id INTEGER REFERENCES organizations(id) ON DELETE SET NULL,
  site_id INTEGER REFERENCES sites(id) ON DELETE SET NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_role_assignments_uid ON role_assignments(uid);
CREATE INDEX IF NOT EXISTS idx_role_assignments_user_uid ON role_assignments(user_uid);
