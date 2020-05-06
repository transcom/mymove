import React from 'react';
import { shallow } from 'enzyme';
import RequestedShipments from './RequestedShipments';

describe('RequestedShipments', () => {
  it('renders the container successfully', () => {
    const wrapper = shallow(<RequestedShipments />);
    expect(wrapper.find('div[data-cy="requested-shipments"]').exists()).toBe(true);
  });

  it('renders a child component passed to it', () => {
    const wrapper = shallow(
      <RequestedShipments>
        <>TESTING</>
      </RequestedShipments>,
    );
    expect(wrapper.find('div[data-cy="requested-shipments"]').text()).toContain('TESTING');
  });

  it('renders multiple child components passed to it', () => {
    const wrapper = shallow(
      <RequestedShipments>
        <>TESTING1</>
        <>TESTING2</>
      </RequestedShipments>,
    );
    expect(wrapper.find('div[data-cy="requested-shipments"]').text()).toContain('TESTING1');
    expect(wrapper.find('div[data-cy="requested-shipments"]').text()).toContain('TESTING2');
  });
});
