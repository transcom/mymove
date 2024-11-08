import React from 'react';
import { render, screen, within } from '@testing-library/react';

import {
  shipments,
  ordersInfo,
  allowancesInfo,
  customerInfo,
  serviceItemsEmpty,
  ppmOnlyShipments,
  closeoutOffice,
} from './RequestedShipmentsTestData';
import SubmittedRequestedShipments from './SubmittedRequestedShipments';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
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
  <SubmittedRequestedShipments
    allowancesInfo={allowancesInfo}
    moveCode="TE5TC0DE"
    mtoShipments={shipments}
    closeoutOffice={closeoutOffice}
    customerInfo={customerInfo}
    ordersInfo={ordersInfo}
    approveMTO={approveMTO}
  />
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
      expect(screen.queryByTestId('services-counseling-completed-text')).not.toBeInTheDocument();
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

    it('does not render the "Add service items to move" section or Counseling option when all shipments are PPM', () => {
      const testPropsServiceItemsEmpty = {
        mtoServiceItems: serviceItemsEmpty,
        mtoShipments: ppmOnlyShipments,
        ...conditionalFormTestProps,
      };
      renderComponent(testPropsServiceItemsEmpty);

      expect(screen.queryByText('Add service items to this move')).not.toBeInTheDocument();
      expect(screen.getByText('Approve selected')).toBeInTheDocument();
      expect(screen.queryByTestId('shipmentManagementFee')).not.toBeInTheDocument();
      expect(screen.queryByTestId('counselingFee')).not.toBeInTheDocument();
    });
  });
});
