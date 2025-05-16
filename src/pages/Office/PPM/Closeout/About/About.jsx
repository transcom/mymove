import React, { useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import ppmPageStyles from 'pages/Office/PPM/PPM.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { shipmentTypes } from 'constants/shipments';
import AboutForm from 'components/Shared/PPM/Closeout/AboutForm/AboutForm';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { servicesCounselingRoutes } from 'constants/routes';
import { usePPMShipmentAndDocsOnlyQueries } from 'hooks/queries';
import { formatDateForSwagger } from 'shared/dates';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { updateMTOShipment } from 'services/ghcApi';
import { MTO_SHIPMENT } from 'constants/queryKeys';
import { APP_NAME } from 'constants/apps';

const About = () => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [isSubmitted, setIsSubmitted] = useState(false);

  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { moveCode, shipmentId } = useParams();
  const { mtoShipment, documents, isLoading, isError } = usePPMShipmentAndDocsOnlyQueries(shipmentId);
  const appName = APP_NAME.OFFICE;

  const { mutate: mutateMTOShipment } = useMutation(updateMTOShipment, {
    onSuccess: (updatedMTOShipment) => {
      queryClient.setQueryData([MTO_SHIPMENT, updatedMTOShipment.moveTaskOrderID, false], updatedMTOShipment);
      queryClient.invalidateQueries([MTO_SHIPMENT, updatedMTOShipment.moveTaskOrderID]);

      let path;
      if (documents?.WeightTickets?.length === 0) {
        path = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
          moveCode,
          shipmentId,
        });
      } else {
        path = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, { moveCode, shipmentId });
      }
      navigate(path);
    },
  });

  if (isError) return <SomethingWentWrong />;

  if (!mtoShipment || isLoading) {
    return <LoadingPlaceholder />;
  }

  const handleBack = () => {
    navigate(-1);
  };

  const handleSetError = (error, defaultError) => {
    if (error?.response?.body?.message !== null && error?.response?.body?.message !== undefined) {
      setErrorMessage(`${error?.response?.body?.message}`);
    } else {
      setErrorMessage(defaultError);
    }
  };

  const handleSubmit = async (values) => {
    setIsSubmitted(true);
    setErrorMessage(null);
    const hasReceivedAdvance = values.hasReceivedAdvance === 'true';
    const payload = {
      shipmentType: mtoShipment.shipmentType,
      ppmShipment: {
        id: mtoShipment.ppmShipment.id,
        actualMoveDate: formatDateForSwagger(values.actualMoveDate),
        pickupAddress: values.pickupAddress,
        hasSecondaryPickupAddress: values.hasSecondaryPickupAddress === 'true',
        secondaryPickupAddress: values.hasSecondaryPickupAddress === 'true' ? values.secondaryPickupAddress : null,
        destinationAddress: values.destinationAddress,
        hasSecondaryDestinationAddress: values.hasSecondaryDestinationAddress === 'true',
        secondaryDestinationAddress:
          values.hasSecondaryDestinationAddress === 'true' ? values.secondaryDestinationAddress : null,
        hasReceivedAdvance,
        advanceAmountReceived: hasReceivedAdvance ? values.advanceAmountReceived * 100 : null,
        w2Address: values.w2Address,
      },
    };
    const updatePPMPayload = {
      moveTaskOrderID: mtoShipment.moveTaskOrderID,
      shipmentID: mtoShipment.id,
      ifMatchETag: mtoShipment.eTag,
      normalize: false,
      body: payload,
    };

    mutateMTOShipment(updatePPMPayload, {
      onError: (error) => {
        setIsSubmitted(false);
        handleSetError(error, `Something went wrong, and your changes were not saved. Please try again.`);
      },
    });
  };

  return (
    <div className={ppmPageStyles.tabContent}>
      <div className={ppmPageStyles.container}>
        <NotificationScrollToTop dependency={errorMessage} />
        <GridContainer className={ppmPageStyles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <div className={ppmPageStyles.closeoutPageWrapper}>
                <ShipmentTag shipmentType={shipmentTypes.PPM} />
                <h1>About your PPM</h1>
                {errorMessage && (
                  <Alert data-testid="errorMessage" type="error" headingLevel="h4" heading="An error occurred">
                    {errorMessage}
                  </Alert>
                )}
                <AboutForm
                  mtoShipment={mtoShipment}
                  onSubmit={handleSubmit}
                  onBack={handleBack}
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

export default About;
