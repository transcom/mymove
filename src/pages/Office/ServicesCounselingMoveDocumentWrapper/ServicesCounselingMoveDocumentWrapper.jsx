import React, { useState } from 'react';
import { matchPath, useLocation, useParams } from 'react-router-dom';
import moment from 'moment';

import styles from 'styles/documentViewerWithSidebar.module.scss';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import { servicesCounselingRoutes } from 'constants/routes';
import { useOrdersDocumentQueries, useAmendedDocumentQueries } from 'hooks/queries';
import ServicesCounselingMoveAllowances from 'pages/Office/ServicesCounselingMoveAllowances/ServicesCounselingMoveAllowances';
import ServicesCounselingOrders from 'pages/Office/ServicesCounselingOrders/ServicesCounselingOrders';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { MOVE_DOCUMENT_TYPE } from 'shared/constants';

const ServicesCounselingMoveDocumentWrapper = () => {
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

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const showOrders = matchPath(
    {
      path: servicesCounselingRoutes.BASE_ORDERS_EDIT_PATH,
      exact: true,
    },
    pathname,
  );

  return (
    <div className={styles.DocumentWrapper}>
      {documentsForViewer && (
        <div data-testid="sc-doc-viewer-container" className={styles.embed}>
          <DocumentViewer files={documentsForViewer} allowDownload />
        </div>
      )}
      {showOrders ? (
        <ServicesCounselingOrders
          moveCode={moveCode}
          files={documentsByTypes}
          amendedDocumentId={amendedDocumentId}
          updateAmendedDocument={updateAmendedDocument}
        />
      ) : (
        <ServicesCounselingMoveAllowances moveCode={moveCode} />
      )}
    </div>
  );
};

export default ServicesCounselingMoveDocumentWrapper;
