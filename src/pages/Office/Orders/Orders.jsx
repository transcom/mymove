import React, { useEffect, useReducer, useCallback } from 'react';
import { Link, useNavigate, useParams, useLocation, generatePath } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import ordersFormValidationSchema from './ordersFormValidationSchema';

import styles from 'styles/documentViewerWithSidebar.module.scss';
import { milmoveLogger } from 'utils/milmoveLog';
import { getTacValid, updateOrder } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { tooRoutes, tioRoutes } from 'constants/routes';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import OrdersDetailForm from 'components/Office/OrdersDetailForm/OrdersDetailForm';
import { formatSwaggerDate, dropdownInputOptions } from 'utils/formatters';
import { DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';
import { ORDERS_PAY_GRADE_OPTIONS, ORDERS_TYPE_DETAILS_OPTIONS, ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { ORDERS } from 'constants/queryKeys';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { TAC_VALIDATION_ACTIONS, reducer, initialState } from 'reducers/tacValidation';
import { LOA_TYPE } from 'shared/constants';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

const deptIndicatorDropdownOptions = dropdownInputOptions(DEPARTMENT_INDICATOR_OPTIONS);
const ordersTypeDropdownOptions = dropdownInputOptions(ORDERS_TYPE_OPTIONS);
const ordersTypeDetailsDropdownOptions = dropdownInputOptions(ORDERS_TYPE_DETAILS_OPTIONS);
const payGradeDropdownOptions = dropdownInputOptions(ORDERS_PAY_GRADE_OPTIONS);

const Orders = () => {
  const navigate = useNavigate();
  const { moveCode } = useParams();
  const [tacValidationState, tacValidationDispatch] = useReducer(reducer, null, initialState);

  const { move, orders, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const { state } = useLocation();
  const orderId = move?.ordersId;
  const from = state?.from;
  const handleClose = useCallback(() => {
    let redirectPath;
    if (from === 'paymentRequestDetails') {
      redirectPath = generatePath(tioRoutes.BASE_PAYMENT_REQUESTS_PATH, { moveCode });
    } else {
      redirectPath = generatePath(tooRoutes.BASE_MOVE_VIEW_PATH, { moveCode });
    }
    navigate(redirectPath);
  }, [navigate, moveCode, from]);
  const queryClient = useQueryClient();
  const { mutate: mutateOrders } = useMutation(updateOrder, {
    onSuccess: (data, variables) => {
      const updatedOrder = data.orders[variables.orderID];
      queryClient.setQueryData([ORDERS, variables.orderID], {
        orders: {
          [`${variables.orderID}`]: updatedOrder,
        },
      });
      queryClient.invalidateQueries(ORDERS);
      handleClose();
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  const handleHHGTacValidation = async (value) => {
    if (value && value.length === 4 && value !== tacValidationState[LOA_TYPE.HHG].tac) {
      const response = await getTacValid({ tac: value });
      tacValidationDispatch({
        type: TAC_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        loaType: LOA_TYPE.HHG,
        isValid: response.isValid,
        tac: value,
      });
    }
  };

  const handleNTSTacValidation = async (value) => {
    if (value && value.length === 4 && value !== tacValidationState[LOA_TYPE.NTS].tac) {
      const response = await getTacValid({ tac: value });
      tacValidationDispatch({
        type: TAC_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        loaType: LOA_TYPE.NTS,
        isValid: response.isValid,
        tac: value,
      });
    }
  };

  const order = Object.values(orders)?.[0];
  const { entitlement, uploadedAmendedOrderID, amendedOrdersAcknowledgedAt } = order;
  // TODO - passing in these fields so they don't get unset. Need to rework the endpoint.
  const {
    proGearWeight,
    proGearWeightSpouse,
    requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment,
  } = entitlement;

  useEffect(() => {
    if (isLoading || isError) {
      return;
    }

    const checkHHGTac = async () => {
      const response = await getTacValid({ tac: order.tac });
      tacValidationDispatch({
        type: TAC_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        loaType: LOA_TYPE.HHG,
        isValid: response.isValid,
        tac: order.tac,
      });
    };

    const checkNTSTac = async () => {
      const response = await getTacValid({ tac: order.ntsTac });
      tacValidationDispatch({
        type: TAC_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        loaType: LOA_TYPE.NTS,
        isValid: response.isValid,
        tac: order.ntsTac,
      });
    };

    if (order?.tac && order.tac.length === 4) {
      checkHHGTac();
    }
    if (order?.ntsTac && order.ntsTac.length === 4) {
      checkNTSTac();
    }
  }, [order?.tac, order?.ntsTac, isLoading, isError]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = (values) => {
    const { originDutyLocation, newDutyLocation, ...fields } = values;
    const body = {
      ...fields,
      originDutyLocationId: values.originDutyLocation.id,
      newDutyLocationId: values.newDutyLocation.id,
      issueDate: formatSwaggerDate(values.issueDate),
      reportByDate: formatSwaggerDate(values.reportByDate),
      proGearWeight,
      proGearWeightSpouse,
      requiredMedicalEquipmentWeight,
      organizationalClothingAndIndividualEquipment,
      grade: values.payGrade,
    };

    mutateOrders({ orderID: orderId, ifMatchETag: order.eTag, body });
  };

  const tacWarningMsg =
    'This TAC does not appear in TGET, so it might not be valid. Make sure it matches what‘s on the orders before you continue.';

  const hasAmendedOrders = !!uploadedAmendedOrderID;

  const initialValues = {
    agency: order?.agency,
    originDutyLocation: order?.originDutyLocation,
    newDutyLocation: order?.destinationDutyLocation,
    issueDate: order?.date_issued,
    reportByDate: order?.report_by_date,
    departmentIndicator: order?.department_indicator,
    ordersNumber: order?.order_number || '',
    ordersType: order?.order_type,
    ordersTypeDetail: order?.order_type_detail,
    tac: order?.tac || '',
    sac: order?.sac || '',
    ntsTac: order?.ntsTac,
    ntsSac: order?.ntsSac,
    ordersAcknowledgement: !!amendedOrdersAcknowledgedAt,
    payGrade: order?.grade,
  };

  return (
    <div className={styles.sidebar}>
      <Formik initialValues={initialValues} validationSchema={ordersFormValidationSchema} onSubmit={onSubmit}>
        {(formik) => {
          // onBlur, if the value has 4 digits, run validator and show warning if invalid
          const hhgTacWarning = tacValidationState[LOA_TYPE.HHG].isValid ? '' : tacWarningMsg;
          const ntsTacWarning = tacValidationState[LOA_TYPE.NTS].isValid ? '' : tacWarningMsg;

          return (
            <form onSubmit={formik.handleSubmit}>
              <div className={styles.content}>
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
                  <h2 className={styles.header}>View Orders</h2>
                  <div>
                    {/* the prop "from" represents the page the Link navigates away from.
                        Passing it via state allows the new page to have access to what the previous page was
                        for purposes of redirecting back to the previous page on cancel
                        For documentation please see https://reactrouter.com/en/main/components/link#state
                     */}
                    <Link
                      className={styles.viewAllowances}
                      data-testid="view-allowances"
                      to="../allowances"
                      state={{ from }}
                    >
                      View Allowances
                    </Link>
                  </div>
                </div>
                <div className={styles.body}>
                  <Restricted
                    to={permissionTypes.updateOrders}
                    fallback={
                      <OrdersDetailForm
                        deptIndicatorOptions={deptIndicatorDropdownOptions}
                        ordersTypeOptions={ordersTypeDropdownOptions}
                        ordersTypeDetailOptions={ordersTypeDetailsDropdownOptions}
                        hhgTacWarning={hhgTacWarning}
                        ntsTacWarning={ntsTacWarning}
                        validateHHGTac={handleHHGTacValidation}
                        validateNTSTac={handleNTSTacValidation}
                        showOrdersAcknowledgement={hasAmendedOrders}
                        ordersType={order.order_type}
                        setFieldValue={formik.setFieldValue}
                        payGradeOptions={payGradeDropdownOptions}
                        formIsDisabled
                      />
                    }
                  >
                    <OrdersDetailForm
                      deptIndicatorOptions={deptIndicatorDropdownOptions}
                      ordersTypeOptions={ordersTypeDropdownOptions}
                      ordersTypeDetailOptions={ordersTypeDetailsDropdownOptions}
                      hhgTacWarning={hhgTacWarning}
                      ntsTacWarning={ntsTacWarning}
                      validateHHGTac={handleHHGTacValidation}
                      validateNTSTac={handleNTSTacValidation}
                      showOrdersAcknowledgement={hasAmendedOrders}
                      ordersType={order.order_type}
                      setFieldValue={formik.setFieldValue}
                      payGradeOptions={payGradeDropdownOptions}
                    />
                  </Restricted>
                </div>
                <Restricted to={permissionTypes.updateOrders}>
                  <div className={styles.bottom}>
                    <div className={styles.buttonGroup}>
                      <Button disabled={formik.isSubmitting} type="submit">
                        Save
                      </Button>
                      <Button type="button" secondary onClick={handleClose}>
                        Cancel
                      </Button>
                    </div>
                  </div>
                </Restricted>
              </div>
            </form>
          );
        }}
      </Formik>
    </div>
  );
};

export default Orders;
