import React from 'react';
import { shallow } from 'enzyme';
import { Radio } from '@trussworks/react-uswds';

import { SelectMoveType } from 'pages/MyMove/SelectMoveType';

describe('SelectMoveType', () => {
  it('should render radio buttons', () => {
    const wrapper = shallow(<SelectMoveType />);
    expect(wrapper.find(Radio).length).toBe(2);

    // PPM button should be checked on page load
    expect(wrapper.find(Radio).at(0).dive().text()).toContain('Arrange it all yourself');
    expect(wrapper.find(Radio).at(0).dive().find('.usa-radio__input').html()).toContain('checked');

    // HHG button should be disabled
    expect(wrapper.find(Radio).at(1).dive().text()).toContain('Have professionals pack and move it all');
    expect(wrapper.find(Radio).at(1).dive().find('.usa-radio__input').html()).toContain('disabled');
  });
});
