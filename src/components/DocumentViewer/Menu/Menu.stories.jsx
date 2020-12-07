import React from 'react';

import DocViewerMenu from './Menu';

export default {
  title: 'Components/Document Viewer/Menu',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/8f32f4ab-cbe5-45f7-a0df-4ad19c0902a8?mode=design',
    },
  },
};

const testFiles = [
  {
    filename: 'Test File.pdf',
    fileType: 'pdf',
    filePath: '',
  },
  {
    filename: 'Test File 2.pdf',
    fileType: 'pdf',
    filePath: '',
  },
  {
    filename: 'Test File 3.pdf',
    fileType: 'pdf',
    filePath: '',
  },
  {
    filename: 'Test File 3 - A really long title that overflows with ellipsis.pdf',
    fileType: 'pdf',
    filePath: '',
  },
];

export const Menu = () => <DocViewerMenu isOpen files={testFiles} />;
