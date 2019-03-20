import CheckBox from './index';
import { shallow } from 'enzyme';
import React from 'react';

describe('Checkbox', () => {
  const props = { checked: true, onChange: jest.fn };
  it('checkbox displays as checked', function() {
    const wrapper = shallow(<CheckBox {...props} />);
    const checkbox = wrapper.find({ type: 'checkbox' });

    expect(checkbox.props().checked).toBe(true);
  });

  it('calls onChangeHandler on change events', function() {
    const onChangeHandler = jest.fn();
    const wrapper = shallow(<CheckBox {...props} onChangeHandler={onChangeHandler} />);
    wrapper.find({ type: 'checkbox' }).simulate('change', { target: { checked: true } });

    expect(onChangeHandler).toHaveBeenCalledWith(true);
  });
});
