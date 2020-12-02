import React from 'react';
import { mount } from 'enzyme';

import ConnectedFlashMessage, { FlashMessage } from './FlashMessage';

import { MockProviders, createMockHistory } from 'testUtils';

jest.mock('react-router-dom', () => ({
  __esModule: true,
  ...jest.requireActual('react-router-dom'),
}));

// Skipping these tests since I couldn't figure out how to change the react-router-dom mock for only some tests
// All of the test cases are repeated in the ConnectedFlashMessage block
describe.skip('FlashMessage component', () => {
  it('doesn’t crash if there is no flash object', () => {
    const wrapper = mount(<FlashMessage clearFlashMessage={jest.fn()} />);
    expect(wrapper.find('FlashMessage').exists()).toBe(true);
  });

  it('renders nothing if there is no flash message in Redux', () => {
    const wrapper = mount(
      <FlashMessage
        flash={{
          type: null,
          message: null,
          key: null,
        }}
        clearFlashMessage={jest.fn()}
      />,
    );

    expect(wrapper.find('FlashMessage').exists()).toBe(true);
    expect(wrapper.find('FlashMessage').children().length).toBe(0);
    expect(wrapper.find('FlashMessage').html()).toBe(null);
  });

  it('renders an Alert if there is a flash message in Redux', () => {
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

  it('does not render children if there is no flash message in Redux', () => {
    const wrapper = mount(
      <FlashMessage
        flash={{
          type: null,
          message: null,
          key: null,
        }}
        clearFlashMessage={jest.fn()}
      >
        This is my custom flash message
      </FlashMessage>,
    );

    expect(wrapper.find('FlashMessage').exists()).toBe(true);
    expect(wrapper.find('FlashMessage').children().length).toBe(0);
    expect(wrapper.find('FlashMessage').html()).toBe(null);
  });
});

describe('ConnectedFlashMessage component', () => {
  it('doesn’t crash if there is no flash object', () => {
    const testState = {};

    const wrapper = mount(
      <MockProviders initialState={testState}>
        <ConnectedFlashMessage />
      </MockProviders>,
    );

    expect(wrapper.find('FlashMessage').exists()).toBe(true);
  });

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

    expect(wrapper.find('FlashMessage').exists()).toBe(true);
    expect(wrapper.find('FlashMessage').children().length).toBe(0);
    expect(wrapper.find('FlashMessage').html()).toBe(null);
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

    expect(wrapper.find('FlashMessage').exists()).toBe(true);
    expect(wrapper.find('FlashMessage').children().length).toBe(0);
    expect(wrapper.find('FlashMessage').html()).toBe(null);
  });

  it('clears the flash if the pathname changes', () => {
    const testState = {
      flash: {
        flashMessage: {
          type: null,
          message: null,
          key: null,
        },
      },
    };

    const testHistory = createMockHistory(['/']);

    const wrapper = mount(
      <MockProviders initialState={testState} history={testHistory}>
        <ConnectedFlashMessage />
      </MockProviders>,
    );

    expect(wrapper.find('FlashMessage').exists()).toBe(true);
    expect(wrapper.find('FlashMessage').children().length).toBe(0);
    expect(wrapper.find('FlashMessage').html()).toBe(null);

    testHistory.push('/new-path');

    expect(wrapper.find('FlashMessage').children().length).toBe(0);
    expect(wrapper.find('FlashMessage').html()).toBe(null);
  });
});
