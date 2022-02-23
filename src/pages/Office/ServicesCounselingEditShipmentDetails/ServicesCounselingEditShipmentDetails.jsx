import React from 'react';
import PropTypes from 'prop-types';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';

import styles from '../ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

import 'styles/office.scss';
import ShipmentForm from 'components/Office/ShipmentForm/ShipmentForm';
import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { MatchShape } from 'types/officeShapes';
import { useEditShipmentQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { updateMTOShipment } from 'services/ghcApi';
import { servicesCounselingRoutes } from 'constants/routes';
import { roleTypes } from 'constants/userRoles';

const ServicesCounselingEditShipmentDetails = ({ match, onUpdate }) => {
  const { moveCode, shipmentId } = useParams();
  const history = useHistory();
  const { move, order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);
  const [mutateMTOShipment] = useMutation(updateMTOShipment, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryCache.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryCache.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);
      history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
      onUpdate('success');
    },
    onError: () => {
      history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
      onUpdate('error');
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

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <ShipmentForm
                match={match}
                history={history}
                submitHandler={mutateMTOShipment}
                isCreatePage={false}
                isForServicesCounseling
                currentResidence={customer.current_address}
                newDutyStationAddress={order.destinationDutyLocation?.address}
                selectedMoveType={matchingShipment.shipmentType}
                mtoShipment={matchingShipment}
                serviceMember={{ weightAllotment }}
                moveTaskOrderID={move.id}
                mtoShipments={mtoShipments}
                TACs={TACs}
                SACs={SACs}
                userRole={roleTypes.SERVICES_COUNSELOR}
                orderType={order.order_type}
              />
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

ServicesCounselingEditShipmentDetails.propTypes = {
  match: MatchShape.isRequired,
  onUpdate: PropTypes.func.isRequired,
};

export default ServicesCounselingEditShipmentDetails;
