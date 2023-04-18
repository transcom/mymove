import React, { useEffect, useRef, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Grid } from '@trussworks/react-uswds';
import { generatePath, useHistory, withRouter } from 'react-router-dom';

import { calculateWeightRequested } from '../../../../hooks/custom';

import styles from './ReviewDocuments.module.scss';

import ReviewDocumentsSidePanel from 'components/Office/PPM/ReviewDocumentsSidePanel/ReviewDocumentsSidePanel';
import { ErrorMessage } from 'components/form';
import { servicesCounselingRoutes } from 'constants/routes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { MatchShape } from 'types/router';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import DocumentViewerSidebar from 'pages/Office/DocumentViewerSidebar/DocumentViewerSidebar';
import { useReviewShipmentWeightsQuery, usePPMShipmentDocsQueries } from 'hooks/queries';
import ReviewWeightTicket from 'components/Office/PPM/ReviewWeightTicket/ReviewWeightTicket';
import ReviewExpense from 'components/Office/PPM/ReviewExpense/ReviewExpense';
import { DOCUMENTS } from 'constants/queryKeys';
import ReviewProGear from 'components/Office/PPM/ReviewProGear/ReviewProGear';

// TODO: This should be in src/constants/ppms.js, but it's causing a lot of errors in unrelated tests, so I'll leave
//  this here for now.
const DOCUMENT_TYPES = {
  WEIGHT_TICKET: 'WEIGHT_TICKET',
  PROGEAR_WEIGHT_TICKET: 'PROGEAR_WEIGHT_TICKET',
  MOVING_EXPENSE: 'MOVING_EXPENSE',
};

export const ReviewDocuments = ({ match }) => {
  const { shipmentId, moveCode } = match.params;
  const { orders, mtoShipments } = useReviewShipmentWeightsQuery(moveCode);
  const { mtoShipment, documents, isLoading, isError } = usePPMShipmentDocsQueries(shipmentId);

  const order = Object.values(orders)?.[0];

  const [documentSetIndex, setDocumentSetIndex] = useState(0);
  const [moveHasExcessWeight, setMoveHasExcessWeight] = useState(false);

  let documentSets = [];
  const weightTickets = documents?.WeightTickets ?? [];
  const proGearWeightTickets = documents?.ProGearWeightTickets ?? [];
  const movingExpenses = documents?.MovingExpenses ?? [];

  const moveWeightTotal = calculateWeightRequested(mtoShipments);
  useEffect(() => {
    setMoveHasExcessWeight(moveWeightTotal > order.entitlement.totalWeight);
  }, [moveWeightTotal, order.entitlement.totalWeight]);

  if (weightTickets.length > 0) {
    weightTickets.sort((a, b) => (a.createdAt < b.createdAt ? -1 : 1));

    documentSets = documentSets.concat(
      weightTickets.map((weightTicket, index) => {
        return {
          documentSetType: DOCUMENT_TYPES.WEIGHT_TICKET,
          documentSet: weightTicket,
          uploads: [
            ...weightTicket.emptyDocument.uploads,
            ...weightTicket.fullDocument.uploads,
            ...weightTicket.proofOfTrailerOwnershipDocument.uploads,
          ],
          tripNumber: index + 1,
        };
      }),
    );
  }

  if (proGearWeightTickets.length > 0) {
    proGearWeightTickets.sort((a, b) => (a.createdAt < b.createdAt ? -1 : 1));

    documentSets = documentSets.concat(
      proGearWeightTickets.map((proGearWeightTicket, index) => {
        return {
          documentSetType: DOCUMENT_TYPES.PROGEAR_WEIGHT_TICKET,
          documentSet: proGearWeightTicket,
          uploads: proGearWeightTicket.document.uploads,
          tripNumber: index + 1,
        };
      }),
    );
  }

  if (movingExpenses.length > 0) {
    movingExpenses.sort((a, b) => (a.createdAt < b.createdAt ? -1 : 1));

    documentSets = documentSets.concat(
      movingExpenses.map((movingExpense, index) => {
        return {
          documentSetType: DOCUMENT_TYPES.MOVING_EXPENSE,
          documentSet: movingExpense,
          uploads: movingExpense.document.uploads,
          tripNumber: index + 1,
        };
      }),
    );
  }

  const history = useHistory();

  const formRef = useRef();
  const mainRef = useRef();

  const [serverError, setServerError] = useState(null);
  const [showOverview, setShowOverview] = useState(false);

  const queryClient = useQueryClient();

  const onClose = () => {
    history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
  };

  const onBack = () => {
    setServerError(null);
    if (showOverview) {
      setShowOverview(false);
    } else if (documentSetIndex > 0) {
      setDocumentSetIndex(documentSetIndex - 1);
    }
  };

  const onSuccess = () => {
    queryClient.invalidateQueries([DOCUMENTS, shipmentId]);

    if (documentSetIndex < documentSets.length - 1) {
      const newDocumentSetIndex = documentSetIndex + 1;
      setDocumentSetIndex(newDocumentSetIndex);
    } else {
      setShowOverview(true);
    }
  };

  const getAllUploads = () => {
    return documentSets.reduce((acc, documentSet) => {
      return acc.concat(documentSet.uploads);
    }, []);
  };

  const onError = () => {
    setServerError('There was an error submitting the form. Please try again later.');
  };

  const onConfirmSuccess = () => {
    history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
  };

  const onContinue = () => {
    setServerError(null);
    if (formRef.current) {
      formRef.current.handleSubmit();
    }
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const currentDocumentSet = documentSets[documentSetIndex];
  const disableBackButton = documentSetIndex === 0 && !showOverview;

  const reviewShipmentWeightsURL = generatePath(servicesCounselingRoutes.REVIEW_SHIPMENT_WEIGHTS_PATH, {
    moveCode,
    shipmentId,
  });

  const reviewShipmentWeightsLink = <a href={reviewShipmentWeightsURL}>Review shipment weights</a>;

  return (
    <div data-testid="ReviewDocuments" className={styles.ReviewDocuments}>
      <div className={styles.embed}>
        <DocumentViewer files={showOverview ? getAllUploads() : currentDocumentSet.uploads} allowDownload />
      </div>
      <DocumentViewerSidebar
        title="Review documents"
        onClose={onClose}
        className={styles.sidebar}
        supertitle={
          showOverview ? 'All Document Sets' : `${documentSetIndex + 1} of ${documentSets.length} Document Sets`
        }
        defaultH3
        hyperlink={reviewShipmentWeightsLink}
      >
        <DocumentViewerSidebar.Content mainRef={mainRef}>
          <NotificationScrollToTop dependency={documentSetIndex || serverError} target={mainRef.current} />
          {moveHasExcessWeight && (
            <Grid className={styles.alertContainer}>
              <Alert headingLevel="h4" slim type="warning">
                <span>This move has excess weight. Edit the PPM net weight to resolve.</span>
              </Alert>
            </Grid>
          )}
          <ErrorMessage display={!!serverError}>{serverError}</ErrorMessage>
          {documentSets &&
            (showOverview ? (
              <ReviewDocumentsSidePanel
                ppmShipment={mtoShipment.ppmShipment}
                weightTickets={weightTickets}
                proGearTickets={proGearWeightTickets}
                expenseTickets={movingExpenses}
                onError={onError}
                onSuccess={onConfirmSuccess}
                formRef={formRef}
              />
            ) : (
              <>
                {currentDocumentSet.documentSetType === DOCUMENT_TYPES.WEIGHT_TICKET && (
                  <ReviewWeightTicket
                    weightTicket={currentDocumentSet.documentSet}
                    ppmNumber={1}
                    tripNumber={currentDocumentSet.tripNumber}
                    mtoShipment={mtoShipment}
                    order={order}
                    mtoShipments={mtoShipments}
                    onError={onError}
                    onSuccess={onSuccess}
                    formRef={formRef}
                  />
                )}
                {currentDocumentSet.documentSetType === DOCUMENT_TYPES.PROGEAR_WEIGHT_TICKET && (
                  <ReviewProGear
                    proGear={currentDocumentSet.documentSet}
                    ppmNumber={1}
                    tripNumber={currentDocumentSet.tripNumber}
                    mtoShipment={mtoShipment}
                    onError={onError}
                    onSuccess={onSuccess}
                    formRef={formRef}
                  />
                )}
                {currentDocumentSet.documentSetType === DOCUMENT_TYPES.MOVING_EXPENSE && (
                  <ReviewExpense
                    expense={currentDocumentSet.documentSet}
                    ppmNumber={1}
                    tripNumber={currentDocumentSet.tripNumber}
                    mtoShipment={mtoShipment}
                    onError={onError}
                    onSuccess={onSuccess}
                    formRef={formRef}
                  />
                )}
              </>
            ))}
        </DocumentViewerSidebar.Content>
        <DocumentViewerSidebar.Footer>
          <Button className="usa-button--secondary" onClick={onBack} disabled={disableBackButton}>
            Back
          </Button>
          <Button type="submit" onClick={onContinue}>
            {showOverview ? 'Confirm' : 'Continue'}
          </Button>
        </DocumentViewerSidebar.Footer>
      </DocumentViewerSidebar>
    </div>
  );
};

ReviewDocuments.propTypes = {
  match: MatchShape.isRequired,
};

export default withRouter(ReviewDocuments);
