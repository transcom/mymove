/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow } from 'enzyme';

import { OfficeWrapper } from './index';

describe('OfficeWrapper tests', () => {
  let wrapper;

  const mockOfficeProps = {
    getCurrentUserInfo: jest.fn(),
    loadInternalSchema: jest.fn(),
    loadPublicSchema: jest.fn(),
  };

  beforeEach(() => {
    wrapper = shallow(<OfficeWrapper {...mockOfficeProps} />);
  });

  it('renders without crashing or erroring', () => {
    const officeWrapper = wrapper.find('div');
    expect(officeWrapper).toBeDefined();
    expect(wrapper.find('SomethingWentWrong')).toHaveLength(0);
  });

  describe('if an error occurs', () => {
    it('renders the fail whale', () => {
      wrapper.setState({ hasError: true });
      expect(wrapper.find('SomethingWentWrong')).toHaveLength(1);
    });
  });
});
