/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { act } from 'react-dom/test-utils';
import { shallow, mount } from 'enzyme';

import ReviewServiceItems from './ReviewServiceItems';

import { toDollarString } from 'utils/formatters';
import { SHIPMENT_OPTIONS, PAYMENT_SERVICE_ITEM_STATUS, PAYMENT_REQUEST_STATUS } from 'shared/constants';
import { serviceItemCodes } from 'content/serviceItems';

const pendingPaymentRequest = {
  status: PAYMENT_REQUEST_STATUS.PENDING,
};

const reviewedPaymentRequest = {
  status: PAYMENT_REQUEST_STATUS.REVIEWED,
  reviewedAt: '2020-11-03T20:01:01.001Z',
};

const serviceItemCards = [
  {
    id: '1',
    mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    mtoShipmentID: '10',
    mtoServiceItemName: serviceItemCodes.DLH,
    mtoServiceItemCode: 'DLH',
    amount: 6423,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    createdAt: '2020-01-01T00:08:00.999Z',
  },
  {
    id: '2',
    mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    mtoShipmentID: '10',
    mtoServiceItemName: serviceItemCodes.FSC,
    mtoServiceItemCode: 'FSC',
    amount: 50.25,
    createdAt: '2020-01-01T00:08:30.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  },
  {
    id: '3',
    mtoShipmentType: SHIPMENT_OPTIONS.NTS,
    mtoShipmentID: '20',
    mtoServiceItemName: serviceItemCodes.DLH,
    mtoServiceItemCode: 'DLH',
    amount: 0.1,
    createdAt: '2020-01-01T00:09:00.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  },
  {
    id: '4',
    mtoShipmentType: null,
    mtoShipmentID: null,
    mtoServiceItemName: serviceItemCodes.CS,
    mtoServiceItemCode: 'CS',
    amount: 1000,
    createdAt: '2020-01-01T00:02:00.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  },
];

const basicServiceItemCards = [
  {
    id: '4',
    mtoShipmentType: null,
    mtoShipmentID: null,
    mtoServiceItemName: serviceItemCodes.CS,
    mtoServiceItemCode: 'CS',
    amount: 1000,
    createdAt: '2020-01-01T00:02:00.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  },
  {
    id: '5',
    mtoShipmentType: null,
    mtoShipmentID: null,
    mtoServiceItemName: serviceItemCodes.MS,
    mtoServiceItemCode: 'MS',
    amount: 1,
    createdAt: '2020-01-01T00:01:00.999Z',
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  },
];

const compareItem = (component, item) => {
  expect(component.find('[data-testid="serviceItemName"]').text()).toEqual(item.mtoServiceItemName);
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

  const shallowComponent = shallow(
    <ReviewServiceItems
      paymentRequest={pendingPaymentRequest}
      serviceItemCards={serviceItemCards}
      {...requiredProps}
    />,
  );
  const mountedComponent = mount(
    <ReviewServiceItems
      paymentRequest={pendingPaymentRequest}
      serviceItemCards={serviceItemCards}
      {...requiredProps}
    />,
  );

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
    expect(mountedComponent.find('button[data-testid="closeSidebar"]').prop('onClick')).toBe(handleClose);
  });

  it('displays the total count', () => {
    expect(mountedComponent.find('[data-testid="itemCount"]').text()).toEqual('1 OF 4 ITEMS');
  });

  it('disables previous button at beginning', () => {
    expect(mountedComponent.find('button[data-testid="prevServiceItem"]').prop('disabled')).toBe(true);
  });

  it('enables next button at beginning', () => {
    expect(mountedComponent.find('button[data-testid="nextServiceItem"]').prop('disabled')).toBe(false);
  });

  it('displays the total approved amount', () => {
    expect(mountedComponent.find('[data-testid="approvedAmount"]').text()).toEqual('$0.00');
  });

  it('renders two basic service item cards', () => {
    const basicWrapper = mount(
      <ReviewServiceItems
        paymentRequest={pendingPaymentRequest}
        serviceItemCards={basicServiceItemCards}
        {...requiredProps}
      />,
    );
    expect(basicWrapper.find('ServiceItemCard').length).toBe(2);
  });

  describe('navigating through service items', () => {
    const cardsWithInitialValues = [
      {
        id: '1',
        mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
        mtoShipmentID: '10',
        mtoServiceItemName: serviceItemCodes.DLH,
        mtoServiceItemCode: 'DLH',
        amount: 6423,
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
      {
        id: '2',
        mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
        mtoShipmentID: '10',
        mtoServiceItemName: serviceItemCodes.FSC,
        mtoServiceItemCode: 'FSC',
        amount: 50.25,
        createdAt: '2020-01-01T00:08:30.999Z',
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
      },
      {
        id: '3',
        mtoShipmentType: SHIPMENT_OPTIONS.NTS,
        mtoShipmentID: '20',
        mtoServiceItemName: serviceItemCodes.DLH,
        mtoServiceItemCode: 'DLH',
        amount: 0.1,
        createdAt: '2020-01-01T00:09:00.999Z',
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
      },
      {
        id: '4',
        mtoShipmentType: null,
        mtoShipmentID: null,
        mtoServiceItemName: serviceItemCodes.CS,
        mtoServiceItemCode: 'CS',
        amount: 1000,
        createdAt: '2020-01-01T00:02:00.999Z',
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
      },
    ];
    const componentWithInitialValues = mount(
      <ReviewServiceItems
        paymentRequest={pendingPaymentRequest}
        serviceItemCards={cardsWithInitialValues}
        {...requiredProps}
      />,
    );
    const nextButton = componentWithInitialValues.find('button[data-testid="nextServiceItem"]');
    const prevButton = componentWithInitialValues.find('button[data-testid="prevServiceItem"]');

    it('renders the service item cards ordered by timestamp ascending', () => {
      compareItem(componentWithInitialValues, serviceItemCards[3]);

      nextButton.simulate('click');
      componentWithInitialValues.update();

      compareItem(componentWithInitialValues, serviceItemCards[0]);

      nextButton.simulate('click');
      componentWithInitialValues.update();

      compareItem(componentWithInitialValues, serviceItemCards[1]);

      nextButton.simulate('click');
      componentWithInitialValues.update();

      compareItem(componentWithInitialValues, serviceItemCards[2]);
    });

    it('shows the Complete Review step after the last item', async () => {
      nextButton.simulate('click');
      componentWithInitialValues.update();

      expect(componentWithInitialValues.find('[data-testid="authorizePaymentBtn"]').exists()).toBe(true);
    });

    it('does not show a Next button on the Complete Review step', async () => {
      expect(componentWithInitialValues.find('[data-testid="nextServiceItem"]').exists()).toBe(false);
    });

    it('can click back to the first item', () => {
      prevButton.simulate('click');
      componentWithInitialValues.update();

      compareItem(componentWithInitialValues, serviceItemCards[2]);

      prevButton.simulate('click');
      componentWithInitialValues.update();

      compareItem(componentWithInitialValues, serviceItemCards[1]);

      prevButton.simulate('click');
      componentWithInitialValues.update();

      compareItem(componentWithInitialValues, serviceItemCards[0]);

      prevButton.simulate('click');
      componentWithInitialValues.update();

      compareItem(componentWithInitialValues, serviceItemCards[3]);
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
      const nextButton = mountedComponent.find('button[data-testid="nextServiceItem"]');

      nextButton.simulate('click');
      mountedComponent.update();

      const rejectInput = mountedComponent.find(`input[name="status"][value="DENIED"]`);
      expect(rejectInput.length).toBe(1);

      await act(async () => {
        rejectInput.simulate('change');
      });
      mountedComponent.update();
      const saveButton = mountedComponent.find('button[data-testid="rejectionSaveButton"]');
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
      const clearSelectionButton = mountedComponent.find('button[data-testid="clearStatusButton"]');
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
        mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
        mtoShipmentID: '10',
        mtoServiceItemName: serviceItemCodes.DLH,
        amount: 6423,
        status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
      {
        id: '2',
        mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
        mtoShipmentID: '10',
        mtoServiceItemName: serviceItemCodes.FSC,
        amount: 50.25,
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
        createdAt: '2020-01-01T00:08:30.999Z',
      },
      {
        id: '3',
        mtoShipmentType: SHIPMENT_OPTIONS.NTS,
        mtoShipmentID: '20',
        mtoServiceItemName: serviceItemCodes.DLH,
        amount: 0.1,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
      {
        id: '4',
        mtoShipmentType: null,
        mtoShipmentID: null,
        mtoServiceItemName: serviceItemCodes.CS,
        amount: 1000,
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
        createdAt: '2020-01-01T00:02:00.999Z',
      },
      {
        id: '5',
        mtoShipmentType: null,
        mtoShipmentID: null,
        mtoServiceItemName: serviceItemCodes.MS,
        amount: 1,
        status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
        rejectionReason: 'Wrong amount specified',
        createdAt: '2020-01-01T00:01:00.999Z',
      },
    ];

    const componentWithInitialValues = mount(
      <ReviewServiceItems
        paymentRequest={pendingPaymentRequest}
        serviceItemCards={cardsWithInitialValues}
        {...requiredProps}
      />,
    );
    const approvedAmount = componentWithInitialValues.find('[data-testid="approvedAmount"]');

    it('calculates the approved sum for items with initial values', () => {
      expect(approvedAmount.text()).toEqual('$1,050.25');
    });

    it('can clear approved amount for an item', async () => {
      const approveInput = mountedComponent.find(`input[name="status"][value="APPROVED"]`);
      expect(approveInput.length).toBe(1);

      await act(async () => {
        approveInput.simulate('change');
      });
      mountedComponent.update();

      const clearSelectionButton = mountedComponent.find('button[data-testid="clearStatusButton"]');
      expect(clearSelectionButton.length).toBe(1);

      await act(async () => {
        clearSelectionButton.simulate('click');
      });
      mountedComponent.update();

      const clearedAmount = mountedComponent.find('[data-testid="approvedAmount"]');
      expect(clearedAmount.text()).toEqual('$0.00');
    });
  });

  describe('completing the review step', () => {
    describe('with no error', () => {
      const cardsWithInitialValues = [
        {
          id: '1',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoShipmentID: '10',
          mtoServiceItemName: serviceItemCodes.DLH,
          mtoServiceItemCode: 'DLH',
          amount: 6423,
          status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
          createdAt: '2020-01-01T00:08:00.999Z',
        },
        {
          id: '2',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoShipmentID: '10',
          mtoServiceItemName: serviceItemCodes.FSC,
          mtoServiceItemCode: 'FSC',
          amount: 50.25,
          status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
          createdAt: '2020-01-01T00:08:30.999Z',
        },
      ];
      const componentWithInitialValues = mount(
        <ReviewServiceItems
          paymentRequest={pendingPaymentRequest}
          serviceItemCards={cardsWithInitialValues}
          {...requiredProps}
        />,
      );

      it('lands on the Complete Review step after reviewing all items', () => {
        const nextButton = componentWithInitialValues.find('button[data-testid="nextServiceItem"]');

        nextButton.simulate('click');
        mountedComponent.update();

        nextButton.simulate('click');
        mountedComponent.update();

        nextButton.simulate('click');
        mountedComponent.update();

        const header = componentWithInitialValues.find('h2');
        expect(header.exists()).toBe(true);
        expect(header.text()).toEqual('Complete request');

        const body = componentWithInitialValues.find('[data-testid="AuthorizePayment"] [data-testid="header"]');
        expect(body.exists()).toBe(true);
        expect(body.text()).toEqual('Do you authorize this payment of $6,473.25?');
      });

      it('can click on Authorize Payment', async () => {
        const authorizeBtn = componentWithInitialValues.find('button[data-testid="authorizePaymentBtn"]');
        expect(authorizeBtn.exists()).toBe(true);

        await act(async () => {
          authorizeBtn.simulate('click');
        });

        componentWithInitialValues.update();
        expect(onCompleteReview).toHaveBeenCalled();
      });
    });

    describe('with one item that needs review', () => {
      const cardWithInitialValues = [
        {
          id: '1',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoShipmentID: '10',
          mtoServiceItemName: serviceItemCodes.DLH,
          amount: 6423,
          status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
          createdAt: '2020-01-01T00:08:00.999Z',
        },
      ];
      const componentWithInitialValues = mount(
        <ReviewServiceItems
          paymentRequest={pendingPaymentRequest}
          serviceItemCards={cardWithInitialValues}
          {...requiredProps}
        />,
      );

      it('lands on the Complete Review step after reviewing one item', () => {
        const nextButton = componentWithInitialValues.find('button[data-testid="nextServiceItem"]');

        nextButton.simulate('click');
        mountedComponent.update();

        const header = componentWithInitialValues.find('h2');
        expect(header.exists()).toBe(true);
        expect(header.text()).toEqual('Complete request');

        const needsReviewHeader = componentWithInitialValues.find(
          '[data-testid="NeedsReview"] > [data-testid="header"]',
        );
        expect(needsReviewHeader.exists()).toBe(true);
        expect(needsReviewHeader.text()).toEqual('1 item still needs your review');

        const needsReviewContent = componentWithInitialValues.find(
          '[data-testid="NeedsReview"] [data-testid="content"]',
        );
        expect(needsReviewContent.exists()).toBe(true);
        expect(needsReviewContent.text()).toEqual('Accept or reject all service items, then authorized payment.');
      });
    });

    describe('with items that needs review', () => {
      const cardsWithInitialValues = [
        {
          id: '1',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoShipmentID: '10',
          mtoServiceItemName: serviceItemCodes.DLH,
          amount: 6423,
          status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
          createdAt: '2020-01-01T00:08:00.999Z',
        },
        {
          id: '2',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoShipmentID: '10',
          mtoServiceItemName: serviceItemCodes.FSC,
          amount: 50.25,
          status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
          createdAt: '2020-01-01T00:08:30.999Z',
        },
        {
          id: '3',
          mtoShipmentType: null,
          mtoShipmentID: null,
          mtoServiceItemName: serviceItemCodes.MS,
          amount: 10.25,
          status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
          createdAt: '2020-01-01T00:01:30.999Z',
        },
      ];
      const componentWithInitialValues = mount(
        <ReviewServiceItems
          paymentRequest={pendingPaymentRequest}
          serviceItemCards={cardsWithInitialValues}
          {...requiredProps}
        />,
      );

      it('lands on the Complete Review step after reviewing all items', () => {
        const nextButton = componentWithInitialValues.find('button[data-testid="nextServiceItem"]');

        nextButton.simulate('click');
        mountedComponent.update();

        nextButton.simulate('click');
        mountedComponent.update();

        nextButton.simulate('click');
        mountedComponent.update();

        const header = componentWithInitialValues.find('h2');
        expect(header.exists()).toBe(true);
        expect(header.text()).toEqual('Complete request');

        const needsReviewHeader = componentWithInitialValues.find(
          '[data-testid="NeedsReview"] > [data-testid="header"]',
        );
        expect(needsReviewHeader.exists()).toBe(true);
        expect(needsReviewHeader.text()).toEqual('2 items still needs your review');

        const needsReviewContent = componentWithInitialValues.find(
          '[data-testid="NeedsReview"] [data-testid="content"]',
        );
        expect(needsReviewContent.exists()).toBe(true);
        expect(needsReviewContent.text()).toEqual('Accept or reject all service items, then authorized payment.');
      });

      it('can click on Finish review button', async () => {
        const finishReviewBtn = componentWithInitialValues.find('button[data-testid="finishReviewBtn"]');
        expect(finishReviewBtn.exists()).toBe(true);

        await act(async () => {
          finishReviewBtn.simulate('click');
        });

        // goes back to first service item that needs to be reviewed
        componentWithInitialValues.update();
        compareItem(componentWithInitialValues, cardsWithInitialValues[2]);
      });
    });

    describe('with all items rejected', () => {
      const cardsWithInitialValues = [
        {
          id: '1',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoShipmentID: '10',
          mtoServiceItemName: serviceItemCodes.DLH,
          mtoServiceItemCode: 'DLH',
          amount: 6423,
          status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
          createdAt: '2020-01-01T00:08:00.999Z',
        },
        {
          id: '2',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoShipmentID: '10',
          mtoServiceItemName: serviceItemCodes.FSC,
          mtoServiceItemCode: 'FSC',
          amount: 50.25,
          status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
          createdAt: '2020-01-01T00:08:30.999Z',
        },
        {
          id: '3',
          mtoShipmentType: null,
          mtoShipmentID: null,
          mtoServiceItemName: serviceItemCodes.MS,
          mtoServiceItemCode: 'MS',
          amount: 10.25,
          status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
          createdAt: '2020-01-01T00:01:30.999Z',
        },
      ];
      const componentWithInitialValues = mount(
        <ReviewServiceItems
          paymentRequest={pendingPaymentRequest}
          serviceItemCards={cardsWithInitialValues}
          {...requiredProps}
        />,
      );

      it('lands on the Complete Review step after reviewing all items', () => {
        const nextButton = componentWithInitialValues.find('button[data-testid="nextServiceItem"]');

        nextButton.simulate('click');
        mountedComponent.update();

        nextButton.simulate('click');
        mountedComponent.update();

        nextButton.simulate('click');
        mountedComponent.update();

        const header = componentWithInitialValues.find('h2');
        expect(header.exists()).toBe(true);
        expect(header.text()).toEqual('Complete request');

        const rejectRequestHeader = componentWithInitialValues.find(
          '[data-testid="RejectRequest"] > [data-testid="content"]',
        );
        expect(rejectRequestHeader.exists()).toBe(true);
        expect(rejectRequestHeader.text()).toEqual(
          `You're rejecting all service items. No payment will be authorized.`,
        );
      });

      it('can click on Reject request button', async () => {
        const rejectRequestBtn = componentWithInitialValues.find('button[data-testid="rejectRequestBtn"]');
        expect(rejectRequestBtn.exists()).toBe(true);

        await act(async () => {
          rejectRequestBtn.simulate('click');
        });

        // this hooks up to the same onCompleteReview function call
        componentWithInitialValues.update();
        expect(onCompleteReview).toHaveBeenCalled();
      });
    });

    describe('with a validation error', () => {
      const componentWithMockError = mount(
        <ReviewServiceItems
          paymentRequest={pendingPaymentRequest}
          serviceItemCards={serviceItemCards}
          {...requiredProps}
          completeReviewError={{ detail: 'A validation error occurred' }}
        />,
      );

      it('lands on the Complete Review step after reviewing all items', () => {
        const nextButton = componentWithMockError.find('button[data-testid="nextServiceItem"]');

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

  describe('viewing authorized payment request', () => {
    describe('with approved service items', () => {
      const cardsWithInitialValues = [
        {
          id: '1',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoShipmentID: '10',
          mtoServiceItemName: serviceItemCodes.DLH,
          mtoServiceItemCode: 'DLH',
          amount: 6423,
          status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
          rejectionReason: 'Duplicate charge',
          createdAt: '2020-01-01T00:08:00.999Z',
        },
        {
          id: '2',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoShipmentID: '10',
          mtoServiceItemName: serviceItemCodes.FSC,
          mtoServiceItemCode: 'FSC',
          amount: 50.25,
          status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
          createdAt: '2020-01-01T00:08:30.999Z',
        },
      ];
      const reviewedComponent = mount(
        <ReviewServiceItems
          paymentRequest={reviewedPaymentRequest}
          serviceItemCards={cardsWithInitialValues}
          {...requiredProps}
        />,
      );

      it('displays the payment reviewed component', () => {
        expect(reviewedComponent.find('[data-testid="PaymentReviewed"]').length).toBe(1);
      });

      it('displays the authorized amount and reviewed on date', () => {
        expect(reviewedComponent.find('[data-testid="paymentAuthorizedAmt"]').text()).toEqual(
          'Payment authorized: $50.25',
        );
        expect(reviewedComponent.find('[data-testid="reviewedOn"]').text()).toEqual('On: 03 Nov 2020');
      });

      it('displays the success alert', () => {
        const alert = reviewedComponent.find('Alert');
        expect(alert.length).toBe(1);
        expect(alert.prop('type')).toBe('success');
        expect(alert.text()).toEqual('The payment request was successfully submitted.');
      });

      it('disables the form elements', () => {
        const backButton = reviewedComponent.find('button[data-testid="prevServiceItem"]');

        backButton.simulate('click');
        reviewedComponent.update();

        let statusSummary = reviewedComponent.find('[data-testid="statusHeading"]');
        expect(statusSummary.find('FontAwesomeIcon').prop('icon')).toEqual('check');
        expect(statusSummary.text()).toBe('Accepted');

        backButton.simulate('click');
        reviewedComponent.update();

        statusSummary = reviewedComponent.find('[data-testid="statusHeading"]');
        expect(statusSummary.find('FontAwesomeIcon').prop('icon')).toEqual('times');
        expect(statusSummary.text()).toBe('Rejected');
        expect(reviewedComponent.find('[data-testid="rejectionReason"]').text()).toBe('Duplicate charge');
      });
    });

    describe('with all rejected service items', () => {
      const cardsWithInitialValues = [
        {
          id: '1',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoShipmentID: '10',
          mtoServiceItemName: serviceItemCodes.DLH,
          mtoServiceItemCode: 'DLH',
          amount: 6423,
          status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
          rejectionReason: 'Duplicate charge',
          createdAt: '2020-01-01T00:08:00.999Z',
        },
        {
          id: '2',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoShipmentID: '10',
          mtoServiceItemName: serviceItemCodes.FSC,
          mtoServiceItemCode: 'FSC',
          amount: 50.25,
          status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
          rejectionReason: 'Not applicable',
          createdAt: '2020-01-01T00:08:30.999Z',
        },
      ];
      const reviewedComponent = mount(
        <ReviewServiceItems
          paymentRequest={reviewedPaymentRequest}
          serviceItemCards={cardsWithInitialValues}
          {...requiredProps}
        />,
      );

      it('displays the payment reviewed component', () => {
        expect(reviewedComponent.find('[data-testid="PaymentReviewed"]').length).toBe(1);
      });

      it('displays authorized amount none', () => {
        expect(reviewedComponent.find('[data-testid="paymentAuthorizedAmt"]').text()).toEqual(
          'Payment authorized: none',
        );
      });

      it('displays the success alert', () => {
        const alert = reviewedComponent.find('Alert');
        expect(alert.length).toBe(1);
        expect(alert.prop('type')).toBe('success');
        expect(alert.text()).toEqual('The payment request was successfully submitted.');
      });

      it('displays service item status summary', () => {
        const backButton = reviewedComponent.find('button[data-testid="prevServiceItem"]');

        backButton.simulate('click');
        reviewedComponent.update();

        let statusSummary = reviewedComponent.find('[data-testid="statusHeading"]');
        expect(statusSummary.find('FontAwesomeIcon').prop('icon')).toEqual('times');
        expect(statusSummary.text()).toBe('Rejected');
        expect(reviewedComponent.find('[data-testid="rejectionReason"]').text()).toBe('Not applicable');

        backButton.simulate('click');
        reviewedComponent.update();

        statusSummary = reviewedComponent.find('[data-testid="statusHeading"]');
        expect(statusSummary.find('FontAwesomeIcon').prop('icon')).toEqual('times');
        expect(statusSummary.text()).toBe('Rejected');
        expect(reviewedComponent.find('[data-testid="rejectionReason"]').text()).toBe('Duplicate charge');
      });

      it('displays the expected status', () => {
        expect(reviewedComponent.find('[data-testid="rejectionReason"]').text()).toBe('Duplicate charge');
      });
    });
  });
});
