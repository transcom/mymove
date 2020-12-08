/* eslint-disable camelcase */
import React from 'react';
import { withRouter } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import moveOrdersStyles from '../MoveOrders/MoveOrders.module.scss';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { HistoryShape, MatchShape } from 'types/router';
import { useOrdersDocumentQueries } from 'hooks/queries';

const MoveAllowances = ({ history, match }) => {
  const { moveOrderId } = match.params;

  const { upload, isLoading, isError } = useOrdersDocumentQueries(moveOrderId);

  const handleClose = () => {
    history.push(`/moves/${moveOrderId}/details`);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = () => {
    handleClose();
  };

  const initialValues = {};

  const documentsForViewer = Object.values(upload);

  return (
    <div className={moveOrdersStyles.MoveOrders}>
      {documentsForViewer && (
        <div className={moveOrdersStyles.embed}>
          <DocumentViewer files={documentsForViewer} />
        </div>
      )}
      <div className={moveOrdersStyles.sidebar}>
        <Formik initialValues={initialValues} onSubmit={onSubmit}>
          {(formik) => (
            <form onSubmit={formik.handleSubmit}>
              <div className={moveOrdersStyles.orderDetails}>
                <div className={moveOrdersStyles.top}>
                  <Button
                    className={moveOrdersStyles.closeButton}
                    data-testid="closeSidebar"
                    type="button"
                    onClick={handleClose}
                    unstyled
                  >
                    <FontAwesomeIcon icon="times" title="Close sidebar" aria-label="Close sidebar" />
                  </Button>
                  <h2 className={moveOrdersStyles.header}>View Allowances</h2>
                  <div>
                    <Button type="button" className={moveOrdersStyles.viewAllowances} unstyled>
                      View Orders
                    </Button>
                  </div>
                </div>

                <div className={moveOrdersStyles.bottom}>
                  <div className={moveOrdersStyles.buttonGroup}>
                    <Button type="submit" disabled={formik.isSubmitting}>
                      Save
                    </Button>
                    <Button type="button" secondary onClick={handleClose}>
                      Cancel
                    </Button>
                  </div>
                </div>
              </div>
            </form>
          )}
        </Formik>
      </div>
    </div>
  );
};

MoveAllowances.propTypes = {
  history: HistoryShape.isRequired,
  match: MatchShape.isRequired,
};

export default withRouter(MoveAllowances);
