# Dark Theme Setup

TestMesh now uses **dark theme by default** for a modern, eye-friendly experience.

## What Changed

### 1. Theme Provider Added
- Installed `next-themes` package
- Created `/lib/providers/theme-provider.tsx`
- Wrapped app with ThemeProvider in `app/layout.tsx`

### 2. Default Configuration
```tsx
<ThemeProvider
  attribute="class"
  defaultTheme="dark"        // ← Default is dark
  enableSystem={false}       // ← Don't use system preference
  disableTransitionOnChange  // ← Smooth transitions
>
```

### 3. CSS Variables
Dark theme colors are defined in `app/globals.css`:
- Background: `oklch(0.145 0 0)` - Very dark gray
- Foreground: `oklch(0.985 0 0)` - Almost white
- Primary: `oklch(0.68 0.15 237)` - Blue accent
- Cards: `oklch(0.205 0 0)` - Slightly lighter than background

## Optional: Theme Toggle

A theme toggle component is available at `components/theme-toggle.tsx` if you want to let users switch themes.

### To Add Theme Toggle to Navigation:

```tsx
import { ThemeToggle } from '@/components/theme-toggle';

// Add to your navigation/header component:
<ThemeToggle />
```

This creates a button that switches between Light, Dark, and System themes.

## Enabling System Theme

If you want to respect the user's system theme preference instead:

```tsx
// In app/layout.tsx, change:
<ThemeProvider
  defaultTheme="system"  // ← Use system preference
  enableSystem={true}    // ← Enable system detection
>
```

## Customizing Colors

To customize dark theme colors, edit `app/globals.css` under the `.dark` selector:

```css
.dark {
  --background: oklch(0.145 0 0);     /* Main background */
  --foreground: oklch(0.985 0 0);     /* Text color */
  --primary: oklch(0.68 0.15 237);    /* Primary buttons */
  --card: oklch(0.205 0 0);           /* Card backgrounds */
  /* ... etc */
}
```

## Testing

Start the dev server and all pages will use dark theme:

```bash
pnpm dev
```

Visit any page - it will automatically use dark theme with no flashing or transitions.
