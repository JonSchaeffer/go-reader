# go-reader — Product Plan

## What it is

A self-hosted reading app combining an RSS reader with a personal reading library.
Think Readwise Reader, but yours.

---

## Current State (as of March 2026)

### Completed
- ✅ **Phase 1** — RSS reader polish
  - RSS feed management (add, delete, categorise, refresh)
  - Background feed fetching every 5 minutes via FiveFilters Full-Text RSS
  - GUID deduplication on fetch cycle
  - Article list with read/unread tracking, pagination, mark all read
  - Unread count badges per feed in sidebar
  - Category management
  - Article sidebar panel + full-screen reading modal with prev/next + keyboard nav
  - Background blur when modal is open

- ✅ **Phase 2** — Library
  - Save RSS articles or arbitrary URLs to a personal library
  - Full-text extraction via FiveFilters for saved URLs
  - Read/unread, archive, pagination
  - "Save URL" form in Library view
  - Save button in article detail panel and modal

- ✅ **Phase 3** — Highlights & Notes
  - Select text → colour picker popover → saved highlight with optional note
  - Highlights persist and re-render on article open (DOM TreeWalker approach)
  - `text_offset` anchoring to handle repeated words correctly
  - Edit colour/note or delete by clicking existing highlight
  - Popover dismisses on scroll
  - Works in both sidebar panel and full-screen modal

- ✅ **Phase 3.5** — Highlights page
  - View all highlights across all articles + library items
  - Colour-coded, shows note, source title, date
  - Click source → opens article in reading modal

- ✅ **Phase 4** — Search
  - Unified full-text search across articles and library items
  - `plainto_tsquery` PostgreSQL FTS, HTML stripped from excerpts
  - Debounced search input, type icons, result counts
  - Click result → opens in reading modal

- ✅ **Phase 5a** — Web Clipper (bookmarklet)
  - Drag-to-bookmarks-bar bookmarklet
  - Toast notification on save (saved / already saved / error)
  - Configurable backend URL
  - Raw code copy for manual bookmark creation

---

## Readwise Reader Feature Audit

Comparing what Readwise Reader has vs what we've built.

### Core reading experience

| Feature | Status | Notes |
|---|---|---|
| RSS feeds | ✅ Done | |
| Save URLs to library | ✅ Done | |
| Full-text extraction | ✅ Done | FiveFilters |
| Read/unread tracking | ✅ Done | |
| Archive | ✅ Done | |
| Highlights + notes | ✅ Done | |
| Search | ✅ Done | |
| Web clipper | ✅ Done (bookmarklet) | Extension planned |
| **Reading progress** | ❌ Missing | Save/restore scroll position per article |
| **Reading time estimates** | ❌ Missing | "5 min read" from word count |
| **Tags** | ❌ Missing | Per-item tags for library items |
| **Reader customisation** | ❌ Missing | Font size, font family, line width, theme |
| **Keyboard shortcuts** | ⚠️ Partial | Modal nav only (←→ Esc), not comprehensive |
| **Document notes** | ❌ Missing | Note on the whole article (not just highlights) |
| **Filtering / sorting** | ⚠️ Partial | Unread/All/Archived only; no date/tag/length sort |
| **Bulk actions** | ❌ Missing | Multi-select → archive/delete/tag |

### Differentiating features

| Feature | Status | Notes |
|---|---|---|
| **Daily highlight review** | ❌ Missing | Readwise's core: resurface past highlights daily |
| **Highlight export** | ❌ Missing | Markdown / CSV / copy-all |
| **Spaced repetition** | ❌ Missing | Schedule-based highlight review |
| **Browser extension** | ⚠️ Planned | Bookmarklet done; full extension is next |
| **Email newsletters** | ❌ Missing | Inbound email → library item |
| **PDF support** | ❌ Missing | Render PDFs in the reader |
| **Full page archiving** | ❌ Missing | Preserve full page layout (monolith approach) |
| **Mobile / PWA** | ❌ Missing | Responsive design + installable |
| **Feed discovery** | ❌ Missing | Browse/search for new feeds |

### Out of scope (for now)
- Multi-user / auth
- AI summarisation (Ghostreader equivalent)
- Podcast support
- Social / sharing highlights publicly
- Twitter/X thread saving

---

## Upcoming Phases

---

## Phase 6 — Reader Customisation + Reading Progress

The biggest gap in the reading *experience* right now.

### 6a — Reader customisation
User-controlled typography and theme settings, persisted to `localStorage`.

**Controls:**
- Font size (sm / md / lg / xl)
- Font family (sans-serif / serif / monospace)
- Line width (narrow / medium / wide)
- Theme (dark / light / sepia)

**Implementation:**
- `ReaderSettings.svelte` — small popover or drawer accessible from the modal/detail toolbar
- CSS custom properties (`--reader-font-size`, `--reader-font-family`, etc.) applied to the `.prose` wrapper
- Settings saved to `localStorage` and loaded on mount in `ArticleModal` and `ArticleDetail`

### 6b — Reading progress
Save and restore scroll position per article.

**Backend:**
```sql
-- Add to saved_item and consider a separate progress table for articles
CREATE TABLE reading_progress (
  id           SERIAL PRIMARY KEY,
  item_type    TEXT NOT NULL,   -- 'article' | 'library'
  item_id      INT NOT NULL,
  scroll_pct   FLOAT NOT NULL DEFAULT 0,  -- 0.0–1.0
  updated_at   TIMESTAMP DEFAULT NOW(),
  UNIQUE (item_type, item_id)
);
```

**API:** `PUT /api/progress` with `{ itemType, itemId, scrollPct }` — upsert.
`GET /api/progress?type=&id=` — fetch on open.

**Frontend:**
- Debounced `scroll` event handler saves position (at most once per 2s)
- On article open, fetch progress and `scrollTo()` after content renders

### 6c — Reading time estimates
Simple word count from stripped HTML content.

- Backend: calculate at save time, store as `word_count INT` on `saved_item` and `article`
- Frontend: show "N min read" on article cards and in reader header
- Formula: `Math.ceil(wordCount / 238)` (average adult reading speed)

---

## Phase 7 — Tags

Per-item tags for library items (and optionally RSS articles).

### Backend
```sql
CREATE TABLE tag (
  id    SERIAL PRIMARY KEY,
  name  TEXT NOT NULL UNIQUE
);

CREATE TABLE item_tag (
  item_type  TEXT NOT NULL,  -- 'library' | 'article'
  item_id    INT NOT NULL,
  tag_id     INT NOT NULL REFERENCES tag(id) ON DELETE CASCADE,
  PRIMARY KEY (item_type, item_id, tag_id)
);
```

**API:**
| Method | Path | Description |
|---|---|---|
| GET | `/api/tags` | All tags |
| POST | `/api/tags` | Create tag |
| DELETE | `/api/tags?id=` | Delete tag |
| GET | `/api/tags/item?type=&id=` | Tags for an item |
| POST | `/api/tags/item` | Add tag to item |
| DELETE | `/api/tags/item?type=&id=&tagId=` | Remove tag from item |

### Frontend
- Tag pills on article cards and in reader toolbar
- Inline tag editor (type to create or select existing, click × to remove)
- Filter Library by tag
- Tags shown on Highlights page per highlight's source

---

## Phase 8 — Daily Highlight Review

Readwise's killer feature: a daily in-app review queue that resurfaces past highlights.

### Concept
- Each day, pick N highlights (default 5) from your library using a spaced-repetition-inspired algorithm
- Present them one at a time in a focused "review" mode
- User can mark as reviewed, edit note, or dismiss
- Track review history to avoid showing the same highlight every day

### Backend
```sql
ALTER TABLE highlight ADD COLUMN last_reviewed_at TIMESTAMP;
ALTER TABLE highlight ADD COLUMN review_count INT NOT NULL DEFAULT 0;
```

**API:**
- `GET /api/review/queue` — return today's review highlights (picks unseen/oldest-reviewed)
- `POST /api/review/complete?id=` — mark a highlight as reviewed

### Frontend
- **ReviewView.svelte** — full-screen focused review mode
  - Shows highlight text in its colour, with source title and note
  - "Next" button to advance through queue
  - Progress indicator (3 / 5)
  - "Done for today" when queue is empty
- **Daily review badge** on sidebar nav item when queue is ready

---

## Phase 9 — Highlight Export

Get your highlights out in useful formats.

### Options
- **Copy all** — copy all highlights for an article as plain text to clipboard
- **Markdown export** — download `.md` file with highlights + notes formatted as blockquotes
- **CSV export** — download `.csv` for import into Notion, Obsidian, Anki, etc.

### Implementation
- Export buttons on the Highlights page (export all) and in the reader toolbar (export for current article)
- All client-side — no backend needed, just build strings from the highlight data and trigger a download

---

## Phase 10 — Browser Extension

Full Chrome/Firefox extension using the same core logic as the bookmarklet.

### Architecture
The bookmarklet JS is the seed. The extension wraps it in:
- `manifest.json` (MV3 for Chrome, MV2 compatible for Firefox)
- `popup.html` — small popup with "Save to YARR" button + status
- `background.js` (service worker) — handles the fetch to the backend
- `options.html` — configure backend URL (same as clipper page)

### Key reuse from bookmarklet
- Same `POST /api/library` endpoint
- Same toast-style feedback logic
- Same CORS setup already in place

### Distribution
- Chrome: upload to Chrome Web Store (or sideload via `chrome://extensions`)
- Firefox: upload to AMO (or sideload via `about:debugging`)
- Self-hosted option: zip the extension folder, users install manually

---

## Phase 11 — Full Page Archiving

Preserve the complete page layout for reference material.

### Approach
Run [monolith](https://github.com/Y2Z/monolith) as a Docker sidecar service.
POST a URL to it → get back a self-contained single-file HTML blob → store in DB.

### Backend
```sql
ALTER TABLE saved_item ADD COLUMN raw_html TEXT;
```
New endpoint: `POST /api/library/archive?id=` — fetches and stores the full page for an existing saved item.

### Frontend
Toggle in reader between "Readable" (extracted content) and "Original" (full page in `<iframe sandbox>`).

---

## Phase 12 — Mobile / PWA

Make the app usable on phones and tablets.

### Approach
- Add responsive Tailwind breakpoints (stack sidebar + content vertically on mobile)
- Sidebar becomes a slide-in drawer on small screens
- Add `manifest.json` + service worker for PWA installability
- Touch-friendly highlight selection (long-press to trigger popover)

---

## Architecture Notes

### Existing infrastructure that carries forward
| Existing piece | Reused for |
|---|---|
| FiveFilters service | URL saving (Phase 2 ✅), archiving (Phase 11) |
| bluemonday sanitizer | All saved content ✅ |
| PostgreSQL FTS | Search ✅, tags filter (Phase 7) |
| ArticleModal reading view | All reading phases ✅ |
| Highlight system | Daily review (Phase 8), export (Phase 9) |
| Bookmarklet core logic | Browser extension (Phase 10) |

### Recommended build order
1. Phase 6 (reader customisation + progress) — biggest UX gap
2. Phase 7 (tags) — core organisation feature
3. Phase 8 (daily review) — the Readwise differentiator
4. Phase 9 (highlight export) — quick win
5. Phase 10 (browser extension) — builds on bookmarklet
6. Phase 11 (archiving) — nice to have
7. Phase 12 (mobile) — depends on usage needs
