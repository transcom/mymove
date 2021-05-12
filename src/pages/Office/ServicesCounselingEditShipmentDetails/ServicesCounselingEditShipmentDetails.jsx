import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { GridContainer } from '@trussworks/react-uswds';

import ServicesCounselingShipmentForm from 'components/Office/ServicesCounselingShipmentForm/ServicesCounselingShipmentForm';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { useMoveDetailsQueries } from 'hooks/queries';
import { MatchShape } from 'types/customerShapes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const ServicesCounselingEditShipmentDetails = ({ match }) => {
  const { moveCode, shipmentId } = useParams();
  const history = useHistory();
  const { order, mtoShipments, isLoading, isError } = useMoveDetailsQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances } = order;
  const matchingShipment = mtoShipments.filter((shipment) => shipment.id === shipmentId)[0];
  const weightAllotment = { ...allowances, total_weight_self: allowances.authorizedWeight };

  /*
  console.log(`order:`);
  console.log(order);
  console.log(`mtoShipments`);
  console.log(mtoShipments);
  console.log(`customer:`);
  console.log(customer);
  console.log(allowances);
  */

  return (
    <GridContainer>
      <ServicesCounselingShipmentForm
        match={match}
        history={history}
        updateMtoShipment={() => {}}
        isCreatePage={false}
        currentResidence={customer.current_address}
        newDutyStationAddress={order.destinationDutyStation}
        selectedMoveType={SHIPMENT_OPTIONS.HHG}
        mtoShipment={matchingShipment}
        serviceMember={{ ...customer, weight_allotment: weightAllotment }}
      />
    </GridContainer>
  );
};

ServicesCounselingEditShipmentDetails.propTypes = {
  match: MatchShape,
};

ServicesCounselingEditShipmentDetails.defaultProps = {
  match: { isExact: false, params: { moveID: '' } },
};

export default ServicesCounselingEditShipmentDetails;
