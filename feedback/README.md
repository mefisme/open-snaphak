# feedback/ — the feedback relay

The in-app **Send feedback** dialog (the "?" button in SnapHak Studio) POSTs the user's report here —
a tiny Cloudflare Worker — and the relay files it as a labeled issue on this repo's tracker. Users never
need a GitHub account; the relay holds the only credential, and files issues **as the org's GitHub App
bot identity** (`<app-name>[bot]`), not as any personal account. See `worker.js` for the full flow
(validation → honeypot → dedup-by-signature → create or append) and `docs/feedback.md` for the
end-to-end pipeline including the repo-side hygiene workflows.

## Deploy (maintainer, one-time + on change)

Prereqs: a free Cloudflare account, Node.js.

```
cd feedback
npx wrangler login          # opens the browser, authorizes wrangler against your Cloudflare account
npx wrangler deploy         # prints the deployed URL: https://snaphak-feedback.<your-subdomain>.workers.dev
```

Sanity check: `curl https://snaphak-feedback.<your-subdomain>.workers.dev/` → `snaphak feedback relay: OK`.

The deployed hostname must match `kReportHost` in `src/ui/webview/snaphak_ui_webview.cpp` — update it
there once after the first deploy (the URL is stable across redeploys).

## The credential — a GitHub App (preferred)

Issues should come from a bot identity owned by the org, not a person, so the relay authenticates as a
**GitHub App**. One-time setup:

1. github.com → the **snaphak** org → Settings → Developer settings → **GitHub Apps** → New GitHub App.
   Name it (e.g. `snaphak-feedback`; issues will show as `snaphak-feedback[bot]`), homepage = this
   repo's URL, **uncheck "Active" under Webhook**, and under Permissions set **Repository permissions →
   Issues → Read and write** (nothing else). "Where can this app be installed" → **Only on this account**.
2. On the new app's page: note the **App ID**, then **Generate a private key** (downloads a `.pem`).
3. **Install App** (left sidebar) → install on the snaphak org → **Only select repositories →
   open-snaphak**.
4. GitHub's key download is PKCS#1; the Worker needs PKCS#8. Convert once:
   `openssl pkcs8 -topk8 -inform PEM -nocrypt -in <downloaded>.pem -out app-pkcs8.pem`
   (Windows: `openssl.exe` ships with Git under `C:\Program Files\Git\usr\bin\`.)
5. Store the secrets:
   `npx wrangler secret put APP_ID` (the number from step 2), then
   `npx wrangler secret put APP_PRIVATE_KEY` (paste the whole `app-pkcs8.pem`, BEGIN/END lines included).

No expiry, no rotation chore: the relay mints short-lived (~1 h) installation tokens from the key per
request. If the key is ever compromised, revoke it on the app page, generate a new one, redo steps 4–5.
Delete the local `.pem` files once the secret is stored.

**Fallback:** with no APP_ID/APP_PRIVATE_KEY set, the relay uses a fine-grained PAT from the
`GITHUB_TOKEN` secret instead (repo: only this one; permissions: Issues read/write; max 1-year expiry —
issues then come from the PAT owner's account, and the annual rotation chore applies).

## Abuse posture

Stateless and deliberately minimal: honeypot field + size/length caps in the Worker. If real spam ever
appears, add a Cloudflare **rate-limiting rule** on the dashboard (Security → WAF → Rate limiting rules,
scope it to `POST /report`) — no code change needed. Worst case is spam *issues*, which are deletable;
the token can't touch anything but Issues on this one repo.
