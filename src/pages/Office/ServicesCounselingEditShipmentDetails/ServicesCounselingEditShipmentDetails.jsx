import React from 'react';
import PropTypes from 'prop-types';
import { useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import styles from '../ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

import 'styles/office.scss';
import ShipmentForm from 'components/Office/ShipmentForm/ShipmentForm';
import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { useEditShipmentQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { updateMTOShipment } from 'services/ghcApi';
import { ORDERS_TYPE } from 'constants/orders';
import { roleTypes } from 'constants/userRoles';

const ServicesCounselingEditShipmentDetails = ({ onUpdate, isAdvancePage }) => {
  const { moveCode, shipmentId } = useParams();
  const queryClient = useQueryClient();
  const { move, order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);
  const { mutate: mutateMTOShipment } = useMutation(updateMTOShipment, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances, grade } = order;

  const matchingShipment = mtoShipments?.filter((shipment) => shipment.id === shipmentId)[0];
  const weightAllotment = { ...allowances, totalWeightSelf: allowances.authorizedWeight };

  const TACs = {
    HHG: order.tac,
    NTS: order.ntsTac,
  };

  const SACs = {
    HHG: order.sac,
    NTS: order.ntsSac,
  };

  const isRetirementOrSeparation =
    order?.order_type === ORDERS_TYPE.RETIREMENT || order?.order_type === ORDERS_TYPE.SEPARATION;

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <ShipmentForm
                submitHandler={mutateMTOShipment}
                onUpdate={onUpdate}
                isCreatePage={false}
                isForServicesCounseling
                currentResidence={customer.current_address}
                originDutyLocationAddress={order.originDutyLocation?.address}
                newDutyLocationAddress={order.destinationDutyLocation?.address}
                shipmentType={matchingShipment.shipmentType}
                mtoShipment={matchingShipment}
                serviceMember={{ weightAllotment, agency: customer.agency, grade }}
                moveTaskOrderID={move.id}
                mtoShipments={mtoShipments}
                TACs={TACs}
                SACs={SACs}
                userRole={roleTypes.SERVICES_COUNSELOR}
                displayDestinationType={isRetirementOrSeparation}
                isAdvancePage={isAdvancePage}
                move={move}
              />
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

ServicesCounselingEditShipmentDetails.propTypes = {
  onUpdate: PropTypes.func.isRequired,
  isAdvancePage: PropTypes.bool,
};

ServicesCounselingEditShipmentDetails.defaultProps = {
  isAdvancePage: false,
};

export default ServicesCounselingEditShipmentDetails;
