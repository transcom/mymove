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
import { createMTOShipment } from 'services/ghcApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const ServicesCounselingAddShipment = ({ match }) => {
  const { moveCode } = useParams();
  const history = useHistory();
  const { move, order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);
  const [mutateMTOShipments] = useMutation(createMTOShipment, {
    onSuccess: (newMTOShipment) => {
      mtoShipments.push(newMTOShipment);
      queryCache.setQueryData([MTO_SHIPMENTS, newMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryCache.invalidateQueries([MTO_SHIPMENTS, newMTOShipment.moveTaskOrderID]);
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances } = order;
  const weightAllotment = { ...allowances, totalWeightSelf: allowances.authorizedWeight };

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
                  submitHandler={mutateMTOShipments}
                  isCreatePage
                  currentResidence={customer.current_address}
                  newDutyStationAddress={order.destinationDutyStation?.address}
                  selectedMoveType={SHIPMENT_OPTIONS.HHG}
                  serviceMember={{ weightAllotment }}
                  moveTaskOrderID={move.id}
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
