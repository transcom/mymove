import React from 'react';
import { withRouter } from 'react-router-dom';

import styles from './ReviewDocuments.module.scss';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { MatchShape } from 'types/router';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import { useMoveDetailsQueries } from 'hooks/queries';

export const PaymentRequestReview = ({ match }) => {
  const { shipmentId, moveCode } = match.params;
  const { mtoShipments, isLoading, isError } = useMoveDetailsQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { ppmShipment } = mtoShipments.filter((shipment) => shipment.id === shipmentId)[0];
  const uploads = ppmShipment.weightTickets ? ppmShipment.weightTickets[0]?.emptyDocument?.uploads : [];
  console.log(ppmShipment.weightTickets[0]?.emptyDocument?.uploads);
  console.log(ppmShipment);
  // const emptyWeightDocs = ppmShipment.weightTicket ? ppmShipment.weightTicket.emptyDocument.uploads : [];
  // const fullWeightDocs = ppmShipment.weightTicket ? ppmShipment.weightTicket.fullDocument.uploads : [];
  // const proofOfTrailerWeightDocs = ppmShipment.weightTicket
  //   ? ppmShipment.weightTicket.proofOfTrailerOwnershipDocument.uploads
  //   : [];
  // const uploads = [emptyWeightDocs, fullWeightDocs, proofOfTrailerWeightDocs].flatMap((doc) => doc);
  return (
    <div data-testid="PaymentRequestReview" className={styles.PaymentRequestReview}>
      <div className={styles.embed}>
        <DocumentViewer files={uploads} />
      </div>
    </div>
  );
};

PaymentRequestReview.propTypes = {
  match: MatchShape.isRequired,
};

export default withRouter(PaymentRequestReview);
