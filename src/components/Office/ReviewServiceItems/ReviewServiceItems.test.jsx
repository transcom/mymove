import React from 'react';
import { shallow } from 'enzyme';

import ReviewServiceItems from './ReviewServiceItems';
import ServiceItemCard from './ServiceItemCard';

import { SHIPMENT_OPTIONS, SERVICE_ITEM_STATUS } from 'shared/constants';

const serviceItemCards = [
  {
    id: '1',
    shipmentType: SHIPMENT_OPTIONS.HHG,
    shipmentId: '10',
    serviceItemName: 'Domestic linehaul',
    amount: 6423,
    status: SERVICE_ITEM_STATUS.SUBMITTED,
    createdAt: '2020-01-01T00:08:00.999Z',
  },
  {
    id: '2',
    shipmentType: SHIPMENT_OPTIONS.HHG,
    shipmentId: '10',
    serviceItemName: 'Fuel Surcharge',
    amount: 50.25,
    createdAt: '2020-01-01T00:08:00.999Z',
  },
  {
    id: '3',
    shipmentType: SHIPMENT_OPTIONS.NTS,
    shipmentId: '20',
    serviceItemName: 'Domestic linehaul',
    amount: 0.1,
    createdAt: '2020-01-01T00:09:00.999Z',
  },
  {
    id: '4',
    shipmentType: null,
    shipmentId: null,
    serviceItemName: 'Counseling Services',
    amount: 1000,
    createdAt: '2020-01-01T00:01:00.999Z',
  },
  {
    id: '5',
    shipmentType: null,
    shipmentId: null,
    serviceItemName: 'Move management',
    amount: 1,
    createdAt: '2020-01-01T00:01:00.999Z',
  },
];

describe('ReviewServiceItems component', () => {
  const component = shallow(<ReviewServiceItems serviceItemCards={serviceItemCards} handleClose={jest.fn()} />);

  it('renders without crashing', () => {
    expect(component.find('[data-cy="ReviewServiceItems"]').length).toBe(1);
  });

  it('renders ServiceItemCard component', () => {
    expect(component.find(ServiceItemCard).length).toBe(1);
  });
});
