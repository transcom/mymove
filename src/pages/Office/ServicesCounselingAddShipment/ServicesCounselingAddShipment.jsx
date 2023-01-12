import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';

import styles from '../ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

import 'styles/office.scss';
import CustomerHeader from 'components/CustomerHeader';
import ShipmentForm from 'components/Office/ShipmentForm/ShipmentForm';
import { MOVES, MTO_SHIPMENTS } from 'constants/queryKeys';
import { MatchShape } from 'types/officeShapes';
import { useEditShipmentQueries } from 'hooks/queries';
import { createMTOShipment, updateMoveCloseoutOffice } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { roleTypes } from 'constants/userRoles';
import { SHIPMENT_OPTIONS, SHIPMENT_OPTIONS_URL } from 'shared/constants';

// this returns an object that contains a promise and another thing
function foobar({ shipment, closeoutOffice }) {
  console.log('addship submitHandler (foobar)', { shipment, closeoutOffice });
  return createMTOShipment(shipment).then((newShipment) => {
    return { newShipment, closeoutOffice };
  });
  // return { newShipment: createMTOShipment(shipment), closeoutOffice };
}
const ServicesCounselingAddShipment = ({ match }) => {
  const params = useParams();
  let { shipmentType } = params;
  const { moveCode } = params;

  if (shipmentType === SHIPMENT_OPTIONS_URL.NTSrelease) {
    shipmentType = SHIPMENT_OPTIONS.NTSR;
  } else {
    shipmentType = SHIPMENT_OPTIONS[shipmentType];
  }

  const history = useHistory();
  const { move, order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);
  // what does this syntax do?
  const [mutateMoveCloseoutOffice] = useMutation(updateMoveCloseoutOffice, {
    onSuccess: (updatedMove) => {
      // invalidate some query data?
      queryCache.invalidateQueries([MOVES, moveCode]);
      console.log('updateMoveCloseoutOffice success');
    },
    onError: (error) => {
      // invalidate some query data?
      console.error('update closeout office mutation', error);
    },
  });
  // I think useMutation might wait for a promise to be resolved from the return value, but if it's not a promise, what will it do?
  // I need to find the expectations for this function
  const [mutateMTOShipments] = useMutation(foobar, {
    onSuccess: (result) => {
      if (result.closeoutOffice) {
        console.log('lets try to submit the closeout office', result.closeoutOffice);
        // TODO this is wrong, need move info in args
        mutateMoveCloseoutOffice({
          locator: moveCode,
          ifMatchETag: move.eTag,
          body: { closeoutOfficeId: result.closeoutOffice.id },
        }).then(() => {
          console.log('mutate closeout done');
        });
      }
      // TODO i'm not sure if we wait for the promise above to resolve before getting to this stuff
      mtoShipments.push(result.newShipment);
      queryCache.setQueryData([MTO_SHIPMENTS, result.newShipment.moveTaskOrderID, false], mtoShipments);
      queryCache.invalidateQueries([MTO_SHIPMENTS, result.newShipment.moveTaskOrderID]);
      return result.newShipment;
    },
    onError: (error) => {
      // TODO invalidate some query data?
      console.error('create shipment mutation error', error);
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
    <>
      <CustomerHeader order={order} customer={customer} moveCode={moveCode} />
      <div className={styles.tabContent}>
        <div className={styles.container}>
          <GridContainer className={styles.gridContainer}>
            <Grid row>
              <Grid col desktop={{ col: 8, offset: 2 }}>
                <ShipmentForm
                  match={match}
                  history={history}
                  submitHandler={mutateMTOShipments}
                  isCreatePage
                  ServicesCounselingShipmentForm
                  currentResidence={customer.current_address}
                  originDutyLocationAddress={order.originDutyLocation?.address}
                  newDutyLocationAddress={order.destinationDutyLocation?.address}
                  shipmentType={shipmentType}
                  serviceMember={{ weightAllotment, agency: customer.agency }}
                  moveTaskOrderID={move.id}
                  mtoShipments={mtoShipments}
                  TACs={TACs}
                  SACs={SACs}
                  userRole={roleTypes.SERVICES_COUNSELOR}
                  displayDestinationType
                  closeoutOffice={move.closeoutOffice}
                  move={move}
                />
              </Grid>
            </Grid>
          </GridContainer>
        </div>
      </div>
    </>
  );
};

ServicesCounselingAddShipment.propTypes = {
  match: MatchShape.isRequired,
};

export default ServicesCounselingAddShipment;
