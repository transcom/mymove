import React from 'react';
import { mount, shallow } from 'enzyme';
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

  it('renders the button', () => {
    const wrapper = mount(
      <RequestedShipments>
        <>TESTING1</>
        <>TESTING2</>
      </RequestedShipments>,
    );
    expect(wrapper.find('button[data-testid="button"]').exists()).toBe(true);
    expect(wrapper.find('button[data-testid="button"]').text()).toContain('Approve selected shipments');
  });

  it('renders the button disabled', () => {
    const wrapper = mount(
      <RequestedShipments>
        <>TESTING1</>
        <>TESTING2</>
      </RequestedShipments>,
    );
    expect(wrapper.find('button[data-testid="button"]').html()).toContain('disabled=""');
  });

  it('renders the checkboxes', () => {
    const wrapper = mount(
      <RequestedShipments>
        <>TESTING1</>
        <>TESTING2</>
      </RequestedShipments>,
    );
    expect(wrapper.find('div[data-testid="checkbox"]').exists()).toBe(true);
    expect(wrapper.find('div[data-testid="checkbox"]').length).toEqual(2);
  });

  it('renders the checkboxes', () => {
    const wrapper = mount(
      <RequestedShipments>
        <>TESTING1</>
        <>TESTING2</>
      </RequestedShipments>,
    );
    expect(wrapper.find('div[data-testid="checkbox"]').exists()).toBe(true);
    expect(wrapper.find('div[data-testid="checkbox"]').length).toEqual(2);
  });
});
