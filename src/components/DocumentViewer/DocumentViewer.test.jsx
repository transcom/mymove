/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { screen, waitFor, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import DocumentViewer from './DocumentViewer';
import samplePDF from './sample.pdf';
import sampleJPG from './sample.jpg';
import samplePNG from './sample2.png';
import sampleGIF from './sample3.gif';

import { bulkDownloadPaymentRequest } from 'services/ghcApi';
import { UPLOAD_SCAN_STATUS, UPLOAD_DOC_STATUS_DISPLAY_MESSAGE } from 'shared/constants';
import { renderWithProviders } from 'testUtils';

const toggleMenuClass = () => {
  const container = document.querySelector('[data-testid="menuButtonContainer"]');
  if (container) {
    container.className = container.className === 'closed' ? 'open' : 'closed';
  }
};
// Mocking necessary functions/module
const mockMutateUploads = jest.fn();

jest.mock('@tanstack/react-query', () => ({
  ...jest.requireActual('@tanstack/react-query'),
  useMutation: () => ({ mutate: mockMutateUploads }),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const mockFiles = [
  {
    id: 1,
    filename: 'Test File.pdf',
    contentType: 'application/pdf',
    url: samplePDF,
    createdAt: '2021-06-14T15:09:26.979879Z',
    rotation: 1,
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

const mockErrorFiles = [
  {
    id: 1,
    filename: 'Test File.pdf',
    contentType: 'application/pdf',
    url: '404',
    createdAt: '2021-06-14T15:09:26.979879Z',
  },
];

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  bulkDownloadPaymentRequest: jest.fn(),
}));

jest.mock('./Content/Content', () => ({
  __esModule: true,
  default: ({
    id,
    filename,
    contentType,
    url,
    createdAt,
    rotation,
    filePath,
    onError,
    disableSaveButton,
    saveRotation,
  }) => {
    if (filePath === '404') {
      onError('content error happening');
      return <div>nothing to see here</div>;
    }
    return (
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
        <button type="button" disabled={disableSaveButton} onClick={saveRotation}>
          Save
        </button>
      </div>
    );
  },
}));

describe('DocumentViewer component', () => {
  it('initial state is closed menu and first file selected', async () => {
    renderWithProviders(<DocumentViewer files={mockFiles} />);

    const selectedFileTitle = await screen.getAllByTestId('documentTitle')[0];
    expect(selectedFileTitle.textContent).toEqual('Test File 4.gif - Added on 16 Jun 2021');

    const menuButtonContainer = await screen.findByTestId('menuButtonContainer');
    expect(menuButtonContainer.className).toContain('closed');
  });

  it('renders the file creation date with the correctly sorted props', async () => {
    renderWithProviders(<DocumentViewer files={mockFiles} />);
    const files = screen.getAllByRole('listitem');

    expect(files[0].textContent).toContain('Test File 4.gif - Added on 2021-06-16T15:09:26.979879Z');
  });

  it('renders the title bar with the correct props', async () => {
    renderWithProviders(<DocumentViewer files={mockFiles} />);

    const title = await screen.getAllByTestId('documentTitle')[0];

    expect(title.textContent).toContain('Test File 4.gif - Added on 16 Jun 2021');
  });

  it('handles the open menu button', async () => {
    renderWithProviders(<DocumentViewer files={mockFiles} />);

    const openMenuButton = await screen.findByTestId('menuButton');

    await userEvent.click(openMenuButton);

    const menuButtonContainer = await screen.findByTestId('menuButtonContainer');
    expect(menuButtonContainer.className).toContain('open');
  });

  it('handles the close menu button', async () => {
    renderWithProviders(<DocumentViewer files={mockFiles} />);

    // defaults to closed so we need to open it first.
    const openMenuButton = await screen.findByTestId('menuButton');

    await userEvent.click(openMenuButton);

    const menuButtonContainer = await screen.findByTestId('menuButtonContainer');
    expect(menuButtonContainer.className).toContain('open');

    await userEvent.click(openMenuButton);

    expect(menuButtonContainer.className).toContain('closed');
  });

  it('shows error if file type is unsupported', async () => {
    renderWithProviders(
      <DocumentViewer files={[{ id: 99, filename: 'archive.zip', contentType: 'zip', url: '/path/to/archive.zip' }]} />,
    );

    expect(screen.getByText('id: undefined')).toBeInTheDocument();
  });

  describe('regarding content errors', () => {
    const errorMessageText = 'If your document does not display, please refresh your browser.';
    const downloadLinkText = 'Download file';
    it('no error message normally', async () => {
      renderWithProviders(<DocumentViewer files={mockFiles} />);
      expect(screen.queryByText(errorMessageText)).toBeNull();
    });

    it('download link normally', async () => {
      renderWithProviders(<DocumentViewer files={mockFiles} allowDownload />);
      expect(screen.getByText(downloadLinkText)).toBeVisible();
    });

    it('show message on content error', async () => {
      renderWithProviders(<DocumentViewer files={mockErrorFiles} />);
      expect(screen.getByText(errorMessageText)).toBeVisible();
    });

    it('download link on content error', async () => {
      renderWithProviders(<DocumentViewer files={mockErrorFiles} allowDownload />);
      expect(screen.getByText(downloadLinkText)).toBeVisible();
    });
  });

  describe('when clicking download Download All Files button', () => {
    it('downloads a bulk packet', async () => {
      const mockResponse = {
        ok: true,
        headers: {
          'content-disposition': 'filename="test.pdf"',
        },
        status: 200,
        data: null,
      };

      renderWithProviders(
        <DocumentViewer
          files={[
            { id: 99, filename: 'archive.zip', contentType: 'zip', url: 'path/to/archive.zip' },
            { id: 99, filename: 'archive.zip', contentType: 'zip', url: 'path/to/archive.zip' },
          ]}
          paymentRequestId="PaymentRequestId"
        />,
      );

      bulkDownloadPaymentRequest.mockImplementation(() => Promise.resolve(mockResponse));

      const downloadButton = screen.getByText('Download All Files (PDF)', { exact: false });
      await userEvent.click(downloadButton);
      await waitFor(() => {
        expect(bulkDownloadPaymentRequest).toHaveBeenCalledTimes(1);
      });
    });
  });

  describe('DocumentViewer', () => {
    it('disables Save button when no rotation change', async () => {
      renderWithProviders(<DocumentViewer files={mockFiles} />);

      const saveBtn = await screen.findByRole('button', { name: /save/i });
      expect(saveBtn).toBeDisabled();
    });
  });
});

// Mock the EventSource
class MockEventSource {
  constructor(url) {
    this.url = url;
    this.onmessage = null;
  }

  close() {
    this.isClosed = true;
  }
}
global.EventSource = MockEventSource;
// Helper function for finding the file status text
const findByTextContent = (text) => {
  return screen.getByText((content, node) => {
    const hasText = (element) => element.textContent.includes(text);
    const nodeHasText = hasText(node);
    const childrenDontHaveText = Array.from(node.children).every((child) => !hasText(child));
    return nodeHasText && childrenDontHaveText;
  });
};

describe('Test DocumentViewer File Upload Statuses', () => {
  let eventSource;
  const renderDocumentViewer = (props) => {
    return renderWithProviders(<DocumentViewer {...props} />);
  };

  beforeEach(() => {
    eventSource = new MockEventSource('');
    jest.spyOn(global, 'EventSource').mockImplementation(() => eventSource);
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('displays Uploading status', () => {
    renderDocumentViewer({ files: mockFiles, isFileUploading: true });
    expect(findByTextContent(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.UPLOADING)).toBeInTheDocument();
  });

  it('displays Scanning status', async () => {
    renderDocumentViewer({ files: mockFiles });
    await act(async () => {
      eventSource.onmessage({ data: UPLOAD_SCAN_STATUS.PROCESSING });
    });
    await waitFor(() => {
      expect(findByTextContent(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.SCANNING)).toBeInTheDocument();
    });
  });

  it('displays Establishing document for viewing  status', async () => {
    renderDocumentViewer({ files: mockFiles });
    await act(async () => {
      eventSource.onmessage({ data: UPLOAD_SCAN_STATUS.CLEAN });
    });
    await waitFor(() => {
      expect(
        findByTextContent(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.ESTABLISHING_DOCUMENT_FOR_VIEWING),
      ).toBeInTheDocument();
    });
  });

  it('displays infected file message', async () => {
    renderDocumentViewer({ files: mockFiles });
    await act(async () => {
      eventSource.onmessage({ data: UPLOAD_SCAN_STATUS.INFECTED });
    });
    await waitFor(() => {
      expect(findByTextContent(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.INFECTED_FILE_MESSAGE)).toBeInTheDocument();
    });
  });

  it('displays File Not Found message when no file is selected', () => {
    renderDocumentViewer({ files: [] });
    expect(findByTextContent(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.FILE_NOT_FOUND)).toBeInTheDocument();
  });
});
