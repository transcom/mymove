import classnames from 'classnames';
import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory, useParams } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

import styles from './FinalCloseout.module.scss';

import FinalCloseoutForm from 'components/Customer/PPM/Closeout/FinalCloseoutForm/FinalCloseoutForm';
import ScrollToTop from 'components/ScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { generalRoutes } from 'constants/routes';
import { shipmentTypes } from 'constants/shipments';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { updateMTOShipment } from 'store/entities/actions';
import { selectMTOShipmentById } from 'store/entities/selectors';

const FinalCloseout = () => {
  const history = useHistory();
  const dispatch = useDispatch();
  const [errorMessage, setErrorMessage] = useState();
  const { mtoShipmentId } = useParams();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));

  if (!mtoShipment) {
    return <LoadingPlaceholder />;
  }

  const handleBack = () => {
    history.push(generalRoutes.HOME_PATH);
  };

  const handleSubmit = (values, { setSubmitting }) => {
    setErrorMessage(null);

    const payload = {
      ppmShipment: {
        id: mtoShipment.ppmShipment.id,
      },
    };

    patchMTOShipment(mtoShipmentId, payload, mtoShipment.eTag)
      .then((response) => {
        setSubmitting(false);

        dispatch(updateMTOShipment(response));

        history.push(generalRoutes.HOME_PATH);
      })
      .catch((err) => {
        setSubmitting(false);
        setErrorMessage(getResponseError(err.response, 'Failed to submit PPM documentation due to server error.'));
      });
  };

  return (
    <div className={classnames(ppmPageStyles.ppmPageStyle, styles.FinalCloseout)}>
      <ScrollToTop otherDep={errorMessage} />

      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />

            <h1>Complete PPM</h1>

            {errorMessage && (
              <Alert headingLevel="h4" slim type="error">
                {errorMessage}
              </Alert>
            )}

            <FinalCloseoutForm mtoShipment={mtoShipment} onBack={handleBack} onSubmit={handleSubmit} />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default FinalCloseout;
