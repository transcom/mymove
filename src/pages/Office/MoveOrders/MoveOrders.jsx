/* eslint-disable camelcase */
import React from 'react';
import { withRouter } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import { Formik } from 'formik';
import { queryCache, useMutation } from 'react-query';

import styles from './MoveOrders.module.scss';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import { updateMoveOrder } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import OrdersDetailForm from 'components/Office/OrdersDetailForm/OrdersDetailForm';
import { MatchShape, HistoryShape } from 'types/router';
import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import { dropdownInputOptions, formatSwaggerDate } from 'shared/formatters';
import { DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';
import { ORDERS_TYPE_OPTIONS, ORDERS_TYPE_DETAILS_OPTIONS } from 'constants/orders';
import { MOVE_ORDERS } from 'constants/queryKeys';
import { useOrdersDocumentQueries } from 'hooks/queries';

const deptIndicatorDropdownOptions = dropdownInputOptions(DEPARTMENT_INDICATOR_OPTIONS);
const ordersTypeDropdownOptions = dropdownInputOptions(ORDERS_TYPE_OPTIONS);
const ordersTypeDetailsDropdownOptions = dropdownInputOptions(ORDERS_TYPE_DETAILS_OPTIONS);

const validationSchema = Yup.object({
  originDutyStation: Yup.object().defined('Required'),
  newDutyStation: Yup.object().required('Required'),
  issueDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY').required('Required'),
  reportByDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY').required('Required'),
  departmentIndicator: Yup.string().required('Required'),
  ordersNumber: Yup.string().required('Required'),
  ordersType: Yup.string().required('Required'),
  ordersTypeDetail: Yup.string().required('Required'),
  tac: Yup.string().required('Required'),
  sac: Yup.string().required('Required'),
});

const MoveOrders = ({ history, match }) => {
  const { moveOrderId } = match.params;

  const { moveOrders, upload, isLoading, isError } = useOrdersDocumentQueries(moveOrderId);

  const handleClose = () => {
    history.push(`/moves/${moveOrderId}/details`);
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
      // eslint-disable-next-line no-console
      console.log(errorMsg);
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const moveOrder = Object.values(moveOrders)?.[0];

  const onSubmit = (values) => {
    const { originDutyStation, newDutyStation, ...fields } = values;
    const body = {
      ...fields,
      originDutyStationId: values.originDutyStation.id,
      newDutyStationId: values.newDutyStation.id,
      issueDate: formatSwaggerDate(values.issueDate),
      reportByDate: formatSwaggerDate(values.reportByDate),
    };
    mutateOrders({ moveOrderID: moveOrderId, ifMatchETag: moveOrder.eTag, body });
  };

  const initialValues = {
    originDutyStation: moveOrder?.originDutyStation,
    newDutyStation: moveOrder?.destinationDutyStation,
    issueDate: moveOrder?.date_issued,
    reportByDate: moveOrder?.report_by_date,
    departmentIndicator: moveOrder?.department_indicator,
    ordersNumber: moveOrder?.order_number,
    ordersType: moveOrder?.order_type,
    ordersTypeDetail: moveOrder?.order_type_detail,
    tac: moveOrder?.tac,
    sac: moveOrder?.sac,
  };

  const documentsForViewer = Object.values(upload);

  return (
    <div className={styles.MoveOrders}>
      {documentsForViewer && (
        <div className={styles.embed}>
          <DocumentViewer files={documentsForViewer} />
        </div>
      )}
      <div className={styles.sidebar}>
        <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
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
                    <XLightIcon />
                  </Button>
                  <h2 className={styles.header}>View Orders</h2>
                  <div>
                    <Button type="button" className={styles.viewAllowances} unstyled>
                      View Allowances
                    </Button>
                  </div>
                </div>
                <div className={styles.body}>
                  <OrdersDetailForm
                    deptIndicatorOptions={deptIndicatorDropdownOptions}
                    ordersTypeOptions={ordersTypeDropdownOptions}
                    ordersTypeDetailOptions={ordersTypeDetailsDropdownOptions}
                  />
                </div>
                <div className={styles.bottom}>
                  <div className={styles.buttonGroup}>
                    <Button primary type="submit" disabled={formik.isSubmitting}>
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

MoveOrders.propTypes = {
  history: HistoryShape.isRequired,
  match: MatchShape.isRequired,
};

export default withRouter(MoveOrders);
