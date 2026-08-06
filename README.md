# universal-rbac-manager

cmd/urm/main.go                    → manual wiring, same shape as their startAllServices()
internal/
  config/                          → env config struct (their configs/)
  server/
    router.go                     → gin.Engine setup, New(), Start()  (their http.go)
    routes.go                      → registerAuthRoutes(), registerUserRoutes(), registerRoleRoutes(), registerPermissionRoutes(), registerServiceRoutes(), registerAgentRoutes() — each takes *gin.RouterGroup + *Handler
    handlers.go                    → Handler struct (holds all service deps) + respondJSON/respondError helpers
    handler_auth.go               → login, refresh, introspect
    handler_user.go                → CRUD
    handler_role.go                → CRUD + grants
    handler_permission.go     → catalog browse
    handler_service.go          → self-registration, credential rotation
    middleware.go                → authMiddleware (calls token package), request logging config
  auth/          → auth.go (login, password, lockout) + store.go
  user/          → user.go + store.go
  role/           → role.go + store.go
  permission/  → permission.go + store.go
  service/       → service-principal registration/rotation + store.go
  token/           → JWT mint/verify + refresh token rotation (pure, no gin/pgx coupling beyond store)
  migrate/        → embed.FS + runner
  model/          → shared structs
schemas/            → one .sql per table (you already have the doc/setup.md ordering convention)