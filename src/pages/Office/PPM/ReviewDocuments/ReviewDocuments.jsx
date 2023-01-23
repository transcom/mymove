import React from 'react';
import { withRouter } from 'react-router-dom';
// import { generatePath } from 'react-router'; // need this for close button on side panel

import styles from './ReviewDocuments.module.scss';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { MatchShape } from 'types/router';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import { usePPMShipmentDocsQueries } from 'hooks/queries';
import ReviewDocumentsSidePanel from 'components/Office/PPM/ReviewDocumentsSidePanel/ReviewDocumentsSidePanel';

export const ReviewDocuments = ({ match }) => {
  const { shipmentId } = match.params;
  const { mtoShipment, weightTickets, isLoading, isError } = usePPMShipmentDocsQueries(shipmentId);

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

  return (
    <div data-testid="ReviewDocuments" className={styles.ReviewDocuments}>
      <div className={styles.embed}>
        <DocumentViewer files={uploads} allowDownload />
      </div>
      <div className={styles.sidebar}>
        <ReviewDocumentsSidePanel
          ppmShipment={mtoShipment.ppmShipment}
          weightTickets={weightTickets}
          expenseTickets={expenses}
          proGearTickets={progearTickets}
        />
      </div>
    </div>
  );
};

ReviewDocuments.propTypes = {
  match: MatchShape.isRequired,
};

export default withRouter(ReviewDocuments);
