import React from 'react';
import { useSelector } from 'react-redux';
import { useNavigate, useParams } from 'react-router-dom';
import { Button, Grid, GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './Feedback.module.scss';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { selectMTOShipmentById } from 'store/entities/selectors';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { formatCents, formatCentsTruncateWhole, formatCustomerDate, formatWeight } from 'utils/formatters';
import { calculateTotalMovingExpensesAmount, getW2Address } from 'utils/ppmCloseout';
import { FEEDBACK_DOCUMENT_TYPES, FEEDBACK_TEMPLATES } from 'constants/ppmFeedback';
import FeedbackItems from 'components/Customer/PPM/Closeout/FeedbackItems/FeedbackItems';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import {
  calculateTotalNetWeightForProGearWeightTickets,
  getTotalNetWeightForWeightTickets,
} from 'utils/shipmentWeights';

export const GetTripWeight = (doc) => {
  return doc.fullWeight - doc.emptyWeight;
};

export const FormatRow = (row) => {
  const formattedRow = { ...row };
  // format the values
  if (formattedRow.format) {
    formattedRow.value = formattedRow.format(row.value);
    if (row.secondaryValue !== undefined) {
      formattedRow.secondaryValue = formattedRow.format(row.secondaryValue);
    }
  }

  return formattedRow;
};

const Feedback = () => {
  const { mtoShipmentId } = useParams();
  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const navigate = useNavigate();

  const ppmShipment = mtoShipment?.ppmShipment;
  const weightTickets = ppmShipment?.weightTickets;
  const proGearWeightTickets = ppmShipment?.proGearWeightTickets;
  const movingExpenses = ppmShipment?.movingExpenses;

  // track if a document was adjusted
  let docWasAdjusted = false;

  const getTripWeight = (doc) => {
    return GetTripWeight(doc);
  };

  // key into the passed in document and set the value of the new row
  const setRowValues = (doc, templateRow) => {
    const row = { ...templateRow };

    row.value = doc[row.key];
    if (row.key === 'tripWeight') row.value = getTripWeight(doc);

    // set the secondary value/customer submitted value if
    // it differs from the final value, and note that the doc was adjusted
    if (row.secondaryKey && doc[row.secondaryKey] !== undefined && doc[row.secondaryKey] !== row.value) {
      docWasAdjusted = true;
      row.secondaryValue = doc[row.secondaryKey];
    }

    // format that status for display
    if (row.key === 'status') {
      if (docWasAdjusted && row.value === 'APPROVED') row.value = 'EDITED';
      if (row.value === 'REJECTED' || row.value === 'EXCLUDED') {
        row.label = `${row.value}: `;
        row.value = doc.reason;
      }
    }

    return row;
  };

  // format a single document
  const formatSingleDocForFeedbackItem = (doc, docType) => {
    docWasAdjusted = false;

    return FEEDBACK_TEMPLATES[docType]?.map((templateRow) => {
      const row = setRowValues(doc, templateRow);
      return FormatRow(row);
    });
  };

  // format an array of documents
  const formatDocuments = (documentSet, type) => {
    if (!documentSet) return [];
    return documentSet?.map((doc) => {
      return formatSingleDocForFeedbackItem(doc, type);
    });
  };

  // store formatted documents to pass down to child component
  const formattedWeightTickets = formatDocuments(weightTickets, FEEDBACK_DOCUMENT_TYPES.WEIGHT);
  const formattedProGearWeightTickets = formatDocuments(proGearWeightTickets, FEEDBACK_DOCUMENT_TYPES.PRO_GEAR);
  const formattedMovingExpenses = formatDocuments(movingExpenses, FEEDBACK_DOCUMENT_TYPES.MOVING_EXPENSE);

  // calculate total weights/dollars for document sets
  const weightTicketsTotal = getTotalNetWeightForWeightTickets(weightTickets);
  const proGearTotal = calculateTotalNetWeightForProGearWeightTickets(proGearWeightTickets);
  const expensesTotal = calculateTotalMovingExpensesAmount(movingExpenses);

  if (!mtoShipment) return <LoadingPlaceholder />;

  const ppmDetails = (
    <>
      <h2>About Your PPM</h2>
      <div>Departure Date: {formatCustomerDate(ppmShipment?.actualMoveDate)}</div>
      <div>Starting ZIP: {ppmShipment?.actualPickupPostalCode}</div>
      <div>Ending ZIP: {ppmShipment?.actualDestinationPostalCode}</div>
      <div>
        Advance:
        {ppmShipment?.hasReceivedAdvance
          ? ` Yes, $${formatCentsTruncateWhole(ppmShipment?.advanceAmountReceived)}`
          : ' No'}
      </div>
      <br />
      <div data-testid="w-2Address">W-2 address: {getW2Address(ppmShipment?.w2Address)}</div>
    </>
  );

  return (
    <div className={classnames(ppmPageStyles.ppmPageStyle, styles.PPMFeedback)}>
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Closeout Feedback</h1>
            <SectionWrapper className={styles.aboutSection}>{ppmDetails}</SectionWrapper>
            <SectionWrapper>
              <h2>Documents</h2>
              <div className={styles.editedFieldsLabel}>
                *Edited fields will show their previous values in parentheses
              </div>
              <div className={styles.headingContainer}>
                <div className={styles.headingContent}>
                  <h3>Weight Moved</h3>
                  <span>-&nbsp;{formatWeight(weightTicketsTotal)}</span>
                </div>
              </div>
              <FeedbackItems
                className={styles.feedbackItems}
                documents={formattedWeightTickets}
                docType={FEEDBACK_DOCUMENT_TYPES.WEIGHT}
              />
              {proGearWeightTickets.length > 0 && (
                <>
                  <div className={styles.headingContainer} data-testid="pro-gear-items">
                    <div className={styles.headingContent}>
                      <h3>Pro-gear</h3>
                      <span>-&nbsp;{formatWeight(proGearTotal)}</span>
                    </div>
                  </div>
                  <FeedbackItems
                    className={styles.feedbackItems}
                    documents={formattedProGearWeightTickets}
                    docType={FEEDBACK_DOCUMENT_TYPES.PRO_GEAR}
                  />
                </>
              )}
              {movingExpenses.length > 0 && (
                <>
                  <div className={styles.headingContainer} data-testid="expenses-items">
                    <div className={styles.headingContent}>
                      <h3>Expenses</h3>
                      <span>-&nbsp;${expensesTotal ? formatCents(expensesTotal) : 0}</span>
                    </div>
                  </div>
                  <FeedbackItems
                    className={styles.feedbackItems}
                    documents={formattedMovingExpenses}
                    docType={FEEDBACK_DOCUMENT_TYPES.MOVING_EXPENSE}
                  />
                </>
              )}
            </SectionWrapper>
            <div className={classnames(ppmPageStyles.buttonContainer, styles.navigationButtons)}>
              <Button onClick={() => navigate(-1)} type="button">
                Back
              </Button>
            </div>
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default Feedback;
