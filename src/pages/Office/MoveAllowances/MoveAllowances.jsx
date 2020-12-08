/* eslint-disable camelcase */
import React from 'react';
import { withRouter } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from '../MoveOrders/MoveOrders.module.scss';

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
    <div className={styles.MoveOrders}>
      {documentsForViewer && (
        <div className={styles.embed}>
          <DocumentViewer files={documentsForViewer} />
        </div>
      )}
      <div className={styles.sidebar}>
        <Formik initialValues={initialValues} onSubmit={onSubmit}>
          {(formik) => (
            <form onSubmit={formik.handleSubmit}>
              <div className={styles.orderDetails}>
                <div className={styles.top}>
                  <Button
                    className={styles.closeButton}
                    data-testid="closeSidebar"
                    type="button"
                    onClick={handleClose}
                    unstyled
                  >
                    <FontAwesomeIcon icon="times" title="Close sidebar" aria-label="Close sidebar" />
                  </Button>
                  <h2 className={styles.header}>View Allowances</h2>
                  <div>
                    <Button type="button" className={styles.viewAllowances} unstyled>
                      View Orders
                    </Button>
                  </div>
                </div>

                <div className={styles.bottom}>
                  <div className={styles.buttonGroup}>
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
