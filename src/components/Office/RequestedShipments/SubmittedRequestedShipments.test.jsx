import React from 'react';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import {
  shipments,
  ordersInfo,
  allowancesInfo,
  customerInfo,
  serviceItemsEmpty,
  ppmOnlyShipments,
  closeoutOffice,
  ordersInfoOCONUS,
  ordersInfoOCONUSLocalMove,
} from './RequestedShipmentsTestData';
import SubmittedRequestedShipments from './SubmittedRequestedShipments';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const moveTaskOrderAvailableToPrimeAt = {
  eTag: 'MjAyMC0wNi0yNlQyMDoyMjo0MS43Mjc4NTNa',
  id: '6e8c5ca4-774c-4170-934a-59d22259e480',
  availableToPrimeAt: '2020-06-10T15:58:02.431995Z',
};

const moveTaskOrderServicesCounselingCompleted = {
  eTag: 'MjAyMC0wNi0yNlQyMDoyMjo0MS43Mjc4NTNa',
  id: '6e8c5ca4-774c-4170-934a-59d22259e480',
  serviceCounselingCompletedAt: '2020-10-02T19:20:08.481139Z',
};

const approveMTO = jest.fn().mockResolvedValue({ response: { status: 200 } });

const submittedRequestedShipmentsComponent = (
  <MockProviders permissions={[permissionTypes.updateShipment]}>
    <SubmittedRequestedShipments
      allowancesInfo={allowancesInfo}
      moveCode="TE5TC0DE"
      mtoShipments={shipments}
      closeoutOffice={closeoutOffice}
      customerInfo={customerInfo}
      ordersInfo={ordersInfo}
      approveMTO={approveMTO}
    />
  </MockProviders>
);

const submittedRequestedShipmentsComponentWithPermission = (
  <MockProviders permissions={[permissionTypes.updateShipment]}>
    <SubmittedRequestedShipments
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoShipments={shipments}
      closeoutOffice={closeoutOffice}
      approveMTO={approveMTO}
      moveCode="TE5TC0DE"
    />
  </MockProviders>
);

const submittedRequestedShipmentsComponentAvailableToPrimeAt = (
  <MockProviders permissions={[permissionTypes.updateShipment]}>
    <SubmittedRequestedShipments
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoShipments={shipments}
      closeoutOffice={closeoutOffice}
      approveMTO={approveMTO}
      moveTaskOrder={moveTaskOrderAvailableToPrimeAt}
      moveCode="TE5TC0DE"
    />
  </MockProviders>
);

const submittedRequestedShipmentsComponentServicesCounselingCompleted = (
  <MockProviders permissions={[permissionTypes.updateShipment]}>
    <SubmittedRequestedShipments
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoShipments={shipments}
      closeoutOffice={closeoutOffice}
      approveMTO={approveMTO}
      moveTaskOrder={moveTaskOrderServicesCounselingCompleted}
      moveCode="TE5TC0DE"
    />
  </MockProviders>
);

const submittedRequestedShipmentsCanCreateNewShipment = (
  <MockProviders permissions={[permissionTypes.createTxoShipment]}>
    <SubmittedRequestedShipments
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoShipments={shipments}
      closeoutOffice={closeoutOffice}
      approveMTO={approveMTO}
      moveTaskOrder={moveTaskOrderServicesCounselingCompleted}
      moveCode="TE5TC0DE"
    />
  </MockProviders>
);

const submittedRequestedShipmentsCanCreateNewShipmentOCONUS = (
  <MockProviders permissions={[permissionTypes.createTxoShipment]}>
    <SubmittedRequestedShipments
      ordersInfo={ordersInfoOCONUS}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoShipments={shipments}
      closeoutOffice={closeoutOffice}
      approveMTO={approveMTO}
      moveTaskOrder={moveTaskOrderServicesCounselingCompleted}
      moveCode="TE5TC0DE"
    />
  </MockProviders>
);

const submittedRequestedShipmentsCanCreateNewShipmentOCONUSLocalMove = (
  <MockProviders permissions={[permissionTypes.createTxoShipment]}>
    <SubmittedRequestedShipments
      ordersInfo={ordersInfoOCONUSLocalMove}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoShipments={shipments}
      closeoutOffice={closeoutOffice}
      approveMTO={approveMTO}
      moveTaskOrder={moveTaskOrderServicesCounselingCompleted}
      moveCode="TE5TC0DE"
    />
  </MockProviders>
);

describe('RequestedShipments', () => {
  describe('Prime-handled shipments', () => {
    it('renders the container successfully without services counseling completed', () => {
      render(submittedRequestedShipmentsComponent);
      expect(screen.getByTestId('requested-shipments')).toBeInTheDocument();
      expect(screen.queryByTestId('services-counseling-completed-text')).not.toBeInTheDocument();
    });

    it('renders the container successfully with services counseling completed', () => {
      render(submittedRequestedShipmentsComponentServicesCounselingCompleted);
      expect(screen.getByTestId('requested-shipments')).toBeInTheDocument();
      expect(screen.queryByTestId('services-counseling-completed-text')).toBeInTheDocument();
    });

    it('renders a shipment passed to it', () => {
      render(submittedRequestedShipmentsComponent);
      const withinContainer = within(screen.getByTestId('requested-shipments'));
      expect(withinContainer.getAllByText('HHG').length).toEqual(2);
      expect(withinContainer.getAllByText('NTS').length).toEqual(1);
    });

    it('renders the button', () => {
      render(submittedRequestedShipmentsComponentWithPermission);
      expect(
        screen.getByRole('button', {
          name: 'Approve selected',
        }),
      ).toBeInTheDocument();
      expect(
        screen.getByRole('button', {
          name: 'Approve selected',
        }),
      ).toBeDisabled();
    });

    it('renders the button when it is available to the prime', () => {
      render(submittedRequestedShipmentsComponentAvailableToPrimeAt);
      expect(screen.getByTestId('shipmentApproveButton')).toBeInTheDocument();
      expect(screen.getByTestId('shipmentApproveButton')).toBeDisabled();
    });

    it('renders the checkboxes', () => {
      render(submittedRequestedShipmentsComponentWithPermission);
      expect(screen.getAllByTestId('checkbox').length).toEqual(5);
    });

    it('renders Add a new shipment Button', async () => {
      render(submittedRequestedShipmentsCanCreateNewShipment);

      expect(await screen.getByRole('combobox', { name: 'Add a new shipment' })).toBeInTheDocument();
    });

    it('renders Add a new shipment Button and displays shipment options when clicked', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      render(submittedRequestedShipmentsCanCreateNewShipmentOCONUS);

      // Get the combobox (dropdown button)
      const combobox = await screen.getByRole('combobox', { name: 'Add a new shipment' });

      expect(combobox).toBeInTheDocument();

      // Simulate a user clicking the dropdown
      await userEvent.click(combobox);

      // Check if all expected options appear
      await waitFor(() => {
        expect(screen.getByRole('option', { name: 'HHG' })).toBeInTheDocument();
        expect(screen.getByRole('option', { name: 'PPM' })).toBeInTheDocument();
        expect(screen.getByRole('option', { name: 'NTS' })).toBeInTheDocument();
        expect(screen.getByRole('option', { name: 'NTS-release' })).toBeInTheDocument();
        expect(screen.getByRole('option', { name: 'Boat' })).toBeInTheDocument();
        expect(screen.getByRole('option', { name: 'Mobile Home' })).toBeInTheDocument();
        expect(screen.getByRole('option', { name: 'UB' })).toBeInTheDocument();
      });
    });

    it('renders Add a new shipment Button and does not show UB when orders type is local move', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      render(submittedRequestedShipmentsCanCreateNewShipmentOCONUSLocalMove);

      // Get the combobox (dropdown button)
      const combobox = await screen.getByRole('combobox', { name: 'Add a new shipment' });

      expect(combobox).toBeInTheDocument();

      // Simulate a user clicking the dropdown
      await userEvent.click(combobox);

      // Check if all expected options appear
      await waitFor(() => {
        expect(screen.getByRole('option', { name: 'HHG' })).toBeInTheDocument();
        expect(screen.getByRole('option', { name: 'PPM' })).toBeInTheDocument();
        expect(screen.getByRole('option', { name: 'NTS' })).toBeInTheDocument();
        expect(screen.getByRole('option', { name: 'NTS-release' })).toBeInTheDocument();
        expect(screen.getByRole('option', { name: 'Boat' })).toBeInTheDocument();
        expect(screen.getByRole('option', { name: 'Mobile Home' })).toBeInTheDocument();
      });
      // UB option does not appear when orders type is local move
      expect(screen.queryByRole('option', { name: 'UB' })).not.toBeInTheDocument();
    });
  });

  describe('Conditional form display', () => {
    const renderComponent = (props) => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <SubmittedRequestedShipments {...props} />
        </MockProviders>,
      );
    };
    const conditionalFormTestProps = {
      ordersInfo,
      allowancesInfo,
      customerInfo,
      approveMTO,
      moveCode: 'TE5TC0DE',
    };

    it('should disable the counseling checkbox when all shipments are PPM', () => {
      const testPropsServiceItemsEmpty = {
        mtoServiceItems: serviceItemsEmpty,
        mtoShipments: ppmOnlyShipments,
        ...conditionalFormTestProps,
      };
      renderComponent(testPropsServiceItemsEmpty);

      expect(screen.queryByText('Add service items to this move')).toBeInTheDocument();
      expect(screen.getByText('Approve selected')).toBeInTheDocument();
      expect(screen.queryByTestId('counselingFee')).toBeDisabled();
    });
  });
});
