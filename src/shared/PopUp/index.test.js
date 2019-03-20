import PopUp from './index';
import { shallow } from 'enzyme';
import React from 'react';

describe('PopUp', () => {
  const props = { alertMessage: 'ALERT!' };
  const alert = window.alert;

  beforeEach(() => {
    window.alert = jest.fn();
  });

  afterEach(() => {
    window.alert = alert;
  });

  it('displays alert with alertMessage on click events', () => {
    const wrapper = shallow(<PopUp {...props} />);
    wrapper.find('a').simulate('click', { preventDefault: jest.fn() });

    expect(window.alert).toHaveBeenCalledWith(props.alertMessage);
  });
});
