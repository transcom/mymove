import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { GridContainer, Grid, Button, Alert } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';
import classnames from 'classnames';

import styles from '../ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

import scMoveDetailsStyles from './ServicesCounselingMoveDetails.module.scss';

import 'styles/office.scss';
import { updateMoveStatusServiceCounselingCompleted } from 'services/ghcApi';
import { useMoveDetailsQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { MOVES } from 'constants/queryKeys';
import { MOVE_STATUSES } from 'shared/constants';
import AllowancesTable from 'components/Office/AllowancesTable/AllowancesTable';

const ServicesCounselingMoveDetails = () => {
  const { moveCode } = useParams();
  const [alertMessage, setAlertMessage] = useState(null);
  const [alertType, setAlertType] = useState('success');

  const { move, order, isLoading, isError } = useMoveDetailsQueries(moveCode);

  // use mutation calls
  const [mutateMoveStatus] = useMutation(updateMoveStatusServiceCounselingCompleted, {
    onSuccess: (data) => {
      queryCache.setQueryData([MOVES, data.locator], data);
      queryCache.invalidateQueries([MOVES, data.locator]);

      setAlertMessage('Move submitted.');
      setAlertType('success');
    },
    onError: () => {
      setAlertMessage('There was a problem submitting the move. Please try again later.');
      setAlertType('error');
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances } = order;

  const allowancesInfo = {
    branch: customer.agency,
    rank: order.grade,
    weightAllowance: allowances.totalWeight,
    authorizedWeight: allowances.authorizedWeight,
    progear: allowances.proGearWeight,
    spouseProgear: allowances.proGearWeightSpouse,
    storageInTransit: allowances.storageInTransit,
    dependents: allowances.dependentsAuthorized,
    requiredMedicalEquipmentWeight: allowances.requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment: allowances.organizationalClothingAndIndividualEquipment,
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        {/* LeftNav here */}

        <GridContainer
          className={classnames(styles.gridContainer, scMoveDetailsStyles.ServicesCounselingMoveDetails)}
          data-testid="sc-move-details"
        >
          <Grid row className={scMoveDetailsStyles.pageHeader}>
            {alertMessage && (
              <Grid col={12} className={scMoveDetailsStyles.alertContainer}>
                <Alert slim type={alertType}>
                  {alertMessage}
                </Alert>
              </Grid>
            )}
            <Grid col={6}>
              <h1>Move details</h1>
            </Grid>
            <Grid col={6} className={scMoveDetailsStyles.submitMoveDetailsContainer}>
              {move.status === MOVE_STATUSES.NEEDS_SERVICE_COUNSELING && (
                <Button
                  data-testid="submitMoveDetailsBtn"
                  type="button"
                  onClick={() => {
                    mutateMoveStatus({ moveTaskOrderID: move.id, ifMatchETag: move.eTag });
                  }}
                >
                  Submit move details
                </Button>
              )}
            </Grid>
          </Grid>

          {/* additional work here */}
          <div className={styles.section} id="allowances">
            <GridContainer>
              <Grid row gap>
                <Grid col>
                  <AllowancesTable info={allowancesInfo} />
                </Grid>
              </Grid>
            </GridContainer>
          </div>
        </GridContainer>
      </div>
    </div>
  );
};

ServicesCounselingMoveDetails.propTypes = {};

export default ServicesCounselingMoveDetails;
