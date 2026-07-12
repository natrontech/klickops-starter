# Design System

## Tokens

All colors are CSS custom properties in `ui/src/app.css`, exposed to
Tailwind via `@theme inline`. Light and dark values are defined there;
dark mode follows the OS (`prefers-color-scheme`) with zero JS.

Use: `bg-background`, `text-foreground`, `bg-card`, `text-muted-foreground`,
`border-border`, `bg-primary`, `text-destructive`, `text-success`,
`text-warning`, `outline-ring`.

Never hardcode a color (`bg-white`, `text-gray-500`, hex). To restyle the
whole app, change the token values once.

The token names match shadcn-svelte, so shadcn components can be added
later without re-theming.

## Patterns

- **Cards**: `rounded-xl border border-border bg-card p-5` - or just use
  the `Card` component.
- **Page title**: `text-[28px] font-semibold leading-tight tracking-[-0.015em]`.
- **Body text**: `text-sm`, secondary text `text-sm text-muted-foreground`.
- **Small labels/timestamps**: `text-[11px] text-muted-foreground`.
- **Buttons**: use the `Button` component (`primary | outline | ghost |
  destructive`); height 8 (`h-8`), text-xs.
- **Monospace** for technical values: `font-mono text-xs`.
- Subtle borders use opacity: `border-border/60`, `divide-border/60`.

## Rules

- Extend the primitives in `ui/src/lib/components/ui/` instead of
  hand-rolling one-off styled elements.
- Spacing rhythm: `gap-6` between page sections, `gap-2`/`gap-3` inside.
- Keep interactive elements accessible: focus-visible outlines (already in
  the primitives), labels on icon-only buttons.
