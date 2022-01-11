import React, { useState } from 'react';
import { Link, useParams, useHistory } from 'react-router-dom';
import { queryCache, useMutation } from 'react-query';
import { generatePath } from 'react-router';
import classnames from 'classnames';
import 'styles/office.scss';
import { Alert, Button, Grid, GridContainer } from '@trussworks/react-uswds';

import styles from '../ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

import scMoveDetailsStyles from './ServicesCounselingMoveDetails.module.scss';

import { MOVES } from 'constants/queryKeys';
import { servicesCounselingRoutes } from 'constants/routes';
import AllowancesList from 'components/Office/DefinitionLists/AllowancesList';
import CustomerInfoList from 'components/Office/DefinitionLists/CustomerInfoList';
import OrdersList from 'components/Office/DefinitionLists/OrdersList';
import DetailsPanel from 'components/Office/DetailsPanel/DetailsPanel';
import FinancialReviewModal from 'components/Office/FinancialReviewModal/FinancialReviewModal';
import FinancialReviewButton from 'components/Office/FinancialReviewButton/FinancialReviewButton';
import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { SubmitMoveConfirmationModal } from 'components/Office/SubmitMoveConfirmationModal/SubmitMoveConfirmationModal';
import { useMoveDetailsQueries } from 'hooks/queries';
import { updateMoveStatusServiceCounselingCompleted, updateFinancialFlag } from 'services/ghcApi';
import { MOVE_STATUSES, SHIPMENT_OPTIONS_URL } from 'shared/constants';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import shipmentCardsStyles from 'styles/shipmentCards.module.scss';
import { AlertStateShape } from 'types/alert';
import formattedCustomerName from 'utils/formattedCustomerName';
import { getShipmentTypeLabel } from 'utils/shipmentDisplay';
import ButtonDropdown from 'components/ButtonDropdown/ButtonDropdown';

const ServicesCounselingMoveDetails = ({ infoSavedAlert }) => {
  const { moveCode } = useParams();
  const history = useHistory();
  const [alertMessage, setAlertMessage] = useState(null);
  const [alertType, setAlertType] = useState('success');
  const [isSubmitModalVisible, setIsSubmitModalVisible] = useState(false);
  const [isFinancialModalVisible, setIsFinancialModalVisible] = useState(false);

  const { order, move, mtoShipments, isLoading, isError } = useMoveDetailsQueries(moveCode);
  const { customer, entitlement: allowances } = order;

  const counselorCanEdit = move.status === MOVE_STATUSES.NEEDS_SERVICE_COUNSELING;

  // ntsr defaults shows preferred delivery date, storage facility address, destination address, flagged items when collapsed
  const showWhenCollapsed = { HHG_OUTOF_NTS_DOMESTIC: ['counselorRemarks'] }; // add any additional fields that we also want to always show
  const warnIfMissing = {
    HHG_OUTOF_NTS_DOMESTIC: ['primeActualWeight', 'serviceOrderNumber', 'counselorRemarks', 'tacType', 'sacType'],
  };
  const errorIfMissing = { HHG_OUTOF_NTS_DOMESTIC: ['storageFacility'] };

  let shipmentsInfo = [];
  let disableSubmit = false;

  if (mtoShipments) {
    const submittedShipments = mtoShipments?.filter((shipment) => !shipment.deletedAt);

    shipmentsInfo = submittedShipments.map((shipment) => {
      const editURL = counselorCanEdit
        ? generatePath(servicesCounselingRoutes.SHIPMENT_EDIT_PATH, {
            moveCode,
            shipmentId: shipment.id,
          })
        : '';

      const displayInfo = {
        heading: getShipmentTypeLabel(shipment.shipmentType),
        destinationAddress: shipment.destinationAddress || {
          postalCode: order.destinationDutyStation.address.postalCode,
        },
        ...shipment,
      };

      if (!disableSubmit && errorIfMissing[shipment.shipmentType]) {
        for (let i = 0; i < errorIfMissing[shipment.shipmentType].length; i += 1) {
          if (!displayInfo[errorIfMissing[shipment.shipmentType][i]]) {
            disableSubmit = true;
          }
        }
      }

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
    phone: `+1 ${customer.phone}`,
    email: customer.email,
    currentAddress: customer.current_address,
    backupContact: customer.backup_contact,
  };

  const allowancesInfo = {
    branch: customer.agency,
    rank: order.grade,
    weightAllowance: allowances.totalWeight,
    authorizedWeight: allowances.authorizedWeight,
    progear: allowances.proGearWeight,
    spouseProgear: allowances.proGearWeightSpouse,
    storageInTransit: allowances.storageInTransit,
    dependents: allowances.dependentsAuthorized,
    requiredMedicalEquipmentWeight: allowances.requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment: allowances.organizationalClothingAndIndividualEquipment,
  };

  const ordersInfo = {
    currentDutyStation: order.originDutyStation,
    newDutyStation: order.destinationDutyStation,
    issuedDate: order.date_issued,
    reportByDate: order.report_by_date,
    ordersType: order.order_type,
    tacMDC: order.tac,
    sacSDN: order.sac,
    NTStac: order.ntsTac,
    NTSsac: order.ntsSac,
  };
  const ordersLOA = {
    tac: order.tac,
    sac: order.sac,
    ntsTac: order.ntsTac,
    ntsSac: order.ntsSac,
  };

  const handleButtonDropdownChange = (e) => {
    const selectedOption = e.target.value;

    const addShipmentPath = generatePath(servicesCounselingRoutes.SHIPMENT_ADD_PATH, {
      moveCode,
      shipmentType: selectedOption,
    });

    history.push(addShipmentPath);
  };

  // use mutation calls
  const [mutateMoveStatus] = useMutation(updateMoveStatusServiceCounselingCompleted, {
    onSuccess: (data) => {
      queryCache.setQueryData([MOVES, data.locator], data);
      queryCache.invalidateQueries([MOVES, data.locator]);
      setAlertMessage('Move submitted.');
      setAlertType('success');
    },
    onError: () => {
      setAlertMessage('There was a problem submitting the move. Please try again later.');
      setAlertType('error');
    },
  });

  const [mutateFinancialReview] = useMutation(updateFinancialFlag, {
    onSuccess: (data) => {
      queryCache.setQueryData([MOVES, data.locator], data);
      queryCache.invalidateQueries([MOVES, data.locator]);

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
          <Grid row className={scMoveDetailsStyles.pageHeader}>
            {alertMessage && (
              <Grid col={12} className={scMoveDetailsStyles.alertContainer}>
                <Alert slim type={alertType}>
                  {alertMessage}
                </Alert>
              </Grid>
            )}
            {infoSavedAlert && (
              <Grid col={12} className={scMoveDetailsStyles.alertContainer}>
                <Alert slim type={infoSavedAlert.alertType}>
                  {infoSavedAlert.message}
                </Alert>
              </Grid>
            )}
            <Grid col={6} className={scMoveDetailsStyles.pageTitle}>
              <h1>Move details</h1>
            </Grid>
            <Grid col={6} className={scMoveDetailsStyles.submitMoveDetailsContainer}>
              {counselorCanEdit && (
                <Button
                  disabled={!mtoShipments.length || allShipmentsDeleted || disableSubmit}
                  type="button"
                  onClick={handleShowCancellationModal}
                >
                  Submit move details
                </Button>
              )}
            </Grid>
          </Grid>

          <div className={styles.section} id="shipments">
            <DetailsPanel
              className={scMoveDetailsStyles.noPaddingBottom}
              editButton={
                counselorCanEdit && (
                  <ButtonDropdown onChange={handleButtonDropdownChange}>
                    <option value="">Add a new shipment</option>
                    <option value={SHIPMENT_OPTIONS_URL.HHG}>HHG</option>
                    <option value={SHIPMENT_OPTIONS_URL.NTS}>NTS</option>
                    <option value={SHIPMENT_OPTIONS_URL.NTSrelease}>NTS-release</option>
                  </ButtonDropdown>
                )
              }
              financialReviewOpen={handleShowFinancialReviewModal}
              title="Shipments"
            >
              <div className={scMoveDetailsStyles.scFinancialReviewContainer}>
                <FinancialReviewButton
                  onClick={handleShowFinancialReviewModal}
                  reviewRequested={move.financialReviewFlag}
                />
              </div>
              <div className={shipmentCardsStyles.shipmentCards}>
                {shipmentsInfo.map((shipment) => (
                  <ShipmentDisplay
                    displayInfo={shipment.displayInfo}
                    editURL={shipment.editURL}
                    isSubmitted={false}
                    key={shipment.id}
                    shipmentId={shipment.id}
                    shipmentType={shipment.shipmentType}
                    showIcon={false}
                    ordersLOA={ordersLOA}
                    warnIfMissing={warnIfMissing[shipment.shipmentType]}
                    errorIfMissing={errorIfMissing[shipment.shipmentType]}
                    showWhenCollapsed={showWhenCollapsed[shipment.shipmentType]}
                  />
                ))}
              </div>
            </DetailsPanel>
          </div>

          <div className={styles.section} id="orders">
            <DetailsPanel
              title="Orders"
              editButton={
                counselorCanEdit && (
                  <Link
                    className="usa-button usa-button--secondary"
                    to={generatePath(servicesCounselingRoutes.ORDERS_EDIT_PATH, { moveCode })}
                  >
                    View and edit orders
                  </Link>
                )
              }
            >
              <OrdersList ordersInfo={ordersInfo} showMissingWarnings={false} />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="allowances">
            <DetailsPanel
              title="Allowances"
              editButton={
                counselorCanEdit && (
                  <Link
                    className="usa-button usa-button--secondary"
                    to={generatePath(servicesCounselingRoutes.ALLOWANCES_EDIT_PATH, { moveCode })}
                  >
                    Edit allowances
                  </Link>
                )
              }
            >
              <AllowancesList info={allowancesInfo} showVisualCues />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="customer-info">
            <DetailsPanel
              title="Customer info"
              editButton={
                counselorCanEdit && (
                  <Link
                    className="usa-button usa-button--secondary"
                    data-testid="edit-customer-info"
                    to={generatePath(servicesCounselingRoutes.CUSTOMER_INFO_EDIT_PATH, { moveCode })}
                  >
                    Edit customer info
                  </Link>
                )
              }
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
};

ServicesCounselingMoveDetails.defaultProps = {
  infoSavedAlert: null,
};

export default ServicesCounselingMoveDetails;
