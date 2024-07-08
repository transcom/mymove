import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import DocumentViewerFileManager from './DocumentViewerFileManager';

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
});
