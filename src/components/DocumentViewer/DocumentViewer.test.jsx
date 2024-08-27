/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';

import DocumentViewer from './DocumentViewer';
import samplePDF from './sample.pdf';
import sampleJPG from './sample.jpg';
import samplePNG from './sample2.png';
import sampleGIF from './sample3.gif';

const toggleMenuClass = () => {
  const container = document.querySelector('[data-testid="menuButtonContainer"]');
  if (container) {
    container.className = container.className === 'closed' ? 'open' : 'closed';
  }
};

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
    rotation: 1,
  },
  {
    id: 4,
    filename: 'Test File 4.gif',
    contentType: 'image/gif',
    url: sampleGIF,
    createdAt: '2021-06-16T15:09:26.979879Z',
    rotation: 3,
  },
];

jest.mock('./Content/Content', () => ({
  __esModule: true,
  default: ({ id, filename, contentType, url, createdAt, rotation }) => (
    <div>
      <div data-testid="documentTitle">
        {filename} Uploaded on {createdAt}
      </div>
      <div>id: {id || 'undefined'}</div>
      <div>fileName: {filename || 'undefined'}</div>
      <div>contentType: {contentType || 'undefined'}</div>
      <div>url: {url || 'undefined'}</div>
      <div>createdAt: {createdAt || 'undefined'}</div>
      <div>rotation: {rotation || 'undefined'}</div>
      <div data-testid="listOfFiles">
        <ul>
          {mockFiles.map((file) => (
            <li key={file.id}>
              {file.filename} - Added on {file.createdAt}
            </li>
          ))}
        </ul>
      </div>
      <div data-testid="menuButtonContainer" className="closed">
        <button
          data-testid="menuButton"
          onClick={() => {
            toggleMenuClass();
          }}
          type="button"
        >
          Toggle
        </button>
      </div>
    </div>
  ),
}));

describe('DocumentViewer component', () => {
  it('initial state is closed menu and first file selected', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} />
      </QueryClientProvider>,
    );

    const selectedFileTitle = await screen.getAllByTestId('documentTitle')[0];
    expect(selectedFileTitle.textContent).toEqual('Test File 4.gif - Added on 16 Jun 2021');

    const menuButtonContainer = await screen.findByTestId('menuButtonContainer');
    expect(menuButtonContainer.className).toContain('closed');
  });

  it('renders the file creation date with the correctly sorted props', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} />
      </QueryClientProvider>,
    );

    const files = screen.getAllByRole('listitem');

    expect(files[0].textContent).toContain('Test File 4.gif - Added on 2021-06-16T15:09:26.979879Z');
  });

  it('renders the title bar with the correct props', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} />
      </QueryClientProvider>,
    );

    const title = await screen.getAllByTestId('documentTitle')[0];

    expect(title.textContent).toContain('Test File 4.gif - Added on 16 Jun 2021');
  });

  it('handles the open menu button', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} />
      </QueryClientProvider>,
    );

    const openMenuButton = await screen.findByTestId('menuButton');

    await userEvent.click(openMenuButton);

    const menuButtonContainer = await screen.findByTestId('menuButtonContainer');
    expect(menuButtonContainer.className).toContain('open');
  });

  it('handles the close menu button', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer files={mockFiles} />
      </QueryClientProvider>,
    );

    // defaults to closed so we need to open it first.
    const openMenuButton = await screen.findByTestId('menuButton');

    await userEvent.click(openMenuButton);

    const menuButtonContainer = await screen.findByTestId('menuButtonContainer');
    expect(menuButtonContainer.className).toContain('open');

    await userEvent.click(openMenuButton);

    expect(menuButtonContainer.className).toContain('closed');
  });

  it('shows error if file type is unsupported', async () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <DocumentViewer
          files={[{ id: 99, filename: 'archive.zip', contentType: 'zip', url: '/path/to/archive.zip' }]}
        />
      </QueryClientProvider>,
    );

    expect(screen.getByText('id: undefined')).toBeInTheDocument();
  });
});
