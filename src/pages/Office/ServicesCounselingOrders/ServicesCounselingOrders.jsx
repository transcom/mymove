/* eslint-disable camelcase */
import React from 'react';
import { generatePath } from 'react-router';
import { Link, useHistory, useParams } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import { Formik } from 'formik';
import { queryCache, useMutation } from 'react-query';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from '../ServicesCounselingMoveDocumentWrapper/ServicesCounselingMoveDocumentWrapper.module.scss';

import OrdersDetailForm from 'components/Office/OrdersDetailForm/OrdersDetailForm';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { ORDERS } from 'constants/queryKeys';
import { servicesCounselingRoutes } from 'constants/routes';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { counselingUpdateOrder } from 'services/ghcApi';
import { dropdownInputOptions, formatSwaggerDate } from 'shared/formatters';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const ordersTypeDropdownOptions = dropdownInputOptions(ORDERS_TYPE_OPTIONS);

const validationSchema = Yup.object({
  originDutyStation: Yup.object().defined('Required'),
  newDutyStation: Yup.object().required('Required'),
  issueDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
  reportByDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
  ordersType: Yup.string().required('Required'),
});

const ServicesCounselingOrders = () => {
  const history = useHistory();
  const { moveCode } = useParams();
  const { move, orders, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const orderId = move?.ordersId;

  const handleClose = () => {
    history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
  };

  const [mutateOrders] = useMutation(counselingUpdateOrder, {
    onSuccess: (data, variables) => {
      const updatedOrder = data.orders[variables.orderID];
      queryCache.setQueryData([ORDERS, variables.orderID], {
        orders: {
          [`${variables.orderID}`]: updatedOrder,
        },
      });
      queryCache.invalidateQueries(ORDERS);
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

  const order = Object.values(orders)?.[0];

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = (values) => {
    const body = {
      originDutyStationId: values.originDutyStation.id,
      newDutyStationId: values.newDutyStation.id,
      issueDate: formatSwaggerDate(values.issueDate),
      reportByDate: formatSwaggerDate(values.reportByDate),
      ordersType: values.ordersType,
    };
    mutateOrders({ orderID: orderId, ifMatchETag: order.eTag, body });
  };

  const initialValues = {
    originDutyStation: order?.originDutyStation,
    newDutyStation: order?.destinationDutyStation,
    issueDate: order?.date_issued,
    reportByDate: order?.report_by_date,
    departmentIndicator: order?.department_indicator,
    ordersType: order?.order_type,
  };

  return (
    <div className={styles.sidebar}>
      <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
        {(formik) => {
          return (
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
                  <h2 data-testid="view-orders-header" className={styles.header}>
                    View orders
                  </h2>
                  <div>
                    <Link className={styles.viewAllowances} data-testid="view-allowances" to="allowances">
                      View allowances
                    </Link>
                  </div>
                </div>
                <div className={styles.body}>
                  <OrdersDetailForm
                    ordersTypeOptions={ordersTypeDropdownOptions}
                    showTac={false}
                    showDepartmentIndicator={false}
                    showOrdersNumber={false}
                    showOrdersTypeDetail={false}
                    showSac={false}
                  />
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
          );
        }}
      </Formik>
    </div>
  );
};

export default ServicesCounselingOrders;
