import React from 'react';
import { act } from 'react-dom/test-utils';
import { mount, shallow } from 'enzyme';
import { render, screen } from '@testing-library/react';
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
      const wrapper = shallow(requestedShipmentsComponent);
      expect(wrapper.find('div[data-testid="requested-shipments"]').exists()).toBe(true);
      expect(wrapper.find('p[data-testid="services-counseling-completed-text"]').exists()).toBe(false);
    });

    it('renders the container successfully with services counseling completed', () => {
      const wrapper = shallow(requestedShipmentsComponentServicesCounselingCompleted);
      expect(wrapper.find('div[data-testid="requested-shipments"]').exists()).toBe(true);
      expect(wrapper.find('p[data-testid="services-counseling-completed-text"]').exists()).toBe(true);
    });

    it('renders a shipment passed to it', () => {
      const wrapper = mount(requestedShipmentsComponent);
      expect(wrapper.find('div[data-testid="requested-shipments"]').text()).toContain('HHG');
      expect(wrapper.find('div[data-testid="requested-shipments"]').text()).toContain('NTS');
    });

    it('renders the button', () => {
      const wrapper = mount(requestedShipmentsComponentWithPermission);
      const approveButton = wrapper.find('button[data-testid="shipmentApproveButton"]');
      expect(approveButton.exists()).toBe(true);
      expect(approveButton.text()).toContain('Approve selected');
      expect(approveButton.html()).toContain('disabled=""');
    });

    it('renders the button when it is available to the prime', () => {
      const wrapper = mount(requestedShipmentsComponentAvailableToPrimeAt);
      const approveButton = wrapper.find('button[data-testid="shipmentApproveButton"]');
      expect(approveButton.html()).toContain('disabled=""');
    });

    it('renders the checkboxes', () => {
      const wrapper = mount(requestedShipmentsComponentWithPermission);
      expect(wrapper.find('div[data-testid="checkbox"]').exists()).toBe(true);
      expect(wrapper.find('div[data-testid="checkbox"]').length).toEqual(5);
    });

    it('uses the duty location postal code if there is no destination address', () => {
      const wrapper = mount(requestedShipmentsComponent);
      // The first shipment has a destination address so will not use the duty location postal code
      const destination = shipments[0].destinationAddress;
      expect(wrapper.find('[data-testid="destinationAddress"]').at(0).text()).toEqual(
        `${destination.streetAddress1},\xa0${destination.streetAddress2},\xa0${destination.city}, ${destination.state} ${destination.postalCode}`,
      );
      expect(wrapper.find('[data-testid="destinationAddress"]').at(1).text()).toEqual(
        ordersInfo.newDutyLocation.address.postalCode,
      );
    });

    it('enables the Approve selected button when a shipment and service item are checked', async () => {
      const wrapper = mount(requestedShipmentsComponentWithPermission);

      await act(async () => {
        wrapper
          .find('input[name="shipments"]')
          .at(0)
          .simulate('change', {
            target: {
              name: 'shipments',
              value: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
            },
          });
      });
      wrapper.update();

      expect(wrapper.find('form button[type="button"]').prop('disabled')).toEqual(true);
      expect(wrapper.find('#approvalConfirmationModal').prop('style')).toHaveProperty('display', 'none');

      await act(async () => {
        wrapper
          .find('input[name="shipmentManagementFee"]')
          .simulate('change', { target: { name: 'shipmentManagementFee', value: true } });
      });
      wrapper.update();

      expect(wrapper.find('form button[type="button"]').prop('disabled')).toBe(false);

      await act(async () => {
        wrapper.find('form button[type="button"]').simulate('click');
      });
      wrapper.update();

      expect(wrapper.find('#approvalConfirmationModal').prop('style')).toHaveProperty('display', 'block');
    });

    it('disables the Approve selected button when there is missing required information', async () => {
      const wrapper = mount(requestedShipmentsComponentMissingRequiredInfo);

      await act(async () => {
        wrapper
          .find('input[name="shipments"]')
          .at(0)
          .simulate('change', {
            target: {
              name: 'shipments',
              value: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
            },
          });
      });
      wrapper.update();

      expect(wrapper.find('form button[type="button"]').prop('disabled')).toEqual(true);

      await act(async () => {
        wrapper
          .find('input[name="shipmentManagementFee"]')
          .simulate('change', { target: { name: 'shipmentManagementFee', value: true } });
      });
      wrapper.update();

      expect(wrapper.find('form button[type="button"]').prop('disabled')).toBe(true);
    });

    it('calls approveMTO onSubmit', async () => {
      const mockOnSubmit = jest.fn((id, eTag) => {
        return new Promise((resolve) => {
          resolve({ response: { status: 200, body: { id, eTag } } });
        });
      });

      const wrapper = mount(
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
      //  wrapper.find('form').simulate('submit');

      // When simulating change events you must pass the target with the id and
      // name for formik to know which value to update
      await act(async () => {
        wrapper
          .find('input[name="shipments"]')
          .at(0)
          .simulate('change', {
            target: {
              name: 'shipments',
              value: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
            },
          });

        wrapper
          .find('input[name="shipmentManagementFee"]')
          .simulate('change', { target: { name: 'shipmentManagementFee', value: true } });

        wrapper
          .find('input[name="counselingFee"]')
          .simulate('change', { target: { name: 'counselingFee', value: true } });

        wrapper.find('form button[type="button"]').simulate('click');

        wrapper.find('button[type="submit"]').simulate('click');
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
      const wrapper = mount(
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
      const approvedServiceItemNames = wrapper.find('[data-testid="basicServiceItemName"]');
      const approvedServiceItemDates = wrapper.find('[data-testid="basicServiceItemDate"]');

      expect(approvedServiceItemNames.length).toBe(2);
      expect(approvedServiceItemDates.length).toBe(2);

      expect(approvedServiceItemNames.at(0).text()).toBe('Move management');
      expect(approvedServiceItemDates.at(0).find('FontAwesomeIcon').prop('icon')).toEqual('check');
      expect(approvedServiceItemDates.at(0).text()).toBe(' 01 Jan 2020');

      expect(approvedServiceItemNames.at(1).text()).toBe('Counseling fee');
      expect(approvedServiceItemDates.at(1).find('FontAwesomeIcon').prop('icon')).toEqual('check');
      expect(approvedServiceItemDates.at(1).text()).toBe(' 01 Jan 2020');
    });

    it.each([['APPROVED'], ['SUBMITTED']])(
      'displays the customer and counselor remarks for a(n) %s shipment',
      (status) => {
        const wrapper = mount(
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

        const customerRemarks = wrapper.find('[data-testid="customerRemarks"]');
        const counselorRemarks = wrapper.find('[data-testid="counselorRemarks"]');

        expect(customerRemarks.at(0).text()).toBe('please treat gently');
        expect(customerRemarks.at(1).text()).toBe('please treat gently');

        expect(counselorRemarks.at(0).text()).toBe('looks good');
        expect(counselorRemarks.at(1).text()).toBe('looks good');
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
