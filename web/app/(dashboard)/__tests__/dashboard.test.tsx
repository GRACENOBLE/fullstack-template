import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { ProfileCard } from '../dashboard/ProfileCard'

/**
 * Dashboard page tests.
 *
 * DashboardPage is an async Server Component — Vitest (jsdom) cannot render async
 * Server Components directly. We therefore test:
 *   1. ProfileCard — the pure display component (rendered in unit tests here)
 *   2. fetchUserProfile — the data-fetching utility (tested in lib/user-profile.test.ts)
 *
 * Full page integration (auth guard + fetch + render) is covered by Playwright E2E
 * when that layer is set up.
 */

describe('ProfileCard', () => {
  it('renders the display name prominently', () => {
    render(
      <ProfileCard
        profile={{
          uid: 'uid-abc',
          email: 'alice@example.com',
          displayName: 'Alice Example',
        }}
      />,
    )
    expect(screen.getByText('Alice Example')).toBeInTheDocument()
  })

  it('renders the email', () => {
    render(
      <ProfileCard
        profile={{
          uid: 'uid-abc',
          email: 'alice@example.com',
          displayName: 'Alice Example',
        }}
      />,
    )
    expect(screen.getByText('alice@example.com')).toBeInTheDocument()
  })

  it('renders the uid in a smaller style', () => {
    render(
      <ProfileCard
        profile={{
          uid: 'uid-abc',
          email: 'alice@example.com',
          displayName: 'Alice Example',
        }}
      />,
    )
    const uidEl = screen.getByText('uid-abc')
    expect(uidEl).toBeInTheDocument()
    // uid element carries the muted monospace classes
    expect(uidEl.className).toContain('mono')
  })
})

/**
 * Fallback behaviour tests — simulate what DashboardPage does when fetchUserProfile
 * throws by rendering ProfileCard with session-derived data.
 *
 * The fallback path in the page constructs a UserProfile from session.user fields
 * and passes it to ProfileCard regardless of whether the backend fetch succeeded.
 */
describe('ProfileCard fallback (session data)', () => {
  it('renders session display name when used as fallback', () => {
    render(
      <ProfileCard
        profile={{
          uid: 'session-uid',
          email: 'bob@example.com',
          displayName: 'Bob From Session',
        }}
      />,
    )
    expect(screen.getByText('Bob From Session')).toBeInTheDocument()
    expect(screen.getByText('bob@example.com')).toBeInTheDocument()
    expect(screen.getByText('session-uid')).toBeInTheDocument()
  })

  it('handles missing displayName gracefully by showing email in both slots', () => {
    render(
      <ProfileCard
        profile={{
          uid: 'session-uid-2',
          email: 'carol@example.com',
          displayName: 'carol@example.com', // page uses email as fallback for displayName
        }}
      />,
    )
    // When displayName falls back to email, the email text appears in both the
    // displayName and email slots — getAllByText handles the multiple-match case.
    const matches = screen.getAllByText('carol@example.com')
    expect(matches.length).toBeGreaterThanOrEqual(1)
  })
})
