import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import {
  shipments,
  ntsExternalVendorShipments,
  ordersInfo,
  allowancesInfo,
  customerInfo,
  serviceItemsMSandCS,
  serviceItemsMS,
  serviceItemsCS,
  serviceItemsEmpty,
  ppmOnlyShipments,
} from './RequestedShipmentsTestData';
import ApprovedRequestedShipments from './ApprovedRequestedShipments';
import SubmittedRequestedShipments from './SubmittedRequestedShipments';

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

const submittedRequestedShipmentsComponent = (
  <SubmittedRequestedShipments
    allowancesInfo={allowancesInfo}
    moveCode="TE5TC0DE"
    mtoShipments={shipments}
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
      approveMTO={approveMTO}
      moveCode="TE5TC0DE"
    />
  </MockProviders>
);

const submittedRequestedExternalVendorShipmentsComponent = (
  <MockProviders permissions={[permissionTypes.updateShipment]}>
    <SubmittedRequestedShipments
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoShipments={ntsExternalVendorShipments}
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
    approveMTO={approveMTO}
    moveTaskOrder={moveTaskOrderServicesCounselingCompleted}
    moveCode="TE5TC0DE"
  />
);

const submittedRequestedShipmentsComponentMissingRequiredInfo = (
  <MockProviders permissions={[permissionTypes.updateShipment]}>
    <SubmittedRequestedShipments
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoShipments={shipments}
      approveMTO={approveMTO}
      missingRequiredOrdersInfo
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

    it('uses the duty location postal code if there is no destination address', () => {
      render(submittedRequestedShipmentsComponent);
      const destination = shipments[0].destinationAddress;
      expect(screen.getAllByTestId('destinationAddress').at(0)).toHaveTextContent(
        `${destination.streetAddress1}, ${destination.streetAddress2}, ${destination.city}, ${destination.state} ${destination.postalCode}`,
      );

      expect(screen.getAllByTestId('destinationAddress').at(1)).toHaveTextContent(
        ordersInfo.newDutyLocation.address.postalCode,
      );
    });

    it('enables the Approve selected button when a shipment and service item are checked', async () => {
      const { container } = render(submittedRequestedShipmentsComponentWithPermission);

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
      const { container } = render(submittedRequestedShipmentsComponentMissingRequiredInfo);

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
          <SubmittedRequestedShipments
            mtoShipments={shipments}
            ordersInfo={ordersInfo}
            allowancesInfo={allowancesInfo}
            customerInfo={customerInfo}
            moveTaskOrder={moveTaskOrder}
            approveMTO={mockOnSubmit}
            moveCode="TE5TC0DE"
          />
        </MockProviders>,
      );

      await userEvent.click(screen.getByRole('button', { name: 'Approve selected' }));

      const shipmentInput = container.querySelector('input[name="shipments"]');
      await userEvent.type(shipmentInput, 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee');

      const shipmentManagementFeeInput = screen.getByRole('checkbox', { name: 'Move management' });
      await userEvent.click(shipmentManagementFeeInput);

      const counselingFeeInput = screen.getByRole('checkbox', { name: 'Counseling' });
      await userEvent.click(counselingFeeInput);

      await userEvent.click(screen.getByText('Approve and send'));

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
        {
          onSuccess: expect.any(Function),
          onError: expect.any(Function),
        },
      ]);
    });

    it('displays approved basic service items for approved shipments', () => {
      render(
        <ApprovedRequestedShipments
          ordersInfo={ordersInfo}
          mtoShipments={shipments}
          mtoServiceItems={serviceItemsMSandCS}
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
        const statusTestProps = {
          APPROVED: {
            ordersInfo,
            mtoShipments: shipments,
            mtoServiceItems: serviceItemsMSandCS,
            moveCode: 'TE5TC0DE',
          },
          SUBMITTED: {
            ordersInfo,
            allowancesInfo,
            customerInfo,
            mtoShipments: shipments,
            approveMTO,
            mtoServiceItems: serviceItemsMSandCS,
            moveCode: 'TE5TC0DE',
          },
        };

        const statusComponents = {
          APPROVED: ApprovedRequestedShipments,
          SUBMITTED: SubmittedRequestedShipments,
        };

        const Component = statusComponents[status];

        render(<Component {...statusTestProps[status]} />);

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
      render(submittedRequestedExternalVendorShipmentsComponent);

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
      customerInfo,
      mtoShipments: shipments,
      approveMTO,
      mtoServiceItems: [],
      moveCode: 'TE5TC0DE',
    };
    it('renders the "Add service items to move" section when user has permission', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <SubmittedRequestedShipments {...testProps} />
        </MockProviders>,
      );

      expect(screen.getByText('Add service items to this move')).toBeInTheDocument();
      expect(screen.getByText('Approve selected')).toBeInTheDocument();
    });

    it('does not render the "Add service items to move" section when user does not have permission', () => {
      render(
        <MockProviders permissions={[]}>
          <SubmittedRequestedShipments {...testProps} />
        </MockProviders>,
      );

      expect(screen.queryByText('Add service items to this move')).not.toBeInTheDocument();
      expect(screen.queryByText('Approve selected')).not.toBeInTheDocument();
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
    it('does not render the "Add service items to move" section when both service items are present', () => {
      const testProps = {
        mtoServiceItems: serviceItemsMSandCS,
        mtoShipments: shipments,
        ...conditionalFormTestProps,
      };
      renderComponent(testProps);

      expect(screen.queryByText('Add service items to this move')).not.toBeInTheDocument();
      expect(screen.getByText('Approve selected')).toBeInTheDocument();
    });

    it('does not render the "Add service items to move" section when counseling is present and all shipments are PPM', () => {
      const testProps = {
        mtoServiceItems: serviceItemsCS,
        mtoShipments: ppmOnlyShipments,
        ...conditionalFormTestProps,
      };
      renderComponent(testProps);

      expect(screen.queryByText('Add service items to this move')).not.toBeInTheDocument();
      expect(screen.getByText('Approve selected')).toBeInTheDocument();
    });

    it('renders the "Add service items to move" section with only counseling when only move management is present in service items', () => {
      const testProps = {
        mtoServiceItems: serviceItemsMS,
        mtoShipments: shipments,
        ...conditionalFormTestProps,
      };
      renderComponent(testProps);

      expect(screen.getByText('Add service items to this move')).toBeInTheDocument();
      expect(screen.getByText('Approve selected')).toBeInTheDocument();
      expect(screen.queryByTestId('shipmentManagementFee')).not.toBeInTheDocument();
      expect(screen.getByTestId('counselingFee')).toBeInTheDocument();
    });

    it('renders the "Add service items to move" section with only move management when only counseling is present in service items', () => {
      const testProps = {
        mtoServiceItems: serviceItemsCS,
        mtoShipments: shipments,
        ...conditionalFormTestProps,
      };
      renderComponent(testProps);

      expect(screen.getByText('Add service items to this move')).toBeInTheDocument();
      expect(screen.getByText('Approve selected')).toBeInTheDocument();
      expect(screen.getByTestId('shipmentManagementFee')).toBeInTheDocument();
      expect(screen.queryByTestId('counselingFee')).not.toBeInTheDocument();
    });

    it('renders the "Add service items to move" section with all fields when neither counseling nor move management is present in service items', () => {
      const testProps = {
        mtoServiceItems: serviceItemsEmpty,
        mtoShipments: shipments,
        ...conditionalFormTestProps,
      };
      renderComponent(testProps);

      expect(screen.getByText('Add service items to this move')).toBeInTheDocument();
      expect(screen.getByText('Approve selected')).toBeInTheDocument();
      expect(screen.getByTestId('shipmentManagementFee')).toBeInTheDocument();
      expect(screen.getByTestId('counselingFee')).toBeInTheDocument();
    });

    it('renders the "Add service items to move" section with only counseling when all shipments are PPM', () => {
      const testProps = {
        mtoServiceItems: serviceItemsEmpty,
        mtoShipments: ppmOnlyShipments,
        ...conditionalFormTestProps,
      };
      renderComponent(testProps);

      expect(screen.getByText('Add service items to this move')).toBeInTheDocument();
      expect(screen.getByText('Approve selected')).toBeInTheDocument();
      expect(screen.queryByTestId('shipmentManagementFee')).not.toBeInTheDocument();
      expect(screen.getByTestId('counselingFee')).toBeInTheDocument();
    });
  });
});
