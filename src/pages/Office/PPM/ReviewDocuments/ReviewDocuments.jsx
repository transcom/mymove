import React, { useRef, useState } from 'react';
import { queryCache } from 'react-query';
import { Button } from '@trussworks/react-uswds';
import { generatePath, useHistory, withRouter } from 'react-router-dom';

import styles from './ReviewDocuments.module.scss';

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
import { WEIGHT_TICKETS } from 'constants/queryKeys';

export const ReviewDocuments = ({ match }) => {
  const { shipmentId, moveCode } = match.params;
  const { mtoShipment, weightTickets, isLoading, isError } = usePPMShipmentDocsQueries(shipmentId);

  const ppmShipment = mtoShipment?.ppmShipment;

  const [documentSetIndex, setDocumentSetIndex] = useState(0);

  let documentSet;

  if (weightTickets) {
    weightTickets.sort((a, b) => (a.createdAt < b.createdAt ? -1 : 1));
    documentSet = weightTickets[documentSetIndex];
  }
  const history = useHistory();

  const formRef = useRef();

  const weightTicketPanelRef = useRef();

  const [serverError, setServerError] = useState(null);

  // placeholder pro-gear tickets & expenses
  // const progearTickets = [];
  // const expenses = [];
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  let uploads = [];
  weightTickets?.forEach((weightTicket) => {
    uploads = uploads.concat(weightTicket.emptyDocument?.uploads);
    uploads = uploads.concat(weightTicket.fullDocument?.uploads);
    uploads = uploads.concat(weightTicket.proofOfTrailerOwnershipDocument?.uploads);
  });

  const onClose = () => {
    history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
  };

  const onBack = () => {
    setServerError(null);
    if (documentSetIndex > 0) {
      setDocumentSetIndex(documentSetIndex - 1);
    }
  };

  const onSuccess = () => {
    queryCache.invalidateQueries([WEIGHT_TICKETS, ppmShipment.id]);
    if (documentSetIndex < weightTickets.length - 1) {
      setDocumentSetIndex(documentSetIndex + 1);
    } else {
      history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
    }
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
        // TODO: set this correctly based on total document sets, including pro gear and expenses
        supertitle={`${documentSetIndex + 1} of ${weightTickets.length} Document Sets`}
        defaultH3
      >
        <DocumentViewerSidebar.Content mainRef={weightTicketPanelRef}>
          <NotificationScrollToTop dependency={documentSetIndex || serverError} target={weightTicketPanelRef.current} />
          <ErrorMessage display={!!serverError}>{serverError}</ErrorMessage>
          {documentSet && (
            <ReviewWeightTicket
              weightTicket={documentSet}
              ppmNumber={1}
              tripNumber={documentSetIndex + 1}
              mtoShipment={mtoShipment}
              onError={onError}
              onSuccess={onSuccess}
              formRef={formRef}
            />
          )}
        </DocumentViewerSidebar.Content>
        <DocumentViewerSidebar.Footer>
          <Button onClick={onBack} disabled={documentSetIndex === 0}>
            Back
          </Button>
          <Button type="submit" onClick={onContinue}>
            Continue
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
