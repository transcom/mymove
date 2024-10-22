import React, { useState } from 'react';
import * as PropTypes from 'prop-types';
import { useNavigate } from 'react-router-dom';
import { Checkbox, Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import ErrorModal from 'shared/ErrorModal/ErrorModal';
import { EditButton, ReviewButton } from 'components/form/IconButtons';
import ShipmentInfoListSelector from 'components/Office/DefinitionLists/ShipmentInfoListSelector';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import styles from 'components/Office/ShipmentDisplay/ShipmentDisplay.module.scss';
import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
import { AddressShape } from 'types/address';
import { AgentShape } from 'types/agent';
import { OrdersLOAShape } from 'types/order';
import { shipmentStatuses, ppmShipmentStatuses } from 'constants/shipments';
import { ShipmentStatusesOneOf } from 'types/shipment';
import { retrieveSAC, retrieveTAC } from 'utils/shipmentDisplay';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import affiliation from 'content/serviceMemberAgencies';
import { fieldValidationShape, objectIsMissingFieldWithCondition } from 'utils/displayFlags';

const ShipmentDisplay = ({
  shipmentType,
  displayInfo,
  onChange,
  shipmentId,
  isSubmitted,
  allowApproval,
  editURL,
  reviewURL,
  viewURL,
  ordersLOA,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
  neverShow,
  isMoveLocked,
}) => {
  const navigate = useNavigate();
  const containerClasses = classnames(styles.container, { [styles.noIcon]: !allowApproval });
  const [isExpanded, setIsExpanded] = useState(false);
  const tac = retrieveTAC(displayInfo.tacType, ordersLOA);
  const sac = retrieveSAC(displayInfo.sacType, ordersLOA);
  const [isErrorModalVisible, setIsErrorModalVisible] = useState(false);

  const disableApproval = errorIfMissing.some((requiredInfo) =>
    objectIsMissingFieldWithCondition(displayInfo, requiredInfo),
  );

  const handleExpandClick = () => {
    setIsExpanded((prev) => !prev);
  };
  const expandableIconClasses = classnames({
    'chevron-up': isExpanded,
    'chevron-down': !isExpanded,
  });

  const toggleErrorModal = () => {
    setIsErrorModalVisible((prev) => !prev);
  };

  const errorModalMessage =
    "Something went wrong downloading PPM paperwork. Please try again later. If that doesn't fix it, contact the ";

  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={containerClasses} shipmentType={shipmentType}>
        <div className={styles.heading}>
          <Restricted to={permissionTypes.updateShipment}>
            {allowApproval && isSubmitted && !displayInfo.usesExternalVendor && (
              <Checkbox
                id={`shipment-display-checkbox-${shipmentId}`}
                data-testid="shipment-display-checkbox"
                onChange={onChange}
                name="shipments"
                label="&nbsp;"
                value={shipmentId}
                aria-labelledby={`shipment-display-label-${shipmentId}`}
                disabled={disableApproval || isMoveLocked}
              />
            )}
          </Restricted>

          {allowApproval && !isSubmitted && (
            <FontAwesomeIcon icon={['far', 'circle-check']} className={styles.approved} />
          )}
          <div className={styles.headerContainer}>
            <div className={styles.shipmentTypeHeader}>
              <h3>
                <label id={`shipment-display-label-${shipmentId}`}>
                  <span className={styles.marketCodeIndicator}>{displayInfo.marketCode}</span>
                  {displayInfo.heading}
                </label>
              </h3>
              <div>
                {displayInfo.isDiversion && <Tag>diversion</Tag>}
                {displayInfo.shipmentStatus === shipmentStatuses.CANCELED && (
                  <Tag className="usa-tag--red">cancelled</Tag>
                )}
                {displayInfo.shipmentStatus === shipmentStatuses.DIVERSION_REQUESTED && <Tag>diversion requested</Tag>}
                {displayInfo.shipmentStatus === shipmentStatuses.CANCELLATION_REQUESTED && (
                  <Tag>cancellation requested</Tag>
                )}
                {displayInfo.usesExternalVendor && <Tag>external vendor</Tag>}
                {(displayInfo.ppmShipment?.status === ppmShipmentStatuses.CLOSEOUT_COMPLETE ||
                  displayInfo.ppmShipment?.status === ppmShipmentStatuses.WAITING_ON_CUSTOMER) && (
                  <Tag className={styles.ppmStatus}>packet ready for download</Tag>
                )}
              </div>
            </div>
            <div className={styles.shipmentLocator}>
              <h5>#{displayInfo.shipmentLocator}</h5>
            </div>
          </div>

          <FontAwesomeIcon className={styles.icon} icon={expandableIconClasses} onClick={handleExpandClick} />
        </div>
        <ShipmentInfoListSelector
          className={styles.shipmentDisplayInfo}
          shipment={{ ...displayInfo, tac, sac }}
          shipmentType={shipmentType}
          isExpanded={isExpanded}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          showWhenCollapsed={showWhenCollapsed}
          neverShow={neverShow}
          onErrorModalToggle={toggleErrorModal}
        />
        <ErrorModal isOpen={isErrorModalVisible} closeModal={toggleErrorModal} errorMessage={errorModalMessage} />
        <Restricted to={permissionTypes.updateShipment}>
          {editURL && (
            <EditButton
              onClick={() => {
                navigate(editURL);
              }}
              className={styles.editButton}
              data-testid={editURL}
              label="Edit shipment"
              secondary
              disabled={isMoveLocked}
            />
          )}
          {reviewURL && (
            <ReviewButton
              onClick={() => {
                navigate(reviewURL);
              }}
              className={styles.editButton}
              data-testid={reviewURL}
              label="Review documents"
              secondary
              disabled={isMoveLocked}
            />
          )}
        </Restricted>
        {viewURL && (
          <ReviewButton
            onClick={() => {
              navigate(viewURL);
            }}
            className={styles.editButton}
            data-testid={viewURL}
            label="View documents"
            secondary
            disabled={isMoveLocked}
          />
        )}
      </ShipmentContainer>
    </div>
  );
};

ShipmentDisplay.propTypes = {
  onChange: PropTypes.func,
  shipmentId: PropTypes.string.isRequired,
  isSubmitted: PropTypes.bool.isRequired,
  shipmentType: PropTypes.oneOf([
    SHIPMENT_OPTIONS.HHG,
    SHIPMENT_OPTIONS.NTS,
    SHIPMENT_OPTIONS.NTSR,
    SHIPMENT_OPTIONS.PPM,
    SHIPMENT_TYPES.BOAT_HAUL_AWAY,
    SHIPMENT_TYPES.BOAT_TOW_AWAY,
    SHIPMENT_OPTIONS.MOBILE_HOME,
  ]),
  displayInfo: PropTypes.oneOfType([
    PropTypes.shape({
      agency: PropTypes.oneOf(Object.values(affiliation)),
      closeoutOffice: PropTypes.string,
      heading: PropTypes.string.isRequired,
      isDiversion: PropTypes.bool,
      shipmentStatus: ShipmentStatusesOneOf,
      requestedPickupDate: PropTypes.string,
      pickupAddress: AddressShape,
      secondaryPickupAddress: AddressShape,
      destinationAddress: AddressShape,
      destinationType: PropTypes.string,
      displayDestinationType: PropTypes.bool,
      secondaryDeliveryAddress: AddressShape,
      counselorRemarks: PropTypes.string,
      shipmentId: PropTypes.string,
      shipmentType: PropTypes.string,
      usesExternalVendor: PropTypes.bool,
      customerRemarks: PropTypes.string,
      serviceOrderNumber: PropTypes.string,
      requestedDeliveryDate: PropTypes.string,
      agents: PropTypes.arrayOf(AgentShape),
      primeActualWeight: PropTypes.number,
      storageFacility: PropTypes.shape({
        address: AddressShape.isRequired,
        facilityName: PropTypes.string,
        lotNumber: PropTypes.string,
      }),
      tacType: PropTypes.string,
      sacType: PropTypes.string,
      ntsRecordedWeight: PropTypes.number,
      shipmentLocator: PropTypes.string,
    }),
    PropTypes.shape({
      heading: PropTypes.string.isRequired,
      shipmentType: PropTypes.string,
      hasRequestedAdvance: PropTypes.bool,
      advanceAmountRequested: PropTypes.number,
      estimatedIncentive: PropTypes.number,
      estimatedWeight: PropTypes.string,
      expectedDepartureDate: PropTypes.string,
      proGearWeight: PropTypes.string,
      spouseProGearWeight: PropTypes.string,
      customerRemarks: PropTypes.string,
      tacType: PropTypes.string,
      sacType: PropTypes.string,
    }),
  ]).isRequired,
  allowApproval: PropTypes.bool,
  editURL: PropTypes.string,
  reviewURL: PropTypes.string,
  ordersLOA: OrdersLOAShape,
  warnIfMissing: PropTypes.arrayOf(fieldValidationShape),
  errorIfMissing: PropTypes.arrayOf(fieldValidationShape),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
  neverShow: PropTypes.arrayOf(PropTypes.string),
};

ShipmentDisplay.defaultProps = {
  onChange: () => {},
  shipmentType: SHIPMENT_OPTIONS.HHG,
  allowApproval: true,
  editURL: '',
  reviewURL: '',
  ordersLOA: {
    tac: '',
    sac: '',
    ntsTac: '',
    ntsSac: '',
  },
  warnIfMissing: [],
  errorIfMissing: [],
  showWhenCollapsed: [],
  neverShow: [],
};

export default ShipmentDisplay;
