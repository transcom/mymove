import React from 'react';
import { mount } from 'enzyme';
import { Radio } from '@trussworks/react-uswds';

import { SelectMoveType } from 'pages/MyMove/SelectMoveType';

describe('SelectMoveType', () => {
  it('should render radio buttons', () => {
    // const wrapper = shallow(<SelectMoveType />);
    const wrapper = mount(<SelectMoveType />);
    expect(wrapper.find(Radio).length).toBe(2);

    // PPM button should be checked on page load
    expect(wrapper.find(Radio).at(0).text()).toContain('Arrange it all yourself');
    expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('checked');

    // HHG button should be disabled
    expect(wrapper.find(Radio).at(1).text()).toContain('Have professionals pack and move it all');
    expect(wrapper.find(Radio).at(1).find('.usa-radio__input').html()).toContain('disabled');
  });
});
