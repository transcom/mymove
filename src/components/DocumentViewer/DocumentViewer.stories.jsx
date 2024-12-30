import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import DocumentViewer from './DocumentViewer';
import pdf from './sample.pdf';
import pdf2 from './sample2.pdf';
import pdf3 from './sample3.pdf';
import jpg from './sample.jpg';
import png from './sample2.png';
import gif from './sample3.gif';

const queryClient = new QueryClient();

const withQueryClient = (Story) => (
  <QueryClientProvider client={queryClient}>
    <Story />
  </QueryClientProvider>
);

export default {
  title: 'Components/Document Viewer/Document Viewer',
  decorators: [withQueryClient],
};

const testPDFFiles = [
  {
    filename:
      'Test File - - A very long document title that overflow with ellipsis and shows a title when hovering; A very long document title that overflow with ellipsis and shows a title when hovering; A very long document title that overflow with ellipsis and shows a title when hovering; A very long document title that overflow with ellipsis and shows a title when hovering.pdf',
    contentType: 'pdf',
    url: pdf,
    createdAt: '2021-06-15T15:09:26.979879Z',
  },
  {
    filename: 'Test File 2.pdf',
    contentType: 'pdf',
    url: pdf2,
    createdAt: '2021-06-15T15:09:26.979879Z',
  },
  {
    filename: 'Test File 3.pdf',
    contentType: 'pdf',
    url: pdf3,
    createdAt: '2021-06-15T15:09:26.979879Z',
  },
];

const testImageFiles = [
  {
    filename: 'PCS Orders TACOMA Page 1',
    contentType: 'jpg',
    url: jpg,
    createdAt: '2021-06-15T15:09:26.979879Z',
  },
  {
    filename: 'PCS Orders TACOMA Page 2.png',
    contentType: 'png',
    url: png,
    createdAt: '2021-06-15T15:09:26.979879Z',
  },
  {
    filename: 'PCS Orders TACOMA Page 3.gif',
    contentType: 'gif',
    url: gif,
    createdAt: '2021-06-15T15:09:26.979879Z',
  },
];

export const PDFViewer = () => (
  <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
    <QueryClientProvider client={new QueryClient()}>
      <DocumentViewer files={testPDFFiles} />
    </QueryClientProvider>
  </div>
);

export const ImageViewer = () => (
  <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
    <QueryClientProvider client={new QueryClient()}>
      <DocumentViewer files={testImageFiles} />
    </QueryClientProvider>
  </div>
);

export const DisplayDownloadOption = () => (
  <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
    <QueryClientProvider client={new QueryClient()}>
      <DocumentViewer files={testImageFiles} allowDownload />
    </QueryClientProvider>
  </div>
);
