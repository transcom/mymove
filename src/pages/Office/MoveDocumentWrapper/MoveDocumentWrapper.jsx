import React from 'react';
import { useParams, matchPath, useLocation } from 'react-router-dom';

import moveOrdersStyles from '../Orders/Orders.module.scss';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useOrdersDocumentQueries } from 'hooks/queries';
import MoveOrders from 'pages/Office/Orders/Orders';
import MoveAllowances from 'pages/Office/MoveAllowances/MoveAllowances';

const MoveDocumentWrapper = () => {
  const { moveCode } = useParams();
  const { pathname } = useLocation();

  const { upload, isLoading, isError } = useOrdersDocumentQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const showOrders = matchPath(pathname, {
    path: '/moves/:moveCode/orders',
    exact: true,
  });

  const documentsForViewer = Object.values(upload);

  return (
    <div className={moveOrdersStyles.MoveOrders}>
      {documentsForViewer && (
        <div className={moveOrdersStyles.embed}>
          <DocumentViewer files={documentsForViewer} />
        </div>
      )}
      {showOrders ? <MoveOrders moveCode={moveCode} /> : <MoveAllowances moveCode={moveCode} />}
    </div>
  );
};

export default MoveDocumentWrapper;
