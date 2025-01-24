import React from 'react';
import { GridContainer } from '@trussworks/react-uswds';

import ExpandableServiceItemRow from './ExpandableServiceItemRow';
import '../ServiceItemCalculations/ServiceItemCalculations.module.scss';

export default {
  title: 'Office Components/ExpandableServiceItemRow',
  decorators: [
    (Story) => {
      return (
        <div style={{ maxWidth: '68.2666666667rem', marginRight: '0' }}>
          <GridContainer
            style={{ maxWidth: '63.8rem' }}
            className="expandableServiceItemRow"
            data-testid="tio-payment-request-details"
          >
            <table className="table--stacked">
              <colgroup>
                <col style={{ width: '40%' }} />
                <col style={{ width: '20%' }} />
                <col style={{ width: '20%' }} />
              </colgroup>
              <thead>
                <tr>
                  <th>Service item</th>
                  <th className="align-right">Amount</th>
                  <th className="align-right">Status</th>
                </tr>
              </thead>
              <tbody style={{ padding: '20px' }}>
                <Story />
              </tbody>
            </table>
          </GridContainer>
        </div>
      );
    },
  ],
};

const serviceItemRejected = {
  createdAt: '2025-01-09T22:08:38.788Z',
  eTag: 'MjAyNS0wMS0xN1QxNTowODo0Mi44MDI4MDZa',
  id: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
  mtoServiceItemCode: 'DLH',
  mtoServiceItemID: '526f705d-dba1-4bae-bf9a-e97cd1931bd4',
  mtoServiceItemName: 'Domestic linehaul',
  mtoShipmentID: 'ad5c56af-9e32-41bf-8283-a6a52938cc6a',
  mtoShipmentType: 'HHG',
  paymentServiceItemParams: [
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC44NDQ5Nlo=',
      id: 'd3fba800-cc16-45e3-975d-3236884fbf8a',
      key: 'WeightOriginal',
      origin: 'PRIME',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'INTEGER',
      value: '2000',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC44NDAyNjFa',
      id: '19192fe0-3e0b-4d5d-98dd-ea834fe9062f',
      key: 'ActualPickupDate',
      origin: 'PRIME',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'DATE',
      value: '2025-01-09',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC44MjcwNTZa',
      id: 'e249d609-96fa-4533-90dc-12d3164aed41',
      key: 'RequestedPickupDate',
      origin: 'PRIME',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'DATE',
      value: '2025-01-02',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC45NDc5NjVa',
      id: '75fc3b3b-517d-4383-9d3e-1493bcd564d9',
      key: 'DistanceZip',
      origin: 'SYSTEM',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'INTEGER',
      value: '1540',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC45NTk0MTZa',
      id: '15edbab6-ba00-47e2-bc92-49c7b16f57e1',
      key: 'ContractYearName',
      origin: 'PRICER',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'STRING',
      value: 'Award Term 1',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC44NDIwMVo=',
      id: '984fb8ea-da8a-4b3a-8560-85a9da1589ab',
      key: 'ZipDestAddress',
      origin: 'PRIME',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'STRING',
      value: '85309',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC44NDk1MzNa',
      id: '6d8b063b-8af4-4d6a-9ad4-5880a2d5fbea',
      key: 'ReferenceDate',
      origin: 'SYSTEM',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'DATE',
      value: '2025-01-02',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC44MzY0MjJa',
      id: '7941f3cf-7bcd-495e-9291-f2965291676c',
      key: 'ServiceAreaOrigin',
      origin: 'SYSTEM',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'STRING',
      value: '456',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC45NjExODda',
      id: '9a69db91-2d5d-4446-b1f6-0c15a140cf7a',
      key: 'EscalationCompounded',
      origin: 'PRICER',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'DECIMAL',
      value: '1.10701',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC45NjI5ODla',
      id: '35ddc999-a546-4076-9a2e-5295e4ae4279',
      key: 'IsPeak',
      origin: 'PRICER',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'BOOLEAN',
      value: 'false',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC45NjQ3ODNa',
      id: 'd35e3639-6a43-4f6a-8bf0-5a82685b80c4',
      key: 'PriceRateOrFactor',
      origin: 'PRICER',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'DECIMAL',
      value: '3.148',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC44Mzg0MTRa',
      id: '163babcf-7874-474e-a36d-f602fdca5c88',
      key: 'ContractCode',
      origin: 'SYSTEM',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'STRING',
      value: 'TRUSS_TEST',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC44Mjk2OTFa',
      id: '3f677c1e-1302-4cc1-9333-6fa34f5bdab5',
      key: 'WeightEstimated',
      origin: 'PRIME',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'INTEGER',
      value: '1500',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC44NDc2NDNa',
      id: '1009cca9-cdc1-4814-97e9-e0ef82dce965',
      key: 'WeightBilled',
      origin: 'SYSTEM',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'INTEGER',
      value: '1650',
    },
    {
      eTag: 'MjAyNS0wMS0wOVQyMjowODozOC44MzE2Nzha',
      id: '66b571ee-4142-4576-b563-cc3c8ea04bfe',
      key: 'ZipPickupAddress',
      origin: 'PRIME',
      paymentServiceItemID: '46e2df6f-4fe9-47ee-9baa-b9de28251da8',
      type: 'STRING',
      value: '62225',
    },
  ],
  priceCents: 8855385,
  referenceID: '4131-9325-46e2df6f',
  rejectionReason:
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
  status: 'DENIED',
};

const additionalServiceItemData = {
  approvedAt: '2025-01-09T20:24:58.522Z',
  convertToCustomerExpense: false,
  createdAt: '2025-01-09T20:24:58.621Z',
  deletedAt: '0001-01-01',
  eTag: 'MjAyNS0wMS0wOVQyMDoyNDo1OC42MjE5NzRa',
  id: '526f705d-dba1-4bae-bf9a-e97cd1931bd4',
  moveTaskOrderID: 'b02c42d7-bd4f-48ff-a5f8-6e7332fa5d03',
  mtoShipmentID: 'ad5c56af-9e32-41bf-8283-a6a52938cc6a',
  reServiceCode: 'DLH',
  reServiceID: '8d600f25-1def-422d-b159-617c7d59156e',
  reServiceName: 'Domestic linehaul',
  status: 'APPROVED',
  submittedAt: '0001-01-01',
  updatedAt: '0001-01-01T00:00:00.000Z',
};

export const rejectedServiceItem = () => (
  <ExpandableServiceItemRow
    serviceItem={serviceItemRejected}
    additionalServiceItemData={additionalServiceItemData}
    index={0}
    disableExpansion={false}
    paymentIsDeprecated={false}
    tppsDataExists={false}
  />
);

const serviceItemAccepted = { ...serviceItemRejected };
serviceItemAccepted.status = 'APPROVED';
serviceItemAccepted.rejectionReason = null;
export const acceptedServiceItem = () => (
  <ExpandableServiceItemRow
    serviceItem={serviceItemAccepted}
    additionalServiceItemData={additionalServiceItemData}
    index={0}
    disableExpansion={false}
    paymentIsDeprecated={false}
    tppsDataExists={false}
  />
);
