import React from 'react';
import { act } from 'react-dom/test-utils';
import { shallow, mount } from 'enzyme';

import { toDollarString } from '../../../shared/formatters';

import ReviewServiceItems from './ReviewServiceItems';

import { SHIPMENT_OPTIONS, PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';

const serviceItemCards = [
  {
    id: '1',
    shipmentType: SHIPMENT_OPTIONS.HHG,
    shipmentId: '10',
    serviceItemName: 'Domestic linehaul',
    amount: 6423,
    status: PAYMENT_SERVICE_ITEM_STATUS.SUBMITTED,
    createdAt: '2020-01-01T00:08:00.999Z',
  },
  {
    id: '2',
    shipmentType: SHIPMENT_OPTIONS.HHG,
    shipmentId: '10',
    serviceItemName: 'Fuel Surcharge',
    amount: 50.25,
    createdAt: '2020-01-01T00:08:30.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.SUBMITTED,
  },
  {
    id: '3',
    shipmentType: SHIPMENT_OPTIONS.NTS,
    shipmentId: '20',
    serviceItemName: 'Domestic linehaul',
    amount: 0.1,
    createdAt: '2020-01-01T00:09:00.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.SUBMITTED,
  },
  {
    id: '4',
    shipmentType: null,
    shipmentId: null,
    serviceItemName: 'Counseling Services',
    amount: 1000,
    createdAt: '2020-01-01T00:02:00.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.SUBMITTED,
  },
];

const basicServiceItemCards = [
  {
    id: '4',
    shipmentType: null,
    shipmentId: null,
    serviceItemName: 'Counseling Services',
    amount: 1000,
    createdAt: '2020-01-01T00:02:00.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.SUBMITTED,
  },
  {
    id: '5',
    shipmentType: null,
    shipmentId: null,
    serviceItemName: 'Move management',
    amount: 1,
    createdAt: '2020-01-01T00:01:00.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.SUBMITTED,
  },
];

const compareItem = (component, item) => {
  expect(component.find('[data-testid="serviceItemName"]').text()).toEqual(item.serviceItemName);
  expect(component.find('[data-testid="serviceItemAmount"]').text()).toEqual(toDollarString(item.amount));
};

describe('ReviewServiceItems component', () => {
  const handleClose = jest.fn();
  const patchPaymentServiceItem = jest.fn();
  const shallowComponent = shallow(
    <ReviewServiceItems
      serviceItemCards={serviceItemCards}
      handleClose={handleClose}
      patchPaymentServiceItem={patchPaymentServiceItem}
    />,
  );
  const mountedComponent = mount(
    <ReviewServiceItems
      serviceItemCards={serviceItemCards}
      handleClose={handleClose}
      patchPaymentServiceItem={patchPaymentServiceItem}
    />,
  );

  it('renders without crashing', () => {
    expect(shallowComponent.find('[data-testid="ReviewServiceItems"]').length).toBe(1);
  });

  it('doesn’t crash if there are no cards', () => {
    const componentWithNoCards = shallow(
      <ReviewServiceItems handleClose={handleClose} patchPaymentServiceItem={patchPaymentServiceItem} />,
    );
    expect(componentWithNoCards.exists()).toBe(true);
  });

  it('renders ServiceItemCard component', () => {
    expect(mountedComponent.find('ServiceItemCard').length).toBe(1);
  });

  it('attaches the close listener', () => {
    expect(mountedComponent.find('[data-testid="closeSidebar"]').prop('onClick')).toBe(handleClose);
  });

  it('displays the total count', () => {
    expect(mountedComponent.find('[data-testid="itemCount"]').text()).toEqual('1 OF 4 ITEMS');
  });

  it('disables previous button at beginning', () => {
    expect(mountedComponent.find('[data-testid="prevServiceItem"]').prop('disabled')).toBe(true);
  });

  it('enables next button at beginning', () => {
    expect(mountedComponent.find('[data-testid="nextServiceItem"]').prop('disabled')).toBe(false);
  });

  it('displays the total approved amount', () => {
    expect(mountedComponent.find('[data-testid="approvedAmount"]').text()).toEqual('$0.00');
  });

  it('renders two basic service item cards', () => {
    const basicWrapper = mount(
      <ReviewServiceItems
        serviceItemCards={basicServiceItemCards}
        handleClose={handleClose}
        patchPaymentServiceItem={patchPaymentServiceItem}
      />,
    );
    expect(basicWrapper.find('ServiceItemCard').length).toBe(2);
  });

  describe('navigating through service items', () => {
    const nextButton = mountedComponent.find('[data-testid="nextServiceItem"]');
    const prevButton = mountedComponent.find('[data-testid="prevServiceItem"]');

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
      expect(mountedComponent.find('[data-testid="nextServiceItem"]').prop('disabled')).toBe(true);
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
    const nextButton = mountedComponent.find('[data-testid="nextServiceItem"]');

    it('can approve an item', async () => {
      const approveInput = mountedComponent.find(`input[name="status"][value="APPROVED"]`);
      expect(approveInput.length).toBe(1);

      await act(async () => {
        approveInput.simulate('change');
      });
      mountedComponent.update();
      expect(patchPaymentServiceItem).toHaveBeenCalled();
    });

    it('can reject an item', async () => {
      nextButton.simulate('click');
      mountedComponent.update();

      const rejectInput = mountedComponent.find(`input[name="status"][value="DENIED"]`);
      expect(rejectInput.length).toBe(1);

      await act(async () => {
        rejectInput.simulate('change');
      });
      mountedComponent.update();
      const saveButton = mountedComponent.find('[data-testid="rejectionSaveButton"]');
      expect(saveButton.length).toBe(1);
      await act(async () => {
        saveButton.simulate('click');
      });
      mountedComponent.update();
      expect(patchPaymentServiceItem).toHaveBeenCalled();
    });

    it('can enter a reason for rejecting an item', async () => {
      const rejectReasonInput = mountedComponent.find(`textarea[name="rejectionReason"]`);
      expect(rejectReasonInput.length).toBe(1);

      await act(async () => {
        rejectReasonInput.simulate('change', {
          target: {
            name: 'rejectionReason',
            value: 'This is why I rejected it',
          },
        });
      });
      mountedComponent.update();
      expect(rejectReasonInput.text()).toEqual('This is why I rejected it');
    });

    it('can clear the selections for an item', async () => {
      const clearSelectionButton = mountedComponent.find('[data-testid="clearStatusButton"]');
      expect(clearSelectionButton.length).toBe(1);
      const rejectReasonInput = mountedComponent.find(`textarea[name="rejectionReason"]`);
      expect(rejectReasonInput.length).toBe(1);

      await act(async () => {
        rejectReasonInput.simulate('change', {
          target: {
            name: 'rejectionReason',
            value: 'This is why I rejected it',
          },
        });
      });
      mountedComponent.update();

      await act(async () => {
        clearSelectionButton.simulate('click');
      });
      mountedComponent.update();
      const clearedRejectionInput = mountedComponent.find(`textarea[name="rejectionReason"]`);
      expect(clearedRejectionInput.length).toBe(0);
    });
  });

  describe('updating the total amount approved', () => {
    const cardsWithInitialValues = [
      {
        id: '1',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        shipmentId: '10',
        serviceItemName: 'Domestic linehaul',
        amount: 6423,
        status: PAYMENT_SERVICE_ITEM_STATUS.SUBMITTED,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
      {
        id: '2',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        shipmentId: '10',
        serviceItemName: 'Fuel Surcharge',
        amount: 50.25,
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
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
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
        createdAt: '2020-01-01T00:02:00.999Z',
      },
      {
        id: '5',
        shipmentType: null,
        shipmentId: null,
        serviceItemName: 'Move management',
        amount: 1,
        status: PAYMENT_SERVICE_ITEM_STATUS.REJECTED,
        rejectionReason: 'Wrong amount specified',
        createdAt: '2020-01-01T00:01:00.999Z',
      },
    ];

    const componentWithInitialValues = mount(
      <ReviewServiceItems
        handleClose={handleClose}
        serviceItemCards={cardsWithInitialValues}
        patchPaymentServiceItem={patchPaymentServiceItem}
      />,
    );
    const approvedAmount = componentWithInitialValues.find('[data-testid="approvedAmount"]');

    it('calculates the sum for items with initial values', () => {
      expect(approvedAmount.text()).toEqual('$1,050.25');
    });
  });
});
