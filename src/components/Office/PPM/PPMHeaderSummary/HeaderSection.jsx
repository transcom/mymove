import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Label, Button, Alert } from '@trussworks/react-uswds';
import { useQueryClient, useMutation, useIsFetching } from '@tanstack/react-query';
import classnames from 'classnames';

import styles from './HeaderSection.module.scss';

import EditPPMHeaderSummaryModal from 'components/Office/PPM/PPMHeaderSummary/EditPPMHeaderSummaryModal';
import { formatDate, formatCents, formatWeight } from 'utils/formatters';
import { MTO_SHIPMENTS, PPMCLOSEOUT } from 'constants/queryKeys';
import { updateMTOShipment } from 'services/ghcApi';
import { useEditShipmentQueries, usePPMShipmentDocsQueries } from 'hooks/queries';

export const sectionTypes = {
  incentives: 'incentives',
  shipmentInfo: 'shipmentInfo',
  incentiveFactors: 'incentiveFactors',
};

const HAUL_TYPES = {
  SHORTHAUL: 'Shorthaul',
  LINEHAUL: 'Linehaul',
};

const getSectionTitle = (sectionInfo) => {
  switch (sectionInfo.type) {
    case sectionTypes.incentives:
      return `Incentives/Costs`;
    case sectionTypes.shipmentInfo:
      return `Shipment Info`;
    case sectionTypes.incentiveFactors:
      return `Incentive Factors`;
    default:
      return <Alert>Error getting section title!</Alert>;
  }
};

const OpenModalButton = ({ onClick, isDisabled }) => (
  <Button
    type="button"
    data-testid="editTextButton"
    className={styles['edit-btn']}
    onClick={onClick}
    disabled={isDisabled}
  >
    <span>
      <FontAwesomeIcon icon="pencil" style={{ marginRight: '5px', color: isDisabled ? 'black' : 'inherit' }} />
    </span>
  </Button>
);

// Returns the markup needed for a specific section
const getSectionMarkup = (sectionInfo, handleEditOnClick, isFetchingItems, updatedItemName, readOnly, grade) => {
  const aoaRequestedValue = `$${formatCents(sectionInfo.advanceAmountRequested)}`;
  const aoaValue = `$${formatCents(sectionInfo.advanceAmountReceived)}`;
  const isCivilian = grade === 'CIVILIAN_EMPLOYEE';

  const renderHaulType = (haulType) => {
    return haulType === HAUL_TYPES.LINEHAUL ? 'Linehaul' : 'Shorthaul';
  };
  // check if the itemName is one of the items recalulated after item edit(updatedItemName).
  const isRecalulatedItem = (itemName) => {
    let recalulatedItems = [];
    const incentiveFactors = ['haulPrice', 'haulFSC', 'packPrice', 'unpackPrice', 'dop', 'ddp', 'sitReimbursement'];

    switch (updatedItemName) {
      case 'actualMoveDate':
        recalulatedItems = ['gcc', 'grossIncentive', ...incentiveFactors];
        break;
      case 'pickupAddress':
      case 'destinationAddress':
        recalulatedItems = ['miles', 'gcc', 'grossIncentive', 'remainingIncentive', ...incentiveFactors];
        break;
      case 'advanceAmountReceived':
        if (incentiveFactors.includes(itemName)) return false;
        recalulatedItems = ['remainingIncentive'];
        break;
      default:
        break;
    }

    return recalulatedItems.includes(itemName);
  };

  switch (sectionInfo.type) {
    case sectionTypes.shipmentInfo:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>Actual Expense Reimbursement</Label>
            <span data-testid="isActualExpenseReimbursement" className={styles.light}>
              {isFetchingItems && updatedItemName === 'isActualExpenseReimbursement' ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                <>
                  {sectionInfo.isActualExpenseReimbursement ? 'Yes' : 'No'}
                  <OpenModalButton
                    onClick={() => handleEditOnClick(sectionInfo.type, 'isActualExpenseReimbursement')}
                    isDisabled={isFetchingItems || readOnly || isCivilian}
                  />
                </>
              )}
            </span>
          </div>
          <div>
            <Label>Planned Move Start Date</Label>
            <span className={styles.light}>{formatDate(sectionInfo.plannedMoveDate, null, 'DD-MMM-YYYY')}</span>
          </div>
          <div>
            <Label>Actual Move Start Date</Label>
            <span data-testid="actualMoveDate" className={styles.light}>
              {isFetchingItems && updatedItemName === 'actualMoveDate' ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                <>
                  {formatDate(sectionInfo.actualMoveDate, null, 'DD-MMM-YYYY')}
                  <OpenModalButton
                    onClick={() => handleEditOnClick(sectionInfo.type, 'actualMoveDate')}
                    isDisabled={isFetchingItems || readOnly}
                  />
                </>
              )}
            </span>
          </div>
          <div>
            <Label>Starting Address</Label>
            <span data-testid="pickupAddress" className={styles.light}>
              {isFetchingItems && updatedItemName === 'pickupAddress' ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                <>
                  {sectionInfo.pickupAddress}
                  <OpenModalButton
                    onClick={() => handleEditOnClick(sectionInfo.type, 'pickupAddress')}
                    isDisabled={isFetchingItems || readOnly}
                  />
                </>
              )}
            </span>
          </div>
          <div>
            <Label>Ending Address</Label>
            <span data-testid="destinationAddress" className={styles.light}>
              {isFetchingItems && updatedItemName === 'destinationAddress' ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                <>
                  {sectionInfo.destinationAddress}
                  <OpenModalButton
                    onClick={() => handleEditOnClick(sectionInfo.type, 'destinationAddress')}
                    isDisabled={isFetchingItems || readOnly}
                  />
                </>
              )}
            </span>
          </div>
          <div>
            <Label>Miles</Label>
            <span className={styles.light}>
              {isFetchingItems && isRecalulatedItem('miles') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                sectionInfo.miles
              )}
            </span>
          </div>
          <div>
            <Label>Estimated Net Weight</Label>
            <span className={styles.light}>{formatWeight(sectionInfo.estimatedWeight)}</span>
          </div>
          <div>
            <Label>Actual Net Weight</Label>
            <span className={styles.light}>{formatWeight(sectionInfo.actualWeight)}</span>
          </div>
        </div>
      );

    case sectionTypes.incentives:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>Government Constructed Cost (GCC)</Label>
            <span data-testid="gcc" className={styles.light}>
              {isFetchingItems && isRecalulatedItem('gcc') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                `$${formatCents(sectionInfo.gcc)}`
              )}
            </span>
          </div>
          <div>
            <Label>Gross Incentive</Label>
            <span data-testid="grossIncentive" className={styles.light}>
              {isFetchingItems && isRecalulatedItem('grossIncentive') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                `$${formatCents(sectionInfo.grossIncentive)}`
              )}
            </span>
          </div>
          <div>
            <Label>Advance Requested</Label>
            <span data-testid="advanceRequested" className={styles.light}>
              {aoaRequestedValue}
            </span>
          </div>
          <div>
            <Label>Advance Received</Label>
            <span data-testid="advanceReceived" className={styles.light}>
              {isFetchingItems && updatedItemName === 'advanceAmountReceived' ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                <>
                  {aoaValue}
                  <OpenModalButton
                    onClick={() => handleEditOnClick(sectionInfo.type, 'advanceAmountReceived')}
                    isDisabled={isFetchingItems || readOnly}
                  />
                </>
              )}
            </span>
          </div>
          <div>
            <Label>Remaining Incentive</Label>
            <span data-testid="remainingIncentive" className={styles.light}>
              {isFetchingItems && isRecalulatedItem('remainingIncentive') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                `$${formatCents(sectionInfo.remainingIncentive)}`
              )}
            </span>
          </div>
        </div>
      );

    case sectionTypes.incentiveFactors:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>{renderHaulType(sectionInfo.haulType)} Price</Label>
            <span data-testid="haulPrice" className={styles.light}>
              {isFetchingItems && isRecalulatedItem('haulPrice') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                `$${formatCents(sectionInfo.haulPrice)}`
              )}
            </span>
          </div>
          <div>
            <Label>{renderHaulType(sectionInfo.haulType)} Fuel Rate Adjustment</Label>
            <span data-testid="haulFSC" className={styles.light}>
              {isFetchingItems && isRecalulatedItem('haulFSC') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                <>
                  {sectionInfo.haulFSC < 0 ? '-$' : '$'}
                  {formatCents(Math.abs(sectionInfo.haulFSC))}
                </>
              )}
            </span>
          </div>
          <div>
            <Label>Packing Charge</Label>
            <span data-testid="packPrice" className={styles.light}>
              {isFetchingItems && isRecalulatedItem('packPrice') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                `$${formatCents(sectionInfo.packPrice)}`
              )}
            </span>
          </div>
          <div>
            <Label>Unpacking Charge</Label>
            <span data-testid="unpackPrice" className={styles.light}>
              {isFetchingItems && isRecalulatedItem('unpackPrice') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                `$${formatCents(sectionInfo.unpackPrice)}`
              )}
            </span>
          </div>
          <div>
            <Label>Origin Price</Label>
            <span data-testid="originPrice" className={styles.light}>
              {isFetchingItems && isRecalulatedItem('dop') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                `$${formatCents(sectionInfo.dop)}`
              )}
            </span>
          </div>
          <div>
            <Label>Destination Price</Label>
            <span data-testid="destinationPrice" className={styles.light}>
              {isFetchingItems && isRecalulatedItem('ddp') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                `$${formatCents(sectionInfo.ddp)}`
              )}
            </span>
          </div>
          <div>
            <Label>SIT Reimbursement</Label>
            <span data-testid="sitReimbursement" className={styles.light}>
              ${formatCents(sectionInfo.sitReimbursement)}
            </span>
          </div>
        </div>
      );

    default:
      return <Alert>An error occured while getting section markup!</Alert>;
  }
};

export default function HeaderSection({
  sectionInfo,
  dataTestId,
  updatedItemName,
  setUpdatedItemName,
  readOnly,
  grade,
}) {
  const requestDetailsButtonTestId = `${sectionInfo.type}-showRequestDetailsButton`;
  const { shipmentId, moveCode } = useParams();
  const { mtoShipment, refetchMTOShipment, isFetching: isFetchingMtoShipment } = usePPMShipmentDocsQueries(shipmentId);
  const queryClient = useQueryClient();
  const [showDetails, setShowDetails] = useState(false);
  const [isEditModalVisible, setIsEditModalVisible] = useState(false);
  const [itemName, setItemName] = useState('');
  const [sectionType, setSectionType] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const showRequestDetailsButton = true;
  const handleToggleDetails = () => setShowDetails((prevState) => !prevState);
  const showDetailsChevron = showDetails ? 'chevron-up' : 'chevron-down';
  const showDetailsText = showDetails ? 'Hide details' : 'Show details';

  const isFetchingCloseout = useIsFetching({ queryKey: [PPMCLOSEOUT, mtoShipment?.ppmShipment?.id] }) > 0;
  const isFetchingItems = isFetchingMtoShipment || isFetchingCloseout;

  const handleEditOnClose = () => {
    setIsEditModalVisible(false);
    setItemName('');
    setSectionType('');
    setIsSubmitting(false);
  };

  useEffect(() => {
    if (isEditModalVisible) {
      refetchMTOShipment();
    }
  }, [isEditModalVisible, refetchMTOShipment]);

  const { mtoShipments } = useEditShipmentQueries(moveCode);

  const { mutate: mutateMTOShipment } = useMutation(updateMTOShipment, {
    onSuccess: (updatedMTOShipments) => {
      const updatedMTOShipment = updatedMTOShipments.mtoShipments[shipmentId];
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);
      queryClient.invalidateQueries([PPMCLOSEOUT, updatedMTOShipment?.ppmShipment?.id]);
      refetchMTOShipment();
      handleEditOnClose();
    },
    onSettled: () => {
      setIsSubmitting(false);
    },
  });

  const handleEditOnClick = (type, name) => {
    setIsEditModalVisible(true);
    setItemName(name);
    setSectionType(type);
    setUpdatedItemName(name);
  };

  const handleEditSubmit = (values) => {
    if (isSubmitting) return;

    setIsSubmitting(true);
    let body = {};

    switch (itemName) {
      case 'actualMoveDate':
        body = { actualMoveDate: formatDate(values.actualMoveDate, 'DD MMM YYYY', 'YYYY-MM-DD') };
        break;
      case 'advanceAmountReceived':
        if (values.advanceAmountReceived === '0') {
          body = {
            advanceAmountReceived: null,
            hasReceivedAdvance: false,
          };
        } else {
          body = {
            advanceAmountReceived: values.advanceAmountReceived * 100,
            hasReceivedAdvance: true,
          };
        }
        break;
      case 'pickupAddress':
        body = {
          pickupAddress: values.pickupAddress,
          actualPickupPostalCode: values.pickupAddress?.postalCode,
        };
        break;
      case 'destinationAddress':
        body = {
          destinationAddress: values.destinationAddress,
          actualDestinationPostalCode: values.destinationAddress?.postalCode,
        };
        break;
      case 'isActualExpenseReimbursement':
        body = {
          isActualExpenseReimbursement: values.isActualExpenseReimbursement === 'true',
        };
        break;

      default:
        break;
    }

    mutateMTOShipment({
      moveTaskOrderID: mtoShipment.moveTaskOrderID,
      shipmentID: mtoShipment.id,
      ifMatchETag: mtoShipment.eTag,
      body: {
        ppmShipment: body,
      },
    });
  };

  return (
    <section className={classnames(styles.HeaderSection)} data-testid={dataTestId}>
      <header>
        <h4>{getSectionTitle(sectionInfo)}</h4>
      </header>
      <div className={styles.toggleDrawer}>
        {showRequestDetailsButton && (
          <Button
            aria-expanded={showDetails}
            data-testid={requestDetailsButtonTestId}
            type="button"
            unstyled
            onClick={handleToggleDetails}
          >
            <FontAwesomeIcon icon={showDetailsChevron} /> {showDetailsText}
          </Button>
        )}
      </div>
      {showDetails &&
        getSectionMarkup(sectionInfo, handleEditOnClick, isFetchingItems, updatedItemName, readOnly, grade)}
      {isEditModalVisible && (
        <EditPPMHeaderSummaryModal
          onClose={handleEditOnClose}
          onSubmit={handleEditSubmit}
          sectionInfo={sectionInfo}
          editItemName={itemName}
          sectionType={sectionType}
        />
      )}
    </section>
  );
}
