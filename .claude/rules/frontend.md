# Frontend Conventions

## Framework

- Svelte 5 with runes (`$state`, `$derived`, `$effect`, `$props`,
  `$bindable`). No legacy `export let` / `$:` syntax.
- TypeScript strict mode everywhere.
- SPA only: `ssr = false` stays in `+layout.ts`. The Go backend serves the
  static build; there is no Node server in production.
- **pnpm** - never npm or yarn.

## Data fetching

- Every API call goes through `src/lib/api/client.ts::api()`. Components
  never call `fetch("/api/...")` directly.
- One module per resource in `src/lib/api/<resource>.ts` with exported
  TypeScript interfaces matching the Go JSON shapes exactly (camelCase).
- When you add or change a Go response struct, update the matching TS
  interface in the same change.

## Components

- Design-system primitives live in `src/lib/components/ui/` (Button, Card,
  Input, Badge). Use and extend them - don't inline raw styled HTML for
  patterns a primitive covers. New primitive → same folder, same style.
- Feature components in `src/lib/components/<feature>/` once reused;
  page-specific markup stays in the page file.
- Colors ONLY via design tokens (`bg-card`, `text-muted-foreground`,
  `border-border`, `text-destructive`, …) defined in `src/app.css`. Never
  `text-gray-500` or hex values. Dark mode then works automatically.

## Errors & UX

- Every mutation handles failure and shows the message near where the user
  acted (inline error text), not just console.error.
- Empty states say what to do next, not just "nothing here".
- Loading/disabled states on buttons that trigger requests.

## Testing

- Vitest. Test files `<source>.test.ts` next to the source.
- Pure logic (utils, data transforms) must be tested; presentational
  markup doesn't need tests.

## Formatting

Prettier (`pnpm format`), tabs, double quotes. CI fails on unformatted code.
