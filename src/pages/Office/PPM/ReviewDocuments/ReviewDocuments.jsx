import React, { useEffect, useRef, useState } from 'react';
import { queryCache } from 'react-query';
import { Button } from '@trussworks/react-uswds';
import { generatePath, useHistory, withRouter } from 'react-router-dom';

import styles from './ReviewDocuments.module.scss';

import { servicesCounselingRoutes } from 'constants/routes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
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

  let documentSet;
  if (weightTickets) {
    weightTickets.sort((a, b) => (a.createdAt < b.createdAt ? -1 : 1));
    documentSet = weightTickets[documentSetIndex];
  }
  const [nextEnabled, setNextEnabled] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const history = useHistory();

  const formRef = useRef();
  // let nextEnabled = false;
  // if (formRef.current?.isValid) {
  //   nextEnabled = true;
  // }

  useEffect(() => {
    // console.log('hi from useEffect in ReviewDocuments');

    // const sortedWeightTickets = weightTickets;
    // sortedWeightTickets.sort((a, b) => (a.createdAt < b.createdAt ? -1 : 1));
    // setDocumentSet(sortedWeightTickets[documentSetIndex]);
    // NB: this setter appears to work correctly, and is not affected by subsequent Formik rerenders:
    setNextEnabled(formRef.current?.isValid);
    // formRef.current?.resetForm();
    // formRef.current?.validateForm();
  }, [formRef, setNextEnabled]);

  // useEffect(() => {
  //   console.log('documentSet', documentSet);
  // }, [documentSet]);

  // useEffect(() => {
  //   console.log('formRef', formRef);
  // }, [formRef]);

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

  const onContinue = () => {
    if (formRef.current) {
      formRef.current.handleSubmit();
    }
  };

  const onValid = (errors) => {
    setNextEnabled(Object.keys(errors).length === 0);
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
        <NotificationScrollToTop dependency={documentSetIndex} />
        <DocumentViewerSidebar.Content>
          {documentSet && (
            <ReviewWeightTicket
              weightTicket={documentSet}
              ppmNumber={1}
              tripNumber={documentSetIndex + 1}
              mtoShipment={mtoShipment}
              onError={onError}
              onSuccess={onSuccess}
              onValid={onValid}
              formRef={formRef}
              setSubmitting={setSubmitting}
            />
          )}
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
