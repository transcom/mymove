/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow, mount } from 'enzyme';

import samplePDF from '../sample.pdf';
import sampleJPG from '../sample.jpg';
import samplePNG from '../sample2.png';
import sampleGIF from '../sample3.gif';

import DocViewerMenu from './Menu';

const mockFiles = [
  {
    id: 1,
    filename: 'Test File.pdf',
    contentType: 'application/pdf',
    url: samplePDF,
    createdAt: '2021-06-17T15:09:26.979879Z',
  },
  {
    id: 2,
    filename: 'Test File 2.jpg',
    contentType: 'image/jpeg',
    url: sampleJPG,
    createdAt: '2021-06-16T15:09:26.979879Z',
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
    createdAt: '2021-06-14T15:09:26.979879Z',
  },
];

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
      expect(component.find('[data-testid="button"]').length).toBe(4);
    });

    it('renders the file creation date', () => {
      expect(component.find('li button').at(0).text()).toBe('Test File.pdf  Uploaded on 17-Jun-2021');
    });

    it('selects a file when clicked', () => {
      component.find('li button').at(1).simulate('click');
      expect(mockProps.handleSelectFile).toHaveBeenCalledWith(1);
      component.find('li button').at(0).simulate('click');
      expect(mockProps.handleSelectFile).toHaveBeenCalledWith(0);
    });
  });
});
