import React from 'react';
import { act } from 'react-dom/test-utils';
import { mount } from 'enzyme';

import TableCSVExportButton from './TableCSVExportButton';

import { getPaymentRequestsQueue } from 'services/ghcApi';

const paymentRequestsResponse = {
  page: 1,
  perPage: 1,
  queuePaymentRequests: [
    {
      age: 8,
      customer: {
        agency: 'ARMY',
        cacValidated: false,
        dodID: '4800408743',
        eTag: 'MjAyNC0wNi0xMVQyMTowOToxMi4zMTE5OTda',
        email: 'leo_spaceman_sm@example.com',
        emailIsPreferred: true,
        first_name: 'Leo',
        id: 'd659cfab-e4e0-4785-b4de-a77ee5afedcf',
        last_name: 'Spacemen',
        phone: '212-123-4567',
        userID: '8c161e9e-da34-46fb-bd5c-ca72cc8ed692',
      },
      departmentIndicator: 'AIR_AND_SPACE_FORCE',
      id: '62114ae3-9ac1-4903-a10d-b764aca008eb',
      locator: 'PARAMS',
      moveID: '11721c3e-98de-47bf-829f-b99009caa1dd',
      orderType: 'PERMANENT_CHANGE_OF_STATION',
      originDutyLocation: {
        address: {
          city: 'Des Moines',
          country: 'US',
          county: 'POLK',
          eTag: 'MjAyNC0wNi0xMVQyMTowOToxMi4yOTgyNzNa',
          id: '848fa47e-54dc-4199-9ad4-f41910dad6c7',
          postalCode: '50309',
          state: 'IA',
          streetAddress1: '987 Other Avenue',
          streetAddress2: 'P.O. Box 1234',
          streetAddress3: 'c/o Another Person',
        },
        address_id: '848fa47e-54dc-4199-9ad4-f41910dad6c7',
        eTag: 'MjAyNC0wNi0xMVQyMTowOToxMi4zMDE4NFo=',
        id: '7b182c86-53aa-44fc-9997-511286e6c255',
        name: 'DXf3I11JSD',
      },
      originGBLOC: 'KKFA',
      status: 'Payment requested',
      submittedAt: '2024-06-11T21:09:14.249Z',
    },
  ],
};

const paymentRequestsNoResultsResponse = {
  page: 1,
  perPage: 10,
};

const paymentRequestColumns = [
  {
    Header: ' ',
    id: 'lock',
  },
  {
    Header: 'ID',
    accessor: 'id',
    id: 'id',
  },
  {
    Header: 'Customer name',
    id: 'lastName',
    isFilterable: true,
    exportValue: (row) => {
      return `${row.customer.last_name}, ${row.customer.first_name}`;
    },
  },
  {
    Header: 'DoD ID',
    accessor: 'customer.dodID',
    id: 'dodID',
    isFilterable: true,
    exportValue: (row) => {
      return row.customer.dodID;
    },
  },
  {
    Header: 'Status',
    id: 'status',
    disableSortBy: true,
    accessor: 'status',
  },
  {
    Header: 'Age',
    id: 'age',
    accessor: 'age',
  },
  {
    Header: 'Submitted',
    id: 'submittedAt',
    isFilterable: true,
    accessor: 'submittedAt',
  },
  {
    Header: 'Move Code',
    accessor: 'locator',
    id: 'locator',
    isFilterable: true,
  },
  {
    Header: 'Branch',
    id: 'branch',
    isFilterable: true,
    accessor: 'agency',
  },
  {
    Header: 'Origin GBLOC',
    accessor: 'originGBLOC',
    disableSortBy: true,
  },
  {
    Header: 'Origin Duty Location',
    accessor: 'originDutyLocation.name',
    id: 'originDutyLocation',
    isFilterable: true,
    exportValue: (row) => {
      return row.originDutyLocation?.name;
    },
  },
];

jest.mock('services/ghcApi', () => ({
  getPaymentRequestsQueue: jest.fn().mockImplementation(() => Promise.resolve(paymentRequestsResponse)),
  getPaymentRequestsNoResultsQueue: jest
    .fn()
    .mockImplementation(() => Promise.resolve(paymentRequestsNoResultsResponse)),
}));

describe('TableCSVExportButton', () => {
  const defaultProps = {
    tableColumns: paymentRequestColumns,
    queueFetcher: getPaymentRequestsQueue,
    queueFetcherKey: 'queuePaymentRequests',
    totalCount: 1,
  };

  it('renders without error', () => {
    const wrapper = mount(<TableCSVExportButton {...defaultProps} />);
    expect(wrapper.find({ 'data-test-id': 'csv-export-btn-hidden' }).at(0).hasClass('hidden')).toBe(true);
    expect(wrapper.find('span[data-test-id="csv-export-btn-text"]').text()).toBe('Export to CSV');
  });

  it('click calls fetcher', () => {
    act(() => {
      const wrapper = mount(<TableCSVExportButton {...defaultProps} />);
      const exportButton = wrapper.find('span[data-test-id="csv-export-btn-text"]');
      exportButton.simulate('click');
      wrapper.update();
    });

    expect(getPaymentRequestsQueue).toBeCalled();
  });

  const noResultsProps = {
    tableColumns: paymentRequestColumns,
    queueFetcher: () => Promise.resolve(paymentRequestsNoResultsResponse),
    queueFetcherKey: 'queuePaymentRequests',
    totalCount: 0,
  };

  it('is diabled when there is nothing to export', () => {
    act(() => {
      const wrapper = mount(<TableCSVExportButton {...noResultsProps} />);
      const exportButton = wrapper.find('span[data-test-id="csv-export-btn-text"]');
      exportButton.simulate('click');
      wrapper.update();
    });

    expect(getPaymentRequestsQueue).toBeCalled();
  });
});
