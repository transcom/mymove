import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';

import styles from '../ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

import 'styles/office.scss';
import CustomerHeader from 'components/CustomerHeader';
import ServicesCounselingShipmentForm from 'components/Office/ServicesCounselingShipmentForm/ServicesCounselingShipmentForm';
import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { MatchShape } from 'types/officeShapes';
import { useEditShipmentQueries } from 'hooks/queries';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { updateMTOShipment } from 'services/ghcApi';

const ServicesCounselingEditShipmentDetails = ({ match }) => {
  const { moveCode, shipmentId } = useParams();
  const history = useHistory();
  const { move, order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);
  const [mutateMTOShipment] = useMutation(updateMTOShipment, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryCache.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryCache.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances } = order;

  const matchingShipment = mtoShipments?.filter((shipment) => shipment.id === shipmentId)[0];
  const weightAllotment = { ...allowances, totalWeightSelf: allowances.authorizedWeight };

  const TACs = {
    HHG: order.tac,
    NTS: order.nts_tac,
  };

  const SACs = {
    HHG: order.sac,
    NTS: order.nts_sac,
  };

  return (
    <>
      <CustomerHeader order={order} customer={customer} moveCode={moveCode} />
      <div className={styles.tabContent}>
        <div className={styles.container}>
          <GridContainer className={styles.gridContainer}>
            <Grid row>
              <Grid col desktop={{ col: 8, offset: 2 }}>
                <ServicesCounselingShipmentForm
                  match={match}
                  history={history}
                  submitHandler={mutateMTOShipment}
                  isCreatePage={false}
                  currentResidence={customer.current_address}
                  newDutyStationAddress={order.destinationDutyStation?.address}
                  selectedMoveType={SHIPMENT_OPTIONS.HHG}
                  mtoShipment={matchingShipment}
                  serviceMember={{ weightAllotment }}
                  moveTaskOrderID={move.id}
                  mtoShipments={mtoShipments}
                  TACs={TACs}
                  SACs={SACs}
                />
              </Grid>
            </Grid>
          </GridContainer>
        </div>
      </div>
    </>
  );
};

ServicesCounselingEditShipmentDetails.propTypes = {
  match: MatchShape.isRequired,
};

export default ServicesCounselingEditShipmentDetails;
