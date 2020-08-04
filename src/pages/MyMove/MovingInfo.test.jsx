import React from 'react';
import { mount } from 'enzyme';

import { MovingInfo } from './MovingInfo';

const wrapper = mount(<MovingInfo />);

describe('MovingInfo component', () => {
  it('renders', () => {
    expect(wrapper.exists('MovingInfo')).to.equal(true);
    // JUST CHECK FOR EXISTENCE OF
    // PAGE HEADING
    // SUBHEADERS
  });
});
