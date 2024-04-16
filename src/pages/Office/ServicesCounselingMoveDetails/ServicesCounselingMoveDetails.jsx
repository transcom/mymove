import React, { useState, useEffect, useMemo } from 'react';
import { Link, useParams, useNavigate, generatePath } from 'react-router-dom';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { func } from 'prop-types';
import classnames from 'classnames';
import 'styles/office.scss';
import { Alert, Button, Grid, GridContainer } from '@trussworks/react-uswds';

import styles from '../ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

import scMoveDetailsStyles from './ServicesCounselingMoveDetails.module.scss';

import { MOVES } from 'constants/queryKeys';
import { ORDERS_TYPE } from 'constants/orders';
import { servicesCounselingRoutes } from 'constants/routes';
import AllowancesList from 'components/Office/DefinitionLists/AllowancesList';
import CustomerInfoList from 'components/Office/DefinitionLists/CustomerInfoList';
import OrdersList from 'components/Office/DefinitionLists/OrdersList';
import DetailsPanel from 'components/Office/DetailsPanel/DetailsPanel';
import FinancialReviewButton from 'components/Office/FinancialReviewButton/FinancialReviewButton';
import FinancialReviewModal from 'components/Office/FinancialReviewModal/FinancialReviewModal';
import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { SubmitMoveConfirmationModal } from 'components/Office/SubmitMoveConfirmationModal/SubmitMoveConfirmationModal';
import { useMoveDetailsQueries } from 'hooks/queries';
import { updateMoveStatusServiceCounselingCompleted, updateFinancialFlag } from 'services/ghcApi';
import { MOVE_STATUSES, SHIPMENT_OPTIONS_URL, SHIPMENT_OPTIONS } from 'shared/constants';
import { ppmShipmentStatuses } from 'constants/shipments';
import shipmentCardsStyles from 'styles/shipmentCards.module.scss';
import LeftNav from 'components/LeftNav/LeftNav';
import LeftNavTag from 'components/LeftNavTag/LeftNavTag';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { AlertStateShape } from 'types/alert';
import formattedCustomerName from 'utils/formattedCustomerName';
import { getShipmentTypeLabel } from 'utils/shipmentDisplay';
import ButtonDropdown from 'components/ButtonDropdown/ButtonDropdown';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { objectIsMissingFieldWithCondition } from 'utils/displayFlags';
import { ReviewButton } from 'components/form/IconButtons';
import { calculateWeightRequested } from 'hooks/custom';

const ServicesCounselingMoveDetails = ({ infoSavedAlert, setUnapprovedShipmentCount }) => {
  const { moveCode } = useParams();
  const navigate = useNavigate();
  const [alertMessage, setAlertMessage] = useState(null);
  const [alertType, setAlertType] = useState('success');
  const [moveHasExcessWeight, setMoveHasExcessWeight] = useState(false);
  const [isSubmitModalVisible, setIsSubmitModalVisible] = useState(false);
  const [isFinancialModalVisible, setIsFinancialModalVisible] = useState(false);
  const [shipmentConcernCount, setShipmentConcernCount] = useState(0);

  const { order, customerData, move, closeoutOffice, mtoShipments, isLoading, isError } =
    useMoveDetailsQueries(moveCode);
  const { customer, entitlement: allowances } = order;

  const moveWeightTotal = calculateWeightRequested(mtoShipments);

  let counselorCanReview;
  let reviewWeightsURL;
  let counselorCanEdit;
  let counselorCanEditNonPPM;

  const sections = useMemo(() => {
    return ['shipments', 'orders', 'allowances', 'customer-info'];
  }, []);

  // nts defaults show preferred pickup date and pickup address, flagged items when collapsed
  // ntsr defaults shows preferred delivery date, storage facility address, destination address, flagged items when collapsed
  const showWhenCollapsed = {
    HHG_INTO_NTS_DOMESTIC: ['counselorRemarks'],
    HHG_OUTOF_NTS_DOMESTIC: ['counselorRemarks'],
  }; // add any additional fields that we also want to always show
  const neverShow = { HHG_INTO_NTS_DOMESTIC: ['usesExternalVendor', 'serviceOrderNumber', 'storageFacility'] };
  const warnIfMissing = {
    HHG: [{ fieldName: 'counselorRemarks' }],
    HHG_INTO_NTS_DOMESTIC: [{ fieldName: 'counselorRemarks' }, { fieldName: 'tacType' }, { fieldName: 'sacType' }],
    HHG_OUTOF_NTS_DOMESTIC: [
      { fieldName: 'ntsRecordedWeight' },
      { fieldName: 'serviceOrderNumber' },
      { fieldName: 'counselorRemarks' },
      { fieldName: 'tacType' },
      { fieldName: 'sacType' },
    ],
    PPM: [{ fieldName: 'counselorRemarks' }],
  };
  const errorIfMissing = {
    HHG_OUTOF_NTS_DOMESTIC: [{ fieldName: 'storageFacility' }],
    PPM: [
      {
        fieldName: 'advanceStatus',
        condition: (shipment) => shipment?.ppmShipment?.hasRequestedAdvance === true,
      },
    ],
  };

  let shipmentsInfo = [];
  let ppmShipmentsInfoNeedsApproval = [];
  let ppmShipmentsOtherStatuses = [];
  let disableSubmit = false;
  let disableSubmitDueToMissingOrderInfo = false;
  let numberOfErrorIfMissingForAllShipments = 0;
  let numberOfWarnIfMissingForAllShipments = 0;

  const [hasInvalidProGearAllowances, setHasInvalidProGearAllowances] = useState(false);

  // check if invalid progear weight allowances
  const checkProGearAllowances = () => {
    mtoShipments?.forEach((mto) => {
      if (!order.entitlement.proGearWeight) order.entitlement.proGearWeight = 0;
      if (!order.entitlement.proGearWeightSpouse) order.entitlement.proGearWeightSpouse = 0;

      if (
        mto?.ppmShipment?.proGearWeight > order.entitlement.proGearWeight ||
        mto?.ppmShipment?.spouseProGearWeight > order.entitlement.proGearWeightSpouse
      ) {
        setHasInvalidProGearAllowances(true);
      }
    });
  };

  useEffect(() => {
    checkProGearAllowances();
  });

  // for now we are only showing dest type on retiree and separatee orders
  const isRetirementOrSeparation =
    order.order_type === ORDERS_TYPE.RETIREMENT || order.order_type === ORDERS_TYPE.SEPARATION;

  if (isRetirementOrSeparation) {
    // destination type must be set for for HHG, NTSR shipments only
    errorIfMissing.HHG = [{ fieldName: 'destinationType' }];
    errorIfMissing.HHG_OUTOF_NTS_DOMESTIC.push({ fieldName: 'destinationType' });
  }

  if (!order.department_indicator || !order.order_number || !order.order_type_detail || !order.tac)
    disableSubmitDueToMissingOrderInfo = true;

  if (mtoShipments) {
    const submittedShipments = mtoShipments?.filter((shipment) => !shipment.deletedAt);
    const submittedShipmentsNonPPM = submittedShipments.filter(
      (shipment) => shipment.ppmShipment?.status !== ppmShipmentStatuses.NEEDS_PAYMENT_APPROVAL,
    );
    const ppmNeedsApprovalShipments = submittedShipments.filter(
      (shipment) => shipment.ppmShipment?.status === ppmShipmentStatuses.NEEDS_PAYMENT_APPROVAL,
    );
    const onlyPpmShipments = submittedShipments.filter((shipment) => shipment.shipmentType === 'PPM');
    ppmShipmentsOtherStatuses = onlyPpmShipments.filter(
      (shipment) => shipment.ppmShipment?.status !== ppmShipmentStatuses.NEEDS_PAYMENT_APPROVAL,
    );

    ppmShipmentsInfoNeedsApproval = ppmNeedsApprovalShipments.map((shipment) => {
      const reviewURL = `../${generatePath(servicesCounselingRoutes.SHIPMENT_REVIEW_PATH, {
        moveCode,
        shipmentId: shipment.id,
      })}`;
      const numberofPPMShipments = ppmNeedsApprovalShipments.length;

      const displayInfo = {
        heading: getShipmentTypeLabel(shipment.shipmentType),
        destinationAddress: shipment.destinationAddress || {
          postalCode: order.destinationDutyLocation.address.postalCode,
        },
        agency: customerData.agency,
        closeoutOffice,
        ...shipment.ppmShipment,
        ...shipment,
        displayDestinationType: isRetirementOrSeparation,
      };

      const errorIfMissingList = errorIfMissing[shipment.shipmentType];
      if (errorIfMissingList) {
        errorIfMissingList.forEach((fieldToCheck) => {
          if (objectIsMissingFieldWithCondition(displayInfo, fieldToCheck)) {
            numberOfErrorIfMissingForAllShipments += 1;
            // Since storage facility gets split into two fields - the name and the address
            // it needs to be counted twice.
            if (fieldToCheck.fieldName === 'storageFacility') {
              numberOfErrorIfMissingForAllShipments += 1;
            }
          }
        });
      }

      const warnIfMissingList = warnIfMissing[shipment.shipmentType];
      if (warnIfMissingList) {
        warnIfMissingList.forEach((fieldToCheck) => {
          if (objectIsMissingFieldWithCondition(displayInfo, fieldToCheck)) {
            numberOfWarnIfMissingForAllShipments += 1;
          }
          // Since storage facility gets split into two fields - the name and the address
          // it needs to be counted twice.
          if (fieldToCheck.fieldName === 'storageFacility') {
            numberOfErrorIfMissingForAllShipments += 1;
          }
        });
      }

      disableSubmit = numberOfErrorIfMissingForAllShipments !== 0;

      return {
        id: shipment.id,
        displayInfo,
        reviewURL,
        numberofPPMShipments,
        shipmentType: shipment.shipmentType,
      };
    });

    counselorCanReview = ppmShipmentsInfoNeedsApproval.length > 0;
    reviewWeightsURL = generatePath(servicesCounselingRoutes.BASE_REVIEW_SHIPMENT_WEIGHTS_PATH, { moveCode });
    counselorCanEdit = move.status === MOVE_STATUSES.NEEDS_SERVICE_COUNSELING && ppmShipmentsOtherStatuses.length > 0;
    counselorCanEditNonPPM =
      move.status === MOVE_STATUSES.NEEDS_SERVICE_COUNSELING && shipmentsInfo.shipmentType !== 'PPM';

    shipmentsInfo = submittedShipmentsNonPPM.map((shipment) => {
      const editURL =
        counselorCanEdit || counselorCanEditNonPPM
          ? `../${generatePath(servicesCounselingRoutes.SHIPMENT_EDIT_PATH, {
              shipmentId: shipment.id,
            })}`
          : '';

      const displayInfo = {
        heading: getShipmentTypeLabel(shipment.shipmentType),
        destinationAddress: shipment.destinationAddress || {
          postalCode: order.destinationDutyLocation.address.postalCode,
        },
        ...shipment.ppmShipment,
        ...shipment,
        displayDestinationType: isRetirementOrSeparation,
      };

      if (shipment.shipmentType === SHIPMENT_OPTIONS.PPM) {
        displayInfo.agency = customerData.agency;
        displayInfo.closeoutOffice = closeoutOffice;
      }
      const errorIfMissingList = errorIfMissing[shipment.shipmentType];

      if (errorIfMissingList) {
        errorIfMissingList.forEach((fieldToCheck) => {
          if (objectIsMissingFieldWithCondition(displayInfo, fieldToCheck)) {
            numberOfErrorIfMissingForAllShipments += 1;
            // Since storage facility gets split into two fields - the name and the address
            // it needs to be counted twice.
            if (fieldToCheck.fieldName === 'storageFacility') {
              numberOfErrorIfMissingForAllShipments += 1;
            }
          }
        });
      }

      const warnIfMissingList = warnIfMissing[shipment.shipmentType];
      if (warnIfMissingList) {
        warnIfMissingList.forEach((fieldToCheck) => {
          if (objectIsMissingFieldWithCondition(displayInfo, fieldToCheck)) {
            numberOfWarnIfMissingForAllShipments += 1;
          }
          // Since storage facility gets split into two fields - the name and the address
          // it needs to be counted twice.
          if (fieldToCheck.fieldName === 'storageFacility') {
            numberOfWarnIfMissingForAllShipments += 1;
          }
        });
      }

      disableSubmit = numberOfErrorIfMissingForAllShipments !== 0;

      return {
        id: shipment.id,
        displayInfo,
        editURL,
        shipmentType: shipment.shipmentType,
      };
    });
  }

  const customerInfo = {
    name: formattedCustomerName(customer.last_name, customer.first_name, customer.suffix, customer.middle_name),
    dodId: customer.dodID,
    phone: customer.phone,
    altPhone: customer.secondaryTelephone,
    email: customer.email,
    currentAddress: customer.current_address,
    backupAddress: customerData.backupAddress,
    backupContact: customer.backup_contact,
  };

  const allowancesInfo = {
    branch: customer.agency,
    grade: order.grade,
    authorizedWeight: allowances.authorizedWeight,
    progear: allowances.proGearWeight,
    spouseProgear: allowances.proGearWeightSpouse,
    storageInTransit: allowances.storageInTransit,
    dependents: allowances.dependentsAuthorized,
    requiredMedicalEquipmentWeight: allowances.requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment: allowances.organizationalClothingAndIndividualEquipment,
    gunSafe: allowances.gunSafe,
  };

  const ordersInfo = {
    currentDutyLocation: order.originDutyLocation,
    newDutyLocation: order.destinationDutyLocation,
    departmentIndicator: order.department_indicator,
    issuedDate: order.date_issued,
    reportByDate: order.report_by_date,
    ordersType: order.order_type,
    ordersNumber: order.order_number,
    ordersTypeDetail: order.order_type_detail,
    tacMDC: order.tac,
    sacSDN: order.sac,
    NTStac: order.ntsTac,
    NTSsac: order.ntsSac,
    payGrade: order.grade,
  };
  const ordersLOA = {
    tac: order.tac,
    sac: order.sac,
    ntsTac: order.ntsTac,
    ntsSac: order.ntsSac,
  };

  const handleButtonDropdownChange = (e) => {
    const selectedOption = e.target.value;

    const addShipmentPath = `../${generatePath(servicesCounselingRoutes.SHIPMENT_ADD_PATH, {
      shipmentType: selectedOption,
    })}`;

    navigate(addShipmentPath);
  };

  const handleReviewWeightsButton = (weightsURL) => {
    navigate(weightsURL);
  };

  // use mutation calls
  const queryClient = useQueryClient();
  const { mutate: mutateMoveStatus } = useMutation(updateMoveStatusServiceCounselingCompleted, {
    onSuccess: (data) => {
      queryClient.setQueryData([MOVES, data.locator], data);
      queryClient.invalidateQueries([MOVES, data.locator]);
      setAlertMessage('Move submitted.');
      setAlertType('success');
    },
    onError: () => {
      setAlertMessage('There was a problem submitting the move. Please try again later.');
      setAlertType('error');
    },
  });

  const { mutate: mutateFinancialReview } = useMutation(updateFinancialFlag, {
    onSuccess: (data) => {
      queryClient.setQueryData([MOVES, data.locator], data);
      queryClient.invalidateQueries([MOVES, data.locator]);

      if (data.financialReviewFlag) {
        setAlertMessage('Move flagged for financial review.');
        setAlertType('success');
      } else {
        setAlertMessage('Move unflagged for financial review.');
        setAlertType('success');
      }
    },
    onError: () => {
      setAlertMessage('There was a problem flagging the move for financial review. Please try again later.');
      setAlertType('error');
    },
  });

  useEffect(() => {
    setMoveHasExcessWeight(moveWeightTotal > order.entitlement.totalWeight);
  }, [moveWeightTotal, order.entitlement.totalWeight]);

  // Keep unapproved shipment count in sync
  useEffect(() => {
    let shipmentConcerns = numberOfErrorIfMissingForAllShipments + numberOfWarnIfMissingForAllShipments;
    if (moveHasExcessWeight) {
      shipmentConcerns += 1;
    }
    setShipmentConcernCount(shipmentConcerns);
    setUnapprovedShipmentCount(shipmentConcerns);
  }, [
    moveHasExcessWeight,
    numberOfErrorIfMissingForAllShipments,
    numberOfWarnIfMissingForAllShipments,
    setUnapprovedShipmentCount,
  ]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const handleShowCancellationModal = () => {
    setIsSubmitModalVisible(true);
  };

  const handleConfirmSubmitMoveDetails = () => {
    mutateMoveStatus({ moveTaskOrderID: move.id, ifMatchETag: move.eTag });
    setIsSubmitModalVisible(false);
  };

  const handleShowFinancialReviewModal = () => {
    setIsFinancialModalVisible(true);
  };

  const handleSubmitFinancialReviewModal = (remarks, flagForReview) => {
    const flagForReviewBool = flagForReview === 'yes';
    mutateFinancialReview({
      moveID: move.id,
      ifMatchETag: move.eTag,
      body: { remarks, flagForReview: flagForReviewBool },
    });
    setIsFinancialModalVisible(false);
  };

  const handleCancelFinancialReviewModal = () => {
    setIsFinancialModalVisible(false);
  };

  const allShipmentsDeleted = mtoShipments.every((shipment) => !!shipment.deletedAt);
  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <LeftNav sections={sections}>
          <LeftNavTag
            associatedSectionName="shipments"
            showTag={shipmentConcernCount !== 0}
            testID="requestedShipmentsTag"
          >
            {shipmentConcernCount}
          </LeftNavTag>
        </LeftNav>
        {isSubmitModalVisible && (
          <SubmitMoveConfirmationModal onClose={setIsSubmitModalVisible} onSubmit={handleConfirmSubmitMoveDetails} />
        )}
        {isFinancialModalVisible && (
          <FinancialReviewModal
            onClose={handleCancelFinancialReviewModal}
            onSubmit={handleSubmitFinancialReviewModal}
            initialRemarks={move?.financialReviewRemarks}
            initialSelection={move?.financialReviewFlag}
          />
        )}
        <GridContainer className={classnames(styles.gridContainer, scMoveDetailsStyles.ServicesCounselingMoveDetails)}>
          <NotificationScrollToTop dependency={alertMessage || infoSavedAlert} />
          <Grid row className={scMoveDetailsStyles.pageHeader}>
            {alertMessage && (
              <Grid col={12} className={scMoveDetailsStyles.alertContainer}>
                <Alert headingLevel="h4" slim type={alertType}>
                  {alertMessage}
                </Alert>
              </Grid>
            )}
            {infoSavedAlert && (
              <Grid col={12} className={scMoveDetailsStyles.alertContainer}>
                <Alert headingLevel="h4" slim type={infoSavedAlert.alertType}>
                  {infoSavedAlert.message}
                </Alert>
              </Grid>
            )}
            {moveHasExcessWeight && (
              <Grid col={12} className={scMoveDetailsStyles.alertContainer}>
                <Alert headingLevel="h4" slim type="warning">
                  <span>This move has excess weight. Review PPM weight ticket documents to resolve.</span>
                </Alert>
              </Grid>
            )}
            <Grid col={12} className={scMoveDetailsStyles.pageTitle}>
              <h1>Move details</h1>
              {ppmShipmentsInfoNeedsApproval.length > 0 ? null : (
                <div>
                  {(counselorCanEdit || counselorCanEditNonPPM) && (
                    <Button
                      disabled={
                        !mtoShipments.length ||
                        allShipmentsDeleted ||
                        disableSubmit ||
                        disableSubmitDueToMissingOrderInfo ||
                        hasInvalidProGearAllowances
                      }
                      type="button"
                      onClick={handleShowCancellationModal}
                    >
                      Submit move details
                    </Button>
                  )}
                </div>
              )}
            </Grid>
          </Grid>

          {hasInvalidProGearAllowances ? (
            <div className={scMoveDetailsStyles.allowanceErrorStyle} data-testid="allowanceError">
              Pro Gear weight allowances are less than the weights entered in move.
            </div>
          ) : null}

          <div className={styles.section} id="shipments">
            <DetailsPanel
              className={scMoveDetailsStyles.noPaddingBottom}
              editButton={
                (counselorCanEdit || counselorCanEditNonPPM) && (
                  <ButtonDropdown data-testid="addShipmentButton" onChange={handleButtonDropdownChange}>
                    <option value="">Add a new shipment</option>
                    <option data-testid="hhgOption" value={SHIPMENT_OPTIONS_URL.HHG}>
                      HHG
                    </option>
                    <option value={SHIPMENT_OPTIONS_URL.PPM}>PPM</option>
                    <option value={SHIPMENT_OPTIONS_URL.NTS}>NTS</option>
                    <option value={SHIPMENT_OPTIONS_URL.NTSrelease}>NTS-release</option>
                  </ButtonDropdown>
                )
              }
              reviewButton={
                counselorCanReview && (
                  <ReviewButton
                    onClick={() => handleReviewWeightsButton(reviewWeightsURL)}
                    data-testid={reviewWeightsURL}
                    label="Review shipment weights"
                    secondary
                  />
                )
              }
              financialReviewOpen={handleShowFinancialReviewModal}
              title="Shipments"
              ppmShipmentInfoNeedsApproval={ppmShipmentsInfoNeedsApproval}
            >
              <Restricted to={permissionTypes.updateFinancialReviewFlag}>
                <div className={scMoveDetailsStyles.scFinancialReviewContainer}>
                  <FinancialReviewButton
                    onClick={handleShowFinancialReviewModal}
                    reviewRequested={move.financialReviewFlag}
                  />
                </div>
              </Restricted>
              <div className={shipmentCardsStyles.shipmentCards}>
                {shipmentsInfo.map((shipment) => (
                  <ShipmentDisplay
                    displayInfo={shipment.displayInfo}
                    editURL={shipment.editURL}
                    isSubmitted={false}
                    key={shipment.id}
                    shipmentId={shipment.id}
                    shipmentType={shipment.shipmentType}
                    allowApproval={false}
                    ordersLOA={ordersLOA}
                    warnIfMissing={warnIfMissing[shipment.shipmentType]}
                    errorIfMissing={errorIfMissing[shipment.shipmentType]}
                    showWhenCollapsed={showWhenCollapsed[shipment.shipmentType]}
                    neverShow={neverShow[shipment.shipmentType]}
                  />
                ))}
                {ppmShipmentsInfoNeedsApproval.length > 0 &&
                  ppmShipmentsInfoNeedsApproval.map((shipment) => (
                    <ShipmentDisplay
                      numberofPPMShipments={shipment.numberofPPMShipments}
                      displayInfo={shipment.displayInfo}
                      reviewURL={shipment.reviewURL}
                      isSubmitted={false}
                      key={shipment.id}
                      shipmentId={shipment.id}
                      shipmentType={shipment.shipmentType}
                      allowApproval={false}
                      ordersLOA={ordersLOA}
                      warnIfMissing={warnIfMissing[shipment.shipmentType]}
                      errorIfMissing={errorIfMissing[shipment.shipmentType]}
                      showWhenCollapsed={showWhenCollapsed[shipment.shipmentType]}
                      neverShow={neverShow[shipment.shipmentType]}
                    />
                  ))}
              </div>
            </DetailsPanel>
          </div>

          <div className={styles.section} id="orders">
            <DetailsPanel
              title="Orders"
              editButton={
                (counselorCanEdit || counselorCanEditNonPPM) && (
                  <Link
                    className="usa-button usa-button--secondary"
                    to={`../${servicesCounselingRoutes.ORDERS_EDIT_PATH}`}
                  >
                    View and edit orders
                  </Link>
                )
              }
              ppmShipmentInfoNeedsApproval={ppmShipmentsInfoNeedsApproval}
            >
              <OrdersList ordersInfo={ordersInfo} />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="allowances">
            <DetailsPanel
              title="Allowances"
              editButton={
                (counselorCanEdit || counselorCanEditNonPPM) && (
                  <Link
                    className="usa-button usa-button--secondary"
                    data-testid="edit-allowances"
                    to={`../${servicesCounselingRoutes.ALLOWANCES_EDIT_PATH}`}
                  >
                    Edit allowances
                  </Link>
                )
              }
              ppmShipmentInfoNeedsApproval={ppmShipmentsInfoNeedsApproval}
            >
              <AllowancesList info={allowancesInfo} showVisualCues />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="customer-info">
            <DetailsPanel
              title="Customer info"
              editButton={
                <Restricted to={permissionTypes.updateCustomer}>
                  <Link
                    className="usa-button usa-button--secondary"
                    data-testid="edit-customer-info"
                    to={`../${servicesCounselingRoutes.CUSTOMER_INFO_EDIT_PATH}`}
                  >
                    Edit customer info
                  </Link>
                </Restricted>
              }
              ppmShipmentInfoNeedsApproval={ppmShipmentsInfoNeedsApproval}
            >
              <CustomerInfoList customerInfo={customerInfo} />
            </DetailsPanel>
          </div>
        </GridContainer>
      </div>
    </div>
  );
};

ServicesCounselingMoveDetails.propTypes = {
  infoSavedAlert: AlertStateShape,
  setUnapprovedShipmentCount: func.isRequired,
};

ServicesCounselingMoveDetails.defaultProps = {
  infoSavedAlert: null,
};

export default ServicesCounselingMoveDetails;
