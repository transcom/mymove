import React from 'react';

import pdf from '../sample.pdf';

import DocViewerContent from './Content';

export default {
  title: 'Components/Document Viewer/Content',
  parameters: {
    happo: false,
  },
};

export const ContentArea = () => <DocViewerContent fileType="pdf" filePath={pdf} />;
