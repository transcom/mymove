import React from 'react';

import DocumentViewer from './DocumentViewer';
import pdf from './sample.pdf';

export default {
  title: 'Components|Document Viewer',
};

export const PDFViewer = () => <DocumentViewer filename="Test File.pdf" fileType="pdf" filePath={pdf} />;
