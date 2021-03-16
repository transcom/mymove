import React from 'react';
import { mount } from 'enzyme';
import { act } from 'react-dom/test-utils';
import { Alert } from '@trussworks/react-uswds';

import ConnectedFlashMessage, { FlashMessage } from './FlashMessage';

import { MockProviders } from 'testUtils';

describe('FlashMessage component', () => {
  it('renders an Alert', () => {
    const wrapper = mount(
      <FlashMessage
        flash={{
          type: 'success',
          message: 'This is a successful message!',
          key: 'TEST_SUCCESS_FLASH',
          slim: true,
        }}
        clearFlashMessage={jest.fn()}
      />,
    );

    expect(wrapper.find('FlashMessage').exists()).toBe(true);
    const alert = wrapper.find(Alert);
    expect(alert.exists()).toBe(true);
    expect(alert.prop('type')).toEqual('success');
    expect(alert.text()).toEqual('This is a successful message!');
    expect(alert.prop('slim')).toBe(true);
  });

  it('clears the flash message when unmounting', () => {
    const mockClearFlash = jest.fn();

    const wrapper = mount(
      <FlashMessage
        flash={{
          type: 'success',
          message: 'This is a successful message!',
          key: 'TEST_SUCCESS_FLASH',
        }}
        clearFlashMessage={mockClearFlash}
      />,
    );

    expect(mockClearFlash).toHaveBeenCalledTimes(0);
    act(() => {
      wrapper.unmount();
    });
    expect(mockClearFlash).toHaveBeenCalledTimes(1);
  });
});

describe('ConnectedFlashMessage component', () => {
  it('renders nothing if there is no flash message key set in Redux', () => {
    const testState = {
      flash: {
        flashMessage: {
          type: 'success',
          message: 'This is a successful message!',
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={testState}>
        <ConnectedFlashMessage />
      </MockProviders>,
    );

    expect(wrapper.find('FlashMessage').exists()).toBe(false);
  });

  it('renders an Alert if there is a flash message in Redux', () => {
    const testState = {
      flash: {
        flashMessage: {
          type: 'success',
          message: 'This is a successful message!',
          key: 'TEST_SUCCESS_FLASH',
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={testState}>
        <ConnectedFlashMessage />
      </MockProviders>,
    );

    expect(wrapper.find('FlashMessage').exists()).toBe(true);
    const alert = wrapper.find(Alert);
    expect(alert.exists()).toBe(true);
    expect(alert.prop('slim')).toBeUndefined();
    expect(wrapper.find('FlashMessage').children().length).toBe(1);
  });
});
