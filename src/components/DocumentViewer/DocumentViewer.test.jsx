/* eslint-disable react/jsx-props-no-spreading */
import React, { act } from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';

import DocumentViewer from './DocumentViewer';
import samplePDF from './sample.pdf';
import sampleJPG from './sample.jpg';
import samplePNG from './sample2.png';
import sampleGIF from './sample3.gif';

import { UPLOAD_DOC_STATUS, UPLOAD_SCAN_STATUS, UPLOAD_DOC_STATUS_DISPLAY_MESSAGE } from 'shared/constants';
import { bulkDownloadPaymentRequest } from 'services/ghcApi';

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

// Mock EventSource
class MockEventSource {
  constructor(url, config) {
    this.url = url;
    this.config = config;
    this.onmessage = null;
    this.onerror = null;
  }

  sendMessage(data) {
    if (this.onmessage) {
      this.onmessage({ data });
    }
  }

  triggerError() {
    if (this.onerror) {
      this.onerror();
    }
  }
}

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

  describe('regarding content errors', () => {
    const errorMessageText = 'If your document does not display, please refresh your browser.';
    const downloadLinkText = 'Download file';
    it('no error message normally', async () => {
      render(
        <QueryClientProvider client={new QueryClient()}>
          <DocumentViewer files={mockFiles} />
        </QueryClientProvider>,
      );
      expect(screen.queryByText(errorMessageText)).toBeNull();
    });

    it('download link normally', async () => {
      render(
        <QueryClientProvider client={new QueryClient()}>
          <DocumentViewer files={mockFiles} allowDownload />
        </QueryClientProvider>,
      );
      expect(screen.getByText(downloadLinkText)).toBeVisible();
    });

    it('show message on content error', async () => {
      render(
        <QueryClientProvider client={new QueryClient()}>
          <DocumentViewer files={mockErrorFiles} />
        </QueryClientProvider>,
      );
      expect(screen.getByText(errorMessageText)).toBeVisible();
    });

    it('download link on content error', async () => {
      render(
        <QueryClientProvider client={new QueryClient()}>
          <DocumentViewer files={mockErrorFiles} allowDownload />
        </QueryClientProvider>,
      );
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

      render(
        <QueryClientProvider client={new QueryClient()}>
          <DocumentViewer
            files={[
              { id: 99, filename: 'archive.zip', contentType: 'zip', url: 'path/to/archive.zip' },
              { id: 99, filename: 'archive.zip', contentType: 'zip', url: 'path/to/archive.zip' },
            ]}
            paymentRequestId="PaymentRequestId"
          />
        </QueryClientProvider>,
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

// describe('File upload status', () => {
//   const setup = async (fileStatus, isFileUploading = false) => {
//     await act(async () => {
//       render(<DocumentViewer files={mockFiles[0]} isFileUploading={isFileUploading} />);
//     });
//     act(() => {
//       switch (fileStatus) {
//         case UPLOAD_SCAN_STATUS.PROCESSING:
//           DocumentViewer.setFileStatus(UPLOAD_DOC_STATUS.SCANNING);
//           break;
//         case UPLOAD_SCAN_STATUS.CLEAN:
//           DocumentViewer.setFileStatus(UPLOAD_DOC_STATUS.ESTABLISHING);
//           break;
//         case UPLOAD_SCAN_STATUS.INFECTED:
//           DocumentViewer.setFileStatus(UPLOAD_DOC_STATUS.INFECTED);
//           break;
//         default:
//           break;
//       }
//     });
//   };

//   it('renders SCANNING status', () => {
//     setup(UPLOAD_SCAN_STATUS.PROCESSING);
//     expect(screen.getByText('Scanning')).toBeInTheDocument();
//   });

//   it('renders ESTABLISHING status', () => {
//     setup(UPLOAD_SCAN_STATUS.CLEAN);
//     expect(screen.getByText('Establishing Document for View')).toBeInTheDocument();
//   });

//   it('renders INFECTED status', () => {
//     setup(UPLOAD_SCAN_STATUS.INFECTED);
//     expect(screen.getByText('Ask for a new file')).toBeInTheDocument();
//   });
// });

// describe('DocumentViewer component', () => {
//   const files = [
//     {
//       id: '1',
//       createdAt: '2022-01-01T00:00:00Z',
//       contentType: 'application/pdf',
//       filename: 'file1.pdf',
//       url: samplePDF,
//     },
//   ];

//   beforeEach(() => {
//     global.EventSource = MockEventSource;
//   });

//   const renderComponent = (fileStatus) => {
//     render(
//       <QueryClientProvider client={new QueryClient()}>
//         <DocumentViewer files={files} allowDownload isFileUploading fileStatus={fileStatus} />
//       </QueryClientProvider>,
//     );
//   };

//   it('displays Uploading alert when fileStatus is UPLOADING', () => {
//     renderComponent(UPLOAD_DOC_STATUS.UPLOADING);
//     expect(screen.getByText(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.UPLOADING)).toBeInTheDocument();
//   });

//   it('displays Scanning alert when fileStatus is SCANNING', () => {
//     renderComponent(UPLOAD_DOC_STATUS.SCANNING);
//     expect(screen.getByText(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.SCANNING)).toBeInTheDocument();
//   });

//   it('displays Establishing Document for View alert when fileStatus is ESTABLISHING', () => {
//     renderComponent(UPLOAD_DOC_STATUS.ESTABLISHING);
//     expect(screen.getByText(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.ESTABLISHING_DOCUMENT_FOR_VIEW)).toBeInTheDocument();
//   });

//   it('displays File Not Found alert when selectedFile is null', () => {
//     render(<DocumentViewer files={[]} allowDownload />);
//     expect(screen.getByText(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.FILE_NOT_FOUND)).toBeInTheDocument();
//   });

//   it('displays an error alert when fileStatus is INFECTED', () => {
//     renderComponent(UPLOAD_SCAN_STATUS.INFECTED);
//     expect(
//       screen.getByText(
//         'Our antivirus software flagged this file as a security risk. Contact the service member. Ask them to upload a photo of the original document instead.',
//       ),
//     ).toBeInTheDocument();
//   });
// });

describe('DocumentViewer component', () => {
  const files = [
    {
      id: '1',
      createdAt: '2022-01-01T00:00:00Z',
      contentType: 'application/pdf',
      filename: 'file1.pdf',
      url: samplePDF,
    },
  ];
  beforeEach(() => {
    global.EventSource = MockEventSource;
  });

  const renderComponent = () => {
    render(<DocumentViewer files={files} allowDownload paymentRequestId={1234} isFileUploading />);
  };

  test('handles file processing status', async () => {
    renderComponent(UPLOAD_DOC_STATUS.UPLOADING);

    const eventSourceInstance = new MockEventSource(`/internal/uploads/${files[0].id}/status`, {
      withCredentials: true,
    });

    // Simulate different statuses
    await act(async () => {
      eventSourceInstance.sendMessage(UPLOAD_SCAN_STATUS.PROCESSING);
    });
    expect(screen.getByText(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.SCANNING)).toBeInTheDocument();

    await act(async () => {
      eventSourceInstance.sendMessage(UPLOAD_SCAN_STATUS.CLEAN);
    });
    expect(screen.getByText(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.ESTABLISHING_DOCUMENT_FOR_VIEW)).toBeInTheDocument();

    await act(async () => {
      eventSourceInstance.sendMessage(UPLOAD_SCAN_STATUS.INFECTED);
    });
    expect(
      screen.getByText(
        'Our antivirus software flagged this file as a security risk. Contact the service member. Ask them to upload a photo of the original document instead.',
      ),
    ).toBeInTheDocument();
  });

  it('displays File Not Found alert when no selectedFile', () => {
    render(<DocumentViewer files={[]} allowDownload />);
    expect(screen.getByText(UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.FILE_NOT_FOUND)).toBeInTheDocument();
  });
});
