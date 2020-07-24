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
    expect(wrapper.find(Radio).at(0).text()).toContain('I’ll move things myself');
    expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('checked');

    expect(wrapper.find(Radio).at(1).text()).toContain('The government packs for me and moves me');
  });
});
