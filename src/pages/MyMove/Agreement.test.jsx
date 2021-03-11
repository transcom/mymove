/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import { Agreement } from './Agreement';

import * as internalApi from 'services/internalApi';

jest.mock('services/internalApi', () => {
  return {
    getPPMsForMove: jest.fn(() => Promise.resolve('PPM response')),
    submitMoveForApproval: jest.fn(),
  };
});

describe('Agreement page', () => {
  const testProps = {
    moveId: 'testMove123',
    setFlashMessage: jest.fn(),
    push: jest.fn(),
    updatePPMs: jest.fn(),
    updateMove: jest.fn(),
    ppmId: 'testPPM345',
  };

  const getPPMsMock = jest.spyOn(internalApi, 'getPPMsForMove');

  beforeEach(() => {
    getPPMsMock.mockClear();
  });

  it('loads PPMs on mount and stores them in Redux', async () => {
    const wrapper = mount(<Agreement {...testProps} />);
    expect(wrapper.exists()).toBe(true);
    expect(getPPMsMock).toHaveBeenCalled();
    await wrapper.update(); // wait for mock promise to resolve
    expect(testProps.updatePPMs).toHaveBeenCalledWith('PPM response');
  });
});
