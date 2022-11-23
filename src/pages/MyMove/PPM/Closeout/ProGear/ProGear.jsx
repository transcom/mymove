import React, { useEffect, useState } from 'react';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

import { selectMTOShipmentById, selectProGearWeightTicketAndIndexById } from 'store/entities/selectors';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import { customerRoutes, generalRoutes } from 'constants/routes';
import {
  createUploadForPPMDocument,
  createProGearWeightTicket,
  deleteUpload,
  patchProGearWeightTicket,
} from 'services/internalApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import closingPageStyles from 'pages/MyMove/PPM/Closeout/Closeout.module.scss';
import ProGearForm from 'components/Customer/PPM/Closeout/ProGearForm/ProGearForm';
import { updateMTOShipment } from 'store/entities/actions';

const ProGear = () => {
  const dispatch = useDispatch();
  const history = useHistory();
  const handleBack = () => {
    history.push(generalRoutes.HOME_PATH);
  };
  const [errorMessage, setErrorMessage] = useState(null);
  const { moveId, mtoShipmentId, proGearId } = useParams();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const { proGearWeightTicket: currentProGearWeightTicket, index: currentIndex } = useSelector((state) =>
    selectProGearWeightTicketAndIndexById(state, mtoShipmentId, proGearId),
  );

  useEffect(() => {
    if (!proGearId) {
      createProGearWeightTicket(mtoShipment?.ppmShipment?.id)
        .then((resp) => {
          if (mtoShipment?.ppmShipment?.proGearWeightTickets) {
            mtoShipment.ppmShipment.proGearWeightTickets.push(resp);
          } else {
            mtoShipment.ppmShipment.proGearWeightTickets = [resp];
          }
          // Update the URL so the back button would work and not create a new weight ticket or on
          // refresh either.
          history.replace(
            generatePath(customerRoutes.SHIPMENT_PPM_PRO_GEAR_EDIT_PATH, {
              moveId,
              mtoShipmentId,
              proGearId: resp.id,
            }),
          );
          dispatch(updateMTOShipment(mtoShipment));
        })
        .catch(() => {
          setErrorMessage('Failed to create trip record');
        });
    }
  }, [proGearId, moveId, mtoShipmentId, history, dispatch, mtoShipment]);

  const handleCreateUpload = async (fieldName, file, setFieldTouched) => {
    const documentId = currentProGearWeightTicket[`${fieldName}Id`];

    createUploadForPPMDocument(mtoShipment.ppmShipment.id, documentId, file)
      .then((upload) => {
        mtoShipment.ppmShipment.proGearWeightTickets[currentIndex][fieldName].uploads.push(upload);
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
        const filteredUploads = mtoShipment.ppmShipment.proGearWeightTickets[currentIndex][fieldName].uploads.filter(
          (upload) => upload.id !== uploadId,
        );
        mtoShipment.ppmShipment.proGearWeightTickets[currentIndex][fieldName].uploads = filteredUploads;
        setFieldValue(fieldName, filteredUploads, true);
        setFieldTouched(fieldName, true, true);
        dispatch(updateMTOShipment(mtoShipment));
      })
      .catch(() => {
        setErrorMessage('Failed to delete the file upload');
      });
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);
    const hasWeightTickets = !values.missingWeightTicket;
    const belongsToSelf = values.belongsToSelf === 'true';
    const payload = {
      ppmShipmentId: mtoShipment.ppmShipment.id,
      proGearWeightTicketId: currentProGearWeightTicket.id,
      description: values.description,
      weight: parseInt(values.weight, 10),
      belongsToSelf,
      hasWeightTickets,
    };

    patchProGearWeightTicket(
      mtoShipment?.ppmShipment?.id,
      currentProGearWeightTicket.id,
      payload,
      currentProGearWeightTicket.eTag,
    )
      .then((resp) => {
        setSubmitting(false);
        mtoShipment.ppmShipment.proGearWeightTickets[currentIndex] = resp;
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

  if (!mtoShipment || !currentProGearWeightTicket) {
    return renderError() || <LoadingPlaceholder />;
  }
  return (
    <div className={ppmPageStyles.ppmPageStyle}>
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Pro-gear</h1>
            {renderError()}
            <div className={closingPageStyles['closing-section']}>
              <p>
                If you moved pro-gear for yourself or your spouse as part of this PPM, document the total weight here.
                Reminder: This pro-gear should be included in your total weight moved.
              </p>
            </div>
            <ProGearForm
              proGear={currentProGearWeightTicket}
              setNumber={currentIndex + 1}
              onBack={handleBack}
              onSubmit={handleSubmit}
              onCreateUpload={handleCreateUpload}
              onUploadComplete={handleUploadComplete}
              onUploadDelete={handleUploadDelete}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default ProGear;
