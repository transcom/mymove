/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import DocsUploaded from '.';

const defaultProps = {
  files: [],
};
function mountDocsUploaded(props = defaultProps) {
  return mount(<DocsUploaded {...props} />);
}
describe('DocsUploaded component', () => {
  it('renders document list with single file', () => {
    const props = {
      files: [{ filename: 'The fellowship of the file' }],
    };
    const wrapper = mountDocsUploaded(props);
    expect(wrapper.find('h6').text()).toBe('1 File uploaded');
    expect(wrapper.find('.doc-list-item').length).toBe(1);
  });

  it('renders document list with multiple files', () => {
    const props = {
      files: [{ filename: 'The twin files' }, { filename: 'The return of the file' }],
    };
    const wrapper = mountDocsUploaded(props);
    expect(wrapper.find('h6').text()).toBe('2 Files uploaded');
    expect(wrapper.find('.doc-list-item').length).toBe(2);
  });
});
