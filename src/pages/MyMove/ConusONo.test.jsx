import React from 'react';
import { mount } from 'enzyme';
import { Radio } from '@trussworks/react-uswds';

import { ConusONo } from 'pages/MyMove/ConusONo';
import { CONUS_STATUS } from 'shared/constants';

describe('ConusONo', () => {
  const minProps = {
    conusStatus: CONUS_STATUS.CONUS,
  };
  it('should render radio buttons', () => {
    const wrapper = mount(<ConusONo conusStatus={minProps.conusStatus} />);
    expect(wrapper.find(Radio).length).toBe(2);

    // PPM button should be checked on page load
    expect(wrapper.find(Radio).at(0).text()).toContain('CONUS (continental US)');
    expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('checked');

    // HHG button should be disabled
    expect(wrapper.find(Radio).at(1).text()).toContain('OCONUS (Alaska, Hawaii, international)');
  });
});
