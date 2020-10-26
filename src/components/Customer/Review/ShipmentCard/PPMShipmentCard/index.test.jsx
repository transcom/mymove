/*  react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import PPMShipmentCard from '.';

import { formatCustomerDate } from 'utils/formatters';

const defaultProps = {
  editPath: '',
  onEditClick: () => {},
  destinationZIP: '11111',
  estimatedWeight: '5,000',
  expectedDepartureDate: new Date('01/01/2020').toISOString(),
  shipmentId: 'ABC123K-001',
  sitDays: '24',
  originZIP: '00000',
};

function mountPPMShipmentCard(props = defaultProps) {
  return mount(<PPMShipmentCard {...props} />);
}
describe('PPMShipmentCard component', () => {
  it('renders component with all fields', () => {
    const wrapper = mountPPMShipmentCard();
    const tableHeaders = [
      'Expected departure',
      'Starting ZIP',
      'Storage (SIT)',
      'Destination ZIP',
      'Estimated weight',
      'Estimated incentive',
    ];
    const tableData = [
      formatCustomerDate(defaultProps.expectedDepartureDate),
      defaultProps.originZIP,
      `Yes, ${defaultProps.sitDays} days`,
      defaultProps.destinationZIP,
      `${defaultProps.estimatedWeight} lbs`,
      'Rate info unavailable',
    ];

    tableHeaders.forEach((label, index) => expect(wrapper.find('dt').at(index).text()).toBe(label));
    tableData.forEach((label, index) => expect(wrapper.find('dd').at(index).text()).toBe(label));
  });
});
