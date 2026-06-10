---
name: react-code-reviewer
description: Review React, TypeScript, Next.js, and Vite frontend changes for production-risk bugs in hooks, state, async effects, routing, authorization UI, SSR/CSR boundaries, hydration, security, accessibility, performance, forms, validation, testing, error boundaries, bundle size, state management, and maintainability — focusing on user-data-leak, XSS, authz-bypass, and incident-causing defects. Use for useEffect/cleanup correctness, Server/Client boundary mistakes, hydration mismatches, XSS in dangerouslySetInnerHTML, and incident-blocking PR review. Do NOT use for: Tailwind/CSS aesthetics, accessibility audits unrelated to a diff, or non-incident review.
metadata:
  short-description: Evidence-driven React frontend review
---

# React Code Reviewer

Review React frontend code for user-visible and production-risk defects. Load detailed prompts only when the code needs them.

## Required Loading

Always load `prompts/reviewer.md`.

Load additional prompts only when relevant:

- `prompts/nextjs-reviewer.md`: Next.js App Router, Server/Client Components, server actions, route handlers, cache, hydration.
- `prompts/security-reviewer.md`: XSS, unsafe HTML/URLs, token storage, open redirects, sensitive data exposure.
- `prompts/performance-reviewer.md`: unnecessary renders, memo misuse, bundle growth, request waterfalls, large lists.
- `prompts/testing-reviewer.md`: Jest, Vitest, React Testing Library, Playwright, Cypress; behaviour-driven assertions, `act()` wrapping, mock scoping.
- `prompts/forms-reviewer.md`: react-hook-form, Formik, Zod, server actions; submit disable, double-submit, schema-driven error mapping.
- `prompts/state-reviewer.md`: Redux Toolkit, Zustand, Jotai, Context; selector memoisation, store splitting, SSR hydration.
- `prompts/a11y-reviewer.md`: keyboard, screen reader, ARIA, focus management, contrast, live regions, reduced motion.
- `prompts/error-boundary-reviewer.md`: top-level vs route-level boundaries, event-handler async errors, fallback UX, log forwarding.
- `prompts/bundle-reviewer.md`: tree-shaking, dynamic imports, route-level code splitting, icon imports, barrel files.

## Review Contract

- Output concrete issues only.
- Bind each finding to code evidence and a user/runtime impact.
- Mark uncertain findings as `需要结合上下文确认`.
- If no clear high-risk issue is found, output exactly:

```text
未发现明确高风险问题。
```

Use the severity and output contract from `prompts/reviewer.md`.

## Examples

Each bad example has a matching `good-<file>` in this same `examples/`
directory that shows the minimal fix for every Critical/High finding. Read
both side by side when triaging a real diff.
