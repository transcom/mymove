import React from 'react';
import { useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import PropTypes from 'prop-types';

import styles from './EditShipmentDetails.module.scss';

import 'styles/office.scss';
import CustomerHeader from 'components/CustomerHeader';
import ShipmentForm from 'components/Office/ShipmentForm/ShipmentForm';
import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { useEditShipmentQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { updateMTOShipment } from 'services/ghcApi';
import { ORDERS_TYPE } from 'constants/orders';
import { roleTypes } from 'constants/userRoles';

const EditShipmentDetails = ({ isMultiRole }) => {
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

  const { customer, entitlement: allowances } = order;

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
    <>
      <CustomerHeader move={move} order={order} customer={customer} moveCode={moveCode} isMultiRole={isMultiRole} />
      <div className={styles.tabContent}>
        <div className={styles.container}>
          <GridContainer className={styles.gridContainer}>
            <Grid row>
              <Grid col desktop={{ col: 8, offset: 2 }}>
                <ShipmentForm
                  submitHandler={mutateMTOShipment}
                  isCreatePage={false}
                  currentResidence={customer.current_address}
                  originDutyLocationAddress={order.originDutyLocation?.address}
                  newDutyLocationAddress={order.destinationDutyLocation?.address}
                  shipmentType={matchingShipment.shipmentType}
                  mtoShipment={matchingShipment}
                  serviceMember={{ weightAllotment }}
                  moveTaskOrderID={move.id}
                  mtoShipments={mtoShipments}
                  TACs={TACs}
                  SACs={SACs}
                  userRole={roleTypes.TOO}
                  displayDestinationType={isRetirementOrSeparation}
                />
              </Grid>
            </Grid>
          </GridContainer>
        </div>
      </div>
    </>
  );
};

EditShipmentDetails.propTypes = {
  isMultiRole: PropTypes.bool,
};

EditShipmentDetails.defaultProps = {
  isMultiRole: false,
};

export default EditShipmentDetails;
