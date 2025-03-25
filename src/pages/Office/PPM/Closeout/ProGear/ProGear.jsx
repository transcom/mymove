import React, { useEffect, useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import { APP_NAME } from 'constants/apps';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import { servicesCounselingRoutes } from 'constants/routes';
import {
  createProGearWeightTicket,
  patchProGearWeightTicket,
  createUploadForPPMDocument,
  deleteUploadForDocument,
  // updateMTOShipment,
} from 'services/ghcApi';
import { DOCUMENTS } from 'constants/queryKeys';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import closingPageStyles from 'pages/MyMove/PPM/Closeout/Closeout.module.scss';
import ProGearForm from 'components/Shared/PPM/Closeout/ProGearForm/ProGearForm';
import { usePPMShipmentAndDocsOnlyQueries, useReviewShipmentWeightsQuery } from 'hooks/queries';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const ProGear = () => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { moveCode, shipmentId, proGearId } = useParams();

  const { mtoShipment, documents, isError } = usePPMShipmentAndDocsOnlyQueries(shipmentId);
  const { orders } = useReviewShipmentWeightsQuery(moveCode);
  const appName = APP_NAME.OFFICE;
  const ppmShipment = mtoShipment?.ppmShipment;
  const proGearWeightTickets = documents?.ProGearWeightTickets ?? [];

  const currentProGearWeightTicket = proGearWeightTickets?.find((item) => item.id === proGearId) ?? null;
  const currentIndex = Array.isArray(proGearWeightTickets)
    ? proGearWeightTickets.findIndex((ele) => ele.id === proGearId)
    : -1;

  const reviewPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, { moveCode, shipmentId });

  const { mutate: mutateProGearCreateWeightTicket } = useMutation(createProGearWeightTicket, {
    onSuccess: (createdProGearWeightTicket) => {
      queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
      navigate(
        generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_PRO_GEAR_EDIT_PATH, {
          moveCode,
          shipmentId,
          proGearId: createdProGearWeightTicket?.id,
        }),
        { replace: true },
      );
    },
    onError: () => {
      setErrorMessage(`Failed to create trip record`);
    },
  });

  const { mutate: mutatePatchProGearWeightTicket } = useMutation(patchProGearWeightTicket, {
    onSuccess: () => {
      queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
      navigate(reviewPath);
    },
    onError: () => {
      setIsSubmitted(false);
      setErrorMessage('Failed to save updated trip record');
    },
  });

  // why is this not workinignigningingingindfkbdkjnfdkn
  // const { mutate: mutateUpdateMtoShipment } = useMutation(updateMTOShipment, {
  //   onSuccess: () => {
  //     navigate(reviewPath);
  //   },
  //   onError: (error) => {
  //     setIsSubmitted(false);
  //     setErrorMessage(`${error} Failed to save updated trip record`);
  //   },
  // });

  useEffect(() => {
    if (!proGearId) {
      mutateProGearCreateWeightTicket(ppmShipment?.id);
    }
  }, [mutateProGearCreateWeightTicket, ppmShipment?.id, proGearId]);

  const handleCreateUpload = async (fieldName, file, setFieldTouched) => {
    const documentId = currentProGearWeightTicket[`${fieldName}Id`];
    // Create a date-time stamp in the format "yyyymmddhh24miss"
    const now = new Date();
    const timestamp =
      now.getFullYear().toString() +
      (now.getMonth() + 1).toString().padStart(2, '0') +
      now.getDate().toString().padStart(2, '0') +
      now.getHours().toString().padStart(2, '0') +
      now.getMinutes().toString().padStart(2, '0') +
      now.getSeconds().toString().padStart(2, '0');
    // Create a new filename with the timestamp prepended
    const newFileName = `${file.name}-${timestamp}`;
    // Create and return a new File object with the new filename
    const newFile = new File([file], newFileName, { type: file.type });
    createUploadForPPMDocument(ppmShipment?.id, documentId, newFile, true)
      .then((upload) => {
        documents?.ProGearWeightTickets[currentIndex][fieldName]?.uploads.push(upload);
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
    deleteUploadForDocument(uploadId, null, ppmShipment?.id)
      .then(() => {
        const filteredUploads = documents?.ProGearWeightTickets[currentIndex][fieldName].uploads.filter(
          (upload) => upload.id !== uploadId,
        );
        documents.ProGearWeightTickets[currentIndex][fieldName].uploads = filteredUploads;
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

    const belongsToSelf = values.belongsToSelf === 'true';
    let proGear;
    let spouseProGear;
    if (belongsToSelf) {
      proGear = values.weight;
    }
    if (!belongsToSelf) {
      spouseProGear = values.weight;
    }
    const payload = {
      belongsToSelf,
      description: values.description,
      weight: Number(values.weight),
      ppmShipment: {
        id: ppmShipment.id,
      },
      shipmentType: mtoShipment.shipmentType,
      actualSpouseProGearWeight: parseInt(spouseProGear, 10),
      actualProGearWeight: parseInt(proGear, 10),
      shipmentLocator: values.shipmentLocator,
      eTag: mtoShipment.eTag,
    };

    mutatePatchProGearWeightTicket({
      ppmShipmentId: currentProGearWeightTicket.ppmShipmentId,
      proGearWeightTicketId: currentProGearWeightTicket.id,
      payload,
      eTag: currentProGearWeightTicket.eTag,
    });

    // const moveTaskOrderID = Object.values(orders)?.[0].moveTaskOrderID;

    // let body2;
    // if (belongsToSelf) {
    //   body2 = {
    //     actualProGearWeight: parseInt(proGear, 10),
    //   };
    // } else {
    //   body2 = {
    //     actualSpouseProGearWeight: parseInt(spouseProGear, 10),
    //   };
    // }
    // mutateUpdateMtoShipment({
    //   moveTaskOrderID,
    //   shipmentID: mtoShipment.id,
    //   ifMatchETag: mtoShipment.eTag,
    //   body: body2,
    // });
  };

  // TODO: patchmtoshipment mutateMTOShipment

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
  const entitlements = Object.values(orders)?.[0].entitlement;

  if (isError) return <SomethingWentWrong />;

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
              entitlements={entitlements}
              proGear={currentProGearWeightTicket}
              setNumber={currentIndex + 1}
              onCreateUpload={handleCreateUpload}
              onUploadComplete={handleUploadComplete}
              onUploadDelete={handleUploadDelete}
              onBack={handleBack}
              onSubmit={handleSubmit}
              isSubmitted={isSubmitted}
              appName={appName}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default ProGear;
