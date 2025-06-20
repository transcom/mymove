import React, { useEffect, useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import ppmPageStyles from 'pages/Office/PPM/PPM.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import { servicesCounselingRoutes } from 'constants/routes';
import {
  createGunSafeWeightTicket,
  patchGunSafeWeightTicket,
  createUploadForPPMDocument,
  deleteUploadForDocument,
  updateMTOShipment,
} from 'services/ghcApi';
import { DOCUMENTS } from 'constants/queryKeys';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import GunSafeForm from 'components/Shared/PPM/Closeout/GunSafeForm/GunSafeForm';
import { usePPMShipmentAndDocsOnlyQueries, useReviewShipmentWeightsQuery } from 'hooks/queries';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import appendTimestampToFilename from 'utils/fileUpload';
import NotificationScrollToTop from 'components/NotificationScrollToTop';

const GunSafe = ({ appName }) => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // moveId, mtoShipmentId and gunSafeId used for customer side, shipmentId and gunSafeId for office
  const { shipmentId, gunSafeId } = useParams();

  const { mtoShipment, refetchMTOShipment, documents, isError } = usePPMShipmentAndDocsOnlyQueries(shipmentId);
  const { moveCode } = useParams();
  const { orders } = useReviewShipmentWeightsQuery(moveCode);
  const ppmShipment = mtoShipment?.ppmShipment;
  const gunSafeWeightTickets = documents?.GunSafeWeightTickets ?? [];
  const moveTaskOrderID = Object.values(orders)?.[0].moveTaskOrderID;

  const currentGunSafeWeightTicket = gunSafeWeightTickets?.find((item) => item.id === gunSafeId) ?? null;
  const currentIndex = Array.isArray(gunSafeWeightTickets)
    ? gunSafeWeightTickets.findIndex((ele) => ele.id === gunSafeId)
    : -1;

  const reviewPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, { moveCode, shipmentId });

  const { mutate: mutateGunSafeCreateWeightTicket } = useMutation(createGunSafeWeightTicket, {
    onSuccess: (createdGunSafeWeightTicket) => {
      queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
      navigate(
        generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_GUN_SAFE_EDIT_PATH, {
          moveCode,
          shipmentId,
          gunSafeId: createdGunSafeWeightTicket?.id,
        }),
        { replace: true },
      );
    },
    onError: () => {
      setErrorMessage(`Failed to create trip record`);
    },
  });

  const { mutate: mutatePatchGunSafeWeightTicket } = useMutation(patchGunSafeWeightTicket, {
    onSuccess: async () => {
      await queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
      navigate(reviewPath);
    },
    onError: () => {
      setIsSubmitted(false);
      setErrorMessage('Failed to save updated trip record');
    },
  });

  const { mutate: mutateUpdateMtoShipment } = useMutation(updateMTOShipment, {
    onSuccess: () => {
      navigate(reviewPath);
    },
    onError: () => {
      setErrorMessage(`Failed to update shipment record`);
    },
  });

  useEffect(() => {
    if (!gunSafeId && ppmShipment?.id) {
      mutateGunSafeCreateWeightTicket(ppmShipment?.id);
    }
  }, [mutateGunSafeCreateWeightTicket, ppmShipment?.id, gunSafeId]);

  const updateShipment = async (values) => {
    const shipmentResp = await refetchMTOShipment();
    if (shipmentResp.isSuccess) {
      const shipmentPayload = {
        ppmShipment: {
          id: mtoShipment.ppmShipment.id,
        },
        shipmentType: mtoShipment.shipmentType,
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
    const { documentId } = currentGunSafeWeightTicket;

    createUploadForPPMDocument(ppmShipment?.id, documentId, appendTimestampToFilename(file), true)
      .then((upload) => {
        documents?.GunSafeWeightTickets[currentIndex].document?.uploads.push(upload);
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
        const filteredUploads = documents?.GunSafeWeightTickets[currentIndex].document?.uploads.filter(
          (upload) => upload.id !== uploadId,
        );
        documents.GunSafeWeightTickets[currentIndex].document.uploads = filteredUploads;
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

  const updateGunSafeWeightTicket = async (values) => {
    const hasWeightTickets = !values.missingWeightTicket;

    const payload = {
      hasWeightTickets,
      ppmShipmentId: mtoShipment.ppmShipment.id,
      shipmentType: mtoShipment.shipmentType,
      shipmentLocator: values.shipmentLocator,
      description: values.description,
      weight: Number(values.weight),
    };

    mutatePatchGunSafeWeightTicket(
      {
        ppmShipmentId: currentGunSafeWeightTicket.ppmShipmentId,
        gunSafeWeightTicketId: currentGunSafeWeightTicket.id,
        payload,
        eTag: currentGunSafeWeightTicket.eTag,
      },
      {
        onSuccess: () => {
          queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
          updateShipment(values);
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
    updateGunSafeWeightTicket(values);
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

  if (!mtoShipment || !currentGunSafeWeightTicket) {
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
                <h1>Gun Safe</h1>
                {renderError()}
                <GunSafeForm
                  entitlements={entitlements}
                  gunSafe={currentGunSafeWeightTicket}
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

export default GunSafe;
