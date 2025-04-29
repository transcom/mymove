import React from 'react';
import { act } from 'react-dom/test-utils';
import { mount } from 'enzyme';

import DebounceButton from './DebounceButton';

import { createMTOShipment } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  createMTOShipment: jest.fn().mockImplementation(() => Promise.resolve({ locator: 'ABC123' })),
}));

describe('DebounceButton component', () => {
  const defaultProps = {
    onClick: () => Promise.resolve(setTimeout(createMTOShipment(), 500)),
    delay: 2000,
  };

  it('renders the DebounceButton component without errors', () => {
    const wrapper = mount(<DebounceButton {...defaultProps}>Button Text</DebounceButton>);
    const btn = wrapper.find('[datatest-id="debounce-button"]');

    expect(btn).toBeTruthy();
  });

  it('multi-click calls fetcher once', () => {
    act(() => {
      const wrapper = mount(<DebounceButton {...defaultProps}>Click Me</DebounceButton>);
      const btn = wrapper.find('[datatest-id="debounce-button"]');
      btn.first().simulate('click');
      btn.first().simulate('click');
      btn.first().simulate('click');
      wrapper.update();
    });

    expect(createMTOShipment).toHaveBeenCalledTimes(1);
  });
});
