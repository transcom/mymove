import classnames from 'classnames';
import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory, useParams } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

import styles from './FinalCloseout.module.scss';

import FinalCloseoutForm from 'components/Customer/PPM/Closeout/FinalCloseoutForm/FinalCloseoutForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { generalRoutes } from 'constants/routes';
import { shipmentTypes } from 'constants/shipments';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { ppmSubmissionCertificationText } from 'scenes/Legalese/legaleseText';
import { getMTOShipmentsForMove, getResponseError, submitPPMShipmentSignedCertification } from 'services/internalApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { updateMTOShipment } from 'store/entities/actions';
import { selectServiceMemberAffiliation, selectMTOShipmentById } from 'store/entities/selectors';
import { selectMove } from 'shared/Entities/modules/moves';
import { formatSwaggerDate } from 'utils/formatters';
import { setFlashMessage } from 'store/flash/actions';

const FinalCloseout = () => {
  const history = useHistory();
  const dispatch = useDispatch();
  const [errorMessage, setErrorMessage] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const { moveId, mtoShipmentId } = useParams();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const affiliation = useSelector((state) => selectServiceMemberAffiliation(state));
  const selectedMove = useSelector((state) => selectMove(state, moveId));

  useEffect(() => {
    getMTOShipmentsForMove(moveId)
      .then((response) => {
        dispatch(updateMTOShipment(response.mtoShipments[mtoShipmentId]));
      })
      .catch(() => {
        setErrorMessage('Failed to fetch shipment information');
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, [moveId, mtoShipmentId, dispatch]);

  if (!mtoShipment || isLoading) {
    return <LoadingPlaceholder />;
  }

  const handleBack = () => {
    history.push(generalRoutes.HOME_PATH);
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
      .then((response) => {
        dispatch(
          updateMTOShipment({
            ...mtoShipment,
            ppmShipment: response,
          }),
        );

        dispatch(
          setFlashMessage('PPM_SUBMITTED', 'success', 'You submitted documentation for review.', undefined, false),
        );

        history.push(generalRoutes.HOME_PATH);
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

            <h1>Complete PPM</h1>

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
              affiliation={affiliation}
              selectedMove={selectedMove}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default FinalCloseout;
