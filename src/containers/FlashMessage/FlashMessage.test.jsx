import React from 'react';
import { mount } from 'enzyme';

import ConnectedFlashMessage, { FlashMessage } from './FlashMessage';

import { MockProviders } from 'testUtils';

describe('FlashMessage component', () => {
  it('renders an Alert if there is no children', () => {
    const wrapper = mount(
      <FlashMessage
        flash={{
          type: 'success',
          message: 'This is a successful message!',
          key: 'TEST_SUCCESS_FLASH',
        }}
        clearFlashMessage={jest.fn()}
      />,
    );

    expect(wrapper.find('FlashMessage').exists()).toBe(true);
    expect(wrapper.find('Alert').exists()).toBe(true);
    expect(wrapper.find('FlashMessage').children().length).toBe(1);
  });

  it('renders children if there is a flash message in Redux and it uses children', () => {
    const wrapper = mount(
      <FlashMessage
        flash={{
          key: 'TEST_CUSTOM_FLASH',
        }}
        clearFlashMessage={jest.fn()}
      >
        This is my custom flash message
      </FlashMessage>,
    );

    expect(wrapper.find('FlashMessage').exists()).toBe(true);
    expect(wrapper.find('Alert').exists()).toBe(false);
    expect(wrapper.find('FlashMessage').text()).toBe('This is my custom flash message');
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
    wrapper.unmount();
    expect(mockClearFlash).toHaveBeenCalledTimes(1);
  });
});

describe('ConnectedFlashMessage component', () => {
  it('renders nothing if there is no flash message in Redux', () => {
    const testState = {
      flash: {
        flashMessage: {
          type: null,
          message: null,
          key: null,
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
    expect(wrapper.find('Alert').exists()).toBe(true);
    expect(wrapper.find('FlashMessage').children().length).toBe(1);
  });

  it('renders children if there is a flash message in Redux and it uses children', () => {
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
        <ConnectedFlashMessage>This is my custom flash message</ConnectedFlashMessage>
      </MockProviders>,
    );

    expect(wrapper.find('FlashMessage').exists()).toBe(true);
    expect(wrapper.find('Alert').exists()).toBe(false);
    expect(wrapper.find('FlashMessage').text()).toBe('This is my custom flash message');
  });

  it('does not render children if there is no flash message in Redux', () => {
    const testState = {
      flash: {
        flashMessage: {
          type: null,
          message: null,
          key: null,
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={testState}>
        <ConnectedFlashMessage>This is my custom flash message</ConnectedFlashMessage>
      </MockProviders>,
    );

    expect(wrapper.find('FlashMessage').exists()).toBe(false);
  });
});
