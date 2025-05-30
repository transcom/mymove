import React, { useEffect, useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import { servicesCounselingRoutes } from 'constants/routes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import ppmPageStyles from 'pages/Office/PPM/PPM.module.scss';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import WeightTicketForm from 'components/Shared/PPM/Closeout/WeightTicketForm/WeightTicketForm';
import { usePPMShipmentAndDocsOnlyQueries } from 'hooks/queries';
import {
  createWeightTicket,
  patchWeightTicket,
  createUploadForPPMDocument,
  deleteUploadForDocument,
} from 'services/ghcApi';
import { DOCUMENTS } from 'constants/queryKeys';
import { APP_NAME } from 'constants/apps';
import ErrorModal from 'shared/ErrorModal/ErrorModal';
import appendTimestampToFilename from 'utils/fileUpload';

const WeightTickets = () => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { moveCode, shipmentId, weightTicketId } = useParams();

  const [isErrorModalVisible, setIsErrorModalVisible] = useState(false);
  const toggleErrorModal = () => {
    setIsErrorModalVisible((prev) => !prev);
  };

  const displayHelpDeskLink = false;

  const errorModalMessage =
    'The only Excel file this uploader accepts is the Weight Estimator file. Please convert any other Excel file to PDF.';

  const { mtoShipment, documents, isError } = usePPMShipmentAndDocsOnlyQueries(shipmentId);
  const appName = APP_NAME.OFFICE;
  const ppmShipment = mtoShipment?.ppmShipment;
  const weightTickets = documents?.WeightTickets ?? [];

  const currentWeightTicket = weightTickets?.find((item) => item.id === weightTicketId) ?? null;
  const currentWeightTicketIdx = Array.isArray(weightTickets)
    ? weightTickets.findIndex((ele) => ele.id === weightTicketId)
    : -1;

  const reviewPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, { moveCode, shipmentId });

  const { mutate: mutateCreateWeightTicket } = useMutation(createWeightTicket, {
    onSuccess: (createdWeightTicket) => {
      queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
      navigate(
        generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
          moveCode,
          shipmentId,
          weightTicketId: createdWeightTicket?.id,
        }),
        { replace: true },
      );
    },
    onError: () => {
      setErrorMessage(`Failed to create trip record`);
    },
  });

  const { mutate: mutatePatchWeightTicket } = useMutation(patchWeightTicket, {
    onSuccess: async () => {
      await queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
      navigate(reviewPath);
    },
    onError: () => {
      setIsSubmitted(false);
      setErrorMessage('Failed to save updated trip record');
    },
  });

  useEffect(() => {
    if (!weightTicketId) {
      mutateCreateWeightTicket(ppmShipment?.id);
    }
  }, [mutateCreateWeightTicket, ppmShipment?.id, weightTicketId]);

  const handleCreateUpload = async (fieldName, file, setFieldTouched) => {
    const documentId = currentWeightTicket[`${fieldName}Id`];

    createUploadForPPMDocument(ppmShipment?.id, documentId, appendTimestampToFilename(file), true)
      .then((upload) => {
        documents?.WeightTickets[currentWeightTicketIdx][fieldName]?.uploads.push(upload);
        setFieldTouched(fieldName, true);
        return upload;
      })
      .catch((err) => {
        if (
          err.response.obj.message ===
          'The uploaded .xlsx file does not match the expected weight estimator file format.'
        ) {
          setIsErrorModalVisible(true);
        } else {
          setErrorMessage('Failed to save the file upload');
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
    deleteUploadForDocument(uploadId, null, ppmShipment?.id)
      .then(() => {
        const filteredUploads = documents?.WeightTickets[currentWeightTicketIdx][fieldName]?.uploads.filter(
          (upload) => upload.id !== uploadId,
        );
        documents.WeightTickets[currentWeightTicketIdx][fieldName].uploads = filteredUploads;
        setFieldValue(fieldName, filteredUploads, true);
        setFieldTouched(fieldName, true, true);
      })
      .catch(() => {
        setErrorMessage('Failed to delete the file upload');
      });
  };

  const handleBack = () => {
    navigate(reviewPath);
  };

  const handleSubmit = async (values) => {
    if (isSubmitted) return;

    setIsSubmitted(true);
    setErrorMessage(null);
    const ownsTrailer = values.ownsTrailer === 'true';
    const trailerMeetsCriteria = ownsTrailer ? values.trailerMeetsCriteria === 'true' : false;
    const payload = {
      ppmShipmentId: ppmShipment?.id,
      vehicleDescription: values.vehicleDescription,
      emptyWeight: parseInt(values.emptyWeight, 10),
      missingEmptyWeightTicket: values.missingEmptyWeightTicket,
      fullWeight: parseInt(values.fullWeight, 10),
      missingFullWeightTicket: values.missingFullWeightTicket,
      ownsTrailer,
      trailerMeetsCriteria,
    };
    mutatePatchWeightTicket({
      ppmShipmentId: currentWeightTicket.ppmShipmentId,
      weightTicketId: currentWeightTicket.id,
      payload,
      eTag: currentWeightTicket.eTag,
    });
  };

  const renderError = () => {
    if (!errorMessage) {
      return null;
    }

    return (
      <Alert data-testid="errorMessage" type="error" headingLevel="h4" heading="An error occurred">
        {errorMessage}
      </Alert>
    );
  };

  if (isError) return <SomethingWentWrong />;

  if (!mtoShipment || !currentWeightTicket) {
    return renderError() || <LoadingPlaceholder />;
  }

  return (
    <div className={ppmPageStyles.tabContent}>
      <div className={ppmPageStyles.container}>
        <NotificationScrollToTop dependency={errorMessage} />
        <GridContainer className={ppmPageStyles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <div className={ppmPageStyles.closeoutPageWrapper}>
                <ShipmentTag shipmentType={shipmentTypes.PPM} />
                <h1>Weight Tickets</h1>
                {renderError()}
                <WeightTicketForm
                  weightTicket={currentWeightTicket}
                  tripNumber={currentWeightTicketIdx + 1}
                  onCreateUpload={handleCreateUpload}
                  onUploadComplete={handleUploadComplete}
                  onUploadDelete={handleUploadDelete}
                  onSubmit={handleSubmit}
                  onBack={handleBack}
                  isSubmitted={isSubmitted}
                  appName={appName}
                />
                <ErrorModal
                  isOpen={isErrorModalVisible}
                  closeModal={toggleErrorModal}
                  errorMessage={errorModalMessage}
                  displayHelpDeskLink={displayHelpDeskLink}
                />
              </div>
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

export default WeightTickets;
