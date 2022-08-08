import React, { useEffect, useState } from 'react';
import { generatePath, useHistory, useParams, useLocation } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import qs from 'query-string';

import { selectMTOShipmentById, selectWeightTicketAndIndexById } from 'store/entities/selectors';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { createUploadForDocument, createWeightTicket, deleteUpload, patchWeightTicket } from 'services/internalApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ScrollToTop from 'components/ScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import closingPageStyles from 'pages/MyMove/PPM/Closeout/Closeout.module.scss';
import WeightTicketForm from 'components/Customer/PPM/Closeout/WeightTicketForm/WeightTicketForm';
import { updateMTOShipment } from 'store/entities/actions';

const WeightTickets = () => {
  const [errorMessage, setErrorMessage] = useState();

  const dispatch = useDispatch();
  const history = useHistory();
  const { moveId, mtoShipmentId, weightTicketId } = useParams();

  const { search } = useLocation();

  const { tripNumber } = qs.parse(search);

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const { weightTicket: currentWeightTicket, index: currentIndex } = useSelector((state) =>
    selectWeightTicketAndIndexById(state, mtoShipmentId, weightTicketId),
  );

  useEffect(() => {
    if (!weightTicketId) {
      createWeightTicket(mtoShipment?.ppmShipment?.id)
        .then((resp) => {
          if (mtoShipment?.ppmShipment?.weightTickets) {
            mtoShipment.ppmShipment.weightTickets.push(resp);
          } else {
            mtoShipment.ppmShipment.weightTickets = [resp];
          }
          // Update the URL so the back button would work and not create a new weight ticket or on
          // refresh either.
          history.replace(
            generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
              moveId,
              mtoShipmentId,
              weightTicketId: resp.id,
            }),
          );
          dispatch(updateMTOShipment(mtoShipment));
        })
        .catch(() => {
          setErrorMessage('Failed to create trip record');
        });
    }
  }, [weightTicketId, moveId, mtoShipmentId, history, dispatch, mtoShipment]);

  const handleCreateUpload = async (fieldName, file, setFieldTouched) => {
    const documentId = currentWeightTicket[`${fieldName}Id`];

    createUploadForDocument(file, documentId)
      .then((upload) => {
        mtoShipment.ppmShipment.weightTickets[currentIndex][fieldName].uploads.push(upload);
        dispatch(updateMTOShipment(mtoShipment));
        setFieldTouched(fieldName, true);
        return upload;
      })
      .catch(() => {
        setErrorMessage('Failed to save the file upload');
      });
  };

  const handleUploadComplete = (err) => {
    if (err) {
      setErrorMessage('Encountered error when completing file upload');
    }
  };

  const handleUploadDelete = (uploadId, fieldName, setFieldTouched, setFieldValue) => {
    deleteUpload(uploadId)
      .then(() => {
        const filteredUploads = mtoShipment.ppmShipment.weightTickets[currentIndex][fieldName].uploads.filter(
          (upload) => upload.id !== uploadId,
        );
        mtoShipment.ppmShipment.weightTickets[currentIndex][fieldName].uploads = filteredUploads;

        setFieldValue(fieldName, filteredUploads, true);
        setFieldTouched(fieldName, true, true);
        dispatch(updateMTOShipment(mtoShipment));
      })
      .catch(() => {
        setErrorMessage('Failed to delete the file upload');
      });
  };

  const handleBack = () => {
    history.push(generalRoutes.HOME_PATH);
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);
    const ownsTrailer = values.ownsTrailer === 'true';
    const trailerMeetsCriteria = ownsTrailer ? values.trailerMeetsCriteria === 'true' : false;
    const payload = {
      ppmShipmentId: mtoShipment.ppmShipment.id,
      vehicleDescription: values.vehicleDescription,
      emptyWeight: parseInt(values.emptyWeight, 10),
      missingEmptyWeightTicket: values.missingEmptyWeightTicket,
      fullWeight: parseInt(values.fullWeight, 10),
      missingFullWeightTicket: values.missingFullWeightTicket,
      ownsTrailer,
      trailerMeetsCriteria,
    };

    patchWeightTicket(mtoShipment?.ppmShipment?.id, currentWeightTicket.id, payload, currentWeightTicket.eTag)
      .then((resp) => {
        setSubmitting(false);
        mtoShipment.ppmShipment.weightTickets[currentIndex] = resp;
        history.push(generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId, mtoShipmentId }));
        dispatch(updateMTOShipment(mtoShipment));
      })
      .catch(() => {
        setSubmitting(false);
        setErrorMessage('Failed to save updated trip record');
      });
  };

  const renderError = () => {
    if (!errorMessage) {
      return null;
    }

    return (
      <Alert slim type="error">
        {errorMessage}
      </Alert>
    );
  };

  if (!mtoShipment || !currentWeightTicket) {
    return renderError() || <LoadingPlaceholder />;
  }

  return (
    <div className={ppmPageStyles.ppmPageStyle}>
      <ScrollToTop otherDep={errorMessage} />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Weight Tickets</h1>
            {renderError()}
            <div className={closingPageStyles['closing-section']}>
              <p>
                Weight tickets should include both an empty or full weight ticket for each segment or trip. If you’re
                missing a weight ticket, you’ll be able to use a government-created spreadsheet to estimate the weight.
              </p>
              <p>Weight tickets must be certified, legible, and unaltered. Files must be 25MB or smaller.</p>
              <p>You must upload at least one set of weight tickets to get paid for your PPM.</p>
            </div>
            <WeightTicketForm
              weightTicket={currentWeightTicket}
              tripNumber={tripNumber}
              onCreateUpload={handleCreateUpload}
              onUploadComplete={handleUploadComplete}
              onUploadDelete={handleUploadDelete}
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
