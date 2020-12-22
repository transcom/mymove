/* eslint-disable camelcase */
import React from 'react';
import { Link, useHistory, useParams } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { queryCache, useMutation } from 'react-query';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import * as Yup from 'yup';

import moveOrdersStyles from '../MoveOrders/MoveOrders.module.scss';
import AllowancesDetailForm from '../../../components/Office/AllowancesDetailForm/AllowancesDetailForm';

import { updateMoveOrder } from 'services/ghcApi';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';
import { dropdownInputOptions } from 'shared/formatters';
import { MOVE_ORDERS } from 'constants/queryKeys';

const rankDropdownOptions = dropdownInputOptions(ORDERS_RANK_OPTIONS);

const branchDropdownOption = dropdownInputOptions(ORDERS_BRANCH_OPTIONS);

const validationSchema = Yup.object({
  authorizedWeight: Yup.number().min(1, 'Authorized weight must be greater than or equal to 1').required('Required'),
});

const MoveAllowances = () => {
  const { moveCode } = useParams();
  const history = useHistory();

  const { move, moveOrders, upload, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const moveOrderId = move?.ordersId;

  const handleClose = () => {
    history.push(`/moves/${moveCode}/details`);
  };

  const [mutateOrders] = useMutation(updateMoveOrder, {
    onSuccess: (data, variables) => {
      const updatedOrder = data.moveOrders[variables.moveOrderID];
      queryCache.setQueryData([MOVE_ORDERS, variables.moveOrderID], {
        moveOrders: {
          [`${variables.moveOrderID}`]: updatedOrder,
        },
      });
      queryCache.invalidateQueries(MOVE_ORDERS);
      handleClose();
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      // TODO: Handle error some how
      // RA Summary: eslint: no-console - System Information Leak: External
      // RA: The linter flags any use of console.
      // RA: This console displays an error message from unsuccessful mutation.
      // RA: TODO: As indicated, this error needs to be handled and needs further investigation and work.
      // RA: POAM story here: https://dp3.atlassian.net/browse/MB-5597
      // RA Developer Status: Known Issue
      // RA Validator Status: Known Issue
      // RA Modified Severity: CAT II
      // eslint-disable-next-line no-console
      console.log(errorMsg);
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const moveOrder = Object.values(moveOrders)?.[0];
  const onSubmit = (values) => {
    const { grade, authorizedWeight, dependentsAuthorized } = values;
    const body = {
      issueDate: moveOrder.date_issued,
      newDutyStationId: moveOrder.destinationDutyStation.id,
      ordersNumber: moveOrder.order_number,
      ordersType: moveOrder.order_type,
      originDutyStationId: moveOrder.originDutyStation.id,
      reportByDate: moveOrder.report_by_date,
      grade,
      authorizedWeight: Number(authorizedWeight),
      dependentsAuthorized,
    };
    mutateOrders({ moveOrderID: moveOrderId, ifMatchETag: moveOrder.eTag, body });
  };

  const documentsForViewer = Object.values(upload);

  const { entitlement, grade } = moveOrder;
  const { authorizedWeight, dependentsAuthorized } = entitlement;

  const initialValues = { authorizedWeight: `${authorizedWeight}`, grade, dependentsAuthorized };
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
                    <Link className={moveOrdersStyles.viewAllowances} data-testid="view-orders" to="orders">
                      View Orders
                    </Link>
                  </div>
                </div>
                <div className={moveOrdersStyles.body}>
                  <AllowancesDetailForm
                    entitlements={moveOrder.entitlement}
                    rankOptions={rankDropdownOptions}
                    branchOptions={branchDropdownOption}
                  />
                </div>
                <div className={moveOrdersStyles.bottom}>
                  <div className={moveOrdersStyles.buttonGroup}>
                    <Button disabled={formik.isSubmitting} type="submit">
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
