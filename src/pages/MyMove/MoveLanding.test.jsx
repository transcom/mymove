import React from 'react';
import { mount } from 'enzyme';

import ConnectedMoveLanding, { MoveLanding } from 'pages/MyMove/MoveLanding';
import { MockProviders } from 'testUtils';

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

describe('ConnectedMoveLanding', () => {
  it('renders without errors', () => {
    const initialState = {
      entities: {
        user: {
          userId123: {
            service_member: 'testServiceMemberId456',
          },
        },
        serviceMembers: {
          testServiceMemberId456: {
            first_name: 'Frida',
          },
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={initialState}>
        <ConnectedMoveLanding />
      </MockProviders>,
    );

    expect(wrapper.find('h2').text()).toContain('Frida');
  });
});
