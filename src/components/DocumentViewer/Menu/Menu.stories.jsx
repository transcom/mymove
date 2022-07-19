import React from 'react';

import DocViewerMenu from './Menu';

export default {
  title: 'Components/Document Viewer/Menu',
};

const testFiles = [
  {
    filename: 'Test File.pdf',
    fileType: 'pdf',
    filePath: '',
    createdAt: '2021-06-17T15:09:26.979879Z',
  },
  {
    filename: 'Test File 2.pdf',
    fileType: 'pdf',
    filePath: '',
    createdAt: '2021-06-16T15:09:26.979879Z',
  },
  {
    filename: 'Test File 3.pdf',
    fileType: 'pdf',
    filePath: '',
    createdAt: '2021-06-14T15:09:26.979879Z',
  },
  {
    filename: 'Test File 3 - A really long title that overflows with ellipsis.pdf',
    fileType: 'pdf',
    filePath: '',
    createdAt: '2021-06-12T15:09:26.979879Z',
  },
];

export const Menu = () => <DocViewerMenu isOpen files={testFiles} />;
