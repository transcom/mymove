import React from 'react';
import { mount } from 'enzyme';

import { MoveLanding } from 'pages/MyMove/MoveLanding';

describe('MoveLanding', () => {
  const minProps = {
    serviceMember: {
      first_name: 'Frida',
    },
  };
  it('Should render', () => {
    const wrapper = mount(<MoveLanding serviceMember={minProps.serviceMember} />);
    expect(wrapper.find('h1').length).toBe(1);
    expect(wrapper.find('h2').text()).toContain('Frida');
  });
});
