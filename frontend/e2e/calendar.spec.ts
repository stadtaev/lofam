import { test, expect } from '@playwright/test'

const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December'
]

test.describe('Calendar', () => {
  test('Today button navigates to current month', async ({ page }) => {
    await page.goto('/')

    const today = new Date()
    const currentMonth = MONTHS[today.getMonth()]
    const currentYear = today.getFullYear().toString()

    // Navigate away from current month
    await page.getByRole('button', { name: '›' }).click()
    await page.getByRole('button', { name: '›' }).click()

    // Verify we're not on current month
    const monthHeading = page.getByRole('heading', { level: 2 })
    await expect(monthHeading).not.toHaveText(currentMonth)

    // Click Today button
    await page.getByRole('button', { name: 'Today' }).click()

    // Verify we're back to current month
    await expect(monthHeading).toHaveText(currentMonth)
    await expect(page.getByText(currentYear)).toBeVisible()
  })

  test('Today button works after navigating to previous months', async ({ page }) => {
    await page.goto('/')

    const today = new Date()
    const currentMonth = MONTHS[today.getMonth()]

    // Navigate to previous months
    await page.getByRole('button', { name: '‹' }).click()
    await page.getByRole('button', { name: '‹' }).click()
    await page.getByRole('button', { name: '‹' }).click()

    // Click Today
    await page.getByRole('button', { name: 'Today' }).click()

    // Verify current month
    await expect(page.getByRole('heading', { level: 2 })).toHaveText(currentMonth)
  })
})
