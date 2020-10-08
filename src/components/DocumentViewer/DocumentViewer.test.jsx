/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { act } from 'react-dom/test-utils';
import { shallow, mount } from 'enzyme';

import DocViewerContent from './Content/Content';
import DocViewerMenu from './Menu/Menu';
import DocumentViewer from './DocumentViewer';
import samplePDF from './sample.pdf';
import sampleJPG from './sample.jpg';
import samplePNG from './sample2.png';
import sampleGIF from './sample3.gif';

const mockFile = {
  contentType: 'pdf',
  url: samplePDF,
};

const mockFiles = [
  {
    filename: 'Test File.pdf',
    contentType: 'application/pdf',
    url: samplePDF,
  },
  {
    filename: 'Test File 2.jpg',
    contentType: 'image/jpeg',
    url: sampleJPG,
  },
  {
    filename: 'Test File 3.png',
    contentType: 'image/png',
    url: samplePNG,
  },
  {
    filename: 'Test File 4.gif',
    contentType: 'image/gif',
    url: sampleGIF,
  },
];

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
      component.find('button[data-testid="closeMenu"]').simulate('click');
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
    expect(content.prop('filePath')).toBe(mockFiles[0].url);
  });

  it('renders DocViewerContent and DocViewerMenu with the correct props', () => {
    expect(content.length).toBe(1);
    expect(menu.length).toBe(1);
    expect(menu.prop('files')).toBe(mockFiles);
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
    expect(component.find('DocViewerContent').prop('fileType')).toBe('jpg');
    expect(component.find('.unsupported-message').exists()).toBe(false);

    act(() => {
      component.find('button[data-testid="openMenu"]').prop('onClick')();
      menu.find('li button').at(2).simulate('click');
    });
    component.update();
    expect(component.find('DocViewerMenu').prop('isOpen')).toBe(false);
    expect(component.find('DocViewerMenu').prop('selectedFileIndex')).toBe(2);
    expect(component.find('DocViewerContent').prop('filePath')).toBe(mockFiles[2].url);
    expect(component.find('DocViewerContent').prop('fileType')).toBe('png');
    expect(component.find('.unsupported-message').exists()).toBe(false);

    act(() => {
      component.find('button[data-testid="openMenu"]').prop('onClick')();
      menu.find('li button').at(3).simulate('click');
    });
    component.update();
    expect(component.find('DocViewerMenu').prop('isOpen')).toBe(false);
    expect(component.find('DocViewerMenu').prop('selectedFileIndex')).toBe(3);
    expect(component.find('DocViewerContent').prop('filePath')).toBe(mockFiles[3].url);
    expect(component.find('DocViewerContent').prop('fileType')).toBe('gif');
    expect(component.find('.unsupported-message').exists()).toBe(false);
  });

  it('shows error if file type is unsupported', () => {
    const wrapper = mount(
      <DocumentViewer files={[{ filename: 'archive.zip', contentType: 'zip', url: '/path/to/archive.zip' }]} />,
    );
    expect(wrapper.find('.unsupported-message').text()).toEqual('.zip is not supported.');
  });

  it('displays file not found for empty files array', () => {
    const wrapper = mount(<DocumentViewer />);
    expect(wrapper.find('h2').text()).toEqual('File Not Found');
  });
});
