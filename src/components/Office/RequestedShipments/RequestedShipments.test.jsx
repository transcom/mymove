import React from 'react';
import { act } from 'react-dom/test-utils';
import { fireEvent, render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import {
  shipments,
  ntsExternalVendorShipments,
  ordersInfo,
  allowancesInfo,
  customerInfo,
  agents,
  serviceItems,
} from './RequestedShipmentsTestData';
import RequestedShipments from './RequestedShipments';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

const moveTaskOrder = {
  eTag: 'MjAyMC0wNi0yNlQyMDoyMjo0MS43Mjc4NTNa',
  id: '6e8c5ca4-774c-4170-934a-59d22259e480',
};

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

const requestedShipmentsComponent = (
  <RequestedShipments
    ordersInfo={ordersInfo}
    allowancesInfo={allowancesInfo}
    mtoAgents={agents}
    customerInfo={customerInfo}
    mtoShipments={shipments}
    approveMTO={approveMTO}
    shipmentsStatus="SUBMITTED"
    moveCode="TE5TC0DE"
  />
);

const requestedShipmentsComponentWithPermission = (
  <MockProviders permissions={[permissionTypes.updateShipment]}>
    <RequestedShipments
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      mtoAgents={agents}
      customerInfo={customerInfo}
      mtoShipments={shipments}
      approveMTO={approveMTO}
      shipmentsStatus="SUBMITTED"
      moveCode="TE5TC0DE"
    />
  </MockProviders>
);
const requestedExternalVendorShipmentsComponent = (
  <MockProviders permissions={[permissionTypes.updateShipment]}>
    <RequestedShipments
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      mtoAgents={agents}
      customerInfo={customerInfo}
      mtoShipments={ntsExternalVendorShipments}
      approveMTO={approveMTO}
      shipmentsStatus="SUBMITTED"
      moveCode="TE5TC0DE"
    />
  </MockProviders>
);

const requestedShipmentsComponentAvailableToPrimeAt = (
  <MockProviders permissions={[permissionTypes.updateShipment]}>
    <RequestedShipments
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      mtoAgents={agents}
      customerInfo={customerInfo}
      mtoShipments={shipments}
      approveMTO={approveMTO}
      shipmentsStatus="SUBMITTED"
      moveTaskOrder={moveTaskOrderAvailableToPrimeAt}
      moveCode="TE5TC0DE"
    />
  </MockProviders>
);

const requestedShipmentsComponentServicesCounselingCompleted = (
  <RequestedShipments
    ordersInfo={ordersInfo}
    allowancesInfo={allowancesInfo}
    mtoAgents={agents}
    customerInfo={customerInfo}
    mtoShipments={shipments}
    approveMTO={approveMTO}
    shipmentsStatus="SUBMITTED"
    moveTaskOrder={moveTaskOrderServicesCounselingCompleted}
    moveCode="TE5TC0DE"
  />
);

const requestedShipmentsComponentMissingRequiredInfo = (
  <MockProviders permissions={[permissionTypes.updateShipment]}>
    <RequestedShipments
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      mtoAgents={agents}
      customerInfo={customerInfo}
      mtoShipments={shipments}
      approveMTO={approveMTO}
      shipmentsStatus="SUBMITTED"
      missingRequiredOrdersInfo
      moveCode="TE5TC0DE"
    />
  </MockProviders>
);

describe('RequestedShipments', () => {
  describe('Prime-handled shipments', () => {
    it('renders the container successfully without services counseling completed', () => {
      render(requestedShipmentsComponent);
      expect(screen.getByTestId('requested-shipments')).toBeInTheDocument();
      expect(screen.queryByTestId('services-counseling-completed-text')).not.toBeInTheDocument();
    });

    it('renders the container successfully with services counseling completed', () => {
      render(requestedShipmentsComponentServicesCounselingCompleted);
      expect(screen.getByTestId('requested-shipments')).toBeInTheDocument();
      expect(screen.queryByTestId('services-counseling-completed-text')).not.toBeInTheDocument();
    });

    it('renders a shipment passed to it', () => {
      render(requestedShipmentsComponent);
      const withinContainer = within(screen.getByTestId('requested-shipments'));
      expect(withinContainer.getAllByText('HHG').length).toEqual(2);
      expect(withinContainer.getAllByText('NTS').length).toEqual(1);
    });

    it('renders the button', () => {
      render(requestedShipmentsComponentWithPermission);
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
      render(requestedShipmentsComponentAvailableToPrimeAt);
      expect(screen.getByTestId('shipmentApproveButton')).toBeInTheDocument();
      expect(screen.getByTestId('shipmentApproveButton')).toBeDisabled();
    });

    it('renders the checkboxes', () => {
      render(requestedShipmentsComponentWithPermission);
      expect(screen.getAllByTestId('checkbox').length).toEqual(5);
    });

    it('uses the duty location postal code if there is no destination address', () => {
      render(requestedShipmentsComponent);
      const destination = shipments[0].destinationAddress;
      expect(screen.getAllByTestId('destinationAddress').at(0)).toHaveTextContent(
        `${destination.streetAddress1}, ${destination.streetAddress2}, ${destination.city}, ${destination.state} ${destination.postalCode}`,
      );

      expect(screen.getAllByTestId('destinationAddress').at(1)).toHaveTextContent(
        ordersInfo.newDutyLocation.address.postalCode,
      );
    });

    it('enables the Approve selected button when a shipment and service item are checked', async () => {
      const { container } = render(requestedShipmentsComponentWithPermission);

      // TODO this doesn't seem right
      await act(async () => {
        await userEvent.type(
          container.querySelector('input[name="shipments"]'),
          'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
        );
      });

      expect(screen.getByRole('button', { name: 'Approve selected' })).toBeDisabled();
      expect(container.querySelector('#approvalConfirmationModal')).toHaveStyle('display: none');

      // TODO
      await act(async () => {
        await userEvent.click(screen.getByRole('checkbox', { name: 'Move management' }));
      });

      expect(screen.getByRole('button', { name: 'Approve selected' })).not.toBeDisabled();

      // TODO
      await act(async () => {
        await userEvent.click(screen.getByRole('button', { name: 'Approve selected' }));
      });
      expect(container.querySelector('#approvalConfirmationModal')).toHaveStyle('display: block');
    });

    it('disables the Approve selected button when there is missing required information', async () => {
      const { container } = render(requestedShipmentsComponentMissingRequiredInfo);

      // TODO
      await act(async () => {
        await userEvent.type(
          container.querySelector('input[name="shipments"]'),
          'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
        );
      });

      expect(screen.getByRole('button', { name: 'Approve selected' })).toBeDisabled();

      await act(async () => {
        await userEvent.click(screen.getByRole('checkbox', { name: 'Move management' }));
      });

      expect(screen.getByRole('button', { name: 'Approve selected' })).toBeDisabled();
    });

    it('calls approveMTO onSubmit', async () => {
      const mockOnSubmit = jest.fn((id, eTag) => {
        return new Promise((resolve) => {
          resolve({ response: { status: 200, body: { id, eTag } } });
        });
      });

      const { container } = render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <RequestedShipments
            mtoShipments={shipments}
            mtoAgents={agents}
            ordersInfo={ordersInfo}
            allowancesInfo={allowancesInfo}
            customerInfo={customerInfo}
            moveTaskOrder={moveTaskOrder}
            approveMTO={mockOnSubmit}
            shipmentsStatus="SUBMITTED"
            moveCode="TE5TC0DE"
          />
        </MockProviders>,
      );

      // You could take the shortcut and call submit directly as well if providing initial values
      fireEvent.submit(container.querySelector('form'));

      // When simulating change events you must pass the target with the id and
      // name for formik to know which value to update

      await act(async () => {
        const shipmentInput = container.querySelector('input[name="shipments"]');
        await userEvent.type(shipmentInput, 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee');

        const shipmentManagementFeeInput = screen.getByRole('checkbox', { name: 'Move management' });
        await userEvent.click(shipmentManagementFeeInput);

        const counselingFeeInput = screen.getByRole('checkbox', { name: 'Counseling' });
        await userEvent.click(counselingFeeInput);

        userEvent.click(container.querySelector('form button[type="button"]'));
        userEvent.click(screen.getAllByRole('button').at(0));
      });

      expect(mockOnSubmit).toHaveBeenCalled();
      expect(mockOnSubmit.mock.calls[0]).toEqual([
        {
          moveTaskOrderID: moveTaskOrder.id,
          ifMatchETag: moveTaskOrder.eTag,
          mtoApprovalServiceItemCodes: {
            serviceCodeCS: true,
            serviceCodeMS: true,
          },
          normalize: false,
        },
      ]);
    });

    it('displays approved basic service items for approved shipments', () => {
      render(
        <RequestedShipments
          ordersInfo={ordersInfo}
          allowancesInfo={allowancesInfo}
          mtoAgents={agents}
          customerInfo={customerInfo}
          mtoShipments={shipments}
          approveMTO={approveMTO}
          shipmentsStatus="APPROVED"
          mtoServiceItems={serviceItems}
          moveCode="TE5TC0DE"
        />,
      );
      const approvedServiceItemNames = screen.getAllByTestId('basicServiceItemName');
      const approvedServiceItemDates = screen.getAllByTestId('basicServiceItemDate');

      expect(approvedServiceItemNames.length).toBe(2);
      expect(approvedServiceItemDates.length).toBe(2);

      expect(approvedServiceItemNames.at(0).textContent).toEqual('Move management');
      expect(approvedServiceItemDates.at(0).textContent).toEqual(' 01 Jan 2020');

      expect(approvedServiceItemNames.at(1).textContent).toEqual('Counseling fee');
      expect(approvedServiceItemDates.at(1).textContent).toEqual(' 01 Jan 2020');
    });

    it.each([['APPROVED'], ['SUBMITTED']])(
      'displays the customer and counselor remarks for a(n) %s shipment',
      (status) => {
        render(
          <RequestedShipments
            ordersInfo={ordersInfo}
            allowancesInfo={allowancesInfo}
            mtoAgents={agents}
            customerInfo={customerInfo}
            mtoShipments={shipments}
            approveMTO={approveMTO}
            shipmentsStatus={status}
            mtoServiceItems={serviceItems}
            moveCode="TE5TC0DE"
          />,
        );

        const customerRemarks = screen.getAllByTestId('customerRemarks');
        const counselorRemarks = screen.getAllByTestId('counselorRemarks');

        expect(customerRemarks.at(0).textContent).toBe('please treat gently');
        expect(customerRemarks.at(1).textContent).toBe('please treat gently');

        expect(counselorRemarks.at(0).textContent).toBe('looks good');
        expect(counselorRemarks.at(1).textContent).toBe('looks good');
      },
    );
  });

  describe('External vendor shipments', () => {
    it('enables the Approve selected button when there is only external vendor shipments and a service item is checked', async () => {
      render(requestedExternalVendorShipmentsComponent);

      expect(screen.getByTestId('shipmentApproveButton')).toBeDisabled();

      await userEvent.click(screen.getByLabelText('Move management'));

      expect(screen.getByLabelText('Move management').checked).toEqual(true);

      expect(screen.getByTestId('shipmentApproveButton')).toBeEnabled();
    });
  });

  describe('Permission dependent rendering', () => {
    const testProps = {
      ordersInfo,
      allowancesInfo,
      mtoAgents: agents,
      customerInfo,
      mtoShipments: shipments,
      approveMTO,
      shipmentsStatus: 'SUBMITTED',
      mtoServiceItems: serviceItems,
      moveCode: 'TE5TC0DE',
    };
    it('renders the "Add service items to move" section when user has permission', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <RequestedShipments {...testProps} />
        </MockProviders>,
      );

      expect(screen.getByText('Add service items to this move')).toBeInTheDocument();
      expect(screen.getByText('Approve selected')).toBeInTheDocument();
    });

    it('does not render the "Add service items to move" section when user does not have permission', () => {
      render(
        <MockProviders permissions={[]}>
          <RequestedShipments {...testProps} />
        </MockProviders>,
      );

      expect(screen.queryByText('Add service items to this move')).not.toBeInTheDocument();
      expect(screen.queryByText('Approve selected')).not.toBeInTheDocument();
    });
  });
});
