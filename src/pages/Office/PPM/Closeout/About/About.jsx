import React, { useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import ppmPageStyles from 'pages/Office/PPM/PPM.module.scss';
import closingPageStyles from 'pages/Office/PPM/Closeout/Closeout.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { shipmentTypes } from 'constants/shipments';
import AboutForm from 'components/Office/PPM/Closeout/AboutForm/AboutForm';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { servicesCounselingRoutes } from 'constants/routes';
import { usePPMShipmentAndDocsOnlyQueries } from 'hooks/queries';
import { formatDateForSwagger } from 'shared/dates';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { updateMTOShipment } from 'services/ghcApi';
import { MTO_SHIPMENT } from 'constants/queryKeys';

const About = () => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [isSubmitted, setIsSubmitted] = useState(false);

  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { moveCode, shipmentId } = useParams();
  const { mtoShipment, isLoading, isError } = usePPMShipmentAndDocsOnlyQueries(shipmentId);

  const { mutate: mutateMTOShipment } = useMutation(updateMTOShipment, {
    onSuccess: (updatedMTOShipment) => {
      queryClient.setQueryData([MTO_SHIPMENT, updatedMTOShipment.moveTaskOrderID, false], updatedMTOShipment);
      queryClient.invalidateQueries([MTO_SHIPMENT, updatedMTOShipment.moveTaskOrderID]);

      const path = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, { moveCode, shipmentId });
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
        actualPickupPostalCode: values.pickupAddress.postalCode,
        actualDestinationPostalCode: values.destinationAddress.postalCode,
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
                <div className={classnames(closingPageStyles['closing-section'], closingPageStyles['about-ppm'])}>
                  <p>Finish moving this PPM before you start documenting it.</p>
                  <h2>How to complete your PPM</h2>
                  <p>To complete your PPM, you will:</p>
                  <ul>
                    <li>Upload weight tickets for each trip</li>
                    <li>Upload receipts to document any expenses</li>
                    <li>Upload receipts if you used short-term storage, so you can request reimbursement</li>
                    <li>
                      Upload any other documentation (such as proof of ownership for a trailer, if you used your own)
                    </li>
                    <li>Complete your PPM to send it to a counselor for review</li>
                  </ul>
                  <h2>About your final payment</h2>
                  <p>Your final payment will be:</p>
                  <ul>
                    <li>based on your final incentive</li>
                    <li>modified by expenses submitted (authorized expenses reduce your tax burden)</li>
                    <li>minus any taxes withheld (the IRS considers your incentive to be taxable income)</li>
                    <li>plus any reimbursements you receive</li>
                  </ul>
                  <p>
                    Verified expenses reduce the taxable income you report to the IRS on form W-2. They may not be
                    claimed again as moving expenses. Federal tax withholding will be deducted from the profit
                    (entitlement less eligible operating expenses.)
                  </p>
                </div>
                <AboutForm
                  mtoShipment={mtoShipment}
                  onSubmit={handleSubmit}
                  onBack={handleBack}
                  isSubmitted={isSubmitted}
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
