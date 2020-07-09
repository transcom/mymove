import React from 'react';

import DocViewerMenu from './Menu';

export default {
  title: 'Components|Document Viewer|Menu',
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
];

export const Menu = () => <DocViewerMenu isOpen files={testFiles} />;
