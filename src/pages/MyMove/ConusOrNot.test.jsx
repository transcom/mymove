import React from 'react';
import { mount } from 'enzyme';
import { Radio } from '@trussworks/react-uswds';

import { ConusOrNot } from 'pages/MyMove/ConusOrNot';
import { CONUS_STATUS } from 'shared/constants';

describe('ConusOrNot', () => {
  const minProps = {
    conusStatus: CONUS_STATUS.CONUS,
    setLocation: () => {},
  };
  it('should render radio buttons', () => {
    //  react/jsx-props-no-spreading
    const wrapper = mount(<ConusOrNot {...minProps} />);
    expect(wrapper.find(Radio).length).toBe(2);

    // PPM button should be checked on page load
    expect(wrapper.find(Radio).at(0).text()).toContain('CONUS');
    expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('checked');

    // HHG button should be disabled
    expect(wrapper.find(Radio).at(1).text()).toContain('OCONUS');
  });
});
