/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, fireEvent, act, waitFor } from '@testing-library/react';
import { mount } from 'enzyme';

import testParams from '../ServiceItemCalculations/serviceItemTestParams';

import ServiceItemCard from './ServiceItemCard';

import { SHIPMENT_OPTIONS, PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { serviceItemCodes } from 'content/serviceItems';

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
  mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
  mtoShipmentID: '2',
  mtoServiceItemName: serviceItemCodes.FSC,
  mtoServiceItemCode: 'FSC',
  amount: 1000,
  createdAt: '2020-01-01T00:02:00.999Z',
  status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
  paymentServiceItemParams: testParams.FuelSurchage,
  patchPaymentServiceItem: jest.fn(),
};

const reviewedServiceItemCard = {
  ...needsReviewServiceItemCard,
  requestComplete: true,
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

    // using react testing library to test dom interactions
    it('save button is disabled when rejection reason is empty', () => {
      const { getByTestId } = render(<ServiceItemCard {...basicServiceItemCard} />);
      const rejectRadio = getByTestId('rejectRadio');

      waitFor(() => {
        fireEvent.click(rejectRadio);
      });

      expect(getByTestId('rejectionSaveButton')).toHaveAttribute('disabled');
    });

    // using react testing library to test dom interactions
    it('edit reason button link show after saving', async () => {
      const { getByTestId } = render(<ServiceItemCard {...basicServiceItemCard} />);
      const rejectRadio = getByTestId('rejectRadio');

      // Click on reject radio and fill in text area
      await waitFor(() => {
        fireEvent.click(rejectRadio);
      });

      await waitFor(() => {
        fireEvent.change(getByTestId('textarea'), {
          target: { value: 'Rejected just because.' },
        });
      });

      expect(getByTestId('rejectionSaveButton').hasAttribute('disabled')).toBeFalsy();

      // Save
      await waitFor(() => {
        fireEvent.click(getByTestId('rejectionSaveButton'));
      });

      expect(getByTestId('editReasonButton')).toBeTruthy();
      expect(getByTestId('rejectionReasonReadOnly').textContent).toBe('Rejected just because.');
    });

    // using react testing library to test dom interactions
    it('edit existing rejection reason', async () => {
      const data = {
        ...basicServiceItemCard,
        status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
        rejectionReason: 'Rejected just because.',
      };
      const { getByTestId } = render(<ServiceItemCard {...data} />);

      expect(getByTestId('editReasonButton')).toBeTruthy();
      expect(getByTestId('rejectionReasonReadOnly').textContent).toBe('Rejected just because.');

      // Click on Edit reason button, edit text area and save
      await waitFor(() => {
        fireEvent.click(getByTestId('editReasonButton'));
      });

      await waitFor(() => {
        fireEvent.change(getByTestId('textarea'), {
          target: { value: 'Edited rejection reason.' },
        });
      });

      await waitFor(() => {
        fireEvent.click(getByTestId('rejectionSaveButton'));
      });

      expect(getByTestId('editReasonButton')).toBeTruthy();
      expect(getByTestId('rejectionReasonReadOnly').textContent).toBe('Edited rejection reason.');
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
  });
});
