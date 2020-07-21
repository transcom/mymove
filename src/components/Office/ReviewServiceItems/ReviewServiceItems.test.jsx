import React from 'react';
import { act } from 'react-dom/test-utils';
import { shallow, mount } from 'enzyme';

import { toDollarString } from '../../../shared/formatters';

import ReviewServiceItems from './ReviewServiceItems';

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
];

// const basicServiceItemsCard = [
//   {
//     id: '4',
//     shipmentType: null,
//     shipmentId: null,
//     serviceItemName: 'Counseling Services',
//     amount: 1000,
//     createdAt: '2020-01-01T00:02:00.999Z',
//   },
//   {
//     id: '5',
//     shipmentType: null,
//     shipmentId: null,
//     serviceItemName: 'Move management',
//     amount: 1,
//     createdAt: '2020-01-01T00:01:00.999Z',
//   },
// ];

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
    expect(shallowComponent.find('[data-cy="ReviewServiceItems"]').length).toBe(1);
  });

  it('renders a Formik form', () => {
    expect(shallowComponent.find('Formik').length).toBe(1);
  });

  it('renders ServiceItemCard component', () => {
    expect(mountedComponent.find('ServiceItemCard').length).toBe(1);
  });

  it('attaches the close listener', () => {
    expect(mountedComponent.find('[data-cy="closeSidebar"]').prop('onClick')).toBe(handleClose);
  });

  it('displays the total count', () => {
    expect(mountedComponent.find('[data-cy="itemCount"]').text()).toEqual('1 OF 4 ITEMS');
  });

  it('disables previous button at beginning', () => {
    expect(mountedComponent.find('[data-cy="prevServiceItem"]').prop('disabled')).toBe(true);
  });

  it('enables next button at beginning', () => {
    expect(mountedComponent.find('[data-cy="nextServiceItem"]').prop('disabled')).toBe(false);
  });

  describe('navigating through service items', () => {
    const nextButton = mountedComponent.find('[data-cy="nextServiceItem"]');
    const prevButton = mountedComponent.find('[data-cy="prevServiceItem"]');

    it('renders the service item cards ordered by timestamp ascending', () => {
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
    });

    it('disables the Next button on the last item', () => {
      expect(mountedComponent.find('[data-cy="nextServiceItem"]').prop('disabled')).toBe(true);
    });

    it('can click back to the first item', () => {
      prevButton.simulate('click');
      mountedComponent.update();

      compareItem(mountedComponent, serviceItemCards[1]);

      prevButton.simulate('click');
      mountedComponent.update();

      compareItem(mountedComponent, serviceItemCards[0]);

      prevButton.simulate('click');
      mountedComponent.update();

      compareItem(mountedComponent, serviceItemCards[3]);
    });
  });

  describe('filling out the service item form', () => {
    const nextButton = mountedComponent.find('[data-cy="nextServiceItem"]');

    it('the item values are blank by default', () => {
      const serviceItemCard = mountedComponent.find('ServiceItemCard');

      expect(serviceItemCard.prop('value')).toEqual({ status: undefined, rejectionReason: undefined });
    });

    it('can approve an item', async () => {
      const serviceItemId = serviceItemCards[3].id;
      const approveInput = mountedComponent.find(`input[name="${serviceItemId}.status"][value="APPROVED"]`);
      expect(approveInput.length).toBe(1);

      await act(async () => {
        approveInput.simulate('change');
      });
      mountedComponent.update();
      const serviceItemCard = mountedComponent.find('ServiceItemCard');
      expect(serviceItemCard.prop('value')).toEqual({ status: 'APPROVED', rejectionReason: undefined });
    });

    it('can reject an item', async () => {
      nextButton.simulate('click');
      mountedComponent.update();

      const serviceItemId = serviceItemCards[0].id;
      const rejectInput = mountedComponent.find(`input[name="${serviceItemId}.status"][value="REJECTED"]`);
      expect(rejectInput.length).toBe(1);

      await act(async () => {
        rejectInput.simulate('change');
      });
      mountedComponent.update();
      const serviceItemCard = mountedComponent.find('ServiceItemCard');
      expect(serviceItemCard.prop('value')).toEqual({ status: 'REJECTED', rejectionReason: undefined });
    });

    it('can enter a reason for rejecting an item', async () => {
      const serviceItemId = serviceItemCards[0].id;
      const rejectReasonInput = mountedComponent.find(`textarea[name="${serviceItemId}.rejectionReason"]`);
      expect(rejectReasonInput.length).toBe(1);

      await act(async () => {
        rejectReasonInput.simulate('change', {
          target: {
            name: `${serviceItemId}.rejectionReason`,
            value: 'This is why I rejected it',
          },
        });
      });
      mountedComponent.update();
      const serviceItemCard = mountedComponent.find('ServiceItemCard');
      expect(serviceItemCard.prop('value')).toEqual({
        status: 'REJECTED',
        rejectionReason: 'This is why I rejected it',
      });
    });

    it('can clear the selections for an item', async () => {
      const clearSelectionButton = mountedComponent.find('[data-testid="clearStatusButton"]');
      expect(clearSelectionButton.length).toBe(1);

      await act(async () => {
        clearSelectionButton.simulate('click');
      });
      mountedComponent.update();
      const serviceItemCard = mountedComponent.find('ServiceItemCard');
      expect(serviceItemCard.prop('value')).toEqual({
        status: undefined,
        rejectionReason: undefined,
      });
    });
  });
});
