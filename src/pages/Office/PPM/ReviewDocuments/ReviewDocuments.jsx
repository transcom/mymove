import React, { useRef, useState } from 'react';
import { queryCache } from 'react-query';
import { Button } from '@trussworks/react-uswds';
import { generatePath, useHistory, withRouter } from 'react-router-dom';

import styles from './ReviewDocuments.module.scss';

import { servicesCounselingRoutes } from 'constants/routes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { MatchShape } from 'types/router';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import DocumentViewerSidebar from 'pages/Office/DocumentViewerSidebar/DocumentViewerSidebar';
import { usePPMShipmentDocsQueries } from 'hooks/queries';
import ReviewWeightTicket from 'components/Office/PPM/ReviewWeightTicket/ReviewWeightTicket';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';

export const ReviewDocuments = ({ match }) => {
  const { shipmentId, moveCode } = match.params;
  const { mtoShipment, weightTickets, isLoading, isError } = usePPMShipmentDocsQueries(shipmentId);

  const [documentSetIndex, setDocumentSetIndex] = useState(0);
  const [nextEnabled, setNextEnabled] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const history = useHistory();

  const formRef = useRef();

  // placeholder pro-gear tickets & expenses
  const progearTickets = [];
  const expenses = [];
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  let uploads = [];
  weightTickets?.forEach((weightTicket) => {
    uploads = uploads.concat(weightTicket.emptyDocument?.uploads);
    uploads = uploads.concat(weightTicket.fullDocument?.uploads);
    uploads = uploads.concat(weightTicket.proofOfTrailerOwnershipDocument?.uploads);
  });

  // TODO: select the documentSet from among weight tickets, pro gear, and expenses
  const documentSet = weightTickets[documentSetIndex];

  const onClose = () => {
    history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
  };

  const onBack = () => {
    if (documentSetIndex > 0) {
      setDocumentSetIndex(documentSetIndex - 1);
    }
  };

  const onError = (error) => {
    setSubmitting(false);
    const errorMsg = error?.response?.body;
    milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
  };

  const onSuccess = () => {
    setSubmitting(false);
    queryCache.invalidateQueries([], moveCode);
    if (documentSetIndex < weightTickets.length - 1) {
      setDocumentSetIndex(documentSetIndex + 1);
    } else {
      history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
    }
  };

  const onValid = (valid) => {
    setNextEnabled(valid);
  };

  const onContinue = () => {
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
      >
        <DocumentViewerSidebar.Content>
          <ReviewWeightTicket
            weightTicket={documentSet}
            ppmNumber={1}
            tripNumber={1}
            mtoShipment={mtoShipment}
            onError={onError}
            onSuccess={onSuccess}
            onValid={onValid}
            formRef={formRef}
            setSubmitting={setSubmitting}
          />
        </DocumentViewerSidebar.Content>
        <DocumentViewerSidebar.Footer>
          <Button onClick={onBack} disabled={documentSetIndex === 0}>
            Back
          </Button>
          <Button type="submit" onClick={onContinue} disabled={!nextEnabled || submitting}>
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
