import React from 'react';
import { shallow } from 'enzyme';

import ReviewServiceItems from 'components/Office/ReviewServiceItems/ReviewServiceItems';
import ServiceItemCard from 'components/Office/ReviewServiceItems/ServiceItemCard';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const serviceItemCards = [
  {
    id: '1',
    shipmentType: SHIPMENT_OPTIONS.HHG,
    serviceItemName: 'Domestic linehaul',
    amount: 6423,
    createdAt: '2020-01-01T00:08:00.999Z',
  },
];

describe('ReviewServiceItems component', () => {
  const component = shallow(<ReviewServiceItems serviceItemCards={serviceItemCards} handleClose={jest.fn()} />);

  it('renders without crashing', () => {
    expect(component.find('[data-testid="ReviewServiceItems"]').length).toBe(1);
  });

  it('renders ServiceItemCard component', () => {
    expect(component.find(ServiceItemCard).length).toBe(1);
  });
});
