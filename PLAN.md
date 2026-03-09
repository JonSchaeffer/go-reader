# go-reader — Product Plan

## What it is

A self-hosted reading app combining an RSS reader with a personal reading library.
Think Readwise Reader, but yours.

---

## Current State (as of March 2026)

### Working
- RSS feed management (add, delete, categorise)
- Background feed fetching every 5 minutes via FiveFilters Full-Text RSS
- Article list with read/unread tracking, search, pagination
- Category management
- SvelteKit frontend: feed list, article sidebar panel, full-screen reading modal with prev/next navigation

### Known tech debt
- `SaveRSSArticles()` in `rss/rss.go` re-processes HTML for every article on every fetch cycle, even duplicates — should check GUIDs first before sanitizing
- `preferences` store (theme, articlesPerPage, showReadArticles) defined but unused
- Bookmarks nav item in sidebar is a placeholder

---

## Phase 1 — RSS Reader Polish (near-term)

Small remaining items before moving on to the Library.

- [ ] GUID deduplication in `SaveRSSArticles()` — skip processing if GUID already exists in DB
- [ ] Pagination UI in the feed list ("load more" or infinite scroll)
- [ ] Manual feed refresh button
- [ ] Unread count badges in the sidebar per-feed
- [ ] Mark all as read for a feed
- [ ] `preferences` store wired up (articles per page, theme toggle)

---

## Phase 2 — Library (save for later)

A dedicated section for saving content to read later, separate from the live RSS feed.

### Concept
- Save any RSS article to your Library with one click
- Save any URL directly (paste a link, content is fetched and stored)
- Saved items persist indefinitely — they don't disappear when a feed is updated
- Items can be tagged, archived, and searched

### Backend — new `saved_item` table

```sql
CREATE TABLE saved_item (
  id          SERIAL PRIMARY KEY,
  source_type TEXT NOT NULL,        -- 'rss_article' | 'url' | 'webpage'
  source_id   INT REFERENCES article(id) ON DELETE SET NULL,
  url         TEXT,
  title       TEXT,
  author      TEXT,
  content     TEXT,                 -- sanitized full HTML
  published_date TEXT,
  saved_at    TIMESTAMP DEFAULT NOW(),
  read        BOOLEAN DEFAULT false,
  archived    BOOLEAN DEFAULT false
);
```

### Backend — new API endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/library` | Save item — body: `{ url }` or `{ articleId }` |
| GET | `/api/library` | List saved items (paginated, filter by archived/read) |
| GET | `/api/library/single?id=` | Single saved item |
| PUT | `/api/library/update?id=` | Update read/archived status |
| DELETE | `/api/library?id=` | Delete saved item |
| GET | `/api/library/search?q=` | Full-text search across title + content |

### Content fetching for URLs
The FiveFilters service is already running and already used for full-text RSS extraction.
The same service can fetch a single URL:
```
GET http://fullfeedrss:80/makefulltextfeed.php?url=<URL>&max=1
```
This returns RSS XML with one item containing the full page content.
Parse it the same way as RSS articles — reuse existing `SaveRSSArticles` logic.

For pages that FiveFilters can't handle well (SPAs, JS-heavy sites), consider adding
[go-readability](https://github.com/go-shiori/go-readability) as a fallback.

### Frontend

- **Library** section in sidebar (below the feed list)
- **LibraryFeed.svelte** — saved items list, same card style as article feed
  - Filter: All / Unread / Archived
  - "Save URL" button at the top
- **Save button** on each article card and in ArticleDetail / ArticleModal toolbars
- Saved items open in the same ArticleModal reading view

---

## Phase 3 — Highlights & Notes

Annotate saved items with text highlights and notes.

### Concept
- Select text in the reading view → highlight toolbar appears → choose colour → saved
- Highlights persist and are shown when you reopen the item
- Each highlight can have an optional note attached
- Notes are searchable

### Backend — new tables

```sql
CREATE TABLE highlight (
  id            SERIAL PRIMARY KEY,
  saved_item_id INT NOT NULL REFERENCES saved_item(id) ON DELETE CASCADE,
  text          TEXT NOT NULL,        -- the highlighted text
  start_offset  INT,                  -- character offset in content
  end_offset    INT,
  color         TEXT DEFAULT 'yellow', -- 'yellow' | 'green' | 'blue' | 'pink'
  note          TEXT,
  created_at    TIMESTAMP DEFAULT NOW()
);
```

### Backend — new API endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/highlights` | Create highlight — body: `{ savedItemId, text, startOffset, endOffset, color }` |
| GET | `/api/highlights?itemId=` | Get all highlights for an item |
| PUT | `/api/highlights?id=` | Update color or note |
| DELETE | `/api/highlights?id=` | Delete highlight |

### Frontend

**Highlight rendering:**
- When a saved item opens, fetch its highlights
- Re-inject highlights into the HTML by wrapping matched text spans with `<mark>` elements
- Color-coded: yellow / green / blue / pink

**Selection toolbar:**
- On `mouseup` in the article body, check `window.getSelection()`
- If a non-empty selection exists, show a small floating toolbar above it
- Toolbar options: highlight colours + dismiss
- Capture `startOffset`/`endOffset` relative to the article content container

**Notes panel:**
- Drawer or side panel listing all highlights for the current item
- Each highlight shows: quoted text, colour chip, note (editable inline)
- Clicking a highlight scrolls to it in the article

---

## Phase 4 — Search

A unified search across both the RSS feed and the Library.

### Backend
PostgreSQL full-text search is already in place for RSS articles (`SearchArticles`).
Apply the same `tsvector` pattern to `saved_item`:

```sql
-- Add to saved_item table
ALTER TABLE saved_item ADD COLUMN search_vector tsvector
  GENERATED ALWAYS AS (
    to_tsvector('english', coalesce(title, '') || ' ' || coalesce(content, '') || ' ' || coalesce(author, ''))
  ) STORED;

CREATE INDEX saved_item_search_idx ON saved_item USING GIN(search_vector);
```

Also search across highlight notes:
```sql
-- In search query, JOIN with highlights and search notes too
```

### Frontend
- Existing search bar (once built) queries both RSS articles and saved items
- Results grouped by source (Feed vs Library)
- Highlight matched terms in results

---

## Phase 5 — Web Clipper (stretch)

Save entire web pages including layout, images, and styling — for documentation,
reference material, long-form articles.

### Concept
Unlike Phase 2 (which extracts readable text), this preserves the page more fully.

### Options
1. **Single-file HTML** — fetch the page, inline all CSS/images as base64, store as a
   single HTML blob. Libraries: [monolith](https://github.com/Y2Z/monolith) (Rust CLI,
   could be a separate Docker service).
2. **Screenshot** — use headless Chromium (Playwright/Puppeteer) to screenshot the page.
   Heavy but accurate.
3. **Wayback-style** — store the raw HTML + separately crawled assets.

**Recommended approach:** Run `monolith` as a small Docker sidecar service.
POST a URL to it, get back a self-contained `.html` file, store it in the DB (or on disk).

### New DB field
```sql
ALTER TABLE saved_item ADD COLUMN raw_html TEXT; -- full single-file HTML blob
```

A toggle in the reader switches between "readable" view (extracted content) and
"original" view (full page in an iframe sandbox).

---

## Architecture Notes

### Existing infrastructure that carries forward
| Existing piece | Reused for |
|---|---|
| FiveFilters service | Fetching full content for saved URLs (Phase 2) |
| bluemonday sanitizer | Cleaning saved content (Phase 2) |
| PostgreSQL full-text search | Library search (Phase 4) |
| ArticleModal reading view | Library reading view (Phase 2) |
| Sidebar navigation pattern | Library nav section (Phase 2) |

### New Docker services (if needed)
- `monolith` — web page archiving (Phase 5 only)
- No new services needed for Phases 2–4

### Frontend nav structure (end state)
```
Sidebar
├── Feed
│   ├── All Articles
│   ├── [Category]
│   │   └── [Feed]
│   └── ...
├── Library            ← Phase 2
│   ├── All Saved
│   ├── Unread
│   └── Archived
├── ─────────────
├── Manage Feeds
└── Categories
```

---

## Non-goals (for now)
- Multi-user / auth
- Mobile app
- Browser extension (the "save URL" flow covers the use case manually)
- AI summarisation
