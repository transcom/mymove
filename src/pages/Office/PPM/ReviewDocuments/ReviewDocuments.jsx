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

export const ReviewDocuments = ({ match }) => {
  const { shipmentId, moveCode } = match.params;
  const { mtoShipment, documents, isLoading, isError } = usePPMShipmentDocsQueries(shipmentId);

  const [documentSetIndex, setDocumentSetIndex] = useState(0);

  let documentSet = [];
  const allDocuments = [];

  const movingExpenses = documents?.MovingExpenses;
  const weightTickets = documents?.WeightTickets;
  const proGearWeightTickets = documents?.ProGearWeightTickets;

  if (movingExpenses?.length !== 0) {
    allDocuments.push(movingExpenses);
  }
  if (weightTickets?.length !== 0) {
    allDocuments.push(weightTickets);
  }
  if (proGearWeightTickets?.length !== 0) {
    allDocuments.push(proGearWeightTickets);
  }

  const fullDocuments = [];
  allDocuments?.map((docSet) => {
    docSet?.map((doc) => {
      return fullDocuments.push(doc);
    });
    return fullDocuments;
  });

  let uploads = [];
  weightTickets?.forEach((weightTicket) => {
    uploads = uploads.concat(weightTicket.emptyDocument?.uploads);
    uploads = uploads.concat(weightTicket.fullDocument?.uploads);
    uploads = uploads.concat(weightTicket.proofOfTrailerOwnershipDocument?.uploads);
  });

  if (weightTickets) {
    weightTickets.sort((a, b) => (a.createdAt < b.createdAt ? -1 : 1));
    documentSet = documentSet.concat(weightTickets[documentSetIndex]);
  }

  proGearWeightTickets?.forEach((proGearWeightTicket) => {
    uploads = uploads.concat(proGearWeightTicket.document?.uploads);
  });

  if (proGearWeightTickets) {
    proGearWeightTickets.sort((a, b) => (a.createdAt < b.createdAt ? -1 : 1));
    documentSet = documentSet.concat(proGearWeightTickets[documentSetIndex]);
  }

  movingExpenses?.forEach((movingExpense) => {
    uploads = uploads.concat(movingExpense.document?.uploads);
  });

  if (movingExpenses) {
    movingExpenses.sort((a, b) => (a.createdAt < b.createdAt ? -1 : 1));
    documentSet = documentSet.concat(movingExpenses[documentSetIndex]);
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
    if (documentSetIndex < fullDocuments.length - 1) {
      setDocumentSetIndex(documentSetIndex + 1);
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

  return (
    <div data-testid="ReviewDocuments" className={styles.ReviewDocuments}>
      <div className={styles.embed}>
        <DocumentViewer files={uploads} allowDownload />
      </div>
      <DocumentViewerSidebar
        title="Review documents"
        onClose={onClose}
        className={styles.sidebar}
        supertitle={`${documentSetIndex + 1} of ${fullDocuments.length} Document Sets`}
        defaultH3
      >
        <DocumentViewerSidebar.Content mainRef={weightTicketPanelRef}>
          <NotificationScrollToTop dependency={documentSetIndex || serverError} target={weightTicketPanelRef.current} />
          <ErrorMessage display={!!serverError}>{serverError}</ErrorMessage>
          {documentSet &&
            (showOverview ? (
              <ReviewDocumentsSidePanel
                ppmShipment={ppmShipment}
                weightTickets={weightTickets}
                proGearTickets={proGearTickets}
                expenseTickets={expenseTickets}
                onError={onError}
                onSuccess={onConfirmSuccess}
                formRef={formRef}
              />
            ) : (
              <ReviewWeightTicket
                weightTicket={documentSet}
                ppmNumber={1}
                tripNumber={documentSetIndex + 1}
                mtoShipment={mtoShipment}
                onError={onError}
                onSuccess={onSuccess}
                formRef={formRef}
              />
            ))}
        </DocumentViewerSidebar.Content>
        <DocumentViewerSidebar.Footer>
          <Button className="usa-button--secondary" onClick={onBack} disabled={documentSetIndex === 0}>
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
