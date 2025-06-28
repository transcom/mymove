/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { mount } from 'enzyme';

import testParams from '../ServiceItemCalculations/serviceItemTestParams';

import ServiceItemCard from './ServiceItemCard';

import { PAYMENT_SERVICE_ITEM_STATUS, SHIPMENT_OPTIONS } from 'shared/constants';
import { serviceItemCodes } from 'content/serviceItems';
import { shipmentModificationTypes } from 'constants/shipments';

const basicServiceItemCard = {
  id: '1',
  mtoServiceItemName: serviceItemCodes.CS,
  mtoServiceItemCode: 'CS',
  amount: 1000,
  createdAt: '2020-01-01T00:02:00.999Z',
  status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  patchPaymentServiceItem: jest.fn(),
};

const reviewedBasicServiceItemCard = {
  ...basicServiceItemCard,
  requestComplete: true,
};

const needsReviewServiceItemCard = {
  id: '1',
  mtoShipmentType: SHIPMENT_OPTIONS.HHG,
  mtoShipmentID: '2',
  mtoShipmentDepartureDate: '04 May 2021',
  mtoShipmentPickupAddress: 'Fairfield, CA 94535',
  mtoShipmentDestinationAddress: 'Beverly Hills, CA 90210',
  mtoServiceItemName: serviceItemCodes.FSC,
  mtoServiceItemCode: 'FSC',
  amount: 1000,
  createdAt: '2020-01-01T00:02:00.999Z',
  status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  paymentServiceItemParams: testParams.FuelSurchage,
  patchPaymentServiceItem: jest.fn(),
};

const additionalDaySITServiceItemCard = {
  ...needsReviewServiceItemCard,
  mtoServiceItemName: serviceItemCodes.DOASIT,
  mtoServiceItemCode: 'DOASIT',
  paymentServiceItemParams: testParams.DomesticOriginAdditionalSIT,
  shipmentSITBalance: {
    previouslyBilledDays: 30,
    previouslyBilledEndDate: '2021-06-08',
    pendingSITDaysInvoiced: 60,
    pendingBilledEndDate: '2021-08-08',
    totalSITDaysAuthorized: 120,
    totalSITDaysRemaining: 30,
    totalSITEndDate: '2021-09-08',
  },
};

const reviewedServiceItemCard = {
  ...needsReviewServiceItemCard,
  requestComplete: true,
};

const canceledShipmentServiceItemCard = {
  ...needsReviewServiceItemCard,
  mtoShipmentModificationType: shipmentModificationTypes.CANCELED,
};

const divertedShipmentServiceItemCard = {
  ...needsReviewServiceItemCard,
  mtoShipmentModificationType: shipmentModificationTypes.DIVERSION,
};

describe('ServiceItemCard component', () => {
  describe('when payment request needs reviewed', () => {
    const wrapper = mount(<ServiceItemCard {...needsReviewServiceItemCard} />);
    it('toggles pricer calculations when button is clicked', () => {
      const toggleButton = wrapper.find('button[data-testid="toggleCalculations"]');
      expect(toggleButton.text()).toEqual('Show calculations');

      act(() => {
        toggleButton.simulate('click');
      });
      wrapper.update();

      expect(toggleButton.text()).toEqual('Hide calculations');
      expect(wrapper.find('ServiceItemCalculations').exists()).toBe(true);

      act(() => {
        toggleButton.simulate('click');
      });
      wrapper.update();

      expect(toggleButton.text()).toEqual('Show calculations');
      expect(wrapper.find('ServiceItemCalculations').exists()).toBe(false);
    });

    it('does not render calculations toggle when the service item calculations are not implemented', () => {
      const component = mount(<ServiceItemCard {...basicServiceItemCard} />);
      expect(component.find('button[data-testid="toggleCalculations"]').exists()).toBe(false);
    });

    it('displays the Days In SIT information for additional day service items', () => {
      render(<ServiceItemCard {...additionalDaySITServiceItemCard} />);
      expect(screen.getByText('SIT days invoiced')).toBeInTheDocument();
      expect(screen.getByTestId('DaysInSITAllowance')).toBeInTheDocument();
    });
  });

  describe('when Reject is selected', () => {
    it('the component displays correctly and shows asterisks for required fields', async () => {
      render(<ServiceItemCard {...basicServiceItemCard} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'BASIC SERVICE ITEMS' })).toBeInTheDocument();
      });

      expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');
      const approveButtonLabel = screen.getByLabelText(/Approve or reject the service item */);
      expect(approveButtonLabel).toBeInTheDocument();

      const approveButton = screen.getByLabelText('Approve');
      const rejectButton = screen.getByLabelText('Reject');
      expect(screen.getByLabelText('Reject')).toBeInTheDocument();
      expect(approveButton).toBeInTheDocument();

      await userEvent.click(rejectButton);
      expect(rejectButton).toBeChecked();
      expect(screen.queryByText('Add a reason why this service item is rejected')).not.toBeInTheDocument();
      expect(screen.getByLabelText('Reason for rejection *')).toBeInTheDocument();
    });
    describe('when a reason is added', () => {
      it('Approve is selected, and Reject is reselected, the reason is cleared, and no error appears', async () => {
        render(<ServiceItemCard {...basicServiceItemCard} />);

        const approveButton = screen.getByLabelText('Approve');
        const rejectButton = screen.getByLabelText('Reject');

        await userEvent.click(rejectButton);
        expect(screen.getByLabelText(/Reason for rejection/)).toBeInTheDocument();
        const reason = 'why it was rejected';
        await userEvent.type(screen.getByLabelText(/Reason for rejection/), reason);
        await userEvent.click(approveButton);
        await userEvent.click(rejectButton);
        expect(screen.queryByText('Add a reason why this service item is rejected')).not.toBeInTheDocument();
        expect(screen.queryByText(reason)).not.toBeInTheDocument();
      });
      it('and removed, and the textbox is blurred, an error is shown', async () => {
        render(<ServiceItemCard {...basicServiceItemCard} />);
        const rejectButton = screen.getByLabelText('Reject');
        await userEvent.click(rejectButton);
        expect(screen.getByLabelText(/Reason for rejection/)).toBeInTheDocument();
        await userEvent.type(screen.getByLabelText(/Reason for rejection/), 'a{backspace}');
        await userEvent.click(rejectButton);
        expect(screen.queryByText('Add a reason why this service item is rejected')).toBeInTheDocument();
      });
    });
  });

  describe('when payment request has been reviewed', () => {
    const wrapper = mount(<ServiceItemCard {...reviewedServiceItemCard} />);
    it('toggles pricer calculations when button is clicked', () => {
      const toggleButton = wrapper.find('button[data-testid="toggleCalculations"]');
      expect(toggleButton.text()).toEqual('Show calculations');

      act(() => {
        toggleButton.simulate('click');
      });
      wrapper.update();

      expect(toggleButton.text()).toEqual('Hide calculations');
      expect(wrapper.find('ServiceItemCalculations').exists()).toBe(true);

      act(() => {
        toggleButton.simulate('click');
      });
      wrapper.update();

      expect(toggleButton.text()).toEqual('Show calculations');
      expect(wrapper.find('ServiceItemCalculations').exists()).toBe(false);
    });

    it('does not render calculations toggle when the service item calculations are not implemented', () => {
      const component = mount(<ServiceItemCard {...reviewedBasicServiceItemCard} />);
      expect(component.find('button[data-testid="toggleCalculations"]').exists()).toBe(false);
    });

    it('does not display days in SIT info for additional day service items', () => {
      const reviewedDOASIT = { ...additionalDaySITServiceItemCard, requestComplete: true };

      render(<ServiceItemCard {...reviewedDOASIT} />);
      expect(screen.queryByText('SIT days invoiced')).not.toBeInTheDocument();
      expect(screen.queryByText('DaysInSITAllowance')).not.toBeInTheDocument();
    });
  });

  describe('When a service item has a shipment that was canceled ', () => {
    const component = mount(<ServiceItemCard {...canceledShipmentServiceItemCard} />);
    it('there is a canceled tag displayed', () => {
      expect(component.find('ShipmentModificationTag').text()).toBe(shipmentModificationTypes.CANCELED);
    });
  });

  describe('When a service item has a shipment that was diverted ', () => {
    const component = mount(<ServiceItemCard {...divertedShipmentServiceItemCard} />);
    it('there is a diversion tag displayed', () => {
      expect(component.find('ShipmentModificationTag').text()).toBe(shipmentModificationTypes.DIVERSION);
    });
  });
});
