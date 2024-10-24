import React from 'react';
import { render, screen, within, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import userEvent from '@testing-library/user-event';

import SupportingDocuments from './SupportingDocuments';

import { permissionTypes } from 'constants/permissions';
import { MockProviders } from 'testUtils';

beforeEach(() => {
  jest.clearAllMocks();
});

// prevents react-fileviewer from throwing errors without mocking relevant DOM elements
jest.mock('components/DocumentViewer/Content/Content', () => {
  const MockContent = () => <div>Content</div>;
  return MockContent;
});

const mockUploads = [
  {
    bytes: 1235030,
    contentType: 'image/png',
    createdAt: '2024-06-25T14:16:01.682Z',
    filename: 'some_doc.png',
    id: 'cfd9f68f-88e7-4bd1-a437-56f38a1e9c7f',
    status: 'PROCESSING',
    updatedAt: '2024-06-25T14:16:01.682Z',
    url: '/a/fake/path',
  },
  {
    bytes: 166930,
    contentType: 'application/pdf',
    createdAt: '2024-06-25T14:50:22.132Z',
    filename: 'some_other_doc.pdf',
    id: '1739c51a-f87e-4f09-a3ed-cf49dff31711',
    status: 'PROCESSING',
    updatedAt: '2024-06-25T14:50:22.132Z',
    url: '/a/fake/path',
  },
];

const mockProps = {
  move: {
    orderId: '1739c51a-f87e-4f09-a3ed-cf49dff31755',
    additionalDocuments: {
      id: '1739c51a-f87e-4f09-a3ed-cf49dff31712',
    },
  },
  uploads: mockUploads,
};

describe('Supporting Documents Viewer', () => {
  describe('displays viewer', () => {
    it('renders document viewer correctly on load', async () => {
      render(
        <QueryClientProvider client={new QueryClient()}>
          <SupportingDocuments {...mockProps} />
        </QueryClientProvider>,
      );
      const docMenuButton = await screen.findByRole('button', { name: /open menu/i });
      expect(docMenuButton).toBeInTheDocument();

      // We don't really have a better way to grab the DocumentViewerMenu to check its visibility because css isn't
      // loaded in the test environment. Instead, we'll grab it by its test id and check that it has the correct class.
      const docViewer = screen.getByTestId('DocViewerMenu');
      expect(docViewer).toHaveClass('collapsed');

      expect(within(docViewer).getByRole('heading', { level: 3, name: 'Documents' })).toBeInTheDocument();

      await userEvent.click(docMenuButton);

      expect(docViewer).not.toHaveClass('collapsed');

      const uploadList = within(docViewer).getByRole('list');
      expect(uploadList).toBeInTheDocument();

      expect(within(uploadList).getAllByRole('listitem').length).toBe(2);
      expect(within(uploadList).getByRole('button', { name: /some_other_doc\.pdf.*/i })).toBeInTheDocument();
      expect(within(uploadList).getByRole('button', { name: /some_doc\.png.*/i })).toBeInTheDocument();

      expect(screen.getByRole('button', { name: /close menu/i })).toBeInTheDocument();
      expect(screen.getByText('Download file')).toBeInTheDocument();
    });

    it('displays message if no files were uploaded', async () => {
      render(
        <QueryClientProvider client={new QueryClient()}>
          <SupportingDocuments {...mockProps} uploads={[]} />
        </QueryClientProvider>,
      );
      expect(screen.getByRole('heading', { name: /No supporting documents have been uploaded/i })).toBeInTheDocument();
      const docMenuButton = await screen.queryByRole('button', { name: /open menu/i });
      expect(docMenuButton).not.toBeInTheDocument();

      expect(screen.queryByTestId('DocViewerMenu')).not.toBeInTheDocument();

      expect(screen.queryByRole('button', { name: /close menu/i })).not.toBeInTheDocument();
      expect(screen.queryByText('Download file')).not.toBeInTheDocument();
    });

    it('displays message if uploads variable is undefined', async () => {
      render(
        <QueryClientProvider client={new QueryClient()}>
          <SupportingDocuments {...mockProps} uploads={undefined} />
        </QueryClientProvider>,
      );
      expect(screen.getByRole('heading', { name: /No supporting documents have been uploaded/i })).toBeInTheDocument();
      const docMenuButton = await screen.queryByRole('button', { name: /open menu/i });
      expect(docMenuButton).not.toBeInTheDocument();

      expect(screen.queryByTestId('DocViewerMenu')).not.toBeInTheDocument();

      expect(screen.queryByRole('button', { name: /close menu/i })).not.toBeInTheDocument();
      expect(screen.queryByText('Download file')).not.toBeInTheDocument();
    });

    it('displays message if uploads variable is not an array', async () => {
      render(
        <QueryClientProvider client={new QueryClient()}>
          <SupportingDocuments {...mockProps} uploads={1} />
        </QueryClientProvider>,
      );
      expect(screen.getByRole('heading', { name: /No supporting documents have been uploaded/i })).toBeInTheDocument();
      const docMenuButton = await screen.queryByRole('button', { name: /open menu/i });
      expect(docMenuButton).not.toBeInTheDocument();

      expect(screen.queryByTestId('DocViewerMenu')).not.toBeInTheDocument();

      expect(screen.queryByRole('button', { name: /close menu/i })).not.toBeInTheDocument();
      expect(screen.queryByText('Download file')).not.toBeInTheDocument();
    });

    it('displays document manager sidebar', async () => {
      render(
        <MockProviders permissions={[permissionTypes.createSupportingDocuments]}>
          <SupportingDocuments {...mockProps} />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(
          screen.queryByText(/PDF, JPG, or PNG only. Maximum file size 25MB. Each page must be clear and legible/),
        ).toBeInTheDocument();
      });
    });

    it('hides document manager sidebar', async () => {
      render(
        <QueryClientProvider client={new QueryClient()}>
          <SupportingDocuments {...mockProps} />
        </QueryClientProvider>,
      );

      expect(
        screen.queryByText(/PDF, JPG, or PNG only. Maximum file size 25MB. Each page must be clear and legible/),
      ).not.toBeInTheDocument();
    });
  });
});
