import React, { useState } from 'react';
import { generatePath, useHistory, useParams, useLocation } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import qs from 'query-string';

import { selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { formatDateForSwagger } from 'shared/dates';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ScrollToTop from 'components/ScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import closingPageStyles from 'pages/MyMove/PPM/Closeout/Closeout.module.scss';
import WeightTicketForm from 'components/Customer/PPM/Closeout/WeightTicketForm/WeightTicketForm';

const WeightTickets = () => {
  const [errorMessage, setErrorMessage] = useState();

  const history = useHistory();
  const { moveId, mtoShipmentId } = useParams();
  const { search } = useLocation();
  const dispatch = useDispatch();

  const { tripNumber } = qs.parse(search);

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  // TODO add selector for selecting weight ticket from Redux store when data changes are solidified

  const handleBack = () => {
    history.push(generalRoutes.HOME_PATH);
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);
    const hasReceivedAdvance = values.hasReceivedAdvance === 'true';
    const payload = {
      shipmentType: mtoShipment.shipmentType,
      ppmShipment: {
        id: mtoShipment.ppmShipment.id,
        actualMoveDate: formatDateForSwagger(values.actualMoveDate),
        actualPickupPostalCode: values.actualPickupPostalCode,
        actualDestinationPostalCode: values.actualDestinationPostalCode,
        hasReceivedAdvance,
        advanceAmountReceived: hasReceivedAdvance ? values.advanceAmountReceived * 100 : null,
      },
    };

    patchMTOShipment(mtoShipment.id, payload, mtoShipment.eTag)
      .then((response) => {
        setSubmitting(false);
        dispatch(updateMTOShipment(response));
        history.push(generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, { moveId, mtoShipmentId }));
      })
      .catch((err) => {
        setSubmitting(false);
        setErrorMessage(getResponseError(err.response, 'Failed to update MTO shipment due to server error.'));
      });
  };

  if (!mtoShipment) {
    return <LoadingPlaceholder />;
  }

  return (
    <div className={ppmPageStyles.ppmPageStyle}>
      <ScrollToTop otherDep={errorMessage} />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Weight Tickets</h1>
            {errorMessage && (
              <Alert slim type="error">
                {errorMessage}
              </Alert>
            )}
            <div className={closingPageStyles['closing-section']}>
              <p>
                Weight tickets should include both an empty or full weight ticket for each segment or trip. If you’re
                missing a weight ticket, you’ll be able to use a government-created spreadsheet to estimate the weight.
              </p>
              <p>Weight tickets must be certified, legible, and unaltered. Files must be 25MB or smaller.</p>
              <p>You must upload at least one set of weight tickets to get paid for your PPM.</p>
            </div>
            <WeightTicketForm
              mtoShipment={mtoShipment}
              tripNumber={tripNumber}
              onSubmit={handleSubmit}
              onBack={handleBack}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default WeightTickets;
