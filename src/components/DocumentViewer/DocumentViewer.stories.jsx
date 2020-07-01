import React from 'react';

import DocumentViewer from './DocumentViewer';
import pdf from './sample.pdf';

export default {
  title: 'Components|Document Viewer',
};

const testFiles = [
  {
    filename: 'Test File.pdf',
    fileType: 'pdf',
    filePath: pdf,
  },
  {
    filename: 'Test File 2.pdf',
    fileType: 'pdf',
    filePath: pdf,
  },
  {
    filename: 'Test File 3.pdf',
    fileType: 'pdf',
    filePath: pdf,
  },
];

export const PDFViewer = () => (
  <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
    <DocumentViewer files={testFiles} />
  </div>
);
