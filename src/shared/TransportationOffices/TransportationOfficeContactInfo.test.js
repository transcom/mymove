import React from 'react';
import { shallow } from 'enzyme';
import { TransportationOfficeContactInfo } from './TransportationOfficeContactInfo';

describe('TransportationOfficeContactInfo tests', () => {
  let wrapper;
  it('renders without crashing', () => {
    const loadFn = jest.fn();
    const div = document.createElement('div');
    wrapper = shallow(
      <TransportationOfficeContactInfo dutyStation={{ id: '123' }} loadDutyStationTransportationOffice={loadFn} />,
      div,
    );
    expect(wrapper.find('div').length).toEqual(3);
  });
});
