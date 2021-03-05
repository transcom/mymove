import React from 'react';
import { mount } from 'enzyme';
import * as reactRedux from 'react-redux';

import requireCustomerState from './requireCustomerState';

import { MockProviders } from 'testUtils';

describe('requireCustomerState HOC', () => {
  const useDispatchMock = jest.spyOn(reactRedux, 'useDispatch');

  beforeEach(() => {
    useDispatchMock.mockClear();
  });

  const TestComponent = () => <div>My test component</div>;
  const TestComponentWithHOC = requireCustomerState(TestComponent);

  it('dispatches the initOnboarding action on mount', () => {
    const mockInitOnboarding = jest.fn();

    useDispatchMock.mockReturnValue(mockInitOnboarding);

    const wrapper = mount(
      <MockProviders>
        <TestComponentWithHOC />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockInitOnboarding).toHaveBeenCalled();
  });
});
