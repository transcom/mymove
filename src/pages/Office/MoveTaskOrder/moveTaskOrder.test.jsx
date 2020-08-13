import React from 'react';
import { mount } from 'enzyme';

import { MoveTaskOrder } from './moveTaskOrder';

describe('MoveTaskOrder', () => {
  const testId = 'test-id-123';
  const requiredProps = {
    match: { params: { moveTaskOrderId: testId } },
    history: { push: jest.fn() },
  };

  // eslint-disable-next-line react/jsx-props-no-spreading
  const wrapper = mount(<MoveTaskOrder {...requiredProps} />);

  it('should render the h1', () => {
    expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
  });
});
