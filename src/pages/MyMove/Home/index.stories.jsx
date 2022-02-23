/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import { Home } from './index';

import { MockProviders } from 'testUtils';

export default {
  title: 'Customer Components / Pages / Home',
};

const uploadOrdersProps = {
  serviceMember: {
    id: 'testServiceMemberId',
    first_name: 'John',
    last_name: 'Lee',
    current_station: {
      name: 'Fort Knox',
      transportation_office: {
        name: 'Test Transportation Office Name',
        phone_lines: ['555-555-5555'],
      },
    },
    weight_allotment: {},
  },
  showLoggedInUser() {},
  loadMTOShipments() {},
  history: { push: () => {}, goBack: () => {} },
  getSignedCertification() {},
  mtoShipments: [],
  mtoShipment: {},
  isLoggedIn: true,
  loggedInUserIsLoading: false,
  loggedInUserSuccess: true,
  isProfileComplete: true,
  currentPpm: {},
  orders: {},
  location: {},
  move: {
    locator: 'XYZ890',
    status: 'DRAFT',
  },
  uploadedOrderDocuments: [],
  uploadedAmendedOrderDocuments: [],
};

const shipmentSelectionProps = {
  ...uploadOrdersProps,
  serviceMember: {
    ...uploadOrdersProps.serviceMember,
    weight_allotment: {
      total_weight_self: 10000,
    },
  },
  orders: {
    ...uploadOrdersProps.orders,
    new_duty_location: {
      name: 'NAS Jacksonville',
    },
    origin_duty_location: {
      name: 'NAS Norfolk',
    },
    report_by_date: '25 December 2020',
  },
  uploadedOrderDocuments: [
    {
      id: 'file1',
      filename: 'Uploaded_Orders.pdf',
    },
    {
      id: 'file2',
      filename: 'Supporting_Documentation_Screenshot.png',
    },
  ],
};

const withShipmentProps = {
  ...shipmentSelectionProps,
  mtoShipments: [
    {
      id: 'testShipment1',
      shipmentType: 'HHG',
      createdAt: '24 December 2020',
    },
  ],
  currentPpm: {
    id: 'testMove',
  },
};

const submittedProps = {
  ...withShipmentProps,
  move: {
    ...withShipmentProps.move,
    status: 'SUBMITTED',
    submitted_at: '24 December 2020',
  },
};

const amendedOrderProps = {
  ...submittedProps,
  move: {
    ...submittedProps.move,
    status: 'APPROVALS REQUESTED',
  },
  uploadedAmendedOrderDocuments: [
    {
      id: 'file3',
      filename: 'Amended_Orders.pdf',
    },
  ],
};

export const Step2 = () => {
  return (
    <MockProviders>
      <div className="grid-container usa-prose">
        <Home {...uploadOrdersProps} />
      </div>
    </MockProviders>
  );
};

export const Step3 = () => {
  return (
    <MockProviders>
      <div className="grid-container usa-prose">
        <Home {...shipmentSelectionProps} />
      </div>
    </MockProviders>
  );
};

export const Step4 = () => {
  return (
    <MockProviders>
      <div className="grid-container usa-prose">
        <Home {...withShipmentProps} />
      </div>
    </MockProviders>
  );
};

export const SubmittedMove = () => {
  return (
    <MockProviders>
      <div className="grid-container usa-prose">
        <Home {...submittedProps} />
      </div>
    </MockProviders>
  );
};

export const AmendedOrders = () => {
  return (
    <MockProviders>
      <div className="grid-container usa-prose">
        <Home {...amendedOrderProps} />
      </div>
    </MockProviders>
  );
};
