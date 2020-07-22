import React from 'react';
import { mount } from 'enzyme';

import { UnsupportedMove } from 'pages/MyMove/UnsupportedMove';

describe('UnsupportedMove', () => {
  const minProps = {
    serviceMember: {
      current_station: {
        name: 'Yuma AFB',
      },
    },
  };
  it('Should render', () => {
    const wrapper = mount(<UnsupportedMove serviceMember={minProps.serviceMember} />);
    expect(wrapper.find('h1').length).toBe(1);
    expect(wrapper.find('p').text()).toContain('Yuma');
  });
});
