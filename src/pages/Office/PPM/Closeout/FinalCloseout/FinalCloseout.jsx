import classnames from 'classnames';
import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { useMutation } from '@tanstack/react-query';

import styles from './FinalCloseout.module.scss';

import FinalCloseoutForm from 'pages/Office/PPM/Closeout/FinalCloseoutForm/FinalCloseoutForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { servicesCounselingRoutes } from 'constants/routes';
import { shipmentTypes } from 'constants/shipments';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { ppmSubmissionCertificationText } from 'scenes/Legalese/legaleseText';
import { getMove, getMTOShipments, getResponseError, submitPPMShipmentSignedCertification } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { updateMTOShipment } from 'store/entities/actions';
import { selectServiceMemberAffiliation } from 'store/entities/selectors';
import { selectMoveByLocator } from 'shared/Entities/modules/moves';
import { formatSwaggerDate } from 'utils/formatters';
import { setFlashMessage } from 'store/flash/actions';

const FinalCloseout = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const [errorMessage, setErrorMessage] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [mtoShipment, setMtoShipment] = useState(false);
  const { moveCode, shipmentId } = useParams();

  const affiliation = useSelector((state) => selectServiceMemberAffiliation(state));
  const selectedMove = useSelector((state) => selectMoveByLocator(state, moveCode));

  const { mutate: mutateUpdateMTOShipment } = useMutation(updateMTOShipment, {
    onSuccess: () => {
      navigate(
        generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_COMPLETE_PATH, {
          moveCode,
          shipmentId,
        }),
        { replace: true },
      );
    },
    onError: () => {
      setErrorMessage(`Failed to create trip record`);
    },
  });

  useEffect(() => {
    getMove(null, moveCode).then((move) => {
      getMTOShipments(null, move.id)
        .then((response) => {
          setMtoShipment(response.mtoShipments[shipmentId]);
        })
        .catch(() => {
          setErrorMessage('Failed to fetch shipment information');
        })
        .finally(() => {
          setIsLoading(false);
        });
    });
  }, [mutateUpdateMTOShipment, moveCode, shipmentId, dispatch]);

  if (!mtoShipment || isLoading) {
    return <LoadingPlaceholder />;
  }

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
