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

function foobar({ shipment, closeoutOffice }) {
  console.log('editship submitHandler', { shipment, closeoutOffice });
  return updateMTOShipment(shipment).then((newShipment) => {
    return { newShipment, closeoutOffice };
  });
}
const ServicesCounselingEditShipmentDetails = ({ match, onUpdate, isAdvancePage }) => {
  const { moveCode, shipmentId } = useParams();
  const history = useHistory();
  const { move, order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);
  const [mutateMoveCloseoutOffice] = useMutation(updateMoveCloseoutOffice, {
    onSuccess: (updatedMove) => {
      console.log('ROUTE redirecting to move view path after updating closeout office (SKIPPED)');
      // history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
      onUpdate('success');
    },
    onError: () => {
      history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
      onUpdate('error');
    },
  });
  const [mutateMTOShipment] = useMutation(foobar, {
    onSuccess: (result) => {
      console.log('updatedMTOShipment', result);
      if (result.closeoutOffice) {
        console.log('lets try to submit the closeout office', result.closeoutOffice);
        mutateMoveCloseoutOffice({
          locator: moveCode,
          ifMatchETag: move.eTag,
          body: { closeoutOfficeId: result.closeoutOffice.id },
        }).then(() => {
          console.log('mutate closeout done');
          mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === result.newShipment.id)] =
            result.newShipment;
          queryCache.setQueryData([MTO_SHIPMENTS, result.newShipment.moveTaskOrderID, false], mtoShipments);
          queryCache.invalidateQueries([MTO_SHIPMENTS, result.newShipment.moveTaskOrderID]);
          // console.log('--------------- redirect to move view path 2');
          // history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
          onUpdate('success');
        });
      } else {
        console.log('no closeout office, skipping that update');
        mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === result.newShipment.id)] = result.newShipment;
        queryCache.setQueryData([MTO_SHIPMENTS, result.newShipment.moveTaskOrderID, false], mtoShipments);
        queryCache.invalidateQueries([MTO_SHIPMENTS, result.newShipment.moveTaskOrderID]);
        // and then what do i pass in here?
        // TODO removing this makes it so we don't move on after the advance page. i dont know why
        // TODO as i put the same code in the mutation that is supposed to run AFTER this one
        // TODO ohhhhh because it is not called on the advance page.
        // TODO but then why doesn't this one break stuff by redirecting before we can start the other
        // TODO mutation?
        console.log('--------------- redirect to move view path 1');
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

  // console.log('ServicesCounselingEditShipmentDetails move', move);
  // console.log('ServicesCounselingEditShipmentDetails closeoutOffice', move.closeoutOffice);
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
