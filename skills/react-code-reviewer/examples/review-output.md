# Critical

## 1. dangerouslySetInnerHTML with unsanitized server-rendered HTML causes stored XSS

Location:
`UserTable`

Problem:
The component passes `row.nameHtml` straight into `dangerouslySetInnerHTML`. The field comes from `/api/users` and the API has no sanitizer or trusted-source check before returning it.

Impact:
Anyone who can control that field (compromised admin tool, SQL-injection downstream, hostile tenant) executes script in every visitor's session, exfiltrates tokens, or pivots into authenticated actions. This is a session-hijack / account-takeover vector.

Suggestion:
Do not render HTML for this field. If business requires rich text, sanitize server-side with an allowlist sanitizer (DOMPurify on the server, or a server-side equivalent) and restrict the allowed tag/attribute set.

Recommended code:
```tsx
// Option 1: drop the HTML interpretation entirely
<td>{row.name}</td>

// Option 2: render only after server-side sanitization
import DOMPurify from 'isomorphic-dompurify';
<td dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(row.nameHtml) }} />
```

# High

## 1. useEffect with empty dependency array ignores query changes

Location:
`UserTable#useEffect`

Problem:
The effect references `query` but its dependency array is `[]`. After the initial mount, every later change to `query` is silently ignored.

Impact:
The visible user list is permanently bound to the first query, so the user is looking at a table that does not match their filters and may act on the wrong rows. This is a data-correctness bug, not a style issue.

Suggestion:
Add `query` to the dependency array. Use `AbortController` (or a request-version counter) so a slow earlier response cannot overwrite a newer one when the user types fast.

Recommended code:
```tsx
useEffect(() => {
  const ctrl = new AbortController();
  let cancelled = false;
  (async () => {
    const res = await fetch(`/api/users?${qs.stringify(query)}`, {
      signal: ctrl.signal,
    });
    const data = await res.json();
    if (!cancelled) setRows(data);
  })();
  return () => { cancelled = true; ctrl.abort(); };
}, [query]);
```

# Medium

## 1. Optimistic update not rolled back when the save request fails

Location:
`UserTable#onSave`

Problem:
The handler calls `setRows(prev => prev.map(r => r.id === id ? updated : r))` before awaiting `saveUser(updated)`. If `saveUser` rejects, the local state still shows the new value.

Impact:
The user sees a successful save in the UI but the server never accepted it. A subsequent reload, another tab, or any action that depends on the persisted state will silently misbehave.

Suggestion:
Keep a snapshot of the previous row, apply the optimistic update, and roll back from the snapshot in the `catch` branch.

Recommended code:
```tsx
const onSave = async (updated: User) => {
  const previous = rows.find(r => r.id === updated.id);
  setRows(prev => prev.map(r => r.id === updated.id ? updated : r));
  try {
    await saveUser(updated);
  } catch (err) {
    setRows(prev => prev.map(r => r.id === updated.id ? previous! : r));
    toast.error('Save failed; change reverted');
  }
};
```

# Low

## 1. Inline `style` object recreated on every render

Location:
`UserTable#row`

Problem:
The `row` function returns `<tr style={{ padding: 8, color: '#333' }}>`. The object is a new reference on every render, so `React.memo`d children of `<tr>` see a "changed" `style` prop on every parent update.

Impact:
Pure perf concern — minor re-render cost, no functional bug — but it shows up in any list-render profile and is easy to fix.

Suggestion:
Hoist the constant to module scope, or use a CSS class for the visual style.

Recommended code:
```tsx
const ROW_STYLE = { padding: 8, color: '#333' } as const;

// ...
<tr style={ROW_STYLE}>
```
