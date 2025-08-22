# SSO Auth gRPC service

## Бизнес-процесс
**Регистрация**: Register → создаём user (argon2id hash), шлём verify-email токен (UUID в таблице email_verifications).

**Верификация почты**: VerifyEmail(token) → проставляем email_verified=true, логируем событие.

**Логин по паролю**: проверка hash → если mfa_enabled → вернуть «нужна MFA» (можно код 401 с detail) → VerifyMFA → создать session, выдать access(JWT, 5–15 мин) + refresh(rotating, 30–90 дней); refresh хранить как хэш.

**Refresh (rotation + reuse-detector)**: по присланному refresh ищем по хэшу; если уже consumed — триггерим «возможная кража» → ревок всей цепочки session_id. Иначе: пометить consumed_at, выдать новый refresh (chain через previous_token_hash).

**Logout**: ревокнуть текущий refresh (и/или всю session).

**Introspect**: проверить подпись и jti в blacklist, вернуть claims.

**Revoke**: добавить jti в access_jti_blacklist (или ревокнуть по session).

**MFA TOTP**: InitiateMFA → генерим secret, даём otpauth URL, VerifyMFA → сохранить secret и mfa_enabled=true.

**WebAuthn**: рег/аутентика через стандартные вызовы (рекомендация — библиотека github.com/go-webauthn/webauthn), хранить credential id/publicKey.

**Reset password**: StartPasswordReset(email) → одноразовый токен; CompletePasswordReset(token,new).

**Service-to-service (mTLS + JWT)**: для межсервиски лучше: mTLS на ingress + отдельный issuer для svc-токенов с audience конкретного сервиса; короткий TTL (1–5 мин), без refresh.

**RBAC/ABAC**: PermissionsService.Check(subject, action, object, ctx) → простая RBAC таблицами выше. Для ABAC можешь добавлять правила (например, разрешать user:read если subject==object.owner в контексте). Интерцептор на gRPC проверяет Check до бизнес-метода.

## Требования
• **Пароли**: argon2id, не хранить в plaintext, проверка через constant-time.

• Refresh хранить хэшем; rotation с reuse detection.

• **Access**: Ed25519 (EdDSA), короткий TTL, aud, jti, kid, tid (tenant), roles.

• JWKS endpoint + key-rotation, revoke по jti.

• Email verify и password reset через одноразовые токены с TTL.

• **MFA**: TOTP минимум 6 цифр/30 сек, дрейфт ±1 шаг; backup codes (список одноразовых).

• **WebAuthn**: хранить credential id/publicKey, проверка signCount.

• **RBAC/ABAC**: быстрый чек (кеш разрешений в памяти на 30–60 сек).

• Аудит-логирование всех security-событий.

• Rate limiting на Login/Refresh/VerifyMFA.

• **Защита от user-enum**: одинаковые ответы на «не найден/неверный пароль».

• **Мультитенантность**: tenant_id везде (user, roles, tokens).

• **Service-to-service**: отдельный issuer/audience, короткий TTL, без refresh.

• **Трассировка**: gRPC interceptors (OpenTelemetry).

• **Обсервабилити**: метрики по кодам/латентности, счётчики логинов/рефрешей/ошибок.

### Добавить позже
• PKCE/Authorization Code для OIDC-потока.

• **Device management**: список активных сессий, logout remote.

• Attribute-based policies через OPA (local decision point).

• **Внешние провайдеры (Google/GitHub)**: хранить oauth_identities, связывать с existing user, если email совпал и verified.

### Архитектура:
```
.
├── cmd
│   ├── migrator/main.go
│   └── sso/main.go
├── configs
│   ├── local.yaml
│   └── test.yaml
├── internal
│   ├── app
│   │   ├── app.go
│   │   └── grpc/app.go
│   ├── authn           # бизнес-логика аутентификации
│   │   ├── password.go
│   │   ├── jwt.go
│   │   ├── refresh.go
│   │   ├── mfa_totp.go
│   │   ├── webauthn.go
│   │   └── email.go
│   ├── authz           # проверки доступа (RBAC/ABAC)
│   │   ├── checker.go
│   │   └── interceptors.go
│   ├── config
│   │   └── config.go
│   ├── grpc
│   │   ├── auth/server.go
│   │   ├── identity/server.go
│   │   ├── permissions/server.go
│   │   └── admin/server.go
│   ├── jwk
│   │   ├── manager.go  # ротация ключей, хранение, JWKS
│   │   └── store.go
│   ├── mail
│   │   └── sender.go
│   ├── oauth
│   │   ├── provider.go     # вы как OIDC provider (/.well-known, /authorize, /token)
│   │   └── clients.go      # внешние IdP (Google/GitHub) как клиенты
│   ├── services
│   │   ├── auth/service.go
│   │   ├── identity/service.go
│   │   └── permissions/service.go
│   ├── storage
│   │   ├── psql
│   │   │   ├── users_repo.go
│   │   │   ├── sessions_repo.go
│   │   │   ├── tokens_repo.go
│   │   │   ├── jwk_repo.go
│   │   │   ├── permissions_repo.go
│   │   │   └── audits_repo.go
│   │   └── storage.go
│   └── lib
│       ├── logger
│       │   ├── handlers
│       │   │   ├── slogdiscard/slogdiscard.go
│       │   │   └── slogpretty/slogpretty.go
│       │   └── sl/sl.go
│       ├── rate/limiter.go
│       ├── tracing/tracing.go
│       └── validator/validator.go
├── migrations
│   ├── 0001_init.down.sql
│   ├── 0001_init.up.sql
│   ├── 0002_permissions.down.sql
│   ├── 0002_permissions.up.sql
│   ├── 0003_mfa.down.sql
│   ├── 0003_mfa.up.sql
│   ├── 0004_oauth.down.sql
│   ├── 0004_oauth.up.sql
│   ├── 0005_jwk.down.sql
│   └── 0005_jwk.up.sql
└── tests
    ├── migrations
    │   └── ...
    ├── e2e_auth_test.go
    └── suite/suite.go
```
