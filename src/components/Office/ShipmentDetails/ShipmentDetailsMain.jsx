import React, { useMemo, useState } from 'react';
import * as PropTypes from 'prop-types';

import { SIT_EXTENSION_STATUS } from '../../../constants/sitExtensions';

import styles from './ShipmentDetails.module.scss';

import { formatDate } from 'shared/dates';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape } from 'types';
import { ShipmentShape } from 'types/shipment';
import SubmitSITExtensionModal from 'components/Office/SubmitSITExtensionModal/SubmitSITExtensionModal';
import ReviewSITExtensionsModal from 'components/Office/ReviewSITExtensionModal/ReviewSITExtensionModal';
import ShipmentSITDisplay from 'components/Office/ShipmentSITDisplay/ShipmentSITDisplay';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates/ImportantShipmentDates';
import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';
import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';
import ShipmentRemarks from 'components/Office/ShipmentRemarks/ShipmentRemarks';

const ShipmentDetailsMain = ({
  className,
  shipment,
  dutyLocationAddresses,
  handleDivertShipment,
  handleRequestReweighModal,
  handleReviewSITExtension,
  handleSubmitSITExtension,
}) => {
  const {
    requestedPickupDate,
    scheduledPickupDate,
    requiredDeliveryDate,
    pickupAddress,
    destinationAddress,
    primeEstimatedWeight,
    primeActualWeight,
    counselorRemarks,
    customerRemarks,
    sitExtensions,
    sitStatus,
    storageInTransit,
    shipmentType,
    storageFacility,
  } = shipment;
  const { originDutyLocationAddress, destinationDutyLocationAddress } = dutyLocationAddresses;

  const [isReviewSITExtensionModalVisible, setIsReviewSITExtensionModalVisible] = useState(false);
  const [isSubmitITExtensionModalVisible, setIsSubmitITExtensionModalVisible] = useState(false);

  const reviewSITExtension = (sitExtensionID, formValues) => {
    setIsReviewSITExtensionModalVisible(false);
    handleReviewSITExtension(sitExtensionID, formValues, shipment);
  };
  const submitSITExtension = (formValues) => {
    setIsSubmitITExtensionModalVisible(false);
    handleSubmitSITExtension(formValues, shipment);
  };

  const pendingSITExtension = sitExtensions?.find((se) => se.status === SIT_EXTENSION_STATUS.PENDING);

  const summarySITComponent = useMemo(
    () => (
      <ShipmentSITDisplay
        sitExtensions={sitExtensions}
        sitStatus={sitStatus}
        storageInTransit={storageInTransit}
        shipment={shipment}
        showReviewSITExtension={setIsReviewSITExtensionModalVisible}
        showSubmitSITExtension={setIsSubmitITExtensionModalVisible}
        hideSITExtensionAction
      />
    ),
    [
      sitExtensions,
      sitStatus,
      storageInTransit,
      shipment,
      setIsReviewSITExtensionModalVisible,
      setIsSubmitITExtensionModalVisible,
    ],
  );

  let displayedPickupAddress;
  let displayedDeliveryAddress;

  switch (shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
      displayedPickupAddress = pickupAddress;
      displayedDeliveryAddress = destinationAddress || destinationDutyLocationAddress?.postalCode;
      break;
    case SHIPMENT_OPTIONS.NTS:
      displayedPickupAddress = pickupAddress;
      displayedDeliveryAddress = storageFacility ? storageFacility.address : null;
      break;
    case SHIPMENT_OPTIONS.NTSR:
      displayedPickupAddress = storageFacility ? storageFacility.address : null;
      displayedDeliveryAddress = destinationAddress;
      break;
    default:
      displayedPickupAddress = pickupAddress;
      displayedDeliveryAddress = destinationAddress || destinationDutyLocationAddress?.postalCode;
  }

  return (
    <div className={className}>
      {isReviewSITExtensionModalVisible && (
        <ReviewSITExtensionsModal
          onClose={() => setIsReviewSITExtensionModalVisible(false)}
          onSubmit={reviewSITExtension}
          sitExtension={pendingSITExtension}
          summarySITComponent={summarySITComponent}
        />
      )}
      {isSubmitITExtensionModalVisible && (
        <SubmitSITExtensionModal
          onClose={() => setIsSubmitITExtensionModalVisible(false)}
          onSubmit={submitSITExtension}
          summarySITComponent={summarySITComponent}
        />
      )}
      {sitStatus && (
        <ShipmentSITDisplay
          sitExtensions={sitExtensions}
          sitStatus={sitStatus}
          storageInTransit={storageInTransit}
          shipment={shipment}
          showReviewSITExtension={setIsReviewSITExtensionModalVisible}
          showSubmitSITExtension={setIsSubmitITExtensionModalVisible}
          className={styles.shipmentSITSummary}
        />
      )}
      <ImportantShipmentDates
        requestedPickupDate={formatDate(requestedPickupDate)}
        scheduledPickupDate={scheduledPickupDate ? formatDate(scheduledPickupDate) : null}
        requiredDeliveryDate={requiredDeliveryDate ? formatDate(requiredDeliveryDate) : null}
      />
      <ShipmentAddresses
        pickupAddress={displayedPickupAddress}
        destinationAddress={displayedDeliveryAddress}
        originDutyLocation={originDutyLocationAddress}
        destinationDutyLocation={destinationDutyLocationAddress}
        shipmentInfo={{
          id: shipment.id,
          eTag: shipment.eTag,
          status: shipment.status,
          shipmentType: shipment.shipmentType,
        }}
        handleDivertShipment={handleDivertShipment}
      />
      <ShipmentWeightDetails
        estimatedWeight={primeEstimatedWeight}
        actualWeight={primeActualWeight}
        shipmentInfo={{
          shipmentID: shipment.id,
          ifMatchEtag: shipment.eTag,
          reweighID: shipment.reweigh?.id,
          reweighWeight: shipment.reweigh?.weight,
          shipmentType: shipment.shipmentType,
        }}
        handleRequestReweighModal={handleRequestReweighModal}
      />
      {counselorRemarks && <ShipmentRemarks title="Counselor remarks" remarks={counselorRemarks} />}
      {customerRemarks && <ShipmentRemarks title="Customer remarks" remarks={customerRemarks} />}
    </div>
  );
};

ShipmentDetailsMain.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  dutyLocationAddresses: PropTypes.shape({
    originDutyLocationAddress: AddressShape,
    destinationDutyLocationAddress: AddressShape,
  }).isRequired,
  handleDivertShipment: PropTypes.func.isRequired,
  handleRequestReweighModal: PropTypes.func.isRequired,
  handleReviewSITExtension: PropTypes.func.isRequired,
  handleSubmitSITExtension: PropTypes.func.isRequired,
};

ShipmentDetailsMain.defaultProps = {
  className: '',
};

export default ShipmentDetailsMain;
