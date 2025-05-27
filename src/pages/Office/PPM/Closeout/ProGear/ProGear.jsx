import React, { useEffect, useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import { APP_NAME } from 'constants/apps';
import ppmPageStyles from 'pages/Office/PPM/PPM.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import { servicesCounselingRoutes } from 'constants/routes';
import {
  createProGearWeightTicket,
  patchProGearWeightTicket,
  createUploadForPPMDocument,
  deleteUploadForDocument,
} from 'services/ghcApi';
import { DOCUMENTS, MTO_SHIPMENT } from 'constants/queryKeys';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import ProGearForm from 'components/Shared/PPM/Closeout/ProGearForm/ProGearForm';
import { usePPMShipmentAndDocsOnlyQueries, useReviewShipmentWeightsQuery } from 'hooks/queries';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import appendTimestampToFilename from 'utils/fileUpload';

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

  const { mutate: mutatePatchProGearWeightTicket } = useMutation(patchProGearWeightTicket);

  useEffect(() => {
    if (!proGearId) {
      mutateProGearCreateWeightTicket(ppmShipment?.id);
    }
  }, [mutateProGearCreateWeightTicket, ppmShipment?.id, proGearId]);

  const handleCreateUpload = async (fieldName, file, setFieldTouched) => {
    const documentId = currentProGearWeightTicket[`${fieldName}Id`];

    createUploadForPPMDocument(ppmShipment?.id, documentId, appendTimestampToFilename(file), true)
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

  const updateProGearWeightTicket = async (values) => {
    const belongsToSelf = values.belongsToSelf === 'true';
    const hasWeightTickets = !values.missingWeightTicket;

    const payload = {
      hasWeightTickets,
      belongsToSelf,
      ppmShipmentId: mtoShipment.ppmShipment.id,
      shipmentType: mtoShipment.shipmentType,
      shipmentLocator: values.shipmentLocator,
      description: values.description,
      weight: Number(values.weight),
    };

    mutatePatchProGearWeightTicket(
      {
        ppmShipmentId: currentProGearWeightTicket.ppmShipmentId,
        proGearWeightTicketId: currentProGearWeightTicket.id,
        payload,
        eTag: currentProGearWeightTicket.eTag,
      },
      {
        onSuccess: async () => {
          await queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
          await queryClient.invalidateQueries([MTO_SHIPMENT, shipmentId]);
          navigate(reviewPath);
        },
        onError: () => {
          setIsSubmitted(false);
          setErrorMessage('Failed to save updated trip record');
        },
      },
    );
  };

  const handleSubmit = async (values) => {
    if (isSubmitted) return;

    setIsSubmitted(true);
    setErrorMessage(null);
    updateProGearWeightTicket(values);
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
  const entitlements = Object.values(orders)?.[0].entitlement;

  if (isError) return <SomethingWentWrong />;

  if (!mtoShipment || !currentProGearWeightTicket) {
    return renderError() || <LoadingPlaceholder />;
  }
  return (
    <div className={ppmPageStyles.tabContent}>
      <div className={ppmPageStyles.container}>
        <GridContainer>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <div className={ppmPageStyles.closeoutPageWrapper}>
                <ShipmentTag shipmentType={shipmentTypes.PPM} />
                <h1>Pro-gear</h1>
                {renderError()}
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
              </div>
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

export default ProGear;
