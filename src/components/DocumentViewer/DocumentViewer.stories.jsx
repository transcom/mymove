import React from 'react';

import DocumentViewer from './DocumentViewer';
import pdf from './sample.pdf';

export default {
  title: 'Components|Document Viewer',
};

export const PDFViewer = () => (
  <div style={{ height: '100vh' }}>
    <DocumentViewer filename="Test File.pdf" fileType="pdf" filePath={pdf} />
  </div>
);
