package models

// RBAC: Role, Permission и связи. В домене роли/права — сущности, а присваивания — связи.

type Role struct {
	ID       RoleID
	TenantID TenantID
	Name     string // уникальна в рамках tenant
	Desc     string
}

type Permission struct {
	ID       PermID
	TenantID TenantID
	Action   string // e.g. "user.read"
	Object   string // e.g. "user", "project"
	Desc     string
}

// Связи (агрегаты в домене, таблицы-связки в БД)
type UserRole struct {
	UserID UserID
	RoleID RoleID
}

type RolePermission struct {
	RoleID RoleID
	PermID PermID
}
