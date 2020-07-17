import React from 'react';
import { shallow, mount } from 'enzyme';

import { toDollarString } from '../../../shared/formatters';

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
    createdAt: '2020-01-01T00:08:30.999Z',
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
    createdAt: '2020-01-01T00:02:00.999Z',
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

const compareItem = (component, item) => {
  expect(component.find('[data-cy="serviceItemName"]').text()).toEqual(item.serviceItemName);
  expect(component.find('[data-cy="serviceItemAmount"]').text()).toEqual(toDollarString(item.amount));
};

describe('ReviewServiceItems component', () => {
  const handleClose = jest.fn();
  const shallowComponent = shallow(
    <ReviewServiceItems serviceItemCards={serviceItemCards} handleClose={handleClose} />,
  );
  const mountedComponent = mount(<ReviewServiceItems serviceItemCards={serviceItemCards} handleClose={handleClose} />);

  it('renders without crashing', () => {
    const serviceItemCard = shallowComponent.find(ServiceItemCard).dive();
    expect(serviceItemCard.find('[data-testid="ServiceItemCard"]').length).toBe(1);
  });

  it('renders ServiceItemCard component', () => {
    expect(shallowComponent.find(ServiceItemCard).length).toBe(1);
  });

  it('attaches the close listener', () => {
    expect(shallowComponent.find('[data-cy="closeSidebar"]').prop('onClick')).toBe(handleClose);
  });

  it('displays the total count', () => {
    expect(shallowComponent.find('[data-cy="itemCount"]').text()).toEqual('1 OF 5 ITEMS');
  });

  it('disables previous button at beginning', () => {
    expect(shallowComponent.find('[data-cy="prevServiceItem"]').prop('disabled')).toBe(true);
  });

  it('enables next button at beginning', () => {
    expect(shallowComponent.find('[data-cy="nextServiceItem"]').prop('disabled')).toBe(false);
  });

  it('navigates service items in timestamp ascending', () => {
    const nextButton = mountedComponent.find('[data-cy="nextServiceItem"]');
    const prevButton = mountedComponent.find('[data-cy="prevServiceItem"]');

    compareItem(mountedComponent, serviceItemCards[4]);

    nextButton.simulate('click');
    mountedComponent.update();

    compareItem(mountedComponent, serviceItemCards[3]);

    nextButton.simulate('click');
    mountedComponent.update();

    compareItem(mountedComponent, serviceItemCards[0]);

    nextButton.simulate('click');
    mountedComponent.update();

    compareItem(mountedComponent, serviceItemCards[1]);

    nextButton.simulate('click');
    mountedComponent.update();

    compareItem(mountedComponent, serviceItemCards[2]);

    expect(mountedComponent.find('[data-cy="nextServiceItem"]').prop('disabled')).toBe(true);

    prevButton.simulate('click');
    mountedComponent.update();

    compareItem(mountedComponent, serviceItemCards[1]);

    prevButton.simulate('click');
    mountedComponent.update();

    compareItem(mountedComponent, serviceItemCards[0]);

    prevButton.simulate('click');
    mountedComponent.update();

    compareItem(mountedComponent, serviceItemCards[3]);

    prevButton.simulate('click');
    mountedComponent.update();

    compareItem(mountedComponent, serviceItemCards[4]);
  });
});
