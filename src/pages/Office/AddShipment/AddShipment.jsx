import React from 'react';
import { useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import styles from '../ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

import 'styles/office.scss';
import ShipmentForm from 'components/Office/ShipmentForm/ShipmentForm';
import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { useEditShipmentQueries } from 'hooks/queries';
import { createMTOShipment } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { roleTypes } from 'constants/userRoles';
import { SHIPMENT_OPTIONS, SHIPMENT_OPTIONS_URL } from 'shared/constants';
import { ORDERS_TYPE } from 'constants/orders';

const AddShipment = () => {
  const params = useParams();
  let { shipmentType } = params;
  const { moveCode } = params;

  if (shipmentType === SHIPMENT_OPTIONS_URL.NTSrelease) {
    shipmentType = SHIPMENT_OPTIONS.NTSR;
  } else if (shipmentType === SHIPMENT_OPTIONS_URL.BOAT) {
    shipmentType = SHIPMENT_OPTIONS.BOAT;
  } else if (shipmentType === SHIPMENT_OPTIONS_URL.MOBILE_HOME) {
    shipmentType = SHIPMENT_OPTIONS.MOBILE_HOME;
  } else if (shipmentType === SHIPMENT_OPTIONS_URL.UNACCOMPANIED_BAGGAGE) {
    shipmentType = SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE;
  } else {
    shipmentType = SHIPMENT_OPTIONS[shipmentType];
  }

  const { move, order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);
  const isRetirementOrSeparation =
    order?.order_type === ORDERS_TYPE.RETIREMENT || order?.order_type === ORDERS_TYPE.SEPARATION;

  const queryClient = useQueryClient();
  const { mutate: mutateMTOShipments } = useMutation(createMTOShipment, {
    onSuccess: (newMTOShipment) => {
      mtoShipments.push(newMTOShipment);
      queryClient.setQueryData([MTO_SHIPMENTS, newMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries([MTO_SHIPMENTS, newMTOShipment.moveTaskOrderID]);
      return newMTOShipment;
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances } = order;
  const weightAllotment = { ...allowances, totalWeightSelf: allowances.authorizedWeight };

  const TACs = {
    HHG: order.tac,
    NTS: order.ntsTac,
  };

  const SACs = {
    HHG: order.sac,
    NTS: order.ntsSac,
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <ShipmentForm
                submitHandler={mutateMTOShipments}
                isCreatePage
                currentResidence={customer.current_address}
                originDutyLocationAddress={order.originDutyLocation?.address}
                newDutyLocationAddress={order.destinationDutyLocation?.address}
                shipmentType={shipmentType}
                serviceMember={{ weightAllotment, agency: customer.agency }}
                moveTaskOrderID={move.id}
                mtoShipments={mtoShipments}
                TACs={TACs}
                SACs={SACs}
                userRole={roleTypes.TOO}
                displayDestinationType={isRetirementOrSeparation}
                move={move}
              />
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

export default AddShipment;
