import React from 'react';
import { withRouter } from 'react-router-dom';

import styles from './MoveOrders.module.scss';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { MatchShape } from 'types';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const MoveOrders = ({ match }) => {
  const { moveOrderId } = match.params;

  const { upload, isLoading, isError } = useOrdersDocumentQueries(moveOrderId);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const documentsForViewer = Object.values(upload);

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
