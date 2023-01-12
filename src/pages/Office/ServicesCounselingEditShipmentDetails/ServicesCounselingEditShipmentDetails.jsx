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
import { updateMTOShipment, updateMoveCloseoutOffice } from 'services/ghcApi';
import { servicesCounselingRoutes } from 'constants/routes';
import { roleTypes } from 'constants/userRoles';

// updateMTOShipmentWrapper allows us to pass in the closeout office and include it
// with the results from updating the shipment, which allows us to chain on the closeout office
// update.
function updateMTOShipmentWrapper({ shipment, closeoutOffice }) {
  return updateMTOShipment(shipment).then((newShipment) => {
    return { newShipment, closeoutOffice };
  });
}
const ServicesCounselingEditShipmentDetails = ({ match, onUpdate, isAdvancePage }) => {
  const { moveCode, shipmentId } = useParams();
  const history = useHistory();
  const { move, order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);
  const [mutateMoveCloseoutOffice] = useMutation(updateMoveCloseoutOffice, {
    onSuccess: () => {
      onUpdate('success');
    },
    onError: () => {
      history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
      onUpdate('error');
    },
  });
  const [mutateMTOShipment] = useMutation(updateMTOShipmentWrapper, {
    onSuccess: (result) => {
      // if we have a closeout office, we must be on the first page of creating a PPM shipment,
      // so we should update the closeout office and redirect to the advance page
      if (result.closeoutOffice) {
        mutateMoveCloseoutOffice({
          locator: moveCode,
          ifMatchETag: move.eTag,
          body: { closeoutOfficeId: result.closeoutOffice.id },
        }).then(() => {
          mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === result.newShipment.id)] =
            result.newShipment;
          queryCache.setQueryData([MTO_SHIPMENTS, result.newShipment.moveTaskOrderID, false], mtoShipments);
          queryCache.invalidateQueries([MTO_SHIPMENTS, result.newShipment.moveTaskOrderID]);
          onUpdate('success');
        });
      } else {
        // if we don't have a closeout office, we're either on the advance page for a PPM, or the first
        // page for another type of shipment. In either case, we're done now and can head back to the move view
        mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === result.newShipment.id)] = result.newShipment;
        queryCache.setQueryData([MTO_SHIPMENTS, result.newShipment.moveTaskOrderID, false], mtoShipments);
        queryCache.invalidateQueries([MTO_SHIPMENTS, result.newShipment.moveTaskOrderID]);
        history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
        onUpdate('success');
      }
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
                originDutyLocationAddress={order.originDutyLocation?.address}
                newDutyLocationAddress={order.destinationDutyLocation?.address}
                shipmentType={matchingShipment.shipmentType}
                mtoShipment={matchingShipment}
                serviceMember={{ weightAllotment, agency: customer.agency }}
                moveTaskOrderID={move.id}
                mtoShipments={mtoShipments}
                TACs={TACs}
                SACs={SACs}
                userRole={roleTypes.SERVICES_COUNSELOR}
                displayDestinationType
                isAdvancePage={isAdvancePage}
                closeoutOffice={move.closeoutOffice}
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
  match: MatchShape.isRequired,
  onUpdate: PropTypes.func.isRequired,
  isAdvancePage: PropTypes.bool,
};

ServicesCounselingEditShipmentDetails.defaultProps = {
  isAdvancePage: false,
};

export default ServicesCounselingEditShipmentDetails;
