/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow } from 'enzyme';

import samplePDF from '../sample.pdf';

import DocViewerContent from './Content';

const mockFile = {
  contentType: 'pdf',
  url: samplePDF,
  createdAt: '2021-06-15T15:09:26.979879Z',
};

describe('DocViewerContent', () => {
  const component = shallow(
    <DocViewerContent filename={mockFile.filename} fileType={mockFile.contentType} filePath={mockFile.url} />,
  );

  it('renders without crashing', () => {
    expect(component.find('[data-testid="DocViewerContent"]').length).toBe(1);
  });

  it('renders the FileViewer with the file props', () => {
    const fileViewer = component.find('FileViewer');
    expect(fileViewer.exists()).toBe(true);
    expect(fileViewer.prop('fileType')).toBe('pdf');
    expect(fileViewer.prop('filePath')).toBe(mockFile.url);
  });
});
