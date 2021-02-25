/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import SubmitMoveForm from './SubmitMoveForm';

describe('SubmitMoveForm component', () => {
  let wrapper;
  let onSubmit;
  let onPrint;

  beforeEach(() => {
    onSubmit = jest.fn();
    onPrint = jest.fn();
    wrapper = mount(<SubmitMoveForm onSubmit={onSubmit} onPrint={onPrint} />);
  });

  it('renders the default state', () => {
    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('input[name="signature"]').length).toBe(1);
    expect(wrapper.find('input[name="date"]').length).toBe(1);
    expect(wrapper.find('input[name="date"]').prop('disabled')).toBe(true);
    expect(wrapper.find('button[data-testid="wizardCompleteButton"]').length).toBe(1);
  });
});
