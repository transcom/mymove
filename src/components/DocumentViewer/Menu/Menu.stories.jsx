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
