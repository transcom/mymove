import React from 'react';
import { shallow } from 'enzyme';

import SomethingWentWrong from '.';

describe('SomethingWentWrong tests', () => {
  let wrapper;
  it('renders without crashing', () => {
    const div = document.createElement('div');
    wrapper = shallow(<SomethingWentWrong />, div);
    expect(wrapper.find('.usa-grid').length).toEqual(1);
  });
});
