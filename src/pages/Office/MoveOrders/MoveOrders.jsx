import React from 'react';
import { withRouter } from 'react-router-dom';

import styles from './MoveOrders.module.scss';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { MatchShape } from 'types';

const MoveOrders = ({ match }) => {
  const { moveOrderId } = match.params;
  const {
    documents,
    // eslint-disable-next-line no-unused-vars
    isLoading,
    // eslint-disable-next-line no-unused-vars
    isError,
  } = useOrdersDocumentQueries(moveOrderId);

  let documentsForViewer;
  if (documents) {
    // eslint-disable-next-line prefer-destructuring
    documentsForViewer = Object.values(documents.undefined);
  }

  return (
    <div className={styles.MoveOrders}>
      {documentsForViewer && (
        <div className={styles.embed}>
          <DocumentViewer files={documentsForViewer} />
        </div>
      )}
      <div className={styles.sidebar}>View orders</div>
    </div>
  );
};

MoveOrders.propTypes = {
  match: MatchShape.isRequired,
};

export default withRouter(MoveOrders);
