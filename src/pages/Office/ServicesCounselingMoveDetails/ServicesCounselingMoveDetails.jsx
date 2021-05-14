import React, { useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { Alert, Button, Grid, GridContainer } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';
import classnames from 'classnames';

import DetailsPanel from '../../../components/Office/DetailsPanel/DetailsPanel';
import ServicesCounselingOrdersList from '../../../components/Office/DefinitionLists/ServicesCounselingOrdersList';
import AllowancesList from '../../../components/Office/DefinitionLists/AllowancesList';
import CustomerInfoList from '../../../components/Office/DefinitionLists/CustomerInfoList';
import { SubmitMoveConfirmationModal } from '../../../components/Office/SubmitMoveConfirmationModal/SubmitMoveConfirmationModal';
import styles from '../ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

import scMoveDetailsStyles from './ServicesCounselingMoveDetails.module.scss';

import 'styles/office.scss';
import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { updateMoveStatusServiceCounselingCompleted } from 'services/ghcApi';
import { useMoveDetailsQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { MOVES } from 'constants/queryKeys';
import { MOVE_STATUSES, SHIPMENT_OPTIONS } from 'shared/constants';
import shipmentCardsStyles from 'styles/shipmentCards.module.scss';

const ServicesCounselingMoveDetails = () => {
  const { moveCode } = useParams();
  const [alertMessage, setAlertMessage] = useState(null);
  const [alertType, setAlertType] = useState('success');
  const [isSubmitModalVisible, setIsSubmitModalVisible] = useState(false);

  const { order, move, mtoShipments, isLoading, isError } = useMoveDetailsQueries(moveCode);
  const { customer, entitlement: allowances } = order;

  let shipmentsInfo = [];

  if (mtoShipments) {
    shipmentsInfo = mtoShipments.map((shipment) => {
      return {
        id: shipment.id,
        heading: SHIPMENT_OPTIONS.HHG,
        requestedMoveDate: shipment.requestedPickupDate,
        currentAddress: shipment.pickupAddress,
        destinationAddress: shipment.destinationAddress,
        counselorRemarks: shipment.counselorRemarks,
      };
    });
  }

  const customerInfo = {
    name: `${customer.last_name}, ${customer.first_name}`,
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

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const handleShowCancellationModal = () => {
    setIsSubmitModalVisible(true);
  };

  const handleConfirmSubmitMoveDetails = () => {
    mutateMoveStatus({ moveTaskOrderID: move.id, ifMatchETag: move.eTag });
    setIsSubmitModalVisible(false);
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        {isSubmitModalVisible && (
          <SubmitMoveConfirmationModal onClose={setIsSubmitModalVisible} onSubmit={handleConfirmSubmitMoveDetails} />
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
            <Grid col={6} className={scMoveDetailsStyles.pageTitle}>
              <h1>Move details</h1>
            </Grid>
            <Grid col={6} className={scMoveDetailsStyles.submitMoveDetailsContainer}>
              {move.status === MOVE_STATUSES.NEEDS_SERVICE_COUNSELING && (
                <Button type="button" onClick={handleShowCancellationModal}>
                  Submit move details
                </Button>
              )}
            </Grid>
          </Grid>

          <div className={styles.section} id="shipments">
            <DetailsPanel title="Shipments">
              <div className={shipmentCardsStyles.shipmentCards}>
                {shipmentsInfo.map((shipment) => (
                  <ShipmentDisplay
                    displayInfo={shipment}
                    isSubmitted={false}
                    key={shipment.id}
                    shipmentId={shipment.id}
                    shipmentType={SHIPMENT_OPTIONS.HHG}
                    showIcon={false}
                  />
                ))}
              </div>
            </DetailsPanel>
          </div>

          <div className={styles.section} id="orders">
            <DetailsPanel
              title="Orders"
              editButton={
                move.status === MOVE_STATUSES.NEEDS_SERVICE_COUNSELING && (
                  <Link className="usa-button usa-button--secondary" to="orders">
                    View and edit orders
                  </Link>
                )
              }
            >
              <ServicesCounselingOrdersList ordersInfo={ordersInfo} />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="allowances">
            <DetailsPanel
              title="Allowances"
              editButton={
                move.status === MOVE_STATUSES.NEEDS_SERVICE_COUNSELING && (
                  <Link className="usa-button usa-button--secondary" to="allowances">
                    Edit allowances
                  </Link>
                )
              }
            >
              <AllowancesList info={allowancesInfo} />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="customer-info">
            <DetailsPanel
              title="Customer info"
              editButton={
                move.status === MOVE_STATUSES.NEEDS_SERVICE_COUNSELING && (
                  <Link className="usa-button usa-button--secondary" to="#">
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

ServicesCounselingMoveDetails.propTypes = {};

export default ServicesCounselingMoveDetails;
