import React from 'react';
import { mount } from 'enzyme';

import { MoveLanding } from 'pages/MyMove/MoveLanding';
import serviceMemberBuilder from 'utils/test/factories/serviceMember';

describe('MoveLanding', () => {
  const serviceMember = serviceMemberBuilder({
    overrides: {
      first_name: 'Frida',
    },
  });
  const minProps = {
    serviceMember,
  };
  it('Should render', () => {
    const wrapper = mount(<MoveLanding serviceMember={minProps.serviceMember} />);
    expect(wrapper.find('h1').length).toBe(1);
    expect(wrapper.find('h2').text()).toContain('Frida');
  });
});
