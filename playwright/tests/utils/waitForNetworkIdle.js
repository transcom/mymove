export async function waitForNetworkIdle(page, timeout = 1000, maxFailedAttempts = 3) {
  let failedAttempts = 0;
  while (failedAttempts < maxFailedAttempts) {
    try {
      await page.waitForTimeout(timeout);
      await page.waitForRequest(() => true, { timeout: 100 }).catch(() => null);
      return true; // If we reach this point, no requests were made for 'timeout' ms
    } catch (error) {
      failedAttempts += 1;
    }
  }
  return false; // Network did not become idle
}

export default waitForNetworkIdle;
