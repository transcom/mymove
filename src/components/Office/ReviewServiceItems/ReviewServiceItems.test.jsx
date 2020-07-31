/* eslint-disable react/jsx-props-no-spreading */
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
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    createdAt: '2020-01-01T00:08:00.999Z',
  },
  {
    id: '2',
    shipmentType: SHIPMENT_OPTIONS.HHG,
    shipmentId: '10',
    serviceItemName: 'Fuel Surcharge',
    amount: 50.25,
    createdAt: '2020-01-01T00:08:30.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  },
  {
    id: '3',
    shipmentType: SHIPMENT_OPTIONS.NTS,
    shipmentId: '20',
    serviceItemName: 'Domestic linehaul',
    amount: 0.1,
    createdAt: '2020-01-01T00:09:00.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  },
  {
    id: '4',
    shipmentType: null,
    shipmentId: null,
    serviceItemName: 'Counseling Services',
    amount: 1000,
    createdAt: '2020-01-01T00:02:00.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
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
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  },
  {
    id: '5',
    shipmentType: null,
    shipmentId: null,
    serviceItemName: 'Move management',
    amount: 1,
    createdAt: '2020-01-01T00:01:00.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  },
];

const compareItem = (component, item) => {
  expect(component.find('[data-testid="serviceItemName"]').text()).toEqual(item.serviceItemName);
  expect(component.find('[data-testid="serviceItemAmount"]').text()).toEqual(toDollarString(item.amount));
};

describe('ReviewServiceItems component', () => {
  const handleClose = jest.fn();
  const patchPaymentServiceItem = jest.fn();
  const onCompleteReview = jest.fn();

  const requiredProps = {
    handleClose,
    patchPaymentServiceItem,
    onCompleteReview,
  };

  const shallowComponent = shallow(<ReviewServiceItems serviceItemCards={serviceItemCards} {...requiredProps} />);
  const mountedComponent = mount(<ReviewServiceItems serviceItemCards={serviceItemCards} {...requiredProps} />);

  it('renders without crashing', () => {
    expect(shallowComponent.find('[data-testid="ReviewServiceItems"]').length).toBe(1);
  });

  it('doesn’t crash if there are no cards', () => {
    const componentWithNoCards = shallow(<ReviewServiceItems {...requiredProps} />);
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
    const basicWrapper = mount(<ReviewServiceItems serviceItemCards={basicServiceItemCards} {...requiredProps} />);
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

    it('shows the Complete Review step after the last item', () => {
      nextButton.simulate('click');
      mountedComponent.update();

      expect(mountedComponent.find('[data-testid="authorizePaymentBtn"]').exists()).toBe(true);
    });

    it('does not show a Next button on the Complete Review step', () => {
      expect(mountedComponent.find('[data-testid="nextServiceItem"]').exists()).toBe(false);
    });

    it('can click back to the first item', () => {
      prevButton.simulate('click');
      mountedComponent.update();

      compareItem(mountedComponent, serviceItemCards[2]);

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
      const nextButton = mountedComponent.find('[data-testid="nextServiceItem"]');

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
        status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
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
        status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
        rejectionReason: 'Wrong amount specified',
        createdAt: '2020-01-01T00:01:00.999Z',
      },
    ];

    const componentWithInitialValues = mount(
      <ReviewServiceItems serviceItemCards={cardsWithInitialValues} {...requiredProps} />,
    );
    const approvedAmount = componentWithInitialValues.find('[data-testid="approvedAmount"]');

    it('calculates the approved sum for items with initial values', () => {
      expect(approvedAmount.text()).toEqual('$1,050.25');
    });
  });

  describe('completing the review step', () => {
    describe('with no error', () => {
      it('lands on the Complete Review step after reviewing all items', () => {
        const nextButton = mountedComponent.find('[data-testid="nextServiceItem"]');

        nextButton.simulate('click');
        mountedComponent.update();

        nextButton.simulate('click');
        mountedComponent.update();

        nextButton.simulate('click');
        mountedComponent.update();

        const header = mountedComponent.find('h2');
        expect(header.exists()).toBe(true);
        expect(header.text()).toEqual('Complete request');

        const body = mountedComponent.find('.body p');
        expect(body.exists()).toBe(true);
        expect(body.text()).toEqual('Do you authorize this payment of $0.00?');
      });

      it('can click on Authorize Payment', async () => {
        const authorizeBtn = mountedComponent.find('[data-testid="authorizePaymentBtn"]');
        expect(authorizeBtn.exists()).toBe(true);

        await act(async () => {
          authorizeBtn.simulate('click');
        });

        mountedComponent.update();
        expect(onCompleteReview).toHaveBeenCalled();
      });
    });

    describe('with a validation error', () => {
      const componentWithMockError = mount(
        <ReviewServiceItems
          serviceItemCards={serviceItemCards}
          {...requiredProps}
          completeReviewError={{ detail: 'A validation error occurred' }}
        />,
      );

      it('lands on the Complete Review step after reviewing all items', () => {
        const nextButton = componentWithMockError.find('[data-testid="nextServiceItem"]');

        nextButton.simulate('click');
        componentWithMockError.update();

        nextButton.simulate('click');
        componentWithMockError.update();

        nextButton.simulate('click');
        componentWithMockError.update();

        nextButton.simulate('click');
        componentWithMockError.update();

        const header = componentWithMockError.find('h2');
        expect(header.exists()).toBe(true);
        expect(header.text()).toEqual('Complete request');
      });

      it('displays the validation error', async () => {
        const errorMsg = componentWithMockError.find('[data-testid="errorMessage"]');
        expect(errorMsg.exists()).toBe(true);
        expect(errorMsg.text()).toEqual('Error: A validation error occurred');
      });
    });
  });
});
