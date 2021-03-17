import React from 'react';
import { act } from 'react-dom/test-utils';
import { mount, shallow } from 'enzyme';

import { ORDERS_TYPE, ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from '../../../constants/orders';
import { DEPARTMENT_INDICATOR_OPTIONS } from '../../../constants/departmentIndicators';
import SERVICE_ITEM_STATUSES from '../../../constants/serviceItems';

import RequestedShipments from './RequestedShipments';

import { SHIPMENT_OPTIONS, MTOAgentType } from 'shared/constants';
import { serviceItemCodes } from 'content/serviceItems';

const shipments = [
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.404029Z',
    customerRemarks: 'please treat gently',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODk0MTJa',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postal_code: '94535',
      state: 'CA',
      street_address_1: '987 Any Avenue',
      street_address_2: 'P.O. Box 9876',
      street_address_3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTkzMlo=',
      id: '15e8f6cc-e1d7-44b2-b1e0-fcb3d6442831',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTM4OTZa',
      id: '9b79e0c3-8ed5-4fb8-aa36-95845707d8ee',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.404031Z',
  },
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.431993Z',
    customerRemarks: 'please treat gently',
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MzE5OTVa',
    id: 'c2f68d97-b960-4c86-a418-c70a0aeba04e',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MTMyNDha',
      id: '14b1d10d-b34b-4ec5-80e6-69d885206a2a',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjYxODVa',
      id: '1a4f6fec-42b9-4dd2-b205-c6770ac7ea27',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjIwNzVa',
      id: 'e188f33f-f84d-4f86-954a-938b52e38741',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.431995Z',
  },
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.404029Z',
    customerRemarks: 'Please treat gently',
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aeee',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    primeActualWeight: 890,
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    shipmentType: SHIPMENT_OPTIONS.NTS,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.431995Z',
  },
];

const ordersInfo = {
  newDutyStation: {
    address: {
      city: 'Augusta',
      country: 'United States',
      eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
      id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
      postal_code: '30813',
      state: 'GA',
      street_address_1: 'Fort Gordon',
    },
    address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
    id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
    name: 'Fort Gordon',
  },
  currentDutyStation: {
    address: {
      city: 'Des Moines',
      country: 'US',
      eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42NjEwODFa',
      id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
      postal_code: '50309',
      state: 'IA',
      street_address_1: '987 Other Avenue',
      street_address_2: 'P.O. Box 1234',
      street_address_3: 'c/o Another Person',
    },
    address_id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
    eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42Njg5MDFa',
    id: '07282a8f-a496-4648-ae24-119775eef57d',
    name: 'vC6w22RPYC',
  },
  issuedDate: '2018-03-15',
  reportByDate: '2018-08-01',
  departmentIndicator: DEPARTMENT_INDICATOR_OPTIONS.COAST_GUARD,
  ordersNumber: 'ORDER3',
  ordersType: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  ordersTypeDetail: 'TBD',
  tacMDC: 'F381',
  sacSDN: '',
};

const customerInfo = {
  name: 'Smith, Kerry',
  dodId: '9999999999',
  phone: '+1 999-999-9999',
  email: 'ksmith@email.com',
  currentAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  backupContactName: 'Quinn Ocampo',
  backupContactPhone: '+1 999-999-9999',
  backupContactEmail: 'quinnocampo@myemail.com',
};

const agents = [
  {
    type: MTOAgentType.RELEASING_AGENT,
    name: 'Dorothy Lagomarsino',
    email: 'dorothyl@email.com',
    phone: '+1 999-999-9999',
    shipmentId: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
  },
  {
    type: MTOAgentType.RECEIVING_AGENT,
    name: 'Dorothy Lagomarsino',
    email: 'dorothyl@email.com',
    phone: '+1 999-999-9999',
    shipmentId: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
  },
];

const allowancesInfo = {
  branch: ORDERS_BRANCH_OPTIONS.NAVY,
  rank: ORDERS_RANK_OPTIONS.E_6,
  weightAllowance: 11000,
  authorizedWeight: 11000,
  progear: 2000,
  spouseProgear: 500,
  storageInTransit: 90,
  dependents: true,
};

const moveTaskOrder = {
  eTag: 'MjAyMC0wNi0yNlQyMDoyMjo0MS43Mjc4NTNa',
  id: '6e8c5ca4-774c-4170-934a-59d22259e480',
};

const serviceItems = [
  {
    approvedAt: '2020-10-02T19:20:08.481139Z',
    createdAt: '2020-10-01T19:20:08.481139Z',
    id: '12345',
    moveTaskOrderID: '6e8c5ca4-774c-4170-934a-59d22259e480',
    mtoShipmentID: null,
    reServiceCode: 'MS',
    reServiceID: '6789',
    reServiceName: serviceItemCodes.MS,
    status: SERVICE_ITEM_STATUSES.APPROVED,
  },
  {
    approvedAt: '2020-10-02T19:20:08.481139Z',
    createdAt: '2020-10-01T19:20:08.481139Z',
    id: '45678',
    moveTaskOrderID: '6e8c5ca4-774c-4170-934a-59d22259e480',
    mtoShipmentID: null,
    reServiceCode: 'CS',
    reServiceID: '6790',
    reServiceName: serviceItemCodes.CS,
    status: SERVICE_ITEM_STATUSES.APPROVED,
  },
  {
    approvedAt: '2020-10-02T19:20:08.481139Z',
    createdAt: '2020-10-01T19:20:08.481139Z',
    id: '9012',
    moveTaskOrderID: '6e8c5ca4-774c-4170-934a-59d22259e480',
    mtoShipmentID: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
    reServiceCode: 'DLH',
    reServiceID: '6791',
    reServiceRName: serviceItemCodes.DLH,
    status: SERVICE_ITEM_STATUSES.SUBMITTED,
  },
];

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
  />
);

const requestedShipmentsComponentMissingRequiredInfo = (
  <RequestedShipments
    ordersInfo={ordersInfo}
    allowancesInfo={allowancesInfo}
    mtoAgents={agents}
    customerInfo={customerInfo}
    mtoShipments={shipments}
    approveMTO={approveMTO}
    shipmentsStatus="SUBMITTED"
    missingRequiredOrdersInfo
  />
);

describe('RequestedShipments', () => {
  it('renders the container successfully', () => {
    const wrapper = shallow(requestedShipmentsComponent);
    expect(wrapper.find('div[data-testid="requested-shipments"]').exists()).toBe(true);
  });

  it('renders a shipment passed to it', () => {
    const wrapper = mount(requestedShipmentsComponent);
    expect(wrapper.find('div[data-testid="requested-shipments"]').text()).toContain('HHG');
    expect(wrapper.find('div[data-testid="requested-shipments"]').text()).toContain('NTS');
  });

  it('renders the button', () => {
    const wrapper = mount(requestedShipmentsComponent);
    const approveButton = wrapper.find('button[data-testid="shipmentApproveButton"]');
    expect(approveButton.exists()).toBe(true);
    expect(approveButton.text()).toContain('Approve selected shipments');
    expect(approveButton.html()).toContain('disabled=""');
  });

  it('renders the checkboxes', () => {
    const wrapper = mount(requestedShipmentsComponent);
    expect(wrapper.find('div[data-testid="checkbox"]').exists()).toBe(true);
    expect(wrapper.find('div[data-testid="checkbox"]').length).toEqual(5);
  });

  it('uses the duty station postal code if there is no destination address', () => {
    const wrapper = mount(requestedShipmentsComponent);
    // The first shipment has a destination address so will not use the duty station postal code
    const destination = shipments[0].destinationAddress;
    expect(wrapper.find('[data-testid="shipmentDestinationAddress"]').at(0).text()).toEqual(
      `${destination.street_address_1},\xa0${destination.city}, ${destination.state} ${destination.postal_code}`,
    );
    expect(wrapper.find('[data-testid="shipmentDestinationAddress"]').at(1).text()).toEqual(
      ordersInfo.newDutyStation.address.postal_code,
    );
  });

  it('enables the modal button when a shipment and service item are checked', async () => {
    const wrapper = mount(requestedShipmentsComponent);

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

  it('disables the modal button when there is missing required information', async () => {
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
      <RequestedShipments
        mtoShipments={shipments}
        mtoAgents={agents}
        ordersInfo={ordersInfo}
        allowancesInfo={allowancesInfo}
        customerInfo={customerInfo}
        moveTaskOrder={moveTaskOrder}
        approveMTO={mockOnSubmit}
        shipmentsStatus="SUBMITTED"
      />,
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
      />,
    );
    const approvedServiceItemNames = wrapper.find('[data-testid="basicServiceItemName"]');
    const approvedServiceItemDates = wrapper.find('[data-testid="basicServiceItemDate"]');

    expect(approvedServiceItemNames.length).toBe(2);
    expect(approvedServiceItemDates.length).toBe(2);

    expect(approvedServiceItemNames.at(0).text()).toBe('Move management');
    expect(approvedServiceItemDates.at(0).find('FontAwesomeIcon').prop('icon')).toEqual('check');
    expect(approvedServiceItemDates.at(0).text()).toBe(' 02 Oct 2020');

    expect(approvedServiceItemNames.at(1).text()).toBe('Counseling');
    expect(approvedServiceItemDates.at(1).find('FontAwesomeIcon').prop('icon')).toEqual('check');
    expect(approvedServiceItemDates.at(1).text()).toBe(' 02 Oct 2020');
  });
});
