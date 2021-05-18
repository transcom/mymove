import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';

import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { MatchShape } from 'types/officeShapes';
import ServicesCounselingShipmentForm from 'components/Office/ServicesCounselingShipmentForm/ServicesCounselingShipmentForm';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { useEditShipmentQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { updateMTOShipment } from 'services/ghcApi';

const ServicesCounselingEditShipmentDetails = ({ match }) => {
  const { moveCode, shipmentId } = useParams();
  const history = useHistory();
  const { order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);
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

  return (
    <GridContainer containerSize="widescreen">
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <ServicesCounselingShipmentForm
            match={match}
            history={history}
            updateMTOShipment={mutateMTOShipment}
            isCreatePage={false}
            currentResidence={customer.current_address}
            newDutyStationAddress={order.destinationDutyStation?.address}
            selectedMoveType={SHIPMENT_OPTIONS.HHG}
            mtoShipment={matchingShipment}
            serviceMember={{ weightAllotment }}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

ServicesCounselingEditShipmentDetails.propTypes = {
  match: MatchShape.isRequired,
};

export default ServicesCounselingEditShipmentDetails;
