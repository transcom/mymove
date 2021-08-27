/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { act } from 'react-dom/test-utils';
import { mount } from 'enzyme';

import DocumentViewer from './DocumentViewer';
import samplePDF from './sample.pdf';
import sampleJPG from './sample.jpg';
import samplePNG from './sample2.png';
import sampleGIF from './sample3.gif';

const mockFiles = [
  {
    id: 1,
    filename: 'Test File.pdf',
    contentType: 'application/pdf',
    url: samplePDF,
    createdAt: '2021-06-14T15:09:26.979879Z',
  },
  {
    id: 2,
    filename: 'Test File 2.jpg',
    contentType: 'image/jpeg',
    url: sampleJPG,
    createdAt: '2021-06-12T15:09:26.979879Z',
  },
  {
    id: 3,
    filename: 'Test File 3.png',
    contentType: 'image/png',
    url: samplePNG,
    createdAt: '2021-06-15T15:09:26.979879Z',
  },
  {
    id: 4,
    filename: 'Test File 4.gif',
    contentType: 'image/gif',
    url: sampleGIF,
    createdAt: '2021-06-16T15:09:26.979879Z',
  },
];

describe('DocumentViewer component', () => {
  const component = mount(<DocumentViewer files={mockFiles} />);
  const content = component.find('DocViewerContent');
  const menu = component.find('DocViewerMenu');

  it('initial state is closed menu and first file selected', () => {
    expect(menu.prop('isOpen')).toBe(false);
    expect(menu.prop('selectedFileIndex')).toBe(0);
    expect(content.prop('filePath')).toBe(mockFiles[0].url);
  });

  it('renders DocViewerContent and DocViewerMenu with the correct props', () => {
    expect(content.length).toBe(1);
    expect(menu.length).toBe(1);
    expect(menu.prop('files')).toBe(mockFiles);
  });

  it('renders the file creation date with the correctly sorted props', () => {
    expect(component.find('li button').at(0).text()).toBe('Test File 4.gif  Uploaded on 16-Jun-2021');
  });

  it('renders the title bar with the correct props', () => {
    expect(component.find('[data-testid="documentTitle"]').text()).toBe('Test File 4.gif - Added on 16 Jun 2021');
  });

  it('handles the open menu button', () => {
    act(() => {
      component.find('button[data-testid="openMenu"]').prop('onClick')();
    });
    component.update();
    expect(component.find('DocViewerMenu').prop('isOpen')).toBe(true);
  });

  it('handles the close menu button', () => {
    act(() => {
      component.find('button[data-testid="closeMenu"]').prop('onClick')();
    });
    component.update();
    expect(component.find('DocViewerMenu').prop('isOpen')).toBe(false);
  });

  it('handles selecting a different file', () => {
    act(() => {
      component.find('button[data-testid="openMenu"]').prop('onClick')();
      menu.find('li button').at(1).simulate('click');
    });
    component.update();
    expect(component.find('DocViewerMenu').prop('isOpen')).toBe(false);
    expect(component.find('DocViewerMenu').prop('selectedFileIndex')).toBe(1);
    expect(component.find('DocViewerContent').prop('filePath')).toBe(mockFiles[1].url);
    expect(component.find('DocViewerContent').prop('fileType')).toBe('png');
    expect(component.find('.unsupported-message').exists()).toBe(false);

    act(() => {
      component.find('button[data-testid="openMenu"]').prop('onClick')();
      menu.find('li button').at(2).simulate('click');
    });
    component.update();
    expect(component.find('DocViewerMenu').prop('isOpen')).toBe(false);
    expect(component.find('DocViewerMenu').prop('selectedFileIndex')).toBe(2);
    expect(component.find('DocViewerContent').prop('filePath')).toBe(mockFiles[2].url);
    expect(component.find('DocViewerContent').prop('fileType')).toBe('pdf');
    expect(component.find('.unsupported-message').exists()).toBe(false);

    act(() => {
      component.find('button[data-testid="openMenu"]').prop('onClick')();
      menu.find('li button').at(3).simulate('click');
    });
    component.update();
    expect(component.find('DocViewerMenu').prop('isOpen')).toBe(false);
    expect(component.find('DocViewerMenu').prop('selectedFileIndex')).toBe(3);
    expect(component.find('DocViewerContent').prop('filePath')).toBe(mockFiles[3].url);
    expect(component.find('DocViewerContent').prop('fileType')).toBe('jpg');
    expect(component.find('.unsupported-message').exists()).toBe(false);
  });

  it('shows error if file type is unsupported', () => {
    const wrapper = mount(
      <DocumentViewer files={[{ id: 99, filename: 'archive.zip', contentType: 'zip', url: '/path/to/archive.zip' }]} />,
    );
    expect(wrapper.find('.unsupported-message').text()).toEqual('.zip is not supported.');
  });

  it('displays file not found for empty files array', () => {
    const wrapper = mount(<DocumentViewer />);
    expect(wrapper.find('h2').text()).toEqual('File Not Found');
  });
});
