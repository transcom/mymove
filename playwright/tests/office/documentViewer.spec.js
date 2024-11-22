/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../utils/office/officeTest';

test.describe('The document viewer', () => {
  test.describe('When not logged in', () => {
    test('shows page not found', async ({ page, baseURL }) => {
      await Promise.all([page.waitForURL(`${baseURL}/moves/foo/documents`), page.goto('/moves/foo/documents')]);
      await expect(page.getByText('Welcome')).toBeVisible();
      // sign in button not in header
      await expect(page.locator('#main').getByRole('button', { name: 'Sign in' })).toBeVisible();
    });
  });

  test.describe('When logged in', () => {
    test('displays a PDF file correctly', async ({ page, officePage }) => {
      test.slow(); // flaky no flaky

      // Debugging, highly recommend to persist permanently due to the uniqueness of this test
      // Log all page requests for the pdfjs-dist webpack chunk
      page.on('request', (request) => {
        if (request.url().includes('pdfjs-dist-webpack.chunk.js')) {
          // eslint-disable-next-line no-console
          console.log(`>> Outgoing chunk request: ${request.method()} ${request.url()}`);
        }
      });
      // Log all page responses for the pdfjs-dist webpack chunk
      page.on('response', (response) => {
        if (response.url().includes('pdfjs-dist-webpack.chunk.js')) {
          // eslint-disable-next-line no-console
          console.log(`<< Incoming chunk response: ${response.status()} ${response.url()}`);
        }
      });

      // Build a move that has a PDF document
      const move = await officePage.testHarness.buildHHGWithAmendedOrders();

      // Sign in as a TOO user (Any office user works)
      await officePage.signInAsNewTOOUser();

      // Navigate to the move
      await officePage.tooNavigateToMove(move.locator);
      await officePage.waitForLoading();

      // Navigate to the document viewer
      await page.getByTestId('edit-orders').click();
      await officePage.waitForLoading();

      // Verify that the document viewer content is displayed
      await expect(page.getByTestId('DocViewerContent')).toBeVisible(); // This can load but the PDF still fail

      // Wait for the PDF canvas to load <- This is the meat and potatoes
      const pdfCanvas = page.locator('.pdf-canvas canvas').first();
      await pdfCanvas.waitFor({ state: 'visible', timeout: 10000 });

      // Verify that the canvas box has dimensions
      const canvasBox = await pdfCanvas.boundingBox();
      expect(canvasBox.width).toBeGreaterThan(0);
      expect(canvasBox.height).toBeGreaterThan(0);

      // Test zoom functionality

      // Verify that the zoom in and zoom out buttons are visible
      await expect(page.getByRole('button', { name: 'Zoom in' })).toBeVisible();
      await expect(page.getByRole('button', { name: 'Zoom out' })).toBeVisible();

      // Capture the current canvas dimensions, we compare this later for zoom
      const initialWidth = canvasBox.width;
      const initialHeight = canvasBox.height;

      // Make it zoom
      await page.getByRole('button', { name: 'Zoom in' }).click();

      // Check that it zoomed in
      const zoomedCanvasBox = await pdfCanvas.boundingBox();
      expect(zoomedCanvasBox.width).toBeGreaterThan(initialWidth);
      expect(zoomedCanvasBox.height).toBeGreaterThan(initialHeight);

      // Test rotation functionality

      // Verify that rotate buttons are visible
      await expect(page.getByRole('button', { name: 'Rotate left' })).toBeVisible();
      await expect(page.getByRole('button', { name: 'Rotate right' })).toBeVisible();

      // Make it rotate
      await page.getByRole('button', { name: 'Rotate right' }).click();

      // Testing the rotation degree change isn't possible within playwright so instead we verify nothing broke
      // The canvas will unload if the rotation breaks something
      const rotatedPdfCanvas = page.locator('.pdf-canvas canvas').first();
      await rotatedPdfCanvas.waitFor({ state: 'visible', timeout: 10000 });
    });
  });
});
