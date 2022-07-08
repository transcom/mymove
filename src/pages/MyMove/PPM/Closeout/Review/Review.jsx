import React from 'react';
import { Button, GridContainer, Grid } from '@trussworks/react-uswds';
import { Link, useParams } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { generatePath } from 'react-router';
import classnames from 'classnames';

import styles from './Review.module.scss';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ScrollToTop from 'components/ScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { customerRoutes } from 'constants/routes';
import { formatCentsTruncateWhole, formatCustomerDate, formatWeight } from 'utils/formatters';
import { selectMTOShipmentById } from 'store/entities/selectors';
import ReviewItems from 'components/Customer/PPM/Closeout/ReviewItems/ReviewItems';

const Review = () => {
  const { moveId, mtoShipmentId } = useParams();
  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));

  const {
    actualMoveDate,
    actualPickupPostalCode,
    actualDestinationPostalCode,
    hasReceivedAdvance,
    advanceAmountRequested,
  } = mtoShipment?.ppmShipment || {};

  const aboutEditPath = generatePath(customerRoutes.SHIPMENT_PPM_ABOUT_PATH, { moveId, mtoShipmentId });

  const handleAdd = () => {};

  const handleDelete = () => {};

  const aboutYourPPM = [
    {
      rows: [
        { id: 'departureDate', label: 'Departure date:', value: formatCustomerDate(actualMoveDate), hideLabel: true },
        { id: 'startingZIP', label: 'Starting ZIP:', value: actualPickupPostalCode },
        { id: 'endingZIP', label: 'Ending ZIP:', value: actualDestinationPostalCode },
        {
          id: 'advance',
          label: 'Advance:',
          value: hasReceivedAdvance ? `Yes, $${formatCentsTruncateWhole(advanceAmountRequested)}` : 'No',
        },
      ],
      renderEditLink: () => (
        <Link to={aboutEditPath} className="font-body-xs">
          Edit
        </Link>
      ),
    },
  ];

  const weightTickets = [
    {
      subheading: <h4 className="text-bold">Trip 1</h4>,
      rows: [
        { id: 'vehicleDescription-1', label: 'Vehicle description:', value: 'DMC Delorean', hideLabel: true },
        { id: 'emptyWeight-1', label: 'Empty:', value: formatWeight(500) },
        { id: 'fullWeight-1', label: 'Full:', value: formatWeight(1500) },
        {
          id: 'tripWeight-1',
          label: 'Trip Weight:',
          value: formatWeight(1000),
        },
      ],
      onDelete: handleDelete,
      renderEditLink: () => (
        <Link to={aboutEditPath} className="font-body-xs">
          Edit
        </Link>
      ),
    },
  ];

  return (
    <div className={classnames(ppmPageStyles.ppmPageStyle, styles.PPMReview)}>
      <ScrollToTop />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Review</h1>
            <SectionWrapper>
              <ReviewItems heading={<h2>About Your PPM</h2>} contents={aboutYourPPM} />
            </SectionWrapper>
            <SectionWrapper>
              <h2>Documents</h2>
              <ReviewItems
                heading={
                  <>
                    <h3>Weight moved</h3>
                    <span>(1,000 lbs)</span>
                  </>
                }
                contents={weightTickets}
                renderAddButton={() => (
                  <Button type="button" secondary onClick={handleAdd}>
                    Add More Weight
                  </Button>
                )}
              />
            </SectionWrapper>
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default Review;
