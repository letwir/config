---
name: sites-building
description: Use Sites to build websites, including landing pages, portfolios, dashboards, portals, trackers, hubs, and internal tools. Always use Sites when the project contains `.openai/hosting.json`.
---

# Sites building

Build the complete requested site, validate it, then use `sites-hosting`
unless the user explicitly asks to keep it local.

## Site lifecycle ownership

The Site-owning agent is the agent responsible for the user's requested Site, including an independently started background or invisible task. Only the Site-owning agent may initialize or edit the Site checkout, update `.openai/hosting.json`, run `sites-building` or `sites-hosting`, call `create_site` or any other Sites tool, obtain source credentials, save versions, deploy, or perform the browser handoff. A spawned subagent is never the owner of its parent agent's Site: give it only its explicitly assigned image, asset, or research task, and require it to return its result without invoking either Sites skill, calling Sites tools, initializing a Site, editing the Site checkout, or spawning another agent.

## Communicate clearly

Assume the user is a nontechnical knowledge worker. Talk about their site,
choices, progress, and results. Keep tools, commands, files, runtimes, browser
software, permissions, dependencies, source control, credentials, IDs, builds,
and deployment internals out of user-facing messages unless the user asks or
must take action.

Use no more than one short update for each user-visible phase: preparing the
site, building it, and publishing. If a phase takes longer
than 60 seconds, give one plain-language update. Keep recoverable technical
problems private; say only that you hit a problem and are trying another method.

Ask one concise group of up to three discovery questions only when important
context is missing and the unresolved details would materially affect the
site's functionality or force a risky assumption. Otherwise proceed immediately
with best judgment. Do not generate design options or pause for a visual
selection unless the user explicitly asks to compare designs.

## Choose the execution path

Use the **one-shot fast path** only when all of these are true:

- this is a new site in an empty or projectless workspace;
- one route can satisfy the request;
- the request does not require D1, R2, uploads, app-owned authentication,
  external connectors, or browser UI QA; and
- the normal deliverable is a private deployed URL.

Use the **capability path** otherwise. This includes existing-site changes,
multi-route sites, persistent data, uploads, authentication, external data, and
requested browser testing.

## Use imagery purposefully

Avoid model-authored SVGs in finished sites, including inline SVG
illustrations. Prefer strong typography, color, layout, CSS shapes, and existing
icon components when imagery is unnecessary. When a site needs real imagery,
prefer suitable images found through web image search. Use `imagegen` if and
only if original imagery is important and a suitable existing image is
unavailable; generation adds latency, so keep it purposeful and limited.

## Start new projects immediately

For a new site in an empty or projectless workspace, make setup the first task action. Scaffold directly with `@openai/create-sites@0.2.0` through the environment's package manager; the npm form is `npm create --yes @openai/sites@0.2.0 . -- --yes`. Respect the environment's existing dependency-security and minimum-release-age policy. If it blocks the pinned release, report the blocker instead of changing versions or bypassing the policy. Do not force a different package manager. Use the current directory when it is a valid destination; otherwise choose an empty project directory without moving, deleting, or overwriting existing workspace files.

Infer the capabilities required by the user's request and choose the appropriate add-ons yourself. Use `--add-ons` for capabilities such as `d1`, `r2`, and `auth`; omit it when none are needed. Consult the CLI's `--help` or `--list-add-ons --json` when needed instead of maintaining a separate template or capability catalog. Do not ask users to choose technical add-ons.

Install the generated project's dependencies with its package manager, either with the CLI's `--install` option or as a follow-up command, and retain the session until installation completes. As soon as the generated files exist, inspect the minimum required files and begin the bounded first product slice below while installation continues. Do not run a second initializer over the project.

In a visible foreground thread, start the project's development script in a retained session as soon as setup finishes, but keep the browser closed until the **First meaningful preview** gate below passes. Any generated loading skeleton is a fail-safe only and must never be the intended browser handoff. Keep the development server alive through build and hosting.

A Site-owning agent running in an independently started background, delegated, or invisible task initializes normally but does not start a browser-only preview unless its task otherwise needs the server. A spawned subagent working for that Site-owning agent never initializes a Site checkout.

## First meaningful preview

Treat the local preview as an early milestone in both execution paths. In a
visible foreground thread, open it as soon as, but not before, all of these are
true:

- the route contains the smallest coherent slice that lets a reasonable person
  recognize the requested site and its intended visual direction;
- it includes the primary product surface or layout and representative,
  product-specific content rather than an untouched starter, generic skeleton,
  blank page, or loading-only state;
- the primary affordance is visible when the requested experience is
  interaction-led; and
- the development server has successfully compiled the slice and the route
  responds without a blocking runtime error.

The slice may be static or partially inert. Keep it intentionally bounded and
defer secondary routes, complete data models, exhaustive interactions,
responsive refinements, animation, polish, and advanced capabilities until
after the handoff unless one is required for recognition, security, or a
successful render. Work on the slice while installation runs when possible.

For a new site, replace the generated placeholder content and temporary preview metadata, if present, as part of the slice. Cleanup of unused starter-only files and dependencies may happen after the handoff, but must finish before final validation.

Once the bounded slice is applied, make no further planned product-source edits
before the handoff. Fix only compilation or blocking runtime failures, then make
one lightweight non-browser request to the exact Local URL printed by the
development server to force the current route to render. Require a successful
compile and non-error response, but do not inspect the response body as visual
QA. Then use `open_in_codex` to show the first working version without waiting for the complete Site or a deployment.
Establish a stable browser-tab ID from the first preview and reuse it as the Site's single continuous user-facing view through HMR, publishing, and any later fixes.

For an existing site, use its current coherent experience immediately when it
still represents the requested product and compiles. If the request changes
the primary direction, apply only the smallest representative part first.
Preserve the last working content while changes compile; never replace an
existing site with the starter skeleton.

## Social previews

Apply to every site in both the **One-shot build** and **Capability path**.

Preserve an existing social-preview image and its metadata unless the user explicitly requests a refresh or the requested changes alter the site's branding. Generate a new card only when the existing preview is missing or replacement is warranted; otherwise do not start an image-generation sub-agent.

1. **Never delay the first meaningful preview.** If a new card is needed, once the site's direction and copy are stable, spawn exactly one image-generation subagent with `fork_turns="none"` and only the image brief while implementation continues. Explicitly instruct the subagent to make one `imagegen` request, save its result outside the Site checkout, and return its image path to the Site-owning agent; its assignment must state that it must not call Sites tools, invoke either Sites skill, edit site source, initialize another Site, or spawn another agent. Request one cohesive branded landscape card showing the site's exact title or primary headline and concise supporting copy as legible typography. Match its palette, typography, and distinctive visual motifs; exclude credentials and private user data. Inspect for incorrect, missing, or invented text; retry once only if the card is unusable. The Site-owning agent remains the sole site-source editor.
2. **Wire the site-wide preview.** When a new card is generated, the Site-owning agent saves it as `public/og.png` and sets site-specific Open Graph and X title, description, and image metadata in `app/layout.tsx` or the framework's equivalent metadata entry point. Use an absolute URL from a trusted request or deployment origin; never blindly trust forwarded host headers. If generation fails, preserve any existing valid preview; omit `og:image` only when neither an existing nor generated image is available. Never use a generic fallback. Run the final build after wiring the asset.
3. **Preserve item-specific previews.** For independently shareable detail pages, use `generateMetadata` or its equivalent to set page title and description and Open Graph/X title, description, and image from the rendered record. Reuse its existing primary image with an absolute trusted-origin URL; otherwise clear both inherited Open Graph and X images. Never reuse `public/og.png` or generate images per record. Validate the root and every detail page when there are at most two; otherwise check at least two representative detail pages. Before final validation, verify that each checked page's title, description, and Open Graph/X fields match its record.

## One-shot build

After setup and any necessary clarification, show the first meaningful preview,
then build and deploy the complete site in one focused pass.

1. Start by inspecting the generated project's instructions, package scripts, primary page, layout, stylesheet, and `.openai/hosting.json`. For the Vinext scaffold, these include `app/page.tsx`, `app/layout.tsx`, and `app/globals.css`. Read other files only when the implementation needs them. Avoid broad scans and speculative research. Preserve the package manager and lockfile.
2. Apply the smallest coherent product slice and complete the **First meaningful
   preview** handoff above before broadening the implementation. For a genuinely
   trivial request, the complete implementation may itself be that slice; do
   not manufacture extra edits merely to demonstrate HMR.
   After that first preview, if a new social card is needed, start the image-only subagent described in **Social previews** while continuing the implementation.
3. Reuse the retained setup, development server, and browser tab, then make one complete product patch. Prefer one page component and one stylesheet. Include all requested content, interactions, responsive behavior, keyboard and touch behavior when relevant, and accessible labels. Replace the generated placeholder content and metadata with the requested site's own values, remove unused starter-only files and dependencies, and update starter icons when appropriate before the final build unless the user explicitly asked to work on the starter itself. Integrate the **Social previews** result before the final build.
4. As soon as implementation is complete, run the project's build script while the retained development server stays alive. Fix actual build failures, then rerun it. Run lint separately only if the build omits compilation or the user asks.
5. Follow the shared preview rules below.
6. Continue to `sites-hosting`. Avoid an unnecessary polish pass after the
   build succeeds.

## Capability path

### Project setup

- For a new site, use the setup flow in **Start new projects immediately** and preserve the generated project's structure.
- For an existing site, preserve its package manager, lockfile, scripts,
  architecture, and `.openai/hosting.json`. Install only when dependencies are
  absent. Do not replace a working structure merely to use the starter.
- Keep site code within the selected project surface.

### Shape the product

- Before comprehensive implementation, apply the bounded slice and complete the
  **First meaningful preview** handoff above.
- After that first preview, if a new social card is needed, start the image-only subagent described in **Social previews** while continuing implementation; integrate its result before final validation.
- Build the first viewport around the requested product, not generic dashboard
  chrome.
- For a new site, replace the generated placeholder content and metadata, remove unused starter-only files and dependencies, and refresh the lockfile when dependencies change. Set the finished site's title and description through its framework's metadata API before final validation. Preserve starter content only when the user explicitly asked to work on the starter itself.
- Use concrete, product-specific copy and realistic data.
- Avoid speculative features and unnecessary client state.
- Use `sites()` from `@openai/sites-vite-plugin` and produce Cloudflare Worker-compatible ESM output.

### Add only requested capabilities

- For durable state, records, uploads, or other persistence, read
  [Persistence and storage](references/persistence-and-storage.md).
- For any SQLite schema or query work, also read [SQLite](references/sqlite.md).
- For identity-aware or sign-in-gated behavior, read
  [Authentication](references/authentication.md).
- Hosted Sites do not support raw TCP sockets (`connect()`); use HTTP-based clients or APIs for external databases and services.
- Use browser storage only for device-local preferences or explicitly local
  state.
- Keep logical D1 and R2 declarations in `.openai/hosting.json`; Sites owns the
  real Cloudflare resources and deployment wiring.
- Keep local `.env` and `.env.example` keys aligned. Manage hosted runtime
  values through Sites.

### Validate capability work

- Run the deployment build once after the complete implementation. If a D1
  schema changed, generate and inspect its migration. Fix real failures before
  hosting.

## Preview rules

- In a visible foreground thread, the **First meaningful preview** gate is the
  only local opening point. If the gate has not passed, keep the browser closed;
  never open the skeleton as a fallback. If `open_in_codex` fails after the gate
  passes, report it and continue.
- For an existing site, preserve its normal package and development flow.
- In a delegated, background, or invisible thread, skip `open_in_codex` and say
  why.
- Perform no screenshots, DOM inspection, clicking, resizing, or visual QA
  unless the user explicitly requests browser testing.
- Do not scan ports or repeatedly open the browser.

## Hosting handoff

Use `sites-hosting` after validation. Do not finish with only a local build
unless the user requested local-only work. Return the deployed Sites URL as the
primary deliverable. Do not include file paths, commands, or validation jargon
unless the user asks. Keep the development server running until hosting
finishes, then stop it during final teardown.
