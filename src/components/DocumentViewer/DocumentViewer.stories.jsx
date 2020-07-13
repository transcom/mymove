import React from 'react';

import DocumentViewer from './DocumentViewer';
import pdf from './sample.pdf';
import pdf2 from './sample2.pdf';
import pdf3 from './sample3.pdf';

export default {
  title: 'Components|Document Viewer|Document Viewer',
  parameters: { loki: { skip: true } },
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
    filePath: pdf2,
  },
  {
    filename: 'Test File 3.pdf',
    fileType: 'pdf',
    filePath: pdf3,
  },
];

export const PDFViewer = () => (
  <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
    <DocumentViewer files={testFiles} />
  </div>
);
