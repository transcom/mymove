import React, { useRef, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Button } from '@trussworks/react-uswds';
import { generatePath, useHistory, withRouter } from 'react-router-dom';

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
import { usePPMShipmentDocsQueries } from 'hooks/queries';
import ReviewWeightTicket from 'components/Office/PPM/ReviewWeightTicket/ReviewWeightTicket';
import { DOCUMENTS } from 'constants/queryKeys';

// TODO: This should be in src/constants/ppms.js, but it's causing a lot of errors in unrelated tests, so I'll leave
//  this here for now.
const DOCUMENT_TYPES = {
  WEIGHT_TICKET: 'WEIGHT_TICKET',
  PROGEAR_WEIGHT_TICKET: 'PROGEAR_WEIGHT_TICKET',
  MOVING_EXPENSE: 'MOVING_EXPENSE',
};

export const ReviewDocuments = ({ match }) => {
  const { shipmentId, moveCode } = match.params;
  const { mtoShipment, documents, isLoading, isError } = usePPMShipmentDocsQueries(shipmentId);

  const [documentSetIndex, setDocumentSetIndex] = useState(0);

  let documentSets = [];

  const weightTickets = documents?.WeightTickets ?? [];
  const proGearWeightTickets = documents?.ProGearWeightTickets ?? [];
  const movingExpenses = documents?.MovingExpenses ?? [];

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

  const weightTicketPanelRef = useRef();

  const [serverError, setServerError] = useState(null);
  const [showOverview, setShowOverview] = useState(false);

  const queryClient = useQueryClient();

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

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

      // TODO: This is a workaround until we add the ability to work with other document types
      if (documentSets[newDocumentSetIndex].documentSetType === DOCUMENT_TYPES.WEIGHT_TICKET) {
        setDocumentSetIndex(newDocumentSetIndex);
      } else {
        setShowOverview(true);
      }
    } else {
      setShowOverview(true);
    }
  };

  const onConfirmSuccess = () => {
    history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
  };

  const onError = () => {
    setServerError('There was an error submitting the form. Please try again later.');
  };

  const onContinue = () => {
    setServerError(null);
    if (formRef.current) {
      formRef.current.handleSubmit();
    }
  };

  const currentDocumentSet = documentSets[documentSetIndex];
  const disableBackButton = documentSetIndex === 0 && !showOverview;

  return (
    <div data-testid="ReviewDocuments" className={styles.ReviewDocuments}>
      <div className={styles.embed}>
        <DocumentViewer files={currentDocumentSet.uploads} allowDownload />
      </div>
      <DocumentViewerSidebar
        title="Review documents"
        onClose={onClose}
        className={styles.sidebar}
        supertitle={`${documentSetIndex + 1} of ${documentSets.length} Document Sets`}
        defaultH3
      >
        <DocumentViewerSidebar.Content mainRef={weightTicketPanelRef}>
          <NotificationScrollToTop dependency={documentSetIndex || serverError} target={weightTicketPanelRef.current} />
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
              currentDocumentSet.documentSetType === DOCUMENT_TYPES.WEIGHT_TICKET && (
                <ReviewWeightTicket
                  weightTicket={currentDocumentSet.documentSet}
                  ppmNumber={1}
                  tripNumber={currentDocumentSet.tripNumber}
                  mtoShipment={mtoShipment}
                  onError={onError}
                  onSuccess={onSuccess}
                  formRef={formRef}
                />
              )
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
