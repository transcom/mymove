/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow, mount } from 'enzyme';

import { SERVICE_ITEM_STATUS } from '../../../shared/constants';

import RequestedServiceItemsTable from './RequestedServiceItemsTable';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

const defaultProps = {
  handleShowRejectionDialog: jest.fn(),
  handleUpdateMTOServiceItemStatus: jest.fn(),
  serviceItemAddressUpdateAlert: {
    makeVisible: false,
    alertMessage: '',
    alertType: '',
  },
};

const serviceItemWithCrating = {
  id: 'abc123',
  createdAt: '2020-11-20',
  serviceItem: 'Domestic crating',
  code: 'DCRT',
  details: {
    description: 'grandfather clock',
    itemDimensions: { length: 7000, width: 2000, height: 3500 },
  },
};

const serviceItemWithContact = {
  id: 'abc1234',
  createdAt: '2020-09-01',
  serviceItem: 'Domestic destination 1st day SIT',
  code: 'DDFSIT',
  details: {
    sitEntryDate: '',
    customerContacts: [
      {
        timeMilitary: '1200Z',
        firstAvailableDeliveryDate: '2020-09-15',
        dateOfContact: '2020-09-15',
      },
      { timeMilitary: '2300Z', firstAvailableDeliveryDate: '2020-09-21', dateOfContact: '2020-09-21' },
    ],
    reason: 'Took a detour',
  },
};

const serviceItemWithDetails = {
  id: 'abc12345',
  createdAt: '2020-10-15',
  serviceItem: 'Domestic origin 1st day SIT',
  code: 'DOPSIT',
  details: {
    pickupPostalCode: '20050',
    SITPostalCode: '12345',
    reason: 'Took a detour',
    estimatedPrice: 243550,
    status: 'APPROVED',
  },
};

const serviceItemUBP = {
  id: 'ubp123',
  createdAt: '2025-01-15',
  serviceItem: 'International UB price',
  code: 'UBP',
  sort: '1',
};

const serviceItemIUBPK = {
  id: 'iubpk123',
  createdAt: '2025-01-15',
  serviceItem: 'International UB pack',
  code: 'IUBPK',
  sort: '3',
};

const serviceItemIUBUPK = {
  id: 'iubupk123',
  createdAt: '2025-01-15',
  serviceItem: 'International UB unpack',
  code: 'IUBUPK',
  sort: '4',
};

const testDetails = (wrapper) => {
  const labelMap = wrapper
    .find('dt')
    .map((node, index) => [node.text().trim(), wrapper.find('dd').at(index).text().trim()]);

  const getValue = (label) => {
    const entry = labelMap.find(([key]) => key === label);
    if (!entry) throw new Error(`Label "${label}" not found`);
    return entry[1];
  };

  const getAllValues = (label) => labelMap.filter(([key]) => key === label).map(([, value]) => value);

  expect(getValue('Description:')).toBe('grandfather clock');
  expect(getValue('Item size:')).toBe('7"x2"x3.5"');

  expect(getValue('Original Pickup Address:')).toContain('-');
  expect(getValue('Actual Pickup Address:')).toContain('-');
  expect(getValue('Delivery miles into SIT:')).toContain('-');
  expect(getValue('Original Delivery Address:')).toContain('-');
  expect(getValue('SIT entry date:')).toContain('-');

  expect(getValue('First available delivery date 1:')).toContain('15 Sep 2020');
  expect(getValue('Customer contact attempt 1:')).toContain('15 Sep 2020, 1200Z');
  expect(getValue('First available delivery date 2:')).toContain('21 Sep 2020');
  expect(getValue('Customer contact attempt 2:')).toContain('21 Sep 2020, 2300Z');

  expect(getValue('Estimated Price:')).toContain('$2,435.50');

  expect(getAllValues('Reason:')).toContain('Took a detour');
};

describe('RequestedServiceItemsTable', () => {
  it('shows the correct number of service items in the table', () => {
    const serviceItems = [serviceItemWithCrating];

    let wrapper = shallow(
      <RequestedServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    expect(wrapper.text().includes('1 item')).toBe(true);

    serviceItems.push(serviceItemWithContact);

    wrapper = shallow(
      <RequestedServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );
    expect(wrapper.text().includes('2 items')).toBe(true);
  });

  it('displays the service item name and submitted date', () => {
    const serviceItems = [serviceItemWithCrating, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <MockProviders>
        <RequestedServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('.codeName').at(0).text()).toBe('Domestic crating');
    expect(wrapper.find('.nameAndDate').at(0).text().includes('20 Nov 2020')).toBe(true);

    expect(wrapper.find('.codeName').at(1).text()).toBe('Domestic origin 1st day SIT');
    expect(wrapper.find('.nameAndDate').at(1).text().includes('15 Oct 2020')).toBe(true);

    expect(wrapper.find('.codeName').at(2).text()).toBe('Domestic destination 1st day SIT');
    expect(wrapper.find('.nameAndDate').at(2).text().includes('1 Sep 2020')).toBe(true);
  });

  it('shows the service item detail text', () => {
    const serviceItems = [serviceItemWithCrating, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <MockProviders>
        <RequestedServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );
    testDetails(wrapper);
  });

  it('displays the approve and reject status buttons', () => {
    serviceItemWithContact.status = 'SUBMITTED';
    serviceItemWithCrating.status = 'SUBMITTED';
    serviceItemWithDetails.status = 'SUBMITTED';

    const serviceItems = [serviceItemWithCrating, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <RequestedServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    const acceptButtons = wrapper.find({ 'data-testid': 'acceptButton' });
    expect(acceptButtons.at(0).text().includes('Accept')).toBe(true);
    expect(acceptButtons.at(1).text().includes('Accept')).toBe(true);
    expect(acceptButtons.at(2).text().includes('Accept')).toBe(true);

    const rejectButtons = wrapper.find({ 'data-testid': 'rejectButton' });
    expect(rejectButtons.at(0).text().includes('Reject')).toBe(true);
    expect(rejectButtons.at(1).text().includes('Reject')).toBe(true);
    expect(rejectButtons.at(2).text().includes('Reject')).toBe(true);
  });

  it('shows the service item detail text when approved and shows the reject button', () => {
    serviceItemWithDetails.status = 'APPROVED';
    serviceItemWithCrating.status = 'APPROVED';
    serviceItemWithContact.status = 'APPROVED';
    const serviceItems = [serviceItemWithCrating, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <RequestedServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
        />
      </MockProviders>,
    );

    testDetails(wrapper);
    const rejectTextButton = wrapper.find({ 'data-testid': 'rejectTextButton' });
    expect(rejectTextButton.at(0).text().includes('Reject')).toBe(true);
    expect(rejectTextButton.at(1).text().includes('Reject')).toBe(true);
    expect(rejectTextButton.at(2).text().includes('Reject')).toBe(true);
  });

  it('shows the service item detail text when rejected and shows the approve text button', () => {
    serviceItemWithDetails.status = 'REJECTED';
    serviceItemWithCrating.status = 'REJECTED';
    serviceItemWithContact.status = 'REJECTED';
    const serviceItems = [serviceItemWithCrating, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
        <RequestedServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
        />
      </MockProviders>,
    );

    testDetails(wrapper);
    const approveTextButton = wrapper.find({ 'data-testid': 'approveTextButton' });
    expect(approveTextButton.at(0).text().includes('Approve')).toBe(true);
    expect(approveTextButton.at(1).text().includes('Approve')).toBe(true);
    expect(approveTextButton.at(2).text().includes('Approve')).toBe(true);
  });

  it('displays sorted service items in order', () => {
    const serviceItems = [serviceItemIUBPK, serviceItemUBP, serviceItemIUBUPK];
    const wrapper = mount(
      <MockProviders>
        <RequestedServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('.codeName').at(0).text()).toBe('International UB price');
    expect(wrapper.find('.codeName').at(1).text()).toBe('International UB pack');
    expect(wrapper.find('.codeName').at(2).text()).toBe('International UB unpack');
  });

  it('displays sorted service items in order along with non-sorted service items', () => {
    const serviceItems = [serviceItemIUBPK, serviceItemUBP, serviceItemWithCrating, serviceItemIUBUPK];
    const wrapper = mount(
      <MockProviders>
        <RequestedServiceItemsTable
          {...defaultProps}
          serviceItems={serviceItems}
          statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
        />
      </MockProviders>,
    );

    expect(wrapper.find('.codeName').at(0).text()).toBe('International UB price');
    expect(wrapper.find('.codeName').at(1).text()).toBe('International UB pack');
    expect(wrapper.find('.codeName').at(2).text()).toBe('International UB unpack');
    expect(wrapper.find('.codeName').at(3).text()).toBe('Domestic crating');
  });
});
