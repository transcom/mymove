/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import DocumentViewer from './DocumentViewer';
import samplePDF from './sample.pdf';
import sampleJPG from './sample.jpg';
import samplePNG from './sample2.png';
import sampleGIF from './sample3.gif';

import { bulkDownloadPaymentRequest } from 'services/ghcApi';
import { UPLOAD_DOC_STATUS, UPLOAD_SCAN_STATUS, UPLOAD_DOC_STATUS_DISPLAY_MESSAGE } from 'shared/constants';
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

global.EventSource = jest.fn().mockImplementation(() => ({
  addEventListener: jest.fn(),
  removeEventListener: jest.fn(),
  close: jest.fn(),
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
  default: ({ id, filename, contentType, url, createdAt, rotation, filePath, onError }) => {
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
});

describe('Test documentViewer file upload statuses', () => {
  const documentStatus = 'Document Status';
  // Trigger status change helper function
  const triggerStatusChange = (status, fileId, onStatusChange) => {
    // Mocking EventSource
    const mockEventSource = jest.fn();

    global.EventSource = mockEventSource;

    // Create a mock EventSource instance and trigger the onmessage event
    const eventSourceMock = {
      onmessage: () => {
        const event = { data: status };
        onStatusChange(event.data); // Pass status to the callback
      },
      close: jest.fn(),
    };

    mockEventSource.mockImplementationOnce(() => eventSourceMock);

    // Trigger the status change (this would simulate the file status update event)
    const sse = new EventSource(`/ghc/v1/uploads/${fileId}/status`, { withCredentials: true });
    sse.onmessage({ data: status });
  };

  it('displays UPLOADING status when file is uploading', async () => {
    renderWithProviders(<DocumentViewer files={mockFiles} allowDownload paymentRequestId={1} isFileUploading />);
    // Trigger UPLOADING status change
    triggerStatusChange(UPLOAD_DOC_STATUS.UPLOADING, mockFiles[0].id, async () => {
      // Wait for the component to update and check that the status is reflected
      await waitFor(() => {
        expect(screen.getByTestId('documentAlertHeading')).toHaveTextContent(documentStatus);
        expect(screen.getByTestId('documentAlertMessage')).toHaveTextContent(
          UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.UPLOADING,
        );
      });
    });
  });

  it('displays SCANNING status when file is scanning', async () => {
    renderWithProviders(
      <DocumentViewer files={mockFiles} allowDownload paymentRequestId={1} isFileUploading={false} />,
    );

    // Trigger SCANNING status change
    triggerStatusChange(UPLOAD_SCAN_STATUS.PROCESSING, mockFiles[0].id, async () => {
      // Wait for the component to update and check that the status is reflected
      await waitFor(() => {
        expect(screen.getByTestId('documentAlertHeading')).toHaveTextContent(documentStatus);
        expect(screen.getByTestId('documentAlertMessage')).toHaveTextContent(
          UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.SCANNING,
        );
      });
    });
  });

  it('displays ESTABLISHING status when file is establishing', async () => {
    renderWithProviders(
      <DocumentViewer files={mockFiles} allowDownload paymentRequestId={1} isFileUploading={false} />,
    );

    // Trigger ESTABLISHING status change
    triggerStatusChange(UPLOAD_SCAN_STATUS.CLEAN, mockFiles[0].id, async () => {
      // Wait for the component to update and check that the status is reflected
      await waitFor(() => {
        expect(screen.getByTestId('documentAlertHeading')).toHaveTextContent(documentStatus);
        expect(screen.getByTestId('documentAlertMessage')).toHaveTextContent(
          UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.ESTABLISHING_DOCUMENT_FOR_VIEW,
        );
      });
    });
  });

  it('displays FILE_NOT_FOUND status when no file is found', async () => {
    const emptyFileList = [];
    renderWithProviders(
      <DocumentViewer files={emptyFileList} allowDownload paymentRequestId={1} isFileUploading={false} />,
    );

    // Trigger FILE_NOT_FOUND status change (via props)
    triggerStatusChange('FILE_NOT_FOUND', '', async () => {
      // Wait for the component to update and check that the status is reflected
      await waitFor(() => {
        expect(screen.getByTestId('documentAlertHeading')).toHaveTextContent(documentStatus);
        expect(screen.getByTestId('documentAlertMessage')).toHaveTextContent(
          UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.FILE_NOT_FOUND,
        );
      });
    });
  });

  it('displays INFECTED status when file is infected', async () => {
    renderWithProviders(
      <DocumentViewer files={mockFiles} allowDownload paymentRequestId={1} isFileUploading={false} />,
    );
    // Trigger INFECTED status change
    triggerStatusChange(UPLOAD_SCAN_STATUS.INFECTED, mockFiles[0].id, async () => {
      // Wait for the component to update and check that the status is reflected
      await waitFor(() => {
        expect(screen.getByTestId('documentAlertHeading')).toHaveTextContent('Ask for a new file');
        expect(screen.getByTestId('documentAlertMessage')).toHaveTextContent(
          UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.INFECTED_FILE_MESSAGE,
        );
      });
    });
  });
});
