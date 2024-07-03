import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import DocumentViewerFileManager from './DocumentViewerFileManager';

import { MOVE_DOCUMENT_TYPE } from 'shared/constants';

const queryClient = new QueryClient();

const withQueryClient = (Story) => (
  <QueryClientProvider client={queryClient}>
    <Story />
  </QueryClientProvider>
);

export default {
  title: 'Components/DocumentViewerFileManager',
  component: DocumentViewerFileManager,
  decorators: [withQueryClient],
};

const Template = (args) => <DocumentViewerFileManager {...args} />;

export const OrdersDocument = Template.bind({});
OrdersDocument.args = {
  orderId: 'order-id',
  documentId: 'document-id',
  files: [{ id: 'file-1', name: 'File 1', filename: 'file1.pdf', bytes: 1024 }],
  documentType: MOVE_DOCUMENT_TYPE.ORDERS,
};

export const AmendedOrdersDocument = Template.bind({});
AmendedOrdersDocument.args = {
  orderId: 'order-id',
  documentId: 'document-id',
  files: [{ id: 'file-2', name: 'File 2', filename: 'file2.pdf', bytes: 2048 }],
  documentType: MOVE_DOCUMENT_TYPE.AMENDMENTS,
  updateAmendedDocument: () => {},
};

export const SupportingDocuments = Template.bind({});
SupportingDocuments.args = {
  move: { id: 'move-id', locator: 'move-locator' },
  orderId: 'order-id',
  documentId: 'document-id',
  files: [
    { id: 'file-3', name: 'File 3', filename: 'file3.jpg', bytes: 512 },
    { id: 'file-4', name: 'File 4', filename: 'file4.png', bytes: 1024 },
  ],
  documentType: MOVE_DOCUMENT_TYPE.SUPPORTING,
};
