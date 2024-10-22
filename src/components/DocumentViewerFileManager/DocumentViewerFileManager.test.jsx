import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import DocumentViewerFileManager from './DocumentViewerFileManager';

import { deleteUploadForDocument } from 'services/ghcApi';

jest.mock('services/ghcApi', () => ({
  createUploadForDocument: jest.fn(),
  createUploadForAmdendedOrders: jest.fn(),
  createUploadForSupportingDocuments: jest.fn(),
  deleteUploadForDocument: jest.fn(),
  getOrder: jest.fn(),
}));

jest.mock('components/FileUpload/FileUpload', () => ({ onChange }) => (
  <div>
    <button type="button" onClick={() => onChange()}>
      Drag files here or click to upload
    </button>
  </div>
));

jest.mock('components/UploadsTable/UploadsTable', () => ({ uploads, onDelete }) => (
  <div>
    {uploads.map((upload) => (
      <div key={upload.id}>
        <span>{upload.name}</span>
        <button type="button" onClick={() => onDelete(upload.id)}>
          Delete
        </button>
      </div>
    ))}
  </div>
));

jest.mock(
  'components/ConfirmationModals/DeleteDocumentFileConfirmationModal',
  () =>
    ({ isOpen, closeModal, submitModal, fileInfo }) =>
      isOpen ? (
        <div>
          <button type="button" onClick={closeModal}>
            No, keep it
          </button>
          <button type="submit" onClick={submitModal}>
            Yes, delete
          </button>
          <div>
            <p>{fileInfo.filename}</p>
            <p>{fileInfo.bytes} bytes</p>
          </div>
        </div>
      ) : null,
);

jest.mock('components/Hint', () => ({ children }) => <div>{children}</div>);

// Helper to render with React Query
const renderWithQueryClient = (ui) => {
  const queryClient = new QueryClient();
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
};

// Mock DataTransfer for the testing environment
global.DataTransfer = class {
  constructor() {
    this.items = [];
  }

  // Add a file to the DataTransfer
  add(file) {
    this.items.push({ kind: 'file', getAsFile: () => file });
  }
};

describe('DocumentViewerFileManager', () => {
  const defaultProps = {
    className: 'test-class',
    move: { id: 'move-id', locator: 'move-locator' },
    orderId: 'order-id',
    documentId: 'document-id',
    files: [{ id: 'file-1', name: 'File 1' }],
    documentType: 'ORDERS',
    updateAmendedDocument: jest.fn(),
  };

  it('renders without crashing', () => {
    renderWithQueryClient(<DocumentViewerFileManager {...defaultProps} />);
    expect(screen.getByText('Manage Orders')).toBeInTheDocument();
  });

  it('shows upload section when Manage Orders button is clicked', () => {
    renderWithQueryClient(<DocumentViewerFileManager {...defaultProps} />);
    fireEvent.click(screen.getByText('Manage Orders'));
    expect(screen.getByText('Drag files here or click to upload')).toBeInTheDocument();
  });

  it('opens delete confirmation modal when delete button is clicked', () => {
    renderWithQueryClient(<DocumentViewerFileManager {...defaultProps} />);
    fireEvent.click(screen.getByText('Manage Orders'));
    fireEvent.click(screen.getByText('Delete'));
    expect(screen.getByText('Yes, delete')).toBeInTheDocument();
  });

  it('closes delete confirmation modal when close button is clicked', () => {
    renderWithQueryClient(<DocumentViewerFileManager {...defaultProps} />);
    fireEvent.click(screen.getByText('Manage Orders'));
    fireEvent.click(screen.getByText('Delete'));
    fireEvent.click(screen.getByText('No, keep it'));
    expect(screen.queryByText('Yes, delete')).not.toBeInTheDocument();
  });

  it('handles file upload change', () => {
    renderWithQueryClient(<DocumentViewerFileManager {...defaultProps} />);
    fireEvent.click(screen.getByText('Manage Orders'));
    fireEvent.click(screen.getByText('Drag files here or click to upload'));
    expect(screen.queryByText('An error occurred')).not.toBeInTheDocument();
  });

  it('should handle file uploads correctly', async () => {
    const mockUpdateAmendedDocument = jest.fn(); // Mock function to verify the update action
    const queryClient = new QueryClient(); // Create a new instance of QueryClient for React Query

    // Render the DocumentViewerFileManager component within the QueryClientProvider
    render(
      <QueryClientProvider client={queryClient}>
        <DocumentViewerFileManager
          orderId="123" // Sample order ID
          documentId="456" // Sample document ID
          files={[]} // Initialize with an empty array for uploaded files
          documentType="ORDERS" // Set the document type to trigger the appropriate UI
          updateAmendedDocument={mockUpdateAmendedDocument} // Pass the mock function to handle updates
        />
      </QueryClientProvider>,
    );

    // Verify that the "Manage Orders" button is rendered in the document
    const manageDocumentsButton = screen.getByRole('button', { name: /manage orders/i });
    expect(manageDocumentsButton).toBeInTheDocument();

    // Simulate a user clicking the "Manage Orders" button to display the file upload section
    fireEvent.click(manageDocumentsButton);

    // Confirm that the upload area (drag-and-drop zone) is present in the document
    const uploadArea = await screen.findByRole('button', { name: /drag files here or click to upload/i });
    expect(uploadArea).toBeInTheDocument();

    // Create a new File object to simulate a file upload
    const file = new File(['content'], 'testfile.txt', { type: 'text/plain' });

    // Create a DataTransfer object to mimic the file being dragged and dropped
    const dataTransfer = new DataTransfer();
    dataTransfer.add(file); // Add the simulated file to the DataTransfer object

    // Simulate the drag-and-drop events
    fireEvent.dragEnter(uploadArea, {
      dataTransfer, // Trigger the drag enter event on the upload area
    });
    fireEvent.dragOver(uploadArea, {
      dataTransfer, // Trigger the drag over event on the upload area
    });
    fireEvent.drop(uploadArea, {
      dataTransfer, // Trigger the drop event to simulate the file upload
    });
  });

  it('calls deleteUploadForDocument when handleDeleteSubmit is triggered', async () => {
    deleteUploadForDocument.mockResolvedValueOnce({});

    renderWithQueryClient(<DocumentViewerFileManager {...defaultProps} />);

    fireEvent.click(screen.getByText('Manage Orders'));
    fireEvent.click(screen.getByText('Delete'));

    fireEvent.click(screen.getByText('Yes, delete'));

    expect(deleteUploadForDocument).toHaveBeenCalledWith('file-1', 'order-id');
  });
});
