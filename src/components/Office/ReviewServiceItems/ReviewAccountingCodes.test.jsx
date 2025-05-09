import React from 'react';
import { screen, render, within } from '@testing-library/react';

import ReviewAccountingCodes from './ReviewAccountingCodes';

import { LOA_TYPE, PAYMENT_SERVICE_ITEM_STATUS, SHIPMENT_OPTIONS } from 'shared/constants';

describe('components/Office/ReviewServiceItems/ReviewAccountingCodes', () => {
  describe('can display nothing if there are no valid service items', () => {
    it('should not display service items when cards is empty', () => {
      render(
        <ReviewAccountingCodes TACs={{ HHG: '1234', NTS: '5678' }} SACs={{ HHG: 'AB12', NTS: 'CD34' }} cards={[]} />,
      );

      expect(screen.queryByRole('heading', { level: 4, name: 'Accounting codes' })).not.toBeInTheDocument();
    });

    it('should not display service items that are rejected', () => {
      render(
        <ReviewAccountingCodes
          TACs={{ HHG: '1234', NTS: '5678' }}
          SACs={{ HHG: 'AB12', NTS: 'CD34' }}
          cards={[
            {
              amount: 0.01,
              mtoShipmentID: 'X',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG,
              status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
            },
          ]}
        />,
      );

      expect(screen.queryByRole('heading', { level: 4, name: 'Accounting codes' })).not.toBeInTheDocument();
    });

    it('should not display shipment card for service items not attached to a shipment', () => {
      render(
        <ReviewAccountingCodes
          TACs={{ HHG: '1234', NTS: '5678' }}
          SACs={{ HHG: 'AB12', NTS: 'CD34' }}
          cards={[
            {
              amount: 0.01,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
            },
          ]}
        />,
      );

      expect(screen.queryByText('HHG')).not.toBeInTheDocument();
    });

    it('should not display move management fee if move management service item is not requested', () => {
      render(
        <ReviewAccountingCodes
          cards={[
            {
              amount: 20.65,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoServiceItemName: 'Counseling',
            },
          ]}
        />,
      );

      expect(screen.queryByText('Move management fee')).not.toBeInTheDocument();
      expect(screen.getByText('Counseling fee')).toBeInTheDocument();
      expect(screen.getByText('$20.65')).toBeInTheDocument();
    });

    it('should not display counseling fee if counseling service item is not requested', () => {
      render(
        <ReviewAccountingCodes
          cards={[
            {
              amount: 44.33,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoServiceItemName: 'Move management',
            },
          ]}
        />,
      );

      expect(screen.queryByText('Counseling fee')).not.toBeInTheDocument();
      expect(screen.getByText('Move management fee')).toBeInTheDocument();
      expect(screen.getByText('$44.33')).toBeInTheDocument();
    });
  });

  describe('can display codes', () => {
    it('can display a single shipment card', () => {
      render(
        <ReviewAccountingCodes
          TACs={{ HHG: '1234', NTS: '5678' }}
          SACs={{ HHG: 'AB12', NTS: 'CD34' }}
          cards={[
            {
              amount: 10,
              mtoShipmentID: '1',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoShipmentTacType: LOA_TYPE.HHG,
            },
            {
              amount: 20,
              mtoShipmentID: '1',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoShipmentTacType: LOA_TYPE.HHG,
            },
          ]}
        />,
      );

      expect(screen.getByRole('heading', { level: 4, name: 'Accounting codes' })).toBeInTheDocument();
      expect(screen.getByText('TAC: 1234 (HHG)')).toBeInTheDocument();
      expect(screen.getByText('$30.00')).toBeInTheDocument();
    });

    it('can display a multiple shipment cards', () => {
      render(
        <ReviewAccountingCodes
          TACs={{ HHG: '1234', NTS: '5678' }}
          SACs={{ HHG: 'AB12', NTS: 'CD34' }}
          cards={[
            {
              amount: 10,
              mtoShipmentID: '1',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoShipmentTacType: LOA_TYPE.HHG,
            },
            {
              amount: 20,
              mtoShipmentID: '2',
              mtoShipmentType: SHIPMENT_OPTIONS.NTS,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoShipmentTacType: LOA_TYPE.NTS,
              mtoShipmentSacType: LOA_TYPE.HHG,
            },
            {
              amount: 30,
              mtoShipmentID: '3',
              mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoShipmentTacType: LOA_TYPE.NTS,
              mtoShipmentSacType: LOA_TYPE.NTS,
            },
          ]}
        />,
      );

      const firstShipment = screen.getByTestId('shipment-1');
      expect(within(firstShipment).getByText('HHG')).toBeInTheDocument();
      expect(within(firstShipment).queryByText('SAC: AB12 (HHG)')).not.toBeInTheDocument();
      expect(within(firstShipment).getByText('$10.00')).toBeInTheDocument();

      const secondShipment = screen.getByTestId('shipment-2');
      expect(within(secondShipment).getByText('NTS')).toBeInTheDocument();
      expect(within(secondShipment).getByText('TAC: 5678 (NTS)')).toBeInTheDocument();
      expect(within(secondShipment).getByText('SAC: AB12 (HHG)')).toBeInTheDocument();
      expect(within(secondShipment).getByText('$20.00')).toBeInTheDocument();

      const thirdShipment = screen.getByTestId('shipment-3');
      expect(within(thirdShipment).getByText('NTS-release')).toBeInTheDocument();
      expect(within(thirdShipment).getByText('TAC: 5678 (NTS)')).toBeInTheDocument();
      expect(within(thirdShipment).getByText('SAC: CD34 (NTS)')).toBeInTheDocument();
      expect(within(thirdShipment).getByText('$30.00')).toBeInTheDocument();
    });

    it('can display a move level service item card', () => {
      render(
        <ReviewAccountingCodes
          cards={[
            {
              amount: 44.33,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoServiceItemName: 'Move management',
            },
            {
              amount: 20,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoServiceItemName: 'Counseling',
            },
          ]}
        />,
      );

      expect(screen.getByRole('heading', { level: 4, name: 'Accounting codes' })).toBeInTheDocument();
      expect(screen.getByText('Move management fee')).toBeInTheDocument();
      expect(screen.getByText('$44.33')).toBeInTheDocument();
      expect(screen.getByText('Counseling fee')).toBeInTheDocument();
      expect(screen.getByText('$20.00')).toBeInTheDocument();
    });

    it('can display a move level service item card and multiple shipment cards', () => {
      render(
        <ReviewAccountingCodes
          TACs={{ HHG: '1234', NTS: '5678' }}
          SACs={{ HHG: 'AB12', NTS: 'CD34' }}
          cards={[
            {
              amount: 44.33,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoServiceItemName: 'Move management',
            },
            {
              amount: 20.65,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoServiceItemName: 'Counseling',
            },
            {
              amount: 10,
              mtoShipmentID: '1',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoShipmentTacType: LOA_TYPE.HHG,
            },
            {
              amount: 20,
              mtoShipmentID: '2',
              mtoShipmentType: SHIPMENT_OPTIONS.NTS,
              status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
              mtoShipmentTacType: LOA_TYPE.NTS,
              mtoShipmentSacType: LOA_TYPE.HHG,
            },
          ]}
        />,
      );

      expect(screen.getByRole('heading', { level: 4, name: 'Accounting codes' })).toBeInTheDocument();
      const firstShipment = screen.getByTestId('shipment-1');
      expect(within(firstShipment).getByText('HHG')).toBeInTheDocument();
      expect(within(firstShipment).queryByText('SAC: AB12 (HHG)')).not.toBeInTheDocument();
      expect(within(firstShipment).getByText('$10.00')).toBeInTheDocument();
      const secondShipment = screen.getByTestId('shipment-2');
      expect(within(secondShipment).getByText('NTS')).toBeInTheDocument();
      expect(within(secondShipment).getByText('TAC: 5678 (NTS)')).toBeInTheDocument();
      expect(within(secondShipment).getByText('SAC: AB12 (HHG)')).toBeInTheDocument();
      expect(within(secondShipment).getByText('$20.00')).toBeInTheDocument();
      expect(screen.getByText('Move management fee')).toBeInTheDocument();
      expect(screen.getByText('$44.33')).toBeInTheDocument();
      expect(screen.getByText('Counseling fee')).toBeInTheDocument();
      expect(screen.getByText('$20.65')).toBeInTheDocument();
    });
  });
});
