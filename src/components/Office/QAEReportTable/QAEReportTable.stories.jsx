import React from 'react';

import QAEReportTable from './QAEReportTable';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/QAEReportTable',
  component: QAEReportTable,
  decorators: [(Story) => <MockProviders>{Story()}</MockProviders>],
};

const reports = [
  {
    createdAt: '2022-07-14T19:21:27.573Z',
    evalEnd: '09:00',
    evalStart: '10:00',
    id: 'a7fdb0b3-f876-4686-b94f-ad20a2c9a63d',
    inspectionDate: '2022-07-14T00:00:00.000Z',
    inspectionType: 'VIRTUAL',
    location: 'DESTINATION',
    moveID: 'bd1bbbdc-1710-4831-aa41-e35ebedff5cd',
    remarks: 'this is a submitted shipment report',
    shipmentID: '38e87840-d385-413e-9746-b2da2c8245bb',
    submittedAt: '2022-07-14T19:21:27.565Z',
    type: 'SHIPMENT',
    violationsObserved: true,
  },
  {
    createdAt: '2022-07-14T19:21:27.579Z',
    id: '1f9d18a8-7688-428d-be8e-3f3c59a0ae5e',
    inspectionDate: '2022-07-14T00:00:00.000Z',
    inspectionType: 'PHYSICAL',
    timeDepart: '08:00',
    evalStart: '09:00',
    evalEnd: '09:45',
    location: null,
    locationDescription: 'Route 66 at crash inspection site 3',
    moveID: 'bd1bbbdc-1710-4831-aa41-e35ebedff5cd',
    remarks: 'this is a draft NTS shipment report',
    shipmentID: 'd46825dd-cf90-442b-96de-c5075ea2f1bf',
    submittedAt: null,
    type: 'SHIPMENT',
    violationsObserved: true,
  },
];

export const empty = () => (
  <div className="officeApp">
    <QAEReportTable reports={[]} emptyText="No QAE reports have been submitted for this shipment." />
  </div>
);

export const single = () => (
  <div className="officeApp">
    <QAEReportTable reports={reports} emptyText="No QAE reports have been submitted for this shipment." />
  </div>
);
