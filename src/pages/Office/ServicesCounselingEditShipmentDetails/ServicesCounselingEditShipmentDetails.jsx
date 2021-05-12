import React from 'react';
import { matchPath, useHistory, useLocation, useParams } from 'react-router-dom';
import { GridContainer } from '@trussworks/react-uswds';

import ServicesCounselingShipmentForm from 'components/Office/ServicesCounselingShipmentForm/ServicesCounselingShipmentForm';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { useMoveDetailsQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const ServicesCounselingEditShipmentDetails = () => {
  const { moveCode, shipmentId } = useParams();
  const history = useHistory();
  const { order, mtoShipments, isLoading, isError } = useMoveDetailsQueries(moveCode);
  const { pathname } = useLocation();

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances } = order;
  const matchingShipment = mtoShipments.filter((shipment) => shipment.id === shipmentId)[0];
  const weightAllotment = { ...allowances, total_weight_self: allowances.authorizedWeight };

  const updateMtoShipmentHandler = () => {};

  return (
    <GridContainer>
      <ServicesCounselingShipmentForm
        match={matchPath(pathname, {
          isExact: true,
          path: '/moves/:moveId/:shipmentId',
        })}
        history={history}
        updateMTOShipment={updateMtoShipmentHandler}
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

export default ServicesCounselingEditShipmentDetails;
