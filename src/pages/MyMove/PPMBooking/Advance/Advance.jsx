import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import ppmBookingPageStyles from 'pages/MyMove/PPMBooking/PPMBooking.module.scss';
import { customerRoutes } from 'constants/routes';
import AdvanceForm from 'components/Customer/PPMBooking/Advance/AdvanceForm';
import { shipmentTypes } from 'constants/shipments';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import { selectMTOShipmentById } from 'store/entities/selectors';

const Advance = () => {
  const [errorMessage, setErrorMessage] = useState();
  const history = useHistory();
  const { moveId, mtoShipmentId, shipmentNumber } = useParams();
  const dispatch = useDispatch();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));

  const handleBack = () => {
    history.push(generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH, { moveId, mtoShipmentId }));
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);

    const advanceRequested = values.advanceRequested === 'true';

    const payload = {
      ppmShipment: {
        id: mtoShipment.ppmShipment.id,
        advance: values.amountRequested ? values.amountRequested * 100 : null,
        advanceRequested,
      },
    };

    patchMTOShipment(mtoShipment.id, payload, mtoShipment.eTag)
      .then((response) => {
        setSubmitting(false);
        dispatch(updateMTOShipment(response));
        history.push(generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId }));
      })
      .catch((err) => {
        setSubmitting(false);

        setErrorMessage(getResponseError(err.response, 'Failed to update MTO shipment due to server error.'));
      });
  };

  return (
    <div className={ppmBookingPageStyles.PPMBookingPage}>
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} shipmentNumber={shipmentNumber} />
            <h1>Advances</h1>
            {errorMessage && (
              <Alert slim type="error">
                {errorMessage}
              </Alert>
            )}
            <AdvanceForm mtoShipment={mtoShipment} onSubmit={handleSubmit} onBack={handleBack} />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default Advance;
