import { test, expect } from '@playwright/test'

test.describe('auth guard', () => {
  test('redirects /dashboard to /login when unauthenticated', async ({ page }) => {
    await page.goto('/dashboard')
    await expect(page).toHaveURL(/\/login/)
  })

  test('redirects /settings to /login when unauthenticated', async ({ page }) => {
    await page.goto('/settings')
    await expect(page).toHaveURL(/\/login/)
  })
})

test.describe('login page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
  })

  test('renders email field, password field, and sign in button', async ({ page }) => {
    await expect(page.getByLabel('Email')).toBeVisible()
    await expect(page.getByLabel('Password')).toBeVisible()
    await expect(page.getByRole('button', { name: /sign in/i })).toBeVisible()
  })

  test('shows link to register page', async ({ page }) => {
    await expect(page.getByRole('link', { name: /sign up/i })).toBeVisible()
  })

  test('shows validation errors on empty submit', async ({ page }) => {
    await page.getByRole('button', { name: /sign in/i }).click()
    await expect(page.getByText('Please enter a valid email address')).toBeVisible()
    await expect(page.getByText('Password is required')).toBeVisible()
  })
})

test.describe('register page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/register')
  })

  test('renders name, email, password, and confirm password fields', async ({ page }) => {
    await expect(page.getByLabel('Name')).toBeVisible()
    await expect(page.getByLabel('Email')).toBeVisible()
    await expect(page.getByLabel('Password')).toBeVisible()
    await expect(page.getByLabel('Confirm password')).toBeVisible()
  })

  test('shows link back to login page', async ({ page }) => {
    await expect(page.getByRole('link', { name: 'Sign in' })).toBeVisible()
  })
})
