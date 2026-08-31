Add opt-in Google Identity Services credential-flow compatibility
Problem
Applications using Google Identity Services (GIS) do not use an OAuth authorization-code redirect flow. They load https://accounts.google.com/gsi/client, render declarative g_id_onload and g_id_signin markup, receive a signed ID token in a JavaScript callback, and post that token to their backend.

The mock OAuth server provides OAuth/OIDC authorization, token, userinfo, discovery, and JWKS endpoints, but it cannot drive this browser credential-callback contract. As a result, such apps cannot exercise their production-style local ID-token validation flow without real Google credentials.

Proposed Change
Add an explicitly opt-in mock GIS compatibility feature. It must emulate only the narrow browser contract required for local development and tests; it is not a replacement for the Google GIS SDK or a production authentication provider.

Add GET /gsi/client, serving a small browser JavaScript compatibility script.
Have the script recognize an element with id="g_id_onload", read data-client_id and data-callback, and render an accessible sign-in control in .g_id_signin containers.
On activation, have the script request a mock credential from a new mock-server credential endpoint and invoke window[data-callback]({ credential: idToken }), matching the GIS callback shape consumed by such apps.
Add a credential endpoint, recommended as POST /gsi/credential, which accepts a client ID and returns JSON containing credential with an RS256 JWT.
Reuse the existing RSA key pair and kid used by /jwks; do not introduce a separate signing key or an unsigned development token.
Configure the token with the active mock user and issuer. It must contain iss, aud, exp, iat, sub, email, email_verified, name, and picture. aud must equal the client ID submitted by the script.
Register the routes only when an explicit configuration flag is enabled, recommended as MOCK_GIS_ENABLED=true, defaulting to false. Disabled routes should return 404.
Use configured MOCK_ISSUER_URL as the JWT issuer and discovery base URL. It must be reachable by the application verifying tokens and exactly match its expected issuer string.
Protect the credential endpoint from accidental cross-origin access: respond only to same-origin requests by default, or make allowed origins an explicit local-development configuration. Do not enable Access-Control-Allow-Origin: * by default.
Document the feature as development/test-only. Do not alter existing /authorize, /token, /userinfo, discovery, or JWKS behavior for standard users of the project.
Acceptance Criteria
With MOCK_GIS_ENABLED=true, GET /gsi/client succeeds with JavaScript content; when disabled, the route is not exposed.
Given markup with data-client_id="[appName]-local" and data-callback="handleCredentialResponse", a user action results in one callback invocation with a non-empty credential string.
The ID token validates using only the public JWK retrieved from the existing /jwks endpoint.
Validated claims have issuer equal to MOCK_ISSUER_URL, audience equal to the submitted client ID, a future expiry, email_verified=true, and values matching MOCK_USER_EMAIL and MOCK_USER_NAME.
A token issued for one client ID fails audience validation for a different client ID.
go test ./... continues to pass. New tests cover enabled, disabled, successful, malformed-request, and claim-validation behavior.
Existing authorization-code integration tests continue to pass unchanged.
Suggested Implementation Areas
cmd/server/main.go: parse MOCK_GIS_ENABLED and conditionally register routes.
internal/server/server.go: make the new route registration available to in-process test-server construction.
internal/handlers/: add GIS JavaScript and credential handlers.
internal/jwt/jwt.go: expose or reuse a narrow ID-token issuance helper if the current helper does not accept all required mock-user claims and issuer/client inputs.
README.md: add configuration, API contract, browser integration example, and security limitation.
Tests
Handler test for the script content/type and disabled-route behavior.
Credential handler tests for request validation and returned token envelope.
JWT/JWKS interoperability test: fetch /jwks, resolve kid, and verify the returned token's RS256 signature and critical claims.
Minimal browser-contract test, potentially using httptest plus static script assertions if the repository intentionally avoids a browser test runner.
Out of Scope
Google account selection, One Tap, FedCM, prompt moments, popup UX, Google branding, or full GIS SDK behavior.
Replacing or changing existing OAuth authorization-code endpoints.
Production use, persistence, user/password management, or broad CORS enablement.
Consumer Contract: PWA apps
apps local mode will load ${MOCK_OAUTH_ISSUER}/gsi/client, retain its existing credential POST payload to /auth/login, and verify the token against ${MOCK_OAUTH_ISSUER}/jwks with the configured issuer. It must work without outgoing Google requests and must use a verified signed ID token, not a local authentication bypass.