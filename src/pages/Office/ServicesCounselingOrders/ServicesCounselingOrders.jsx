/* eslint-disable camelcase */
import React, { useEffect, useReducer } from 'react';
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
import { ORDERS } from 'constants/queryKeys';
import { servicesCounselingRoutes } from 'constants/routes';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { getTacValid, getLoa, counselingUpdateOrder, getOrder } from 'services/ghcApi';
import { formatSwaggerDate, dropdownInputOptions } from 'utils/formatters';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { LineOfAccountingDfasElementOrder } from 'types/lineOfAccounting';
import { LOA_VALIDATION_ACTIONS, reducer as loaReducer, initialState as initialLoaState } from 'reducers/loaValidation';
import { TAC_VALIDATION_ACTIONS, reducer as tacReducer, initialState as initialTacState } from 'reducers/tacValidation';
import { LOA_TYPE, MOVE_DOCUMENT_TYPE } from 'shared/constants';
import DocumentViewerFileManager from 'components/DocumentViewerFileManager/DocumentViewerFileManager';
import { scrollToViewFormikError } from 'utils/validation';

const deptIndicatorDropdownOptions = dropdownInputOptions(DEPARTMENT_INDICATOR_OPTIONS);
const ordersTypeDropdownOptions = dropdownInputOptions(ORDERS_TYPE_OPTIONS);
const ordersTypeDetailsDropdownOptions = dropdownInputOptions(ORDERS_TYPE_DETAILS_OPTIONS);
const payGradeDropdownOptions = dropdownInputOptions(ORDERS_PAY_GRADE_OPTIONS);

const ServicesCounselingOrders = ({ files, amendedDocumentId, updateAmendedDocument }) => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { moveCode } = useParams();
  const [tacValidationState, tacValidationDispatch] = useReducer(tacReducer, null, initialTacState);
  const [loaValidationState, loaValidationDispatch] = useReducer(loaReducer, null, initialLoaState);
  const { move, orders, isLoading, isError } = useOrdersDocumentQueries(moveCode);

  const orderId = move?.ordersId;
  const orderDocumentId = orders[orderId]?.uploaded_order_id;
  const amendedOrderDocumentId = orders[orderId]?.uploadedAmendedOrderID || amendedDocumentId;

  const ordersDocuments = files[MOVE_DOCUMENT_TYPE.ORDERS];
  const amendedDocuments = files[MOVE_DOCUMENT_TYPE.AMENDMENTS];
  const hasOrdersDocuments = ordersDocuments?.length > 0;
  const handleClose = () => {
    navigate(`../${servicesCounselingRoutes.MOVE_VIEW_PATH}`);
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

  const buildFullLineOfAccountingString = (loa) => {
    const dfasMap = LineOfAccountingDfasElementOrder.map((key) => {
      if (key === 'loaEndFyTx') {
        // Specific logic for DFAS element A3, the loaEndFyTx.
        // This is a combination of both the BgFyTx and EndFyTx
        // and if one are null, then typically we would resort to "XXXXXXXX"
        // but for this one we'll just leave it empty as this is not for the EDI858.
        if (loa.loaBgFyTx != null && loa.loaEndFyTx != null) {
          return `${loa.loaBgFyTx}${loa.loaEndFyTx}`;
        }
        if (loa.loaBgFyTx === null || loa.loaByFyTx === undefined) {
          // Catch the scenario of loaBgFyTx being null but loaEndFyTx not being null
          return '';
        }
      }
      return loa[key] || '';
    });
    let longLoa = dfasMap.join('*');
    // remove any number of spaces following an asterisk in a LOA string
    longLoa = longLoa.replace(/\* +/g, '*');
    // remove any number of spaces preceding an asterisk in a LOA string
    longLoa = longLoa.replace(/ +\*/g, '*');

    return longLoa;
  };

  const { mutate: validateLoa } = useMutation(getLoa, {
    onSuccess: (data, variables) => {
      const { loaType } = variables;
      // The server decides if this is a valid LOA or not
      const isValid = (data?.validHhgProgramCodeForLoa ?? false) && (data?.validLoaForTac ?? false);
      // Construct the long line of accounting string
      const longLineOfAccounting = data ? buildFullLineOfAccountingString(data) : '';
      loaValidationDispatch({
        type: LOA_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        payload: {
          loa: data,
          longLineOfAccounting,
          isValid,
          loaType,
        },
      });
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

  const handleHHGLoaValidation = (formikValues) => {
    // LOA is not a field that can be interacted with
    // Validation is based on the scope of the form
    const { tac, issueDate, departmentIndicator } = formikValues;
    if (tac && tac.length === 4 && departmentIndicator && issueDate) {
      // Only run validation if a 4 length TAC is present, and department and issue date are also present
      validateLoa({
        tacCode: tac,
        departmentIndicator,
        effectiveDate: formatSwaggerDate(issueDate),
        loaType: LOA_TYPE.HHG,
      });
    }
  };

  const handleNTSLoaValidation = (formikValues) => {
    // LOA is not a field that can be interacted with
    // Validation is based on the scope of the form
    const { ntsTac, departmentIndicator } = formikValues;
    if (ntsTac && ntsTac.length === 4 && departmentIndicator) {
      // Only run validation if a 4 length NTS TAC and department are present
      // The effective date for an NTS LOA should be either the approved_at date of the
      // move, or the current time of review (Post review it will save as the approved_at)
      const effectiveDate = move?.approvedAt || Date.now();
      validateLoa({
        tacCode: ntsTac,
        departmentIndicator,
        effectiveDate: formatSwaggerDate(effectiveDate),
        loaType: LOA_TYPE.NTS,
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

    const checkHHGLoa = () => {
      // Only run validation if a 4 length TAC is present, and department and issue date are also present
      validateLoa({
        tacCode: order?.tac,
        departmentIndicator: order?.department_indicator,
        effectiveDate: formatSwaggerDate(order?.date_issued),
        loaType: LOA_TYPE.HHG,
      });
    };

    const checkNTSLoa = () => {
      // Only run validation if a 4 length NTS TAC and department are present
      // The effective date for an NTS LOA should be either the approved_at date of the
      // move, or the current time of review (Post review it will save as the approved_at)
      const effectiveDate = move?.approvedAt || Date.now();
      validateLoa({
        tacCode: order?.ntsTac,
        departmentIndicator: order?.department_indicator,
        effectiveDate: formatSwaggerDate(effectiveDate),
        loaType: LOA_TYPE.NTS,
      });
    };

    if (order?.tac && order.tac.length === 4) {
      checkHHGTac();
    }
    if (order?.ntsTac && order.ntsTac.length === 4) {
      checkNTSTac();
    }
    if (order?.tac && order.tac.length === 4 && order?.department_indicator && order?.date_issued) {
      checkHHGLoa();
    }
    if (order?.ntsTac && order?.ntsTac.length === 4 && order?.department_indicator) {
      checkNTSLoa();
    }
  }, [
    order?.tac,
    order?.ntsTac,
    order?.date_issued,
    order?.department_indicator,
    move?.approvedAt,
    isLoading,
    isError,
    validateLoa,
  ]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = async (values) => {
    const orderResponse = await getOrder(null, orderId);
    let newOrderEtag = order.eTag;
    if (orderResponse) {
      newOrderEtag = orderResponse.orders[orderId].eTag;
    }
    const orderBody = {
      ...values,
      originDutyLocationId: values.originDutyLocation.id,
      newDutyLocationId: values.newDutyLocation.id,
      issueDate: formatSwaggerDate(values.issueDate),
      reportByDate: formatSwaggerDate(values.reportByDate),
      ordersType: values.ordersType,
      grade: values.payGrade,
    };
    mutateOrders({ orderID: orderId, ifMatchETag: newOrderEtag, body: orderBody });
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
  const hhgLoaMissingWarningMsg =
    'Unable to find a LOA based on the provided details. Please ensure an orders issue date, department indicator, and TAC are present on this form.';
  const ntsLoaMissingWarningMsg =
    'Unable to find a LOA based on the provided details. Please ensure a department indicator and TAC are present on this form.';
  const loaInvalidWarningMsg = 'The LOA identified based on the provided details appears to be invalid.';

  return (
    <div className={styles.sidebar}>
      <Formik initialValues={initialValues} validationSchema={ordersFormValidationSchema} onSubmit={onSubmit}>
        {(formik) => {
          const hhgTacWarning = tacValidationState[LOA_TYPE.HHG].isValid ? '' : tacWarningMsg;
          const ntsTacWarning = tacValidationState[LOA_TYPE.NTS].isValid ? '' : tacWarningMsg;
          // Conditionally set the LOA warning message based on off if it is missing or just invalid
          const isHHGLoaMissing =
            loaValidationState[LOA_TYPE.HHG].loa === null || loaValidationState[LOA_TYPE.HHG].loa === undefined;
          const isNTSLoaMissing =
            loaValidationState[LOA_TYPE.NTS].loa === null || loaValidationState[LOA_TYPE.NTS].loa === undefined;
          let hhgLoaWarning = '';
          let ntsLoaWarning = '';
          // Making a nested ternary here goes against linter rules
          // The primary warning should be if it is missing, the other warning should be if it is invalid
          if (isHHGLoaMissing) {
            hhgLoaWarning = hhgLoaMissingWarningMsg;
          } else if (!loaValidationState[LOA_TYPE.HHG].isValid) {
            hhgLoaWarning = loaInvalidWarningMsg;
          }
          if (isNTSLoaMissing) {
            ntsLoaWarning = ntsLoaMissingWarningMsg;
          } else if (!loaValidationState[LOA_TYPE.NTS].isValid) {
            ntsLoaWarning = loaInvalidWarningMsg;
          }

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
                  <DocumentViewerFileManager
                    fileUploadRequired={!hasOrdersDocuments}
                    orderId={orderId}
                    documentId={orderDocumentId}
                    files={ordersDocuments}
                    documentType={MOVE_DOCUMENT_TYPE.ORDERS}
                  />
                  <DocumentViewerFileManager
                    orderId={orderId}
                    documentId={amendedOrderDocumentId}
                    files={amendedDocuments}
                    documentType={MOVE_DOCUMENT_TYPE.AMENDMENTS}
                    updateAmendedDocument={updateAmendedDocument}
                  />
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
                    hhgLoaWarning={hhgLoaWarning}
                    ntsLoaWarning={ntsLoaWarning}
                    validateHHGTac={handleHHGTacValidation}
                    validateHHGLoa={() => handleHHGLoaValidation(formik.values)}
                    validateNTSLoa={() => handleNTSLoaValidation(formik.values)}
                    validateNTSTac={handleNTSTacValidation}
                    payGradeOptions={payGradeDropdownOptions}
                    hhgLongLineOfAccounting={loaValidationState[LOA_TYPE.HHG].longLineOfAccounting}
                    ntsLongLineOfAccounting={loaValidationState[LOA_TYPE.NTS].longLineOfAccounting}
                  />
                </div>
                <div className={styles.bottom}>
                  <div className={styles.buttonGroup}>
                    <Button type="submit" disabled={formik.isSubmitting} onClick={scrollToViewFormikError(formik)}>
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
