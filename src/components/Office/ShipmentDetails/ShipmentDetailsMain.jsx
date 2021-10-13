import React, { useMemo, useState } from 'react';
import * as PropTypes from 'prop-types';

import { SIT_EXTENSION_STATUS } from '../../../constants/sitExtensions';

import styles from './ShipmentDetails.module.scss';

import { formatDate } from 'shared/dates';
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
  dutyStationAddresses,
  handleDivertShipment,
  handleRequestReweighModal,
  handleReviewSITExtension,
  handleSubmitSITExtension,
}) => {
  const {
    requestedPickupDate,
    scheduledPickupDate,
    pickupAddress,
    destinationAddress,
    primeEstimatedWeight,
    primeActualWeight,
    counselorRemarks,
    customerRemarks,
    sitExtensions,
    sitStatus,
    storageInTransit,
  } = shipment;
  const { originDutyStationAddress, destinationDutyStationAddress } = dutyStationAddresses;

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
      />
      <ShipmentAddresses
        pickupAddress={pickupAddress}
        destinationAddress={destinationAddress || destinationDutyStationAddress?.postal_code}
        originDutyStation={originDutyStationAddress}
        destinationDutyStation={destinationDutyStationAddress}
        shipmentInfo={{ shipmentID: shipment.id, ifMatchEtag: shipment.eTag, shipmentStatus: shipment.status }}
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
  dutyStationAddresses: PropTypes.shape({
    originDutyStationAddress: AddressShape,
    destinationDutyStationAddress: AddressShape,
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
