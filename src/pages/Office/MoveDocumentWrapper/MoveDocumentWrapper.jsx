import React, { useState } from 'react';
import { useParams, matchPath, useLocation } from 'react-router-dom';
import moment from 'moment';

import styles from 'styles/documentViewerWithSidebar.module.scss';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useOrdersDocumentQueries, useAmendedDocumentQueries } from 'hooks/queries';
import Orders from 'pages/Office/Orders/Orders';
import MoveAllowances from 'pages/Office/MoveAllowances/MoveAllowances';
import { MOVE_DOCUMENT_TYPE } from 'shared/constants';

const MoveDocumentWrapper = () => {
  const { moveCode } = useParams();
  const { pathname } = useLocation();

  const { upload, amendedOrderDocumentId, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  // some moves do not have amendedOrderDocumentId created and is null.
  // this is to update the id when it is created to store amendedUpload data.
  const [amendedDocumentId, setAmendedDocumentId] = useState(amendedOrderDocumentId);
  const { amendedUpload } = useAmendedDocumentQueries(amendedDocumentId);

  const updateAmendedDocument = (newId) => {
    setAmendedDocumentId(newId);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const showOrders = matchPath(
    {
      path: '/moves/:moveCode/orders',
      end: true,
    },
    pathname,
  );
  // filter out deleted and sort
  const documentsForViewer = Object.values(upload || {})
    .concat(Object.values(amendedUpload || {}))
    ?.filter((file) => {
      return !file.deletedAt;
    })
    ?.sort((a, b) => moment(b.createdAt) - moment(a.createdAt));

  const ordersFilteredAndSorted = Object.values(upload || {})
    ?.filter((file) => {
      return !file.deletedAt;
    })
    ?.sort((a, b) => moment(b.createdAt) - moment(a.createdAt));
  const amendedFilteredAndSorted = Object.values(amendedUpload || {})
    ?.filter((file) => {
      return !file.deletedAt;
    })
    ?.sort((a, b) => moment(b.createdAt) - moment(a.createdAt));

  const documentsByTypes = {
    [MOVE_DOCUMENT_TYPE.ORDERS]: ordersFilteredAndSorted,
    [MOVE_DOCUMENT_TYPE.AMENDMENTS]: amendedFilteredAndSorted,
  };

  return (
    <div data-testid="doc-wrapper" className={styles.DocumentWrapper}>
      {documentsForViewer && (
        <div className={styles.embed}>
          <DocumentViewer files={documentsForViewer} allowDownload />
        </div>
      )}
      {showOrders ? (
        <Orders
          moveCode={moveCode}
          files={documentsByTypes}
          amendedDocumentId={amendedDocumentId}
          updateAmendedDocument={updateAmendedDocument}
        />
      ) : (
        <MoveAllowances moveCode={moveCode} />
      )}
    </div>
  );
};

export default MoveDocumentWrapper;
