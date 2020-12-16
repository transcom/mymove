/* eslint-disable camelcase */
import React from 'react';
import { Link, useHistory, useParams } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import * as Yup from 'yup';

import moveOrdersStyles from '../MoveOrders/MoveOrders.module.scss';
import AllowancesDetailForm from '../../../components/Office/AllowancesDetailForm/AllowancesDetailForm';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useOrdersDocumentQueries } from 'hooks/queries';

const MoveAllowances = () => {
  const { moveOrderId } = useParams();
  const history = useHistory();

  const { moveOrders, upload, isLoading, isError } = useOrdersDocumentQueries(moveOrderId);

  const handleClose = () => {
    history.push(`/moves/${moveOrderId}/details`);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = () => {
    handleClose();
  };

  const documentsForViewer = Object.values(upload);

  const moveOrder = Object.values(moveOrders)?.[0];

  const { authorizedWeight } = moveOrder.entitlement;

  const initialValues = { authorizedWeight: `${authorizedWeight}` };

  const validationSchema = Yup.object({
    authorizedWeight: Yup.number().min(1, 'Authorized weight must be greater than or equal to 1').required('Required'),
  });

  return (
    <div className={moveOrdersStyles.MoveOrders}>
      {documentsForViewer && (
        <div className={moveOrdersStyles.embed}>
          <DocumentViewer files={documentsForViewer} />
        </div>
      )}
      <div className={moveOrdersStyles.sidebar}>
        <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
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
                  <h2 className={moveOrdersStyles.header} data-testid="allowances-header">
                    View Allowances
                  </h2>
                  <div>
                    <Link className={moveOrdersStyles.viewAllowances} data-testid="view-orders-btn" to="orders">
                      View Orders
                    </Link>
                  </div>
                </div>
                <div className={moveOrdersStyles.body}>
                  <AllowancesDetailForm entitlements={moveOrder.entitlement} />
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

export default MoveAllowances;
