import React from 'react';
import { shallow } from 'enzyme';
import FailWhale from '.';

describe('FailWhale tests', () => {
  let wrapper;
  it('renders without crashing', () => {
    const div = document.createElement('div');
    wrapper = shallow(<FailWhale />, div);
    expect(wrapper.find('.usa-grid').length).toEqual(1);
  });
});
