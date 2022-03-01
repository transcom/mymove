import React from 'react';

import {
  ordersLOA,
  hhgInfo,
  ntsInfo,
  ntsMissingInfo,
  ntsReleaseInfo,
  ntsReleaseMissingInfo,
  postalOnlyInfo,
  diversionInfo,
  cancelledInfo,
} from './ShipmentDisplayTestData';

import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/Shipment Display',
  component: ShipmentDisplay,
  decorators: [
    (Story) => (
      <MockProviders>
        <Story />
      </MockProviders>
    ),
  ],
};

const showWhenCollapsed = ['counselorRemarks'];

export const HHGShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay displayInfo={hhgInfo} ordersLOA={ordersLOA} shipmentType={SHIPMENT_OPTIONS.HHG} isSubmitted />
  </div>
);

export const HHGShipmentServiceCounselor = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={hhgInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
      allowApproval={false}
    />
  </div>
);

export const HHGShipmentWithCounselorRemarks = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...hhgInfo, counselorRemarks: 'counselor approved' }}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      ordersLOA={ordersLOA}
      isSubmitted
    />
  </div>
);

export const HHGShipmentEditable = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...hhgInfo, counselorRemarks: 'counselor approved' }}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      ordersLOA={ordersLOA}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.NTS}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSShipmentMissingInfo = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTS}
      shipmentId={ntsMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={['counselorRemarks']}
      errorIfMissing={['storageFacility', 'serviceOrderNumber', 'tacType']}
      showWhenCollapsed={showWhenCollapsed}
      editURL="/"
    />
  </div>
);

export const NTSShipmentExternalVendor = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...ntsInfo, usesExternalVendor: true }}
      shipmentType={SHIPMENT_OPTIONS.NTS}
      ordersLOA={ordersLOA}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSReleaseShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsReleaseInfo}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseInfo.shipmentId}
      ordersLOA={ordersLOA}
      showWhenCollapsed={showWhenCollapsed}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSReleaseShipmentExternalVendor = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...ntsReleaseInfo, usesExternalVendor: true }}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseInfo.shipmentId}
      ordersLOA={ordersLOA}
      showWhenCollapsed={showWhenCollapsed}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSReleaseShipmentMissingInfo = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsReleaseMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={['counselorRemarks']}
      errorIfMissing={['storageFacility', 'ntsRecordedWeight', 'serviceOrderNumber', 'tacType']}
      showWhenCollapsed={showWhenCollapsed}
      editURL="/"
    />
  </div>
);

export const ApprovedShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={hhgInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted={false}
    />
  </div>
);

export const PostalOnlyDestination = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={postalOnlyInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const DivertedShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      shipmentId="1"
      displayInfo={diversionInfo}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      ordersLOA={ordersLOA}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const CancelledShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      shipmentId="1"
      displayInfo={cancelledInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
    />
  </div>
);
