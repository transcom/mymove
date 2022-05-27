import React, { useState } from 'react';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { generatePath } from 'react-router';
import { useHistory, useParams, withRouter } from 'react-router-dom';
import { queryCache, useMutation } from 'react-query';
import { connect } from 'react-redux';
import { func } from 'prop-types';

import PrimeUIShipmentUpdateReweighForm from './PrimeUIShipmentUpdateReweighForm';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import scrollToTop from 'shared/scrollToTop';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { primeSimulatorRoutes } from 'constants/routes';
import { updatePrimeMTOShipmentReweigh } from 'services/primeApi';
import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';

const PrimeUIShipmentUpdateReweigh = ({ setFlashMessage }) => {
  const [errorMessage, setErrorMessage] = useState();
  const { moveCodeOrID, shipmentId, reweighId } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);
  const mtoShipments = moveTaskOrder?.mtoShipments;
  const shipment = mtoShipments?.find((mtoShipment) => mtoShipment?.id === shipmentId);
  const history = useHistory();

  const handleClose = () => {
    history.push(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  const [mutateMTOShipmentReweigh] = useMutation(updatePrimeMTOShipmentReweigh, {
    onSuccess: (updatedReweigh) => {
      const updatedMTOShipment = {
        ...shipment,
        reweigh: updatedReweigh,
      };

      mtoShipments[mtoShipments.findIndex((s) => s.id === updatedReweigh.shipmentID)] = updatedMTOShipment;

      queryCache.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryCache.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);

      setFlashMessage(
        `MSG_UPDATE_REWEIGH_SUCCESS${shipmentId}`,
        'success',
        `Successfully updated shipment reweigh`,
        '',
        true,
      );

      handleClose();
    },
    onError: (error) => {
      const { response: { body } = {} } = error;
      if (body) {
        setErrorMessage({
          title: `Prime API: ${body.title} `,
          detail: `${body.detail}`,
        });
      } else {
        setErrorMessage({
          title: 'Unexpected error',
          detail: 'An unknown error has occurred, please check the address values used',
        });
      }
      scrollToTop();
    },
  });

  const handleSubmit = (values) => {
    mutateMTOShipmentReweigh({
      mtoShipmentID: shipmentId,
      reweighID: reweighId,
      ifMatchETag: shipment.reweigh.eTag,
      body: {
        weight: Number(values.reweighWeight),
        verificationReason: values.reweighRemarks,
      },
    });
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const initialValues = {
    reweighWeight: shipment.reweigh ? String(shipment.reweigh.weight) : '0',
    reweighRemarks:
      shipment.reweigh && shipment.reweigh.verificationReason !== null ? shipment.reweigh.verificationReason : '',
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              {errorMessage?.detail && (
                <div className={primeStyles.errorContainer}>
                  <Alert headingLevel="h4" type="error">
                    <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                    <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
                  </Alert>
                </div>
              )}
              <h1>Edit Reweigh</h1>
              <PrimeUIShipmentUpdateReweighForm
                onSubmit={handleSubmit}
                handleClose={handleClose}
                initialValues={initialValues}
              />
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

PrimeUIShipmentUpdateReweigh.propTypes = {
  setFlashMessage: func,
};

PrimeUIShipmentUpdateReweigh.defaultProps = {
  setFlashMessage: () => {},
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default withRouter(connect(() => ({}), mapDispatchToProps)(PrimeUIShipmentUpdateReweigh));
