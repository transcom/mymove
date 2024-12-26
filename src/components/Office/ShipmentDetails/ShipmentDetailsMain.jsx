import React, { useState } from 'react';
import * as PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './ShipmentDetails.module.scss';

import { SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import { formatDateWithUTC } from 'shared/dates';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape } from 'types';
import { ShipmentShape } from 'types/shipment';
import SubmitSITExtensionModal from 'components/Office/SubmitSITExtensionModal/SubmitSITExtensionModal';
import ReviewSITExtensionsModal from 'components/Office/ReviewSITExtensionModal/ReviewSITExtensionModal';
import ConvertSITToCustomerExpenseModal from 'components/Office/ConvertSITToCustomerExpenseModal/ConvertSITToCustomerExpenseModal';
import ShipmentSITDisplay from 'components/Office/ShipmentSITDisplay/ShipmentSITDisplay';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates/ImportantShipmentDates';
import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';
import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';
import ShipmentRemarks from 'components/Office/ShipmentRemarks/ShipmentRemarks';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

/** @function OpenModalButton
 * The button that opens the modal in SIT Display component
 * @param {string} permission
 * @param {function} onClick
 * @param {string} className
 * @param {string} title
 * @returns {React.ReactElement}
 */
const OpenModalButton = ({ permission, onClick, className, title, isMoveLocked }) => (
  <Restricted to={permission}>
    <Button type="button" onClick={onClick} unstyled className={className} disabled={isMoveLocked}>
      {title}
    </Button>
  </Restricted>
);

const ShipmentDetailsMain = ({
  className,
  shipment,
  dutyLocationAddresses,
  handleShowDiversionModal,
  handleRequestReweighModal,
  handleReviewSITExtension,
  handleSubmitSITExtension,
  handleUpdateSITServiceItemCustomerExpense,
  isMoveLocked,
}) => {
  const {
    requestedPickupDate,
    scheduledPickupDate,
    actualPickupDate,
    requestedDeliveryDate,
    scheduledDeliveryDate,
    actualDeliveryDate,
    requiredDeliveryDate,
    pickupAddress,
    destinationAddress,
    ppmShipment,
    primeEstimatedWeight,
    primeActualWeight,
    counselorRemarks,
    customerRemarks,
    sitExtensions,
    sitStatus,
    storageInTransit,
    shipmentType,
    storageFacility,
    poeLocation,
    podLocation,
  } = shipment;
  const { originDutyLocationAddress, destinationDutyLocationAddress } = dutyLocationAddresses;

  const [isReviewSITExtensionModalVisible, setIsReviewSITExtensionModalVisible] = useState(false);
  const [isSubmitITExtensionModalVisible, setIsSubmitITExtensionModalVisible] = useState(false);
  const [isConvertSITToCustomerExpenseModalVisible, setIsConvertSITToCustomerExpenseModalVisible] = useState(false);
  const [, setSubmittedChangeTime] = useState(Date.now());

  const reviewSITExtension = (sitExtensionID, formValues) => {
    setIsReviewSITExtensionModalVisible(false);
    handleReviewSITExtension(sitExtensionID, formValues, shipment);
    setSubmittedChangeTime(Date.now());
  };

  const submitSITExtension = (formValues) => {
    setIsSubmitITExtensionModalVisible(false);
    handleSubmitSITExtension(formValues, shipment);
    setSubmittedChangeTime(Date.now());
  };

  const convertSITToCustomerExpense = (convertToCustomerExpense, remark) => {
    setIsConvertSITToCustomerExpenseModalVisible(false);
    handleUpdateSITServiceItemCustomerExpense(shipment.id, convertToCustomerExpense, remark, shipment.eTag);
    setSubmittedChangeTime(Date.now());
  };

  const pendingSITExtension = sitExtensions?.find((se) => se.status === SIT_EXTENSION_STATUS.PENDING);

  /**
   * Displays correct button to open the modal on the SIT Display component to open with either Sumbit or Review SIT modal.
   */
  const openModalButton = pendingSITExtension ? (
    <OpenModalButton
      permission={permissionTypes.createSITExtension}
      onClick={setIsReviewSITExtensionModalVisible}
      title="Review request"
      isMoveLocked={isMoveLocked}
    />
  ) : (
    <OpenModalButton
      permission={permissionTypes.updateSITExtension}
      onClick={setIsSubmitITExtensionModalVisible}
      title="Edit"
      className={styles.submitSITEXtensionLink}
      isMoveLocked={isMoveLocked}
    />
  );

  /**
   * Displays button to open the modal on the SIT Display component to open with Convert to customer expense modal.
   */
  const openConvertModalButton = (
    <OpenModalButton
      permission={permissionTypes.updateSITExtension}
      onClick={setIsConvertSITToCustomerExpenseModalVisible}
      title="Convert to customer expense"
      isMoveLocked={isMoveLocked}
    />
  );

  let displayedPickupAddress;
  let displayedDeliveryAddress;
  let weightResult;
  let pickupRequestedDate;
  let pickupScheduledDate;
  let pickupActualDate;
  let plannedMoveDate;
  let actualMoveDate;
  let displayPoeLocation;
  let displayPodLocation;

  switch (shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
      pickupRequestedDate = requestedPickupDate;
      pickupScheduledDate = scheduledPickupDate;
      pickupActualDate = actualPickupDate;
      weightResult = primeEstimatedWeight;
      displayedPickupAddress = pickupAddress;
      displayedDeliveryAddress = destinationAddress || destinationDutyLocationAddress;
      displayPoeLocation = poeLocation;
      displayPodLocation = podLocation;
      break;
    case SHIPMENT_OPTIONS.NTS:
      pickupRequestedDate = requestedPickupDate;
      pickupScheduledDate = scheduledPickupDate;
      pickupActualDate = actualPickupDate;
      weightResult = primeEstimatedWeight;
      displayedPickupAddress = pickupAddress;
      displayedDeliveryAddress = storageFacility ? storageFacility.address : null;
      break;
    case SHIPMENT_OPTIONS.NTSR:
      pickupRequestedDate = requestedPickupDate;
      pickupScheduledDate = scheduledPickupDate;
      pickupActualDate = actualPickupDate;
      weightResult = primeEstimatedWeight;
      displayedPickupAddress = storageFacility ? storageFacility.address : null;
      displayedDeliveryAddress = destinationAddress;
      break;
    case SHIPMENT_OPTIONS.PPM:
      plannedMoveDate = ppmShipment.expectedDepartureDate;
      actualMoveDate = ppmShipment.actualMoveDate;
      weightResult = ppmShipment.estimatedWeight;
      displayedPickupAddress = ppmShipment.pickupAddress;
      displayedDeliveryAddress = ppmShipment.destinationAddress || destinationDutyLocationAddress;
      break;
    default:
      pickupRequestedDate = requestedPickupDate;
      pickupScheduledDate = scheduledPickupDate;
      pickupActualDate = actualPickupDate;
      weightResult = primeEstimatedWeight;
      displayedPickupAddress = pickupAddress;
      displayedDeliveryAddress = destinationAddress || destinationDutyLocationAddress;
      displayPoeLocation = poeLocation;
      displayPodLocation = podLocation;
      break;
  }

  return (
    <div className={className}>
      {isReviewSITExtensionModalVisible && (
        <ReviewSITExtensionsModal
          onClose={() => setIsReviewSITExtensionModalVisible(false)}
          onSubmit={reviewSITExtension}
          shipment={shipment}
          sitExtension={pendingSITExtension}
          sitStatus={sitStatus}
        />
      )}
      {isConvertSITToCustomerExpenseModalVisible && (
        <ConvertSITToCustomerExpenseModal
          onClose={() => setIsConvertSITToCustomerExpenseModalVisible(false)}
          onSubmit={convertSITToCustomerExpense}
          shipment={shipment}
          sitStatus={sitStatus}
        />
      )}
      {isSubmitITExtensionModalVisible && (
        <SubmitSITExtensionModal
          onClose={() => setIsSubmitITExtensionModalVisible(false)}
          onSubmit={submitSITExtension}
          shipment={shipment}
          sitExtensions={sitExtensions}
          sitStatus={sitStatus}
        />
      )}
      {sitStatus && (
        <ShipmentSITDisplay
          sitExtensions={sitExtensions}
          sitStatus={sitStatus}
          storageInTransit={storageInTransit}
          shipment={shipment}
          className={styles.shipmentSITSummary}
          openModalButton={openModalButton}
          openConvertModalButton={openConvertModalButton}
        />
      )}
      {shipmentType === SHIPMENT_OPTIONS.PPM && (
        <ImportantShipmentDates
          plannedMoveDate={plannedMoveDate ? formatDateWithUTC(plannedMoveDate) : null}
          requestedDeliveryDate={requestedDeliveryDate ? formatDateWithUTC(requestedDeliveryDate) : null}
          scheduledDeliveryDate={scheduledDeliveryDate ? formatDateWithUTC(scheduledDeliveryDate) : null}
          actualMoveDate={actualMoveDate ? formatDateWithUTC(actualMoveDate) : null}
          isPPM={shipmentType === SHIPMENT_OPTIONS.PPM}
        />
      )}
      {shipmentType !== SHIPMENT_OPTIONS.PPM && (
        <ImportantShipmentDates
          requestedPickupDate={requestedPickupDate ? formatDateWithUTC(pickupRequestedDate) : null}
          scheduledPickupDate={scheduledPickupDate ? formatDateWithUTC(pickupScheduledDate) : null}
          actualPickupDate={actualPickupDate ? formatDateWithUTC(pickupActualDate) : null}
          requestedDeliveryDate={requestedDeliveryDate ? formatDateWithUTC(requestedDeliveryDate) : null}
          scheduledDeliveryDate={scheduledDeliveryDate ? formatDateWithUTC(scheduledDeliveryDate) : null}
          actualDeliveryDate={actualDeliveryDate ? formatDateWithUTC(actualDeliveryDate) : null}
          requiredDeliveryDate={requiredDeliveryDate ? formatDateWithUTC(requiredDeliveryDate) : null}
          isPPM={shipmentType === SHIPMENT_OPTIONS.PPM}
          shipmentInfo={{
            id: shipment.id,
            status: shipment.status,
            shipmentType: shipment.shipmentType,
            shipmentLocator: shipment.shipmentLocator,
            isDiversion: shipment.diversion,
            diversionReason: shipment.diversionReason,
          }}
        />
      )}
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
          shipmentLocator: shipment.shipmentLocator,
        }}
        handleShowDiversionModal={handleShowDiversionModal}
        isMoveLocked={isMoveLocked}
        poeLocation={displayPoeLocation}
        podLocation={displayPodLocation}
      />
      <ShipmentWeightDetails
        estimatedWeight={weightResult}
        initialWeight={primeActualWeight}
        shipmentInfo={{
          shipmentID: shipment.id,
          ifMatchEtag: shipment.eTag,
          reweighID: shipment.reweigh?.id,
          reweighWeight: shipment.reweigh?.weight,
          shipmentType: shipment.shipmentType,
          shipmentActualProGearWeight: shipment.actualProGearWeight,
          shipmentActualSpouseProGearWeight: shipment.actualSpouseProGearWeight,
        }}
        handleRequestReweighModal={handleRequestReweighModal}
        isMoveLocked={isMoveLocked}
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
  handleShowDiversionModal: PropTypes.func.isRequired,
  handleRequestReweighModal: PropTypes.func.isRequired,
  handleReviewSITExtension: PropTypes.func.isRequired,
  handleSubmitSITExtension: PropTypes.func.isRequired,
  handleUpdateSITServiceItemCustomerExpense: PropTypes.func.isRequired,
};

ShipmentDetailsMain.defaultProps = {
  className: '',
};

export default ShipmentDetailsMain;
