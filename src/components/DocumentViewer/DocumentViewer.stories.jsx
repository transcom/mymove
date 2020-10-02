import React from 'react';

import DocumentViewer from './DocumentViewer';
import pdf from './sample.pdf';
import pdf2 from './sample2.pdf';
import pdf3 from './sample3.pdf';
import jpg from './sample.jpg';
import png from './sample2.png';
import gif from './sample3.gif';

export default {
  title: 'Components|Document Viewer|Document Viewer',
};

const testPDFFiles = [
  {
    filename:
      'Test File - - A very long document title that overflow with ellipsis and shows a title when hovering; A very long document title that overflow with ellipsis and shows a title when hovering; A very long document title that overflow with ellipsis and shows a title when hovering; A very long document title that overflow with ellipsis and shows a title when hovering.pdf',
    contentType: 'pdf',
    url: pdf,
  },
  {
    filename: 'Test File 2.pdf',
    contentType: 'pdf',
    url: pdf2,
  },
  {
    filename: 'Test File 3.pdf',
    contentType: 'pdf',
    url: pdf3,
  },
];

const testImageFiles = [
  {
    filename: 'PCS Orders TACOMA Page 1',
    contentType: 'jpg',
    url: jpg,
  },
  {
    filename: 'PCS Orders TACOMA Page 2.png',
    contentType: 'png',
    url: png,
  },
  {
    filename: 'PCS Orders TACOMA Page 3.gif',
    contentType: 'gif',
    url: gif,
  },
];

export const PDFViewer = () => (
  <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
    <DocumentViewer files={testPDFFiles} />
  </div>
);

export const ImageViewer = () => (
  <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
    <DocumentViewer files={testImageFiles} />
  </div>
);
