/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import FileUpload from './FileUpload';

describe('FileUpload component', () => {
  const testProps = {
    createUpload: jest.fn(() => Promise.resolve({ id: 'testFileId' })),
  };

  it('renders the FilePond component without errors', () => {
    const component = mount(<FileUpload {...testProps} />);
    expect(component.find('FilePond').length).toBe(1);
  });
});
