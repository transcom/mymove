import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Label, Button, Alert } from '@trussworks/react-uswds';
import { useQueryClient, useMutation, useIsFetching } from '@tanstack/react-query';
import classnames from 'classnames';
import propTypes from 'prop-types';

import styles from './HeaderSection.module.scss';

import ToolTip from 'shared/ToolTip/ToolTip';
import EditPPMHeaderSummaryModal from 'components/Office/PPM/PPMHeaderSummary/EditPPMHeaderSummaryModal';
import { formatDate, formatCents, formatWeight, calculateTotal } from 'utils/formatters';
import { MTO_SHIPMENTS, PPMCLOSEOUT } from 'constants/queryKeys';
import { updateMTOShipment } from 'services/ghcApi';
import { useEditShipmentQueries, usePPMShipmentDocsQueries } from 'hooks/queries';
import { getPPMTypeLabel, PPM_TYPES } from 'shared/constants';
import { getTotalPackageWeightSPR, hasProGearSPR, hasSpouseProGearSPR } from 'utils/ppmCloseout';
import { ORDERS_PAY_GRADE_TYPE } from 'constants/orders';
import { renderMultiplier } from 'constants/ppms';

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
    case sectionTypes.shipmentInfo:
      return `Shipment Info`;
    case sectionTypes.incentives:
      return `Incentives/Costs`;
    case sectionTypes.incentiveFactors:
      return `Incentive Factors`;
    default:
      return <Alert>Error getting section title!</Alert>;
  }
};

const OpenModalButton = ({ onClick, isDisabled, dataTestId, ariaLabel }) => (
  <Button
    type="button"
    data-testid={dataTestId || 'editTextButton'}
    className={styles['edit-btn']}
    onClick={onClick}
    disabled={isDisabled}
    aria-label={ariaLabel}
    title="Edit"
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
  const isCivilian = grade === ORDERS_PAY_GRADE_TYPE.CIVILIAN_EMPLOYEE;

  const renderHaulType = (haulType) => {
    if (haulType === '') {
      return null;
    }
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
            <Label>Expense Type</Label>
            <span data-testid="expenseType" className={styles.light}>
              {isFetchingItems && updatedItemName === 'expenseType' ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                <>
                  {getPPMTypeLabel(sectionInfo.ppmType)}
                  <OpenModalButton
                    onClick={() => handleEditOnClick(sectionInfo.type, 'expenseType')}
                    isDisabled={isFetchingItems || readOnly || isCivilian}
                  />
                </>
              )}
            </span>
          </div>
          {sectionInfo.ppmType !== PPM_TYPES.SMALL_PACKAGE ? (
            <>
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
                        ariaLabel="Edit actual move start date"
                      />
                    </>
                  )}
                </span>
              </div>
            </>
          ) : (
            <div>
              <Label>Shipped Date</Label>
              <span className={styles.light}>{formatDate(sectionInfo.plannedMoveDate, null, 'DD-MMM-YYYY')}</span>
            </div>
          )}
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
                    ariaLabel="Edit pickup address"
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
                    ariaLabel="Edit destination address"
                  />
                </>
              )}
            </span>
          </div>
          {sectionInfo.ppmType !== PPM_TYPES.SMALL_PACKAGE ? (
            <>
              {sectionInfo.miles && (
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
              )}
              <div>
                <Label>Estimated Net Weight</Label>
                <span className={styles.light}>{formatWeight(sectionInfo.estimatedWeight)}</span>
              </div>
              <div>
                <Label>Actual Net Weight</Label>
                <span>{formatWeight(sectionInfo.actualWeight)}</span>
              </div>
              <div>
                <Label>Allowable Weight</Label>
                {isFetchingItems && updatedItemName === 'allowableWeight' ? (
                  <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
                ) : (
                  <>
                    <b>{formatWeight(sectionInfo.allowableWeight)}</b>
                    <OpenModalButton
                      onClick={() => handleEditOnClick(sectionInfo.type, 'allowableWeight')}
                      isDisabled={isFetchingItems || readOnly}
                      dataTestId="editAllowableWeightButton"
                      ariaLabel="Edit allowable weight"
                    />
                  </>
                )}
                <ToolTip
                  icon="info-circle"
                  style={{ display: 'inline-block', height: '15px', margin: '0' }}
                  textAreaSize="large"
                  text="The total PPM weight moved (all trips combined). The Counselor may edit this field to reflect the customer's remaining weight entitlement if the combined weight of all shipments exceeds the customer's weight entitlement."
                />
              </div>
            </>
          ) : (
            <>
              <div>
                <Label>Allowable Weight</Label>
                {isFetchingItems && updatedItemName === 'allowableWeight' ? (
                  <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
                ) : (
                  <>
                    <span>{formatWeight(sectionInfo.allowableWeight)}</span>
                    <OpenModalButton
                      onClick={() => handleEditOnClick(sectionInfo.type, 'allowableWeight')}
                      isDisabled={isFetchingItems || readOnly}
                      dataTestId="editAllowableWeightButton"
                      ariaLabel="Edit allowable weight"
                    />
                  </>
                )}
                <ToolTip
                  icon="info-circle"
                  style={{ display: 'inline-block', height: '15px', margin: '0' }}
                  textAreaSize="large"
                  text="The total PPM weight sent via Small Package (all shipments combined). The Counselor may edit this field to reflect the customer's remaining weight entitlement if the combined weight of all shipments exceeds the customer's remaining weight entitlement."
                />
              </div>
              <div>
                <Label>Total Weight Shipped</Label>
                <span>{formatWeight(getTotalPackageWeightSPR(sectionInfo.movingExpenses))}</span>
              </div>
              <div>
                <Label>Pro-gear</Label>
                <span>{hasProGearSPR(sectionInfo.movingExpenses)}</span>
              </div>
              <div>
                <Label>Spouse Pro-gear</Label>
                <span>{hasSpouseProGearSPR(sectionInfo.movingExpenses)}</span>
              </div>
            </>
          )}
        </div>
      );

    case sectionTypes.incentives:
      return (
        <div className={classnames(styles.Details)}>
          <div className={styles.row}>
            <Label className={styles.label}>Government Constructed Cost (GCC)</Label>
            <span data-testid="gcc" className={styles.light}>
              {isFetchingItems && isRecalulatedItem('gcc') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                `$${formatCents(sectionInfo.gcc)}`
              )}
            </span>
          </div>
          <div className={styles.row}>
            <Label className={styles.label}>Gross Incentive</Label>
            <span data-testid="grossIncentive" className={styles.light}>
              {isFetchingItems && isRecalulatedItem('grossIncentive') ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                `$${formatCents(sectionInfo.grossIncentive)}`
              )}
            </span>
          </div>
          <div className={styles.row}>
            <Label className={styles.label}>Advance Requested</Label>
            <span data-testid="advanceRequested" className={styles.light}>
              {aoaRequestedValue}
            </span>
          </div>
          <div className={styles.row}>
            <Label className={styles.label}>Advance Received</Label>
            <span data-testid="advanceReceived" className={styles.light}>
              {isFetchingItems && updatedItemName === 'advanceAmountReceived' ? (
                <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
              ) : (
                <>
                  <OpenModalButton
                    onClick={() => handleEditOnClick(sectionInfo.type, 'advanceAmountReceived')}
                    isDisabled={isFetchingItems || readOnly}
                    ariaLabel="Edit advance amount received"
                  />
                  {aoaValue}
                </>
              )}
            </span>
          </div>
          <div className={styles.row}>
            <Label className={styles.label}>Remaining Incentive</Label>
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
          {sectionInfo.haulPrice ? (
            <div className={styles.row}>
              <Label className={styles.label}>
                {renderHaulType(sectionInfo.haulType)} Price <span>{renderMultiplier(sectionInfo.gccMultiplier)}</span>
              </Label>
              <span data-testid="haulPrice" className={classnames(styles.value, styles.rightAlign)}>
                {isFetchingItems && isRecalulatedItem('haulPrice') ? (
                  <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
                ) : (
                  `$${formatCents(sectionInfo.haulPrice)}`
                )}
              </span>
            </div>
          ) : null}
          {sectionInfo.haulFSC ? (
            <div className={styles.row}>
              <Label className={styles.label}>
                {renderHaulType(sectionInfo.haulType)} Fuel Rate Adjustment{' '}
                <span>{sectionInfo.haulFSC > 0 && renderMultiplier(sectionInfo.gccMultiplier)}</span>
              </Label>
              <span data-testid="haulFSC" className={classnames(styles.value, styles.rightAlign)}>
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
          ) : null}
          {sectionInfo.packPrice ? (
            <div className={styles.row}>
              <Label className={styles.label}>
                Packing Charge <span>{renderMultiplier(sectionInfo.gccMultiplier)}</span>
              </Label>
              <span data-testid="packPrice" className={classnames(styles.value, styles.rightAlign)}>
                {isFetchingItems && isRecalulatedItem('packPrice') ? (
                  <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
                ) : (
                  `$${formatCents(sectionInfo.packPrice)}`
                )}
              </span>
            </div>
          ) : null}
          {sectionInfo.unpackPrice ? (
            <div className={styles.row}>
              <Label className={styles.label}>
                Unpacking Charge <span>{renderMultiplier(sectionInfo.gccMultiplier)}</span>
              </Label>
              <span data-testid="unpackPrice" className={classnames(styles.value, styles.rightAlign)}>
                {isFetchingItems && isRecalulatedItem('unpackPrice') ? (
                  <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
                ) : (
                  `$${formatCents(sectionInfo.unpackPrice)}`
                )}
              </span>
            </div>
          ) : null}
          {sectionInfo.dop ? (
            <div className={styles.row}>
              <Label className={styles.label}>
                Origin Price <span>{renderMultiplier(sectionInfo.gccMultiplier)}</span>
              </Label>
              <span data-testid="originPrice" className={classnames(styles.value, styles.rightAlign)}>
                {isFetchingItems && isRecalulatedItem('dop') ? (
                  <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
                ) : (
                  `$${formatCents(sectionInfo.dop)}`
                )}
              </span>
            </div>
          ) : null}
          {sectionInfo.ddp ? (
            <div className={styles.row}>
              <Label className={styles.label}>
                Destination Price <span>{renderMultiplier(sectionInfo.gccMultiplier)}</span>
              </Label>
              <span data-testid="destinationPrice" className={classnames(styles.value, styles.rightAlign)}>
                {isFetchingItems && isRecalulatedItem('ddp') ? (
                  <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
                ) : (
                  `$${formatCents(sectionInfo.ddp)}`
                )}
              </span>
            </div>
          ) : null}
          {sectionInfo.intlPackPrice ? (
            <div className={styles.row}>
              <Label className={styles.label}>
                International Packing <span>{renderMultiplier(sectionInfo.gccMultiplier)}</span>
              </Label>
              <span data-testid="intlPackPrice" className={classnames(styles.value, styles.rightAlign)}>
                {isFetchingItems && isRecalulatedItem('intlPackPrice') ? (
                  <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
                ) : (
                  `$${formatCents(sectionInfo.intlPackPrice)}`
                )}
              </span>
            </div>
          ) : null}
          {sectionInfo.intlUnpackPrice ? (
            <div className={styles.row}>
              <Label className={styles.label}>
                International Unpacking <span>{renderMultiplier(sectionInfo.gccMultiplier)}</span>
              </Label>
              <span data-testid="intlUnpackPrice" className={classnames(styles.value, styles.rightAlign)}>
                {isFetchingItems && isRecalulatedItem('intlUnpackPrice') ? (
                  <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
                ) : (
                  `$${formatCents(sectionInfo.intlUnpackPrice)}`
                )}
              </span>
            </div>
          ) : null}
          {sectionInfo.intlLinehaulPrice ? (
            <div className={styles.row}>
              <Label className={styles.label}>
                International Shipping & Linehaul <span>{renderMultiplier(sectionInfo.gccMultiplier)}</span>
              </Label>
              <span data-testid="intlLinehaulPrice" className={classnames(styles.value, styles.rightAlign)}>
                {isFetchingItems && isRecalulatedItem('intlLinehaulPrice') ? (
                  <FontAwesomeIcon icon="spinner" spin pulse size="1x" />
                ) : (
                  `$${formatCents(sectionInfo.intlLinehaulPrice)}`
                )}
              </span>
            </div>
          ) : null}
          {sectionInfo.sitReimbursement ? (
            <div className={styles.row}>
              <Label className={styles.label}>SIT Reimbursement</Label>
              <span data-testid="sitReimbursement" className={classnames(styles.value, styles.rightAlign)}>
                ${formatCents(sectionInfo.sitReimbursement)}
              </span>
            </div>
          ) : null}
          <div className={styles.row}>
            <Label className={styles.label}>TOTAL</Label>
            <span data-testid="total" className={classnames(styles.value, styles.rightAlign)}>
              ${calculateTotal(sectionInfo)}
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
  expanded,
}) {
  const requestDetailsButtonTestId = `${sectionInfo.type}-showRequestDetailsButton`;
  const { shipmentId, moveCode } = useParams();
  const { mtoShipment, refetchMTOShipment, isFetching: isFetchingMtoShipment } = usePPMShipmentDocsQueries(shipmentId);
  const queryClient = useQueryClient();
  const [showDetails, setShowDetails] = useState(expanded);
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
      case 'allowableWeight':
        body = { allowableWeight: Number(values?.allowableWeight) };
        break;
      case 'pickupAddress':
        body = {
          pickupAddress: values.pickupAddress,
        };
        break;
      case 'destinationAddress':
        body = {
          destinationAddress: values.destinationAddress,
        };
        break;
      case 'expenseType':
        body = {
          ppmType: values.ppmType,
          isActualExpenseReimbursement: values.ppmType === PPM_TYPES.ACTUAL_EXPENSE,
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
          grade={grade}
        />
      )}
    </section>
  );
}

HeaderSection.propTypes = {
  sectionInfo: propTypes.object.isRequired,
  dataTestId: propTypes.string.isRequired,
  updatedItemName: propTypes.string.isRequired,
  setUpdatedItemName: propTypes.func.isRequired,
  readOnly: propTypes.bool.isRequired,
  expanded: propTypes.bool,
};

HeaderSection.defaultProps = {
  expanded: false,
};
