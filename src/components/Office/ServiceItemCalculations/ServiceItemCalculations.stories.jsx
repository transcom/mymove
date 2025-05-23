import React from 'react';

import ServiceItemCalculations from './ServiceItemCalculations';
import testParams from './serviceItemTestParams';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Office Components/ServiceItemCalculations',
  decorators: [
    (Story) => {
      return (
        <div style={{ padding: '20px' }}>
          <Story />
        </div>
      );
    },
  ],
  argTypes: {
    tableSize: {
      defaultValue: 'large',
      control: {
        type: 'select',
        options: ['small', 'large'],
      },
    },
  },
};

export const DDDSIT = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticDestinationSITDeliveryLonghaul}
    totalAmountRequested={642}
    itemCode="DDDSIT"
    tableSize={data.tableSize}
  />
);

export const DDDSITShorthaulDifferentZIP3 = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticDestinationSITDelivery}
    totalAmountRequested={642}
    itemCode="DDDSIT"
    tableSize={data.tableSize}
  />
);

export const DLH = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticLongHaul}
    totalAmountRequested={642}
    itemCode="DLH"
    tableSize={data.tableSize}
  />
);

export const DLHNTSrelease = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticLongHaul}
    totalAmountRequested={642}
    itemCode="DLH"
    tableSize={data.tableSize}
    shipmentType={SHIPMENT_OPTIONS.NTSR}
  />
);
DLHNTSrelease.storyName = 'DLH (NTS-release)';

export const DOASIT = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticOriginAdditionalSIT}
    totalAmountRequested={642}
    itemCode="DOASIT"
    tableSize={data.tableSize}
  />
);

export const DOFSIT = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticOrigin1stSIT}
    totalAmountRequested={642}
    itemCode="DOFSIT"
    tableSize={data.tableSize}
  />
);

export const DDFSIT = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticDestination1stSIT}
    totalAmountRequested={642}
    itemCode="DDFSIT"
    tableSize={data.tableSize}
  />
);

export const DOP = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticOriginPrice}
    totalAmountRequested={642}
    itemCode="DOP"
    tableSize={data.tableSize}
  />
);

export const DPK = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticPacking}
    totalAmountRequested={642}
    itemCode="DPK"
    tableSize={data.tableSize}
  />
);

export const DNPK = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticNTSPacking}
    totalAmountRequested={642}
    itemCode="DNPK"
    tableSize={data.tableSize}
  />
);

export const DSH = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticShortHaul}
    totalAmountRequested={642}
    itemCode="DSH"
    tableSize={data.tableSize}
  />
);

export const DUPK = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticUnpacking}
    totalAmountRequested={642}
    itemCode="DUPK"
    tableSize={data.tableSize}
  />
);

export const DOPSIT = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticOriginSITPickup}
    totalAmountRequested={642}
    itemCode="DOPSIT"
    tableSize={data.tableSize}
  />
);

export const DOSHUT = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticOriginShuttleService}
    totalAmountRequested={642}
    itemCode="DOSHUT"
    tableSize={data.tableSize}
  />
);

export const DDP = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticDestinationPrice}
    totalAmountRequested={642}
    itemCode="DDP"
    tableSize={data.tableSize}
  />
);

export const DDASIT = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticDestinationAdditionalSIT}
    totalAmountRequested={642}
    itemCode="DDASIT"
    tableSize={data.tableSize}
  />
);

export const DDSHUT = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticDestinationShuttleService}
    totalAmountRequested={642}
    itemCode="DDSHUT"
    tableSize={data.tableSize}
  />
);

export const DCRT = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticCrating}
    additionalServiceItemData={testParams.additionalCratingDataDCRT}
    totalAmountRequested={642}
    itemCode="DCRT"
    tableSize={data.tableSize}
  />
);

export const DUCRT = (data) => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticUncrating}
    additionalServiceItemData={testParams.additionalCratingDataDCRT}
    totalAmountRequested={642}
    itemCode="DUCRT"
    tableSize={data.tableSize}
  />
);
