import React from 'react';
import { matchPath, useHistory, useLocation, useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';

import { MTO_SHIPMENTS, ORDERS } from 'constants/queryKeys';
import ServicesCounselingShipmentForm from 'components/Office/ServicesCounselingShipmentForm/ServicesCounselingShipmentForm';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { useEditShipmentQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { updateMTOShipment } from 'services/ghcApi';

const ServicesCounselingEditShipmentDetails = () => {
  const { moveCode, shipmentId } = useParams();
  const history = useHistory();
  const { order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);
  const { pathname } = useLocation();
  const [mutateMTOShipment] = useMutation(updateMTOShipment, {
    onSuccess: (data) => {
      queryCache.setQueryData([ORDERS, data.locator], data);
      queryCache.setQueryData([MTO_SHIPMENTS, data]);
      queryCache.invalidateQueries([ORDERS, data.locator]);
      queryCache.invalidateQueries([MTO_SHIPMENTS, data.id]);
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances } = order;
  const matchingShipment = mtoShipments.filter((shipment) => shipment.id === shipmentId)[0];
  const weightAllotment = { ...allowances, total_weight_self: allowances.authorizedWeight };

  return (
    <GridContainer containerSize="widescreen">
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <ServicesCounselingShipmentForm
            match={matchPath(pathname, {
              isExact: true,
              path: '/moves/:moveCode/:shipmentId/edit',
            })}
            history={history}
            updateMTOShipment={mutateMTOShipment}
            isCreatePage={false}
            currentResidence={customer.current_address}
            newDutyStationAddress={order.destinationDutyStation?.address}
            selectedMoveType={SHIPMENT_OPTIONS.HHG}
            mtoShipment={matchingShipment}
            serviceMember={{ weight_allotment: weightAllotment }}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

export default ServicesCounselingEditShipmentDetails;
