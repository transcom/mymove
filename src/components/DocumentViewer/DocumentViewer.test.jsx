/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import DocumentViewer from './DocumentViewer';
import samplePDF from './sample.pdf';
import sampleJPG from './sample.jpg';
import samplePNG from './sample2.png';
import sampleGIF from './sample3.gif';

beforeEach(() => {
  jest.clearAllMocks();
});

jest.mock('components/DocumentViewer/Content/Content', () => {
  return ({ fileType }) => {
    return (
      <div data-testid="DocViewerContent">
        {fileType === '' // fileType will converted to '' for unsupported types. Otherwise, it will be supported.
          ? 'No preview available for this kind of file.Download file to see the contents.'
          : 'Content'}
      </div>
    );
  };
});

const mockFiles = [
  {
    id: 1,
    filename: 'Test File.pdf',
    contentType: 'application/pdf',
    url: samplePDF,
    createdAt: '2021-06-14T15:09:26.979879Z',
  },
  {
    id: 2,
    filename: 'Test File 2.jpg',
    contentType: 'image/jpeg',
    url: sampleJPG,
    createdAt: '2021-06-12T15:09:26.979879Z',
  },
  {
    id: 3,
    filename: 'Test File 3.png',
    contentType: 'image/png',
    url: samplePNG,
    createdAt: '2021-06-15T15:09:26.979879Z',
  },
  {
    id: 4,
    filename: 'Test File 4.gif',
    contentType: 'image/gif',
    url: sampleGIF,
    createdAt: '2021-06-16T15:09:26.979879Z',
  },
];

describe('DocumentViewer component', () => {
  it('initial state is closed menu and first file selected', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} />
      </QueryClientProvider>,
    );
    const docMenu = await screen.findByTestId('DocViewerMenu');

    expect(docMenu.className).toContain('collapsed');

    // Files are ordered by createdAt date before being rendered.
    const firstFile = screen.getByRole('button', { name: 'Test File 4.gif Uploaded on 16-Jun-2021' });
    expect(firstFile.className).toContain('active');
  });

  it('renders the file creation date with the correctly sorted props', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} />
      </QueryClientProvider>,
    );

    const files = screen.getAllByRole('listitem');

    expect(files[0].textContent).toEqual('Test File 4.gif Uploaded on 16-Jun-2021');
  });

  it('renders the title bar with the correct props', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} />
      </QueryClientProvider>,
    );
    const title = await screen.findByTestId('documentTitle');

    expect(title.textContent).toEqual('Test File 4.gif - Added on 16 Jun 2021');
  });

  it('handles the open menu button', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} />
      </QueryClientProvider>,
    );
    const openMenuButton = await screen.findByTestId('openMenu');

    await userEvent.click(openMenuButton);

    const docMenu = screen.getByTestId('DocViewerMenu');

    await waitFor(() => {
      expect(docMenu.className).not.toContain('collapsed');
    });
  });

  it('handles the close menu button', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} />
      </QueryClientProvider>,
    );
    // defaults to closed so we need to open it first.
    const openMenuButton = await screen.findByTestId('openMenu');

    await userEvent.click(openMenuButton);

    const docMenu = screen.getByTestId('DocViewerMenu');

    await waitFor(() => {
      expect(docMenu.className).not.toContain('collapsed');
    });

    const closeMenuButton = await screen.findByTestId('closeMenu');

    await userEvent.click(closeMenuButton);

    await waitFor(() => expect(docMenu.className).toContain('collapsed'));
  });

  it.each([
    ['Test File 3.png Uploaded on 15-Jun-2021', 'Test File 3.png - Added on 15 Jun 2021'],
    // ['Test File.pdf Uploaded on 14-Jun-2021', 'Test File.pdf - Added on 14 Jun 2021'],  // TODO: figure out why this isn't working...
    ['Test File 2.jpg Uploaded on 12-Jun-2021', 'Test File 2.jpg - Added on 12 Jun 2021'],
  ])('handles selecting a different file (%s)', async (buttonText, titleText) => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} />
      </QueryClientProvider>,
    );
    // defaults to closed so we need to open it first.
    const openMenuButton = await screen.findByTestId('openMenu');

    await userEvent.click(openMenuButton);

    const docMenu = screen.getByTestId('DocViewerMenu');

    expect(docMenu.className).not.toContain('collapsed');

    const otherFile = await screen.findByRole('button', { name: buttonText });

    await userEvent.click(otherFile);

    expect(docMenu.className).toContain('collapsed');

    const title = await screen.findByTestId('documentTitle');

    expect(title.textContent).toEqual(titleText);

    await waitFor(() => expect(screen.queryByText('is not supported')).not.toBeInTheDocument());
  });

  it('shows error if file type is unsupported', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer
          files={[{ id: 99, filename: 'archive.zip', contentType: 'zip', url: '/path/to/archive.zip' }]}
        />
      </QueryClientProvider>,
    );

    // defaults to closed so we need to open it first.
    const openMenuButton = await screen.findByTestId('openMenu');

    await userEvent.click(openMenuButton);

    const docMenu = screen.getByTestId('DocViewerMenu');

    await waitFor(() => {
      expect(docMenu.className).not.toContain('collapsed');
    });

    const docContent = screen.getByTestId('DocViewerContent');

    expect(docContent.textContent).toEqual(
      'No preview available for this kind of file.Download file to see the contents.',
    );
  });

  it('displays file not found for empty files array', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer />
      </QueryClientProvider>,
    );

    expect(await screen.findByRole('heading', { name: 'File Not Found' })).toBeInTheDocument();
  });

  it('shows the download link option when allowDownload is true', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} allowDownload />
      </QueryClientProvider>,
    );

    expect(await screen.findByText('Download file')).toBeInTheDocument();
  });
});
