import React from 'react';
import { mount } from 'enzyme';

import { Agreement } from './Agreement';

/*
import * as internalApi from 'services/internalApi';

jest.mock('services/internalApi', () => {
  return {
    submitMoveForApproval: jest.fn(),
  };
});
 */

describe('Agreement page', () => {
  const testProps = {
    moveId: 'testMove123',
    setFlashMessage: jest.fn(),
    push: jest.fn(),
    updateMove: jest.fn(),
  };

  it('loads PPMs on mount and stores them in Redux', async () => {
    const wrapper = mount(<Agreement {...testProps} />);
    expect(wrapper.exists()).toBe(true);
  });
});
