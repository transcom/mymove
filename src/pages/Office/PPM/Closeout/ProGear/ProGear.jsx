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
  updateMTOShipment,
} from 'services/ghcApi';
import { DOCUMENTS } from 'constants/queryKeys';
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

  const { mtoShipment, refetchMTOShipment, documents, isError } = usePPMShipmentAndDocsOnlyQueries(shipmentId);
  const { orders } = useReviewShipmentWeightsQuery(moveCode);
  const appName = APP_NAME.OFFICE;
  const ppmShipment = mtoShipment?.ppmShipment;
  const proGearWeightTickets = documents?.ProGearWeightTickets ?? [];
  const moveTaskOrderID = Object.values(orders)?.[0].moveTaskOrderID;

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
  const { mutate: mutateUpdateMtoShipment } = useMutation(updateMTOShipment, {
    onSuccess: () => {
      navigate(reviewPath);
    },
    onError: () => {
      setErrorMessage(`Failed to update shipment record`);
    },
  });

  useEffect(() => {
    if (!proGearId) {
      mutateProGearCreateWeightTicket(ppmShipment?.id);
    }
  }, [mutateProGearCreateWeightTicket, ppmShipment?.id, proGearId]);

  const updateShipment = async (values) => {
    const shipmentResp = await refetchMTOShipment();
    if (shipmentResp.isSuccess) {
      const belongsToSelf = values.belongsToSelf === 'true';
      let proGear;
      let spouseProGear;
      if (belongsToSelf) {
        proGear = values.weight;
      }
      if (!belongsToSelf) {
        spouseProGear = values.weight;
      }

      const shipmentPayload = {
        belongsToSelf,
        ppmShipment: {
          id: mtoShipment.ppmShipment.id,
        },
        shipmentType: mtoShipment.shipmentType,
        actualSpouseProGearWeight: parseInt(spouseProGear, 10),
        actualProGearWeight: parseInt(proGear, 10),
        shipmentLocator: values.shipmentLocator,
        eTag: shipmentResp?.data?.eTag,
      };

      mutateUpdateMtoShipment({
        moveTaskOrderID,
        shipmentID: mtoShipment.id,
        ifMatchETag: shipmentPayload.eTag,
        body: shipmentPayload,
      });
    } else {
      setErrorMessage('Failed to fetch shipment record');
    }
  };

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
        onSuccess: () => {
          queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
          updateShipment(values);
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
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

export default ProGear;
