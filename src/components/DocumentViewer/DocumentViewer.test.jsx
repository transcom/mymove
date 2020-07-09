/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { act } from 'react-dom/test-utils';
import { shallow, mount } from 'enzyme';

import DocViewerContent from './Content/Content';
import DocViewerMenu from './Menu/Menu';
import DocumentViewer from './DocumentViewer';
import samplePDF from './sample.pdf';
import samplePDF2 from './sample2.pdf';
import samplePDF3 from './sample3.pdf';

const mockFile = {
  fileType: 'pdf',
  filePath: samplePDF,
};

const mockFiles = [
  {
    filename: 'Test File.pdf',
    fileType: 'pdf',
    filePath: samplePDF,
  },
  {
    filename: 'Test File 2.pdf',
    fileType: 'pdf',
    filePath: samplePDF2,
  },
  {
    filename: 'Test File 3.pdf',
    fileType: 'pdf',
    filePath: samplePDF3,
  },
];

describe('DocViewerContent', () => {
  const component = shallow(<DocViewerContent {...mockFile} />);

  it('renders without crashing', () => {
    expect(component.find('[data-testid="DocViewerContent"]').length).toBe(1);
  });

  it('renders the FileViewer with the file props', () => {
    const fileViewer = component.find('FileViewer');
    expect(fileViewer.exists()).toBe(true);
    expect(fileViewer.prop('fileType')).toBe(mockFile.fileType);
    expect(fileViewer.prop('filePath')).toBe(mockFile.filePath);
  });
});

describe('DocViewerMenu', () => {
  const mockProps = {
    handleClose: jest.fn(),
    handleSelectFile: jest.fn(),
  };

  describe('closed state', () => {
    const component = shallow(<DocViewerMenu isOpen={false} files={mockFiles} {...mockProps} />);

    it('has the collapsed class', () => {
      expect(component.find('[data-testid="DocViewerMenu"]').hasClass('collapsed')).toBe(true);
    });
  });

  describe('open state', () => {
    const component = shallow(<DocViewerMenu isOpen files={mockFiles} {...mockProps} />);

    it('does not have the collapsed class', () => {
      expect(component.find('[data-testid="DocViewerMenu"]').hasClass('collapsed')).toBe(false);
    });
  });

  describe('close button', () => {
    const component = mount(<DocViewerMenu files={mockFiles} {...mockProps} />);
    it('implements the close handler', () => {
      component.find('[data-testid="closeMenu"]').simulate('click');
      expect(mockProps.handleClose).toHaveBeenCalled();
    });
  });

  describe('file list', () => {
    const component = mount(<DocViewerMenu files={mockFiles} {...mockProps} />);

    it('renders without crashing', () => {
      expect(component.find('[data-testid="DocViewerMenu"]').length).toBe(1);
    });

    it('renders the list of files', () => {
      mockFiles.forEach((file) => {
        expect(component.contains(<p>{file.filename}</p>)).toBe(true);
      });
    });

    it('selects a file when clicked', () => {
      component.find('li button').at(1).simulate('click');
      expect(mockProps.handleSelectFile).toHaveBeenCalledWith(1);
      component.find('li button').at(0).simulate('click');
      expect(mockProps.handleSelectFile).toHaveBeenCalledWith(0);
    });
  });
});

describe('DocumentViewer component', () => {
  const component = mount(<DocumentViewer files={mockFiles} />);
  const content = component.find('DocViewerContent');
  const menu = component.find('DocViewerMenu');

  it('initial state is closed menu and first file selected', () => {
    expect(menu.prop('isOpen')).toBe(false);
    expect(menu.prop('selectedFileIndex')).toBe(0);
    expect(content.prop('filePath')).toBe(mockFiles[0].filePath);
  });

  it('renders DocViewerContent and DocViewerMenu with the correct props', () => {
    expect(content.length).toBe(1);
    expect(menu.length).toBe(1);
    expect(menu.prop('files')).toBe(mockFiles);
  });

  it('handles the open menu button', () => {
    act(() => {
      component.find('[data-testid="openMenu"]').prop('onClick')();
    });
    component.update();
    expect(component.find('DocViewerMenu').prop('isOpen')).toBe(true);
  });

  it('handles the close menu button', () => {
    act(() => {
      component.find('[data-testid="closeMenu"]').prop('onClick')();
    });
    component.update();
    expect(component.find('DocViewerMenu').prop('isOpen')).toBe(false);
  });

  it('handles selecting a different file', () => {
    act(() => {
      component.find('[data-testid="openMenu"]').prop('onClick')();
      menu.find('li button').at(2).simulate('click');
    });
    component.update();
    expect(component.find('DocViewerMenu').prop('isOpen')).toBe(false);
    expect(component.find('DocViewerMenu').prop('selectedFileIndex')).toBe(2);
    expect(component.find('DocViewerContent').prop('filePath')).toBe(mockFiles[2].filePath);
  });
});
