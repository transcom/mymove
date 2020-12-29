/* eslint-disable camelcase */
import React, { lazy } from 'react';
import PropTypes from 'prop-types';
import { useParams } from 'react-router-dom';

import moveOrdersStyles from '../MoveOrders/MoveOrders.module.scss';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useOrdersDocumentQueries } from 'hooks/queries';

const MoveOrders = lazy(() => import('pages/Office/MoveOrders/MoveOrders'));
const MoveAllowances = lazy(() => import('pages/Office/MoveAllowances/MoveAllowances'));

const MoveDocumentWrapper = (props) => {
  const { moveCode } = useParams();
  const { formName } = props;
  // console.log(formName);

  const { upload, isLoading, isError } = useOrdersDocumentQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const documentsForViewer = Object.values(upload);

  return (
    <div className={moveOrdersStyles.MoveOrders}>
      {documentsForViewer && (
        <div className={moveOrdersStyles.embed}>
          <DocumentViewer files={documentsForViewer} />
        </div>
      )}
      {formName === 'allowances' && <MoveAllowances moveCode={moveCode} />}
      {formName === 'orders' && <MoveOrders moveCode={moveCode} />}
    </div>
  );
};

MoveDocumentWrapper.propTypes = {
  formName: PropTypes.string.isRequired,
};
export default MoveDocumentWrapper;
