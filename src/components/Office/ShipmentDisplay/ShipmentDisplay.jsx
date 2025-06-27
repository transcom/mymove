import React, { useEffect, useState } from 'react';
import * as PropTypes from 'prop-types';
import { useNavigate } from 'react-router-dom';
import { Checkbox, Tag, Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';
import { connect } from 'react-redux';

import ErrorModal from 'shared/ErrorModal/ErrorModal';
import { EditButton, ReviewButton } from 'components/form/IconButtons';
import ShipmentInfoListSelector from 'components/Office/DefinitionLists/ShipmentInfoListSelector';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import styles from 'components/Office/ShipmentDisplay/ShipmentDisplay.module.scss';
import { FEATURE_FLAG_KEYS, getPPMTypeLabel, PPM_TYPES, SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
import { AddressShape } from 'types/address';
import { AgentShape } from 'types/agent';
import { OrdersLOAShape } from 'types/order';
import { shipmentStatuses, ppmShipmentStatuses, ppmShipmentStatusLabels } from 'constants/shipments';
import { ShipmentStatusesOneOf } from 'types/shipment';
import { retrieveSAC, retrieveTAC } from 'utils/shipmentDisplay';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import affiliation from 'content/serviceMemberAgencies';
import { fieldValidationShape, objectIsMissingFieldWithCondition } from 'utils/displayFlags';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { SubmitMoveConfirmationModal } from 'components/Office/SubmitMoveConfirmationModal/SubmitMoveConfirmationModal';

const ShipmentDisplay = ({
  shipmentType,
  displayInfo,
  onChange,
  shipmentId,
  isSubmitted,
  allowApproval,
  editURL,
  reviewURL,
  sendPpmToCustomer,
  counselorCanEdit,
  completePpmForCustomerURL,
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
  const [isSubmitPPMShipmentModalVisible, setIsSubmitPPMShipmentModalVisible] = useState(false);
  const [isErrorModalVisible, setIsErrorModalVisible] = useState(false);
  const [enableCompletePPMCloseoutForCustomer, setEnableCompletePPMCloseoutForCustomer] = useState(false);
  const [ppmSprFF, setPpmSprFF] = useState(false);
  const isDisabled = isMoveLocked || displayInfo.shipmentStatus === shipmentStatuses.TERMINATED_FOR_CAUSE;

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

  const handleShowSubmitPPMShipmentModal = () => {
    setIsSubmitPPMShipmentModalVisible(true);
  };

  const handleSubmitPPMShipmentModal = () => {
    sendPpmToCustomer({ ppmShipmentId: displayInfo.ppmShipment.id, eTag: displayInfo.ppmShipment.eTag });
    setIsSubmitPPMShipmentModalVisible();
  };

  const errorModalMessage =
    "Something went wrong downloading PPM paperwork. Please try again later. If that doesn't fix it, contact the ";

  useEffect(() => {
    const fetchData = async () => {
      setEnableCompletePPMCloseoutForCustomer(
        await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.COMPLETE_PPM_CLOSEOUT_FOR_CUSTOMER),
      );
      setPpmSprFF(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.PPM_SPR));
    };
    fetchData();
  }, []);

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
                disabled={disableApproval || isDisabled}
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
              <div className={styles.tagWrapper}>
                {displayInfo.shipmentStatus === shipmentStatuses.TERMINATED_FOR_CAUSE && (
                  <Tag data-testid="terminatedTag" className="usa-tag--cancellation">
                    terminated for cause
                  </Tag>
                )}
                {ppmSprFF && displayInfo.ppmShipment?.ppmType === PPM_TYPES.SMALL_PACKAGE && (
                  <Tag data-testid="smallPackageTag">{getPPMTypeLabel(displayInfo.ppmShipment.ppmType)}</Tag>
                )}
                {displayInfo.ppmShipment?.ppmType === PPM_TYPES.ACTUAL_EXPENSE && (
                  <Tag data-testid="actualReimbursementTag">{getPPMTypeLabel(displayInfo.ppmShipment.ppmType)}</Tag>
                )}
                {displayInfo.isDiversion && <Tag className="usa-tag--diversion">diversion</Tag>}
                {(displayInfo.shipmentStatus === shipmentStatuses.CANCELED ||
                  displayInfo.status === shipmentStatuses.CANCELED) && (
                  <Tag className="usa-tag--cancellation">canceled</Tag>
                )}
                {displayInfo.shipmentStatus === shipmentStatuses.DIVERSION_REQUESTED && (
                  <Tag className="usa-tag--diversion">diversion requested</Tag>
                )}
                {displayInfo.shipmentStatus === shipmentStatuses.CANCELLATION_REQUESTED && (
                  <Tag className="usa-tag--cancellation">cancellation requested</Tag>
                )}
                {displayInfo.usesExternalVendor && <Tag>external vendor</Tag>}
                {displayInfo.ppmShipment?.status && (
                  <Tag
                    className={
                      displayInfo.ppmShipment.status !== ppmShipmentStatuses.CANCELED
                        ? 'usa-tag--ppmStatus'
                        : 'usa-tag--cancellation'
                    }
                    data-testid="ppmStatusTag"
                  >
                    {ppmShipmentStatusLabels[displayInfo.ppmShipment?.status]}
                  </Tag>
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
        {isSubmitPPMShipmentModalVisible && (
          <SubmitMoveConfirmationModal
            onClose={setIsSubmitPPMShipmentModalVisible}
            onSubmit={handleSubmitPPMShipmentModal}
            isShipment
          />
        )}
        <Restricted to={permissionTypes.updateShipment}>
          <div className={styles.flexRight}>
            {editURL && (
              <EditButton
                onClick={() => {
                  navigate(editURL);
                }}
                className={styles.editButton}
                data-testid={editURL}
                label="Edit shipment"
                secondary
                disabled={isDisabled}
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
                disabled={isDisabled}
              />
            )}
            {completePpmForCustomerURL && enableCompletePPMCloseoutForCustomer && (
              <Button
                onClick={() => {
                  navigate(completePpmForCustomerURL);
                }}
                className={styles.editButton}
                data-testid="completePpmForCustomerBtn"
                secondary
                disabled={isDisabled}
              >
                Complete PPM on behalf of the Customer
              </Button>
            )}
            {sendPpmToCustomer &&
              displayInfo.ppmShipment?.status === ppmShipmentStatuses.SUBMITTED &&
              !counselorCanEdit && (
                <Button
                  onClick={() => {
                    handleShowSubmitPPMShipmentModal();
                  }}
                  className={styles.editButton}
                  data-testid="sendPpmToCustomerButton"
                  secondary
                  disabled={isMoveLocked}
                >
                  Send PPM to the Customer
                </Button>
              )}
          </div>
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
            disabled={isDisabled}
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
    SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
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
  sendPpmToCustomer: PropTypes.func,
  counselorCanEdit: PropTypes.bool,
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
  sendPpmToCustomer: null,
  counselorCanEdit: false,
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

export default connect(() => ({}))(ShipmentDisplay);
