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
  ppmInfo,
  ppmInfoApprovedOrExcluded,
  ppmInfoRejected,
  ppmInfoMultiple,
  ppmInfoMultiple2,
} from './ShipmentDisplayTestData';

import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders, MockRouterProvider } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

export default {
  title: 'Office Components/Shipment Display',
  component: ShipmentDisplay,
  decorators: [
    (Story, context) => {
      // Dont wrap with permissions for the read only tests
      if (context.name.includes('Read Only')) {
        return (
          <MockRouterProvider>
            <Story />
          </MockRouterProvider>
        );
      }

      // By default, show component with permissions
      return (
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <Story />
        </MockProviders>
      );
    },
  ],
};

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

export const NTSShipmentMissingInfoAsTOO = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTS}
      shipmentId={ntsMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={[]}
      errorIfMissing={[{ fieldName: 'storageFacility' }, { fieldName: 'serviceOrderNumber' }, { fieldName: 'tacType' }]}
      showWhenCollapsed={['tacType']}
      editURL="/"
    />
  </div>
);

export const NTSShipmentMissingInfoAsSC = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTS}
      shipmentId={ntsMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={[{ fieldName: 'counselorRemarks' }, { fieldName: 'tacType' }, { fieldName: 'sacType' }]}
      errorIfMissing={[]}
      showWhenCollapsed={['counselorRemarks']}
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
      showWhenCollapsed={['counselorRemarks']}
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
      showWhenCollapsed={['counselorRemarks']}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSReleaseShipmentMissingInfoAsTOO = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsReleaseMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={[]}
      errorIfMissing={[
        { fieldName: 'storageFacility' },
        { fieldName: 'ntsRecordedWeight' },
        { fieldName: 'serviceOrderNumber' },
        { fieldName: 'tacType' },
      ]}
      showWhenCollapsed={['tacType', 'ntsRecordedWeight', 'serviceOrderNumber']}
      editURL="/"
    />
  </div>
);

export const NTSReleaseShipmentMissingInfoAsSC = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsReleaseMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={[
        { fieldName: 'ntsRecordedWeight' },
        { fieldName: 'serviceOrderNumber' },
        { fieldName: 'counselorRemarks' },
        { fieldName: 'tacType' },
        { fieldName: 'sacType' },
      ]}
      errorIfMissing={[{ fieldName: 'storageFacility' }]}
      showWhenCollapsed={['counselorRemarks']}
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

export const PPMShipmentServiceCounselor = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ppmInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.PPM}
      isSubmitted
      allowApproval={false}
      warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
    />
  </div>
);

export const PPMShipmentWithCounselorRemarks = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...ppmInfo, counselorRemarks: 'counselor approved' }}
      shipmentType={SHIPMENT_OPTIONS.PPM}
      ordersLOA={ordersLOA}
      isSubmitted
      allowApproval={false}
      warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
    />
  </div>
);

export const PPMShipmentServiceCounselorWithReviewButton = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ppmInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.PPM}
      isSubmitted
      allowApproval={false}
      warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
      reviewURL="/"
    />
  </div>
);

export const MultiplePPMShipmentsServiceCounselorWithReviewButton = () => (
  <div className="shipmentCards_shipmentCards__ok-yC">
    <div style={{ padding: '20px' }}>
      <ShipmentDisplay
        displayInfo={ppmInfoMultiple}
        ordersLOA={ordersLOA}
        shipmentType={SHIPMENT_OPTIONS.PPM}
        isSubmitted
        allowApproval={false}
        warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
        reviewURL="/"
      />
    </div>
    <div style={{ padding: '20px' }}>
      <ShipmentDisplay
        displayInfo={ppmInfoMultiple2}
        ordersLOA={ordersLOA}
        shipmentType={SHIPMENT_OPTIONS.PPM}
        isSubmitted
        allowApproval={false}
        warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
        reviewURL="/"
      />
    </div>
  </div>
);

export const HHGShipmentReadOnly = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay displayInfo={hhgInfo} ordersLOA={ordersLOA} shipmentType={SHIPMENT_OPTIONS.HHG} isSubmitted />
  </div>
);

export const HHGShipmentServiceCounselorReadOnly = () => (
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

export const HHGShipmentWithCounselorRemarksReadOnly = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...hhgInfo, counselorRemarks: 'counselor approved' }}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      ordersLOA={ordersLOA}
      isSubmitted
    />
  </div>
);

export const HHGShipmentEditableReadOnly = () => (
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

export const NTSShipmentReadOnly = () => (
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

export const NTSShipmentMissingInfoAsTOOReadOnly = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTS}
      shipmentId={ntsMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={[]}
      errorIfMissing={[{ fieldName: 'storageFacility' }, { fieldName: 'serviceOrderNumber' }, { fieldName: 'tacType' }]}
      showWhenCollapsed={['tacType']}
      editURL="/"
    />
  </div>
);

export const NTSShipmentMissingInfoAsSCReadOnly = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTS}
      shipmentId={ntsMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={[{ fieldName: 'counselorRemarks' }, { fieldName: 'tacType' }, { fieldName: 'sacType' }]}
      errorIfMissing={[]}
      showWhenCollapsed={['counselorRemarks']}
      editURL="/"
    />
  </div>
);

export const NTSShipmentExternalVendorReadOnly = () => (
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

export const NTSReleaseShipmentReadOnly = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsReleaseInfo}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseInfo.shipmentId}
      ordersLOA={ordersLOA}
      showWhenCollapsed={['counselorRemarks']}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSReleaseShipmentExternalVendorReadOnly = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...ntsReleaseInfo, usesExternalVendor: true }}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseInfo.shipmentId}
      ordersLOA={ordersLOA}
      showWhenCollapsed={['counselorRemarks']}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSReleaseShipmentMissingInfoAsTOOReadOnly = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsReleaseMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={[]}
      errorIfMissing={[
        { fieldName: 'storageFacility' },
        { fieldName: 'ntsRecordedWeight' },
        { fieldName: 'serviceOrderNumber' },
        { fieldName: 'tacType' },
      ]}
      showWhenCollapsed={['tacType', 'ntsRecordedWeight', 'serviceOrderNumber']}
      editURL="/"
    />
  </div>
);

export const NTSReleaseShipmentMissingInfoAsSCReadOnly = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsReleaseMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={[
        { fieldName: 'ntsRecordedWeight' },
        { fieldName: 'serviceOrderNumber' },
        { fieldName: 'counselorRemarks' },
        { fieldName: 'tacType' },
        { fieldName: 'sacType' },
      ]}
      errorIfMissing={[{ fieldName: 'storageFacility' }]}
      showWhenCollapsed={['counselorRemarks']}
      editURL="/"
    />
  </div>
);

export const ApprovedShipmentReadOnly = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={hhgInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted={false}
    />
  </div>
);

export const PostalOnlyDestinationReadOnly = () => (
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

export const DivertedShipmentReadOnly = () => (
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

export const CancelledShipmentReadOnly = () => (
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

export const PPMShipmentServiceCounselorReadOnly = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ppmInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.PPM}
      isSubmitted
      allowApproval={false}
      warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
    />
  </div>
);

export const PPMShipmentWithCounselorRemarksReadOnly = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...ppmInfo, counselorRemarks: 'counselor approved' }}
      shipmentType={SHIPMENT_OPTIONS.PPM}
      ordersLOA={ordersLOA}
      isSubmitted
      allowApproval={false}
      warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
    />
  </div>
);

export const PPMShipmentServiceCounselorApproved = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ppmInfoApprovedOrExcluded}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.PPM}
      isSubmitted
      allowApproval={false}
      warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
      reviewURL="/"
    />
  </div>
);

export const PPMShipmentServiceCounselorRejected = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ppmInfoRejected}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.PPM}
      isSubmitted
      allowApproval={false}
      warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
      reviewURL="/"
    />
  </div>
);

export const PPMShipmentServiceCounselorExcluded = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ppmInfoApprovedOrExcluded}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.PPM}
      isSubmitted
      allowApproval={false}
      warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
      reviewURL="/"
    />
  </div>
);
