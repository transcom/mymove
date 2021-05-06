/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, fireEvent, act, waitFor, screen } from '@testing-library/react';
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
      render(<ServiceItemCard {...basicServiceItemCard} />);

      fireEvent.click(screen.getByTestId('rejectRadio'));

      waitFor(() => {
        expect(screen.getByTestId('rejectionSaveButton')).toHaveAttribute('disabled');
      });
    });

    // using react testing library to test dom interactions
    it('edit reason button link show after saving', async () => {
      render(<ServiceItemCard {...basicServiceItemCard} />);

      // Click on reject radio and fill in text area
      fireEvent.click(screen.getByTestId('rejectRadio'));
      fireEvent.change(screen.getByTestId('textarea'), {
        target: { value: 'Rejected just because.' },
      });

      expect(screen.getByTestId('rejectionSaveButton').hasAttribute('disabled')).toBeFalsy();

      // Save
      fireEvent.click(screen.getByTestId('rejectionSaveButton'));

      expect(screen.getByTestId('editReasonButton')).toBeTruthy();
      expect(screen.getByTestId('rejectionReasonReadOnly').textContent).toBe('Rejected just because.');
    });

    // using react testing library to test dom interactions
    it('edit existing rejection reason', async () => {
      const data = {
        ...basicServiceItemCard,
        status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
        rejectionReason: 'Rejected just because.',
      };
      render(<ServiceItemCard {...data} />);

      expect(screen.getByTestId('editReasonButton')).toBeTruthy();
      expect(screen.getByTestId('rejectionReasonReadOnly').textContent).toBe('Rejected just because.');

      // Click on Edit reason button, edit text area and save
      fireEvent.click(screen.getByTestId('editReasonButton'));
      fireEvent.change(screen.getByTestId('textarea'), {
        target: { value: 'Edited rejection reason.' },
      });
      fireEvent.click(screen.getByTestId('rejectionSaveButton'));

      expect(screen.getByTestId('editReasonButton')).toBeTruthy();
      expect(screen.getByTestId('rejectionReasonReadOnly').textContent).toBe('Edited rejection reason.');
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
