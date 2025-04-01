import classnames from 'classnames';
import React, { useState } from 'react';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

import styles from './FinalCloseout.module.scss';

import { useEditShipmentQueries } from 'hooks/queries';
import FinalCloseoutForm from 'pages/Office/PPM/Closeout/FinalCloseoutForm/FinalCloseoutForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { servicesCounselingRoutes } from 'constants/routes';
import { shipmentTypes } from 'constants/shipments';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { ppmSubmissionCertificationText } from 'scenes/Legalese/legaleseText';
import { getResponseError, submitPPMShipmentSignedCertification } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { formatSwaggerDate } from 'utils/formatters';
import { setFlashMessage } from 'store/flash/actions';

const FinalCloseout = () => {
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState(null);
  // const [isLoading, setIsLoading] = useState(true);
  // const [selectedMove, setSelectedMove] = useState(false);
  const { moveCode, shipmentId } = useParams();
  const { selectedMove, order, mtoShipments, isLoading, isError } = useEditShipmentQueries(moveCode);

  if (isLoading) {
    return <LoadingPlaceholder />;
  }

  if (isError) return <SomethingWentWrong />;

  const mtoShipment = mtoShipments.find((shipment) => shipment.id === shipmentId);

  const handleBack = () => {
    navigate(generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, { moveCode, shipmentId }));
  };

  const handleSubmit = (values) => {
    setErrorMessage(null);
    const ppmShipmentId = mtoShipment.ppmShipment.id;

    const payload = {
      certification_text: ppmSubmissionCertificationText,
      signature: values.signature,
      date: values.date,
    };

    submitPPMShipmentSignedCertification(ppmShipmentId, payload)
      .then(() => {
        setFlashMessage('PPM_SUBMITTED', 'success', 'You submitted documentation for review.', undefined, false);
        navigate(generatePath(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, { moveCode }));
      })
      .catch((err) => {
        setErrorMessage(getResponseError(err.response, 'Failed to submit PPM documentation due to server error.'));
      });
  };

  const initialValues = {
    signature: '',
    date: formatSwaggerDate(new Date()),
  };

  return (
    <div className={classnames(ppmPageStyles.ppmPageStyle, styles.FinalCloseout)}>
      <NotificationScrollToTop dependency={errorMessage} />

      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />

            <h1 data-testid="scCompletePPMHeader">Complete PPM</h1>

            {errorMessage && (
              <Alert headingLevel="h4" slim type="error">
                {errorMessage}
              </Alert>
            )}

            <FinalCloseoutForm
              initialValues={initialValues}
              mtoShipment={mtoShipment}
              onBack={handleBack}
              onSubmit={handleSubmit}
              affiliation={order.agency}
              selectedMove={selectedMove}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default FinalCloseout;
