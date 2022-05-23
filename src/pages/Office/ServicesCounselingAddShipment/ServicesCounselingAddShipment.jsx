import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';

import styles from '../ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

import 'styles/office.scss';
import CustomerHeader from 'components/CustomerHeader';
import ShipmentForm from 'components/Office/ShipmentForm/ShipmentForm';
import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { MatchShape } from 'types/officeShapes';
import { useEditShipmentQueries } from 'hooks/queries';
import { createMTOShipment } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { roleTypes } from 'constants/userRoles';
import { SHIPMENT_OPTIONS, SHIPMENT_OPTIONS_URL } from 'shared/constants';

const ServicesCounselingAddShipment = ({ match }) => {
  const { moveCode, shipmentType } = useParams();
  let selectedMoveType = SHIPMENT_OPTIONS[shipmentType];
  if (shipmentType === SHIPMENT_OPTIONS_URL.NTSrelease) {
    selectedMoveType = SHIPMENT_OPTIONS.NTSR;
  }
  const history = useHistory();
  const { move, order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);
  const [mutateMTOShipments] = useMutation(createMTOShipment, {
    onSuccess: (newMTOShipment) => {
      mtoShipments.push(newMTOShipment);
      queryCache.setQueryData([MTO_SHIPMENTS, newMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryCache.invalidateQueries([MTO_SHIPMENTS, newMTOShipment.moveTaskOrderID]);
      return newMTOShipment;
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
                  newDutyLocationAddress={order.destinationDutyLocation?.address}
                  selectedMoveType={selectedMoveType}
                  serviceMember={{ weightAllotment }}
                  moveTaskOrderID={move.id}
                  mtoShipments={mtoShipments}
                  TACs={TACs}
                  SACs={SACs}
                  userRole={roleTypes.SERVICES_COUNSELOR}
                  displayDestinationType
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
