import React, { useEffect, useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { Alert, Grid, GridContainer, Link } from '@trussworks/react-uswds';

import {
  selectMTOShipmentById,
  selectServiceMemberFromLoggedInUser,
  selectWeightTicketAndIndexById,
} from 'store/entities/selectors';
import { customerRoutes } from 'constants/routes';
import {
  createUploadForPPMDocument,
  createWeightTicket,
  deleteUpload,
  getAllMoves,
  getResponseError,
  patchWeightTicket,
} from 'services/internalApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import WeightTicketForm from 'components/Shared/PPM/Closeout/WeightTicketForm/WeightTicketForm';
import { updateAllMoves, updateMTOShipment } from 'store/entities/actions';
import ErrorModal from 'shared/ErrorModal/ErrorModal';
import { CUSTOMER_ERROR_MESSAGES } from 'constants/errorMessages';
import { APP_NAME } from 'constants/apps';
import appendTimestampToFilename from 'utils/fileUpload';

const WeightTickets = () => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [displayHelpDeskLink, setDisplayHelpDeskLink] = useState(false);

  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { moveId, mtoShipmentId, weightTicketId } = useParams();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const { weightTicket: currentWeightTicket, index: currentIndex } = useSelector((state) =>
    selectWeightTicketAndIndexById(state, mtoShipmentId, weightTicketId),
  );

  const serviceMember = useSelector((state) => selectServiceMemberFromLoggedInUser(state));
  const serviceMemberId = serviceMember.id;

  const [isErrorModalVisible, setIsErrorModalVisible] = useState(false);
  const toggleErrorModal = () => {
    setIsErrorModalVisible((prev) => !prev);
  };

  const appName = APP_NAME.MYMOVE;

  const [errorModalMessage, setErrorModalMessage] = useState(
    "Something went wrong uploading your weight ticket. Please try again. If that doesn't fix it, contact the ",
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
          navigate(
            generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
              moveId,
              mtoShipmentId,
              weightTicketId: resp.id,
            }),
            { replace: true },
          );
          dispatch(updateMTOShipment(mtoShipment));
        })
        .catch(() => {
          setErrorMessage('Failed to create trip record');
        });
    }
  }, [weightTicketId, moveId, mtoShipmentId, navigate, dispatch, mtoShipment]);

  useEffect(() => {
    const moves = getAllMoves(serviceMemberId);
    dispatch(updateAllMoves(moves));
  }, [weightTicketId, moveId, mtoShipmentId, navigate, dispatch, mtoShipment, serviceMemberId]);

  const handleErrorClick = () => {
    const path = generatePath(customerRoutes.SHIPMENT_PPM_ABOUT_PATH, {
      moveId,
      mtoShipmentId,
    });

    navigate(path);
  };

  const zipError = (
    <p>
      We are unable to calculate your distance. It may be that you have entered an invalid ZIP code. Please check
      the&nbsp;
      <Link className="usa-link" href="#" onClick={handleErrorClick} data-testid="ZipLink">
        pickup and delivery ZIP codes
      </Link>
      &nbsp;to ensure they were entered correctly and are not PO boxes. If you do not have a different ZIP code, then
      please contact the&nbsp;
      <Link className="usa-link" href="mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil">
        Technical Help Desk
      </Link>
      .
    </p>
  );

  const handleErrorMessage = (error) => {
    if (error?.response?.status === 412) {
      setErrorMessage(CUSTOMER_ERROR_MESSAGES.PRECONDITION_FAILED);
    } else if (
      // this 'else if' can be removed when E-06516 is implemented
      // along with zipError and handleErrorClick
      error?.response?.body?.detail ===
      'We are unable to calculate your distance. It may be that you have entered an invalid ZIP Code. Please check your ZIP Code to ensure it was entered correctly and is not a PO Box.'
    ) {
      setErrorMessage(zipError);
    } else {
      setErrorMessage(getResponseError(error.response, 'Failed to save updated trip record'));
    }
  };

  const handleCreateUpload = async (fieldName, file, setFieldTouched) => {
    const documentId = currentWeightTicket[`${fieldName}Id`];
    createUploadForPPMDocument(mtoShipment.ppmShipment.id, documentId, appendTimestampToFilename(file), true)
      .then((upload) => {
        mtoShipment.ppmShipment.weightTickets[currentIndex][fieldName].uploads.push(upload);
        dispatch(updateMTOShipment(mtoShipment));
        setFieldTouched(fieldName, true);
        setIsErrorModalVisible(false);
        return upload;
      })
      .catch((err) => {
        if (err.response.obj.title === 'Incorrect Xlsx Template') {
          setErrorModalMessage(
            'The only Excel file this uploader accepts is the Weight Estimator file. Please convert any other Excel file to PDF.',
          );
          setIsErrorModalVisible(true);
        } else {
          setDisplayHelpDeskLink(true);
          setIsErrorModalVisible(true);
        }
      });
  };

  const handleUploadComplete = (err) => {
    if (err) {
      setErrorMessage('Encountered error when completing file upload');
    }
  };

  const handleUploadDelete = (uploadId, fieldName, setFieldTouched, setFieldValue) => {
    deleteUpload(uploadId, null, mtoShipment?.ppmShipment?.id)
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
    const path = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
      moveId,
      mtoShipmentId,
    });
    navigate(path);
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
        navigate(generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId, mtoShipmentId }));
        dispatch(updateMTOShipment(mtoShipment));
      })
      .catch((error) => {
        setSubmitting(false);
        handleErrorMessage(error);
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
      <NotificationScrollToTop dependency={errorMessage} />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Weight Tickets</h1>
            {renderError()}
            <WeightTicketForm
              weightTicket={currentWeightTicket}
              tripNumber={currentIndex + 1}
              onCreateUpload={handleCreateUpload}
              onUploadComplete={handleUploadComplete}
              onUploadDelete={handleUploadDelete}
              onSubmit={handleSubmit}
              onBack={handleBack}
              appName={appName}
            />
            <ErrorModal
              isOpen={isErrorModalVisible}
              closeModal={toggleErrorModal}
              errorMessage={errorModalMessage}
              displayHelpDeskLink={displayHelpDeskLink}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default WeightTickets;
