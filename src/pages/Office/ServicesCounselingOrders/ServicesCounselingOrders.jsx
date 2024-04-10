/* eslint-disable camelcase */
import React, { useEffect, useReducer, useRef, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import ordersFormValidationSchema from '../Orders/ordersFormValidationSchema';

import styles from 'styles/documentViewerWithSidebar.module.scss';
import { milmoveLogger } from 'utils/milmoveLog';
import OrdersDetailForm from 'components/Office/OrdersDetailForm/OrdersDetailForm';
import { DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';
import { ORDERS_TYPE_DETAILS_OPTIONS, ORDERS_TYPE_OPTIONS, ORDERS_PAY_GRADE_OPTIONS } from 'constants/orders';
import { ORDERS, ORDERS_DOCUMENTS } from 'constants/queryKeys';
import { servicesCounselingRoutes } from 'constants/routes';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { getTacValid, counselingUpdateOrder, createUploadForDocument } from 'services/ghcApi';
import { formatSwaggerDate, dropdownInputOptions } from 'utils/formatters';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { TAC_VALIDATION_ACTIONS, reducer, initialState } from 'reducers/tacValidation';
import { LOA_TYPE } from 'shared/constants';
import FileUpload from 'components/FileUpload/FileUpload';
import Hint from 'components/Hint';

const deptIndicatorDropdownOptions = dropdownInputOptions(DEPARTMENT_INDICATOR_OPTIONS);
const ordersTypeDropdownOptions = dropdownInputOptions(ORDERS_TYPE_OPTIONS);
const ordersTypeDetailsDropdownOptions = dropdownInputOptions(ORDERS_TYPE_DETAILS_OPTIONS);
const payGradeDropdownOptions = dropdownInputOptions(ORDERS_PAY_GRADE_OPTIONS);

const ServicesCounselingOrders = ({ hasDocuments }) => {
  const filePondEl = useRef();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { moveCode } = useParams();
  const [tacValidationState, tacValidationDispatch] = useReducer(reducer, null, initialState);
  const { move, orders, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const [showUpload, setShowUpload] = useState(false);
  const [isDoneButtonDisabled, setIsDoneButtonDisabled] = useState(true);
  const orderId = move?.ordersId;
  const documentId = orders[orderId]?.uploaded_order_id;

  const handleClose = () => {
    navigate(`../${servicesCounselingRoutes.MOVE_VIEW_PATH}`);
  };

  const handleUploadFile = (file) => {
    return createUploadForDocument(file, documentId);
  };

  // enable done button when upload completes
  // will need update when implementing deletion
  const handleChange = () => {
    setIsDoneButtonDisabled(false);
  };

  const toggleUploadVisibility = () => {
    setShowUpload((show) => !show);
  };

  // when the user clicks done, invalidate the query to trigger re render
  // of parent to display uploaded orders and hide the button
  const uploadComplete = () => {
    queryClient.invalidateQueries([ORDERS_DOCUMENTS, documentId]);
    toggleUploadVisibility();
  };

  const { mutate: mutateOrders } = useMutation(counselingUpdateOrder, {
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
    const body = {
      ...values,
      originDutyLocationId: values.originDutyLocation.id,
      newDutyLocationId: values.newDutyLocation.id,
      issueDate: formatSwaggerDate(values.issueDate),
      reportByDate: formatSwaggerDate(values.reportByDate),
      ordersType: values.ordersType,
      grade: values.payGrade,
    };
    mutateOrders({ orderID: orderId, ifMatchETag: order.eTag, body });
  };

  const initialValues = {
    originDutyLocation: order?.originDutyLocation,
    newDutyLocation: order?.destinationDutyLocation,
    issueDate: order?.date_issued,
    reportByDate: order?.report_by_date,
    departmentIndicator: order?.department_indicator,
    ordersType: order?.order_type,
    ordersNumber: order?.order_number || '',
    ordersTypeDetail: order?.order_type_detail,
    tac: order?.tac,
    sac: order?.sac,
    ntsTac: order?.ntsTac,
    ntsSac: order?.ntsSac,
    payGrade: order?.grade,
  };

  const tacWarningMsg =
    'This TAC does not appear in TGET, so it might not be valid. Make sure it matches what‘s on the orders before you continue.';

  return (
    <div className={styles.sidebar}>
      <Formik initialValues={initialValues} validationSchema={ordersFormValidationSchema} onSubmit={onSubmit}>
        {(formik) => {
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
                  <h2 data-testid="view-orders-header" className={styles.header}>
                    View orders
                  </h2>
                  <div>
                    <Link className={styles.viewAllowances} data-testid="view-allowances" to="../allowances">
                      View allowances
                    </Link>
                  </div>
                  {!hasDocuments && !showUpload && <Button onClick={toggleUploadVisibility}>Add Orders</Button>}
                  <div>
                    {showUpload && (
                      <div className={styles.upload}>
                        <FileUpload
                          ref={filePondEl}
                          createUpload={handleUploadFile}
                          onChange={handleChange}
                          labelIdle="Drag files here or click to upload"
                        />
                        <Hint>PDF, JPG, or PNG only. Maximum file size 25MB. Each page must be clear and legible</Hint>
                        <Button disabled={isDoneButtonDisabled} onClick={uploadComplete}>
                          Done
                        </Button>
                      </div>
                    )}
                  </div>
                </div>
                <div className={styles.body}>
                  <OrdersDetailForm
                    deptIndicatorOptions={deptIndicatorDropdownOptions}
                    ordersTypeOptions={ordersTypeDropdownOptions}
                    ordersTypeDetailOptions={ordersTypeDetailsDropdownOptions}
                    ordersType={order.order_type}
                    setFieldValue={formik.setFieldValue}
                    hhgTacWarning={hhgTacWarning}
                    ntsTacWarning={ntsTacWarning}
                    validateHHGTac={handleHHGTacValidation}
                    validateNTSTac={handleNTSTacValidation}
                    payGradeOptions={payGradeDropdownOptions}
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
