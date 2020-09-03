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
  parameters: { loki: { skip: true } },
};

const testPDFFiles = [
  {
    filename: 'Test File.pdf',
    fileType: 'pdf',
    filePath: pdf,
  },
  {
    filename: 'Test File 2.pdf',
    fileType: 'pdf',
    filePath: pdf2,
  },
  {
    filename: 'Test File 3.pdf',
    fileType: 'pdf',
    filePath: pdf3,
  },
];

const testImageFiles = [
  {
    filename: 'PCS Orders TACOMA.jpg',
    fileType: 'jpg',
    filePath: jpg,
  },
  {
    filename: 'PCS Orders TACOMA Page 2.png',
    fileType: 'png',
    filePath: png,
  },
  {
    filename: 'PCS Orders TACOMA Page 3.gif',
    fileType: 'gif',
    filePath: gif,
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
